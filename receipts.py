#!/usr/bin/env python

import os
import sys
import json
import hashlib
import logging
import datetime
import subprocess
from functools import partial
from concurrent.futures import ThreadPoolExecutor

import peewee
from peewee import (
    Model,
    CompositeKey,
    DateTimeField,
    DecimalField,
    ForeignKeyField,
    PrimaryKeyField,
    TextField,
    JOIN_LEFT_OUTER,
)
from playhouse.sqlite_ext import FTSModel, SqliteExtDatabase

import tornado.ioloop
import tornado.iostream
import tornado.web
from tornado import gen
from tornado.options import define, options

import pyinsane.abstract as pyinsane


###############################################################################
## Configuration
APP_ROOT = os.path.dirname(os.path.realpath(__file__))
DATABASE = os.path.join(APP_ROOT, 'receipts.db')

db = SqliteExtDatabase(DATABASE, threadlocals=True)
log = logging.getLogger(__name__)
scan_pool = ThreadPoolExecutor(max_workers=1)
ocr_pool  = ThreadPoolExecutor(max_workers=10)

define("debug", default=False, help="run the application in debug mode")
define("device", default=None, help="name of the device to use as a scanner")
define("storage", default=None, help="location to store scanned receipts in")


###############################################################################
## Utils

# Taken from flask-peewee
def get_dictionary_from_model(model, fields=None, exclude=None):
    model_class = type(model)
    data = {}

    fields = fields or {}
    exclude = exclude or {}
    curr_exclude = exclude.get(model_class, [])
    curr_fields = fields.get(model_class, model._meta.get_field_names())

    for field_name in curr_fields:
        if field_name in curr_exclude:
            continue
        field_obj = model_class._meta.fields[field_name]
        field_data = model._data.get(field_name)
        if isinstance(field_obj, ForeignKeyField) and field_data and field_obj.rel_model in fields:
            rel_obj = getattr(model, field_name)
            data[field_name] = get_dictionary_from_model(rel_obj, fields, exclude)
        else:
            data[field_name] = field_data

    return data


class Serializer(object):
    date_format = '%Y-%m-%d'
    time_format = '%H:%M:%S'
    datetime_format = ' '.join([date_format, time_format])

    def convert_value(self, value):
        if isinstance(value, datetime.datetime):
            return value.strftime(self.datetime_format)
        elif isinstance(value, datetime.date):
            return value.strftime(self.date_format)
        elif isinstance(value, datetime.time):
            return value.strftime(self.time_format)
        elif isinstance(value, Model):
            # TODO: keep track of nested models, sideload them
            return value.get_id()
        else:
            return value

    def clean_data(self, data):
        for key, value in data.items():
            if isinstance(value, dict):
                self.clean_data(value)
            elif isinstance(value, (list, tuple)):
                data[key] = map(self.clean_data, value)
            else:
                data[key] = self.convert_value(value)
        return data

    def serialize_object(self, obj, fields=None, exclude=None):
        data = get_dictionary_from_model(obj, fields, exclude)
        return self.clean_data(data)

SERIALIZE = Serializer()


# Tesseract wrapper - tries OCRing an image.  Returns None if something failed,
# the OCR-d text otherwise.
def ocr_image(image_file, language='eng'):
    proc = subprocess.Popen(['tesseract', 'stdin', 'stdout', '-l', language],
                            stdin=image_file,
                            stdout=subprocess.PIPE,
                            stderr=subprocess.PIPE
                            )

    try:
        stdout, stderr = proc.communicate(timeout=15)
    except TimeoutExpired:
        log.warn("OCR process timed out - killing...")
        proc.kill()
        stdout, stderr = proc.communicate()

    if proc.returncode != 0:
        log.warn("OCR process exited with return code %d", proc.returncode)
        return None

    # TODO: convert from bytes to string?
    return stdout


###############################################################################
## Models
class BaseModel(Model):
    class Meta:
        database = db


class Image(BaseModel):
    """The model representing an uploaded or scanned image"""
    id       = PrimaryKeyField()
    hash     = TextField(null=False)
    created  = DateTimeField(null=False, default=datetime.datetime.now)

    @classmethod
    def get_expired(self):
        """
        Return all images that are more than 24 hours old and are not
        associated with a receipt.
        """
        expiry_time = datetime.timedelta(hours=24)
        abs_expiry = datetime.datetime.now() - expiry_time
        return (
            Image.select().
                where(Image.created < abs_expiry).
                join(Receipt, JOIN_LEFT_OUTER).
                where(Receipt.image >> None)
        )


class Receipt(BaseModel):
    """The model representing a single receipt"""
    id       = PrimaryKeyField()
    amount   = DecimalField(null=False)
    ocr_data = TextField(null=False)
    image    = ForeignKeyField(Image, null=False)
    created  = DateTimeField(null=False, default=datetime.datetime.now)


