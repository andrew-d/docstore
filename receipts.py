#!/usr/bin/env python

import os
import sys
import json
import logging
import datetime
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
import tornado.web
from tornado.options import define, options


###############################################################################
## Configuration
APP_ROOT = os.path.dirname(os.path.realpath(__file__))
DATABASE = os.path.join(APP_ROOT, 'receipts.db')

db = SqliteExtDatabase(DATABASE, threadlocals=True)
log = logging.getLogger(__name__)
threadpool = ThreadPoolExecutor(max_workers=4)

define("debug", default=False, help="run the application in debug mode")


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


###############################################################################
## Models
class BaseModel(Model):
    class Meta:
        database = db


class Image(BaseModel):
    """The model representing an uploaded or scanned image"""
    id       = PrimaryKeyField()
    hash     = TextField(null=False)
    path     = TextField(null=False)
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

    def serialize(self, obj, *, name=None, force_serialize=False):
        """
        Serialize a peewee model or query to the appropriate format,
        as expected by Ember-Data.
        """
        if isinstance(obj, Model):
            # Single object - should be serialized at the root.
            obj = SERIALIZE.serialize_object(obj)
        elif (isinstance(obj, peewee.Query) or
              isinstance(obj, (list, tuple)) and force_serialize):
            # Set of objects, should be serialized in the format:
            #   {
            #       "models": [ list of models ]
            #   }
            #
            # Currently, a name must be provided.
            if name is None:
                raise RuntimeError("must provide a name when serializing multiple models")

            obj = {name: [SERIALIZE.serialize_object(x) for x in obj]}
        else:
            raise RuntimeError("unknown type to serialize: %r" % (obj,))

        return self.json(obj)

    def json(self, obj):
        """
        Write an object out in JSON format, setting the correct Content-Type
        """
        self.set_header('Content-Type', 'application/json; charset=utf-8')
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
        self.serialize(Tag.select(), name='tags')

    def post(self):
        params = json.loads(self.request.body.decode('latin1'))
        try:
            with db.transaction():
                t = Tag.create(name=params['name'])
        except peewee.IntegrityError:
            return self.send_error(409, message="tag already exists")

        self.set_status(201)
        self.serialize(t)


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

        self.serialize(tag)

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
        self.serialize(Receipt.select(), name='receipts')

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
        self.serialize(t)


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

        self.serialize(receipt)

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
        self.serialize(Image.select(), name='images')


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

        self.serialize(image)


class ScanImageHandler(BaseHandler):
    """
    Handles the routes:
        POST    /api/images/scan
    """
    @tornado.web.asynchronous
    def post(self):
        # TODO
        pass

# TODO:
#   - ScannedImage model
#       - Endpoint to scan these
#       - Clean up after 24 hours have passed (sched module?)
#       - OCR the input receipt


def check_dependencies():
    # TODO: check for tesseract
    pass


################################################################################
## App setup
if __name__ == "__main__":
    tornado.options.parse_command_line()

    check_dependencies()

    application = tornado.web.Application([
        (r"/",                      MainHandler),
        (r"/api/tags",              TagsHandler),
        (r"/api/tags/(\d+)",        TagHandler),
        (r"/api/receipts",          ReceiptsHandler),
        (r"/api/receipts/(\d+)",    ReceiptHandler),
        (r"/api/images",            ImagesHandler),
        (r"/api/images/(\d+)",      ImageHandler),
        (r"/api/images/scan",       ScanImageHandler),
    ], debug=options.debug)

    # Create database tables
    db.create_tables([Receipt, Tag, ReceiptToTag, FTSReceipt], safe=True)

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