class FTSReceipt(FTSModel):
    """The full-text search model for our Receipt"""
    receipt = ForeignKeyField(Receipt, primary_key=True)
    content = TextField()

    class Meta: # Fix highlighting ):
        database = db

    @classmethod
    def store_note(cls, receipt):
        try:
            fts_receipt = FTSReceipt.get(FTSReceipt.receipt == receipt)
        except FTSReceipt.DoesNotExist:
            fts_receipt = FTSReceipt(receipt=receipt)
            force_insert = True
        else:
            force_insert = False

        fts_receipt.content = receipt.ocr_data
        fts_receipt.save(force_insert=force_insert)


class Tag(BaseModel):
    """The model representing a tag."""
    id = PrimaryKeyField()
    name = TextField(null=False, unique=True)


class ReceiptToTag(BaseModel):
    """A simple "through" table for many-to-many relationship."""
    receipt = ForeignKeyField(Receipt, null=False)
    tag = ForeignKeyField(Tag, null=False)

    class Meta: # Fix highlighting ):
        primary_key = CompositeKey('receipt', 'tag')


###############################################################################
## Handlers
class BaseHandler(tornado.web.RequestHandler):
    def write_error(self, status_code, **kwargs):
        """
        Write an error message in JSON format.  This is called by the
        send_error() function.
        """
        if 'message' not in kwargs:
            if status_code == 405:
                kwargs['message'] = 'Invalid HTTP method'
            else:
                kwargs['message'] = 'Unknown error'

        kwargs["error"] = True
        self.json(kwargs)

    def serialize(self, name, obj, *, force_serialize=False):
        """
        Serialize a peewee model or query to the appropriate format,
        as expected by Ember-Data.
        """
        if isinstance(obj, Model):
            # Single object - should be serialized at the root.
            obj = {name: SERIALIZE.serialize_object(obj)}

        elif (isinstance(obj, peewee.Query) or
              isinstance(obj, (list, tuple)) and force_serialize):
            # Set of objects, should be serialized in the format:
            #   {
            #       "models": [ list of models ]
            #   }
            obj = {name + 's': [SERIALIZE.serialize_object(x) for x in obj]}

        else:
            raise RuntimeError("unknown type to serialize: %r" % (obj,))

        return self.json(obj)

    def json(self, obj):
        """
        Write an object out in JSON format, setting the correct Content-Type
        """
        self.set_header('Content-Type', 'application/json; charset=utf-8')
        self.set_header('Access-Control-Allow-Origin', '*')
        self.write(json.dumps(obj))


class MainHandler(BaseHandler):
    def get(self):
        self.write("Hello, world")


class TagsHandler(BaseHandler):
    """
    Handles the routes:
        GET     /api/tags
        POST    /api/tags
    """
    def get(self):
        self.serialize('tag', Tag.select())

    def post(self):
        params = json.loads(self.request.body.decode('latin1'))
        try:
            with db.transaction():
                t = Tag.create(name=params['name'])
        except peewee.IntegrityError:
            return self.send_error(409, message="tag already exists")

        self.set_status(201)
        self.serialize('tag', t)


class TagHandler(BaseHandler):
    """
    Handles the routes:
        GET     /api/tags/<id>
        DELETE  /api/tags/<id>
    """
    def get(self, tag_id):
        try:
            tag = Tag.get(Tag.id == tag_id)
        except Tag.DoesNotExist:
            return self.send_error(404, message="tag does not exist")

        self.serialize('tag', tag)

    def delete(self, tag_id):
        try:
            tag = Tag.get(Tag.id == tag_id)
        except Tag.DoesNotExist:
            return self.send_error(404, message="tag does not exist")

        tag.delete_instance()
        self.json(True)


class ReceiptsHandler(BaseHandler):
    """
    Handles the routes:
        GET     /api/receipts
        POST    /api/receipts
    """
    def get(self):
        self.serialize('receipt', Receipt.select())

    def post(self):
        params = json.loads(self.request.body.decode('latin1'))
        try:
            # TODO: real parameters here
            with db.transaction():
                t = Receipt.create(
                    path=params['path'],
                    amount=params['amount'],
                )
        except peewee.IntegrityError:
            return self.send_error(409, message="receipt already exists")

        self.set_status(201)
        self.serialize('receipt', t)


class ReceiptHandler(BaseHandler):
    """
    Handles the routes:
        GET     /api/receipts/<id>
        DELETE  /api/receipts/<id>
    """
    def get(self, receipt_id):
        try:
            receipt = Receipt.get(Receipt.id == receipt_id)
        except Receipt.DoesNotExist:
            return self.send_error(404, message="receipt does not exist")

        self.serialize('receipt', receipt)

    def delete(self, receipt_id):
        try:
            receipt = Receipt.get(Receipt.id == receipt_id)
        except Receipt.DoesNotExist:
            return self.send_error(404, message="receipt does not exist")

        receipt.delete_instance()
        self.json(True)


class ImagesHandler(BaseHandler):
    """
    Handles the routes:
        GET     /api/images
    """
    def get(self):
        self.serialize('image', Image.select())


class ImageHandler(BaseHandler):
    """
    Handles the routes:
        GET     /api/images/<id>
    """
    def get(self, image_id):
        try:
            image = Image.get(Image.id == image_id)
        except Image.DoesNotExist:
            return self.send_error(404, message="image does not exist")

        self.serialize('image', image)


class ImageDataHandler(BaseHandler):
    """
    Handles the routes:
        GET     /api/images/<id>/data
    """
    @gen.coroutine
    def get(self, image_id):
        try:
            image = Image.get(Image.id == image_id)
        except Image.DoesNotExist:
            return self.send_error(404, message="image does not exist")

        # Get the path of the file
        fpath = os.path.join(
            self.application.settings['storage'],
            image.hash[0:2],
            image.hash + '.png',
        )

        try:
            size = os.path.getsize(fpath)
        except OSError as e:
            return self.send_error(500, message="could not get size of image")

        self.set_header('Content-Length', size)
        self.set_header('Content-Type', 'image/png')

        with open(fpath, 'rb') as f:
            while True:
                data = f.read(8192)
                if len(data) == 0:
                    break

                try:
                    self.write(data)
                    yield self.flush()
                except tornado.iostream.StreamClosedError:
                    return


class ScanImageHandler(BaseHandler):
    """
    Handles the routes:
        POST    /api/images/scan
    """
    @tornado.web.asynchronous
    def post(self):
        if self.application.settings['device'] is None:
            return self.send_error(501, message='no scanner is available')

        # Submit the request to the threadpool
        args = ()
        kwargs = {}
        job = scan_pool.submit(partial(self.scan_image, *args, **kwargs))

        # When we're finished, call back onto the IO loop
        job.add_done_callback(
            lambda future: tornado.ioloop.IOLoop.instance().add_callback(
                partial(self.scan_finished, future)))

    def scan_image(self):
        device = self.application.settings['device']
        scan_session = device.scan(multiple=False)

        try:
            while True:
                scan_session.scan.read()
        except EOFError:
            pass

        # Return the image to the caller.
        return scan_session.images[0]

    def scan_finished(self, future):
        # Get and hash the image - we use the SHA256 hash as a key
        image = future.result()
        hash = hashlib.sha256(image.tostring()).hexdigest()

        log.debug("got image of size %r with hash %s", image.size, hash)

        # Create an entry in the database.
        # Note: do this before we save the image to disk, so if it fails, then
        # we don't have leftover files lying around.
        with db.transaction():
            dbImg = Image.create(hash=hash)

        # Build the storage path
        dirname = os.path.join(
            self.application.settings['storage'],
            hash[0:2],
        )

        # Ensure the directory exists
        # Note: before 3.4.1 this will raise if the mode isn't the same.
        os.makedirs(dirname, mode=0o700, exist_ok=True)

        # Save the file as PNG.
        image.save(os.path.join(dirname, hash + '.png'))

        # Return the newly-created image object.
        self.serialize('image', dbImg)
        self.finish()


# TODO:
#   - ScannedImage model
#       - Clean up after 24 hours have passed (sched module?)
#       - OCR the input receipt
#   - Allow uploading a document of some sort
#       - If we want to allow OCRing PDFs, then we need to convert them to
#         something that tesseract can read (e.g. TIFF)
#   - Create thumbnails?
#   - Strip black space on the sides of the scanned image?
#       - Should have a script for this around somewhere...
#   - Allow the user to modify tags on a receipt


def check_dependencies():
    # TODO: check for tesseract
    pass


################################################################################
## Main function
def main():
    tornado.options.parse_command_line()

    if options.storage is None:
        log.error("no storage directory given")
        return
    storage = os.path.realpath(os.path.abspath(options.storage))
    log.info("using storage path: %s", storage)

    check_dependencies()

    if options.debug:
        log.debug("scanning for devices")
        devices = pyinsane.get_devices()
        for i, dev in enumerate(devices):
            log.info("found device %d: %s/%s/%s",
                     i, dev.name, dev.vendor, dev.model)

    if options.device is None:
        log.warn("no device specified, scanning will be unavailable")
        device = None
    else:
        device = pyinsane.Scanner(name=options.device)
        log.debug("using device: name = %s, vendor = %s, model = %s",
                  device.name, device.vendor, device.model)

    settings = {
        'debug': options.debug,
        'device': device,
        'storage': storage,
    }

    application = tornado.web.Application([
        (r"/",                      MainHandler),
        (r"/api/tags",              TagsHandler),
        (r"/api/tags/(\d+)",        TagHandler),
        (r"/api/receipts",          ReceiptsHandler),
        (r"/api/receipts/(\d+)",    ReceiptHandler),
        (r"/api/images",            ImagesHandler),
        (r"/api/images/(\d+)",      ImageHandler),
        (r"/api/images/(\d+)/data", ImageDataHandler),
        (r"/api/images/scan",       ScanImageHandler),
    ], **settings)

    # Create database tables
    db.create_tables([Image, Receipt, Tag, ReceiptToTag, FTSReceipt], safe=True)

    # Start app
    log.info("Listening on port %d", 8888)
    application.listen(8888)

    try:
        tornado.ioloop.IOLoop.instance().start()
    except KeyboardInterrupt:
        log.info("Finished")
    finally:
        sys.stdout.flush()
        sys.stderr.flush()


if __name__ == "__main__":
    main()
