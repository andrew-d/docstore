import os
import json
import hashlib
import inspect
import logging
from functools import partial

import tornado.web
from tornado import gen

import peewee
from peewee import (
    Model,
)

from docstore.app import db, scan_pool
from docstore.models import (
    Image,
    Document,
    Tag,
    # DocumentToTag,
    # FTSDocument
)


log = logging.getLogger(__name__)


class BaseHandler(tornado.web.RequestHandler):
    def initialize(self):
        self.__json_body = None

    @property
    def settings(self):
        return self.application.settings

    @property
    def serializer(self):
        return self.application.settings['serializer']

    @property
    def json_body(self):
        if self.__json_body is None:
            # TODO: look at real charset of request
            body_str = self.request.body.decode('latin1')

            try:
                self.__json_body = json.loads(body_str)
            except ValueError:
                raise tornado.web.HTTPError(400, reason="invalid JSON")

        return self.__json_body

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
        self.write_json(kwargs)

    def serialize(self, name, obj, *, force_serialize=False):
        """
        Serialize a peewee model or query to the appropriate format,
        as expected by Ember-Data.
        """
        if isinstance(obj, Model):
            # Single object - should be serialized at the root.
            obj = {name: self.serializer.serialize_object(obj)}

        elif (isinstance(obj, peewee.Query) or
              isinstance(obj, (list, tuple)) and force_serialize):
            # Set of objects, should be serialized in the format:
            #   {
            #       "models": [ list of models ]
            #   }
            serialized = [self.serializer.serialize_object(x) for x in obj]
            obj = {name + 's': serialized}

        else:
            raise RuntimeError("unknown type to serialize: %r" % (obj,))

        return self.write_json(obj)

    def write_json(self, obj):
        """
        Write an object out in JSON format, setting the correct Content-Type
        """
        self.set_header('Content-Type', 'application/json; charset=utf-8')
        self.write(json.dumps(obj))

    def set_default_headers(self):
        self.set_header('Access-Control-Allow-Origin', '*')

    @tornado.web.asynchronous
    def options(self):
        # Walk through the list of supported methods, checking if we have a
        # function with that name defined anywhere except the base
        # RequestHandler.
        method_names = [m.lower() for m in self.SUPPORTED_METHODS]
        supported = []
        for m in method_names:
            func = getattr(self, m)
            if not func:
                continue

            # Walk the inheritance chain, finding the first class that defines
            # this method.
            defining_class = None
            for cls in inspect.getmro(func.__self__.__class__):
                if func.__name__ in cls.__dict__:
                    defining_class = cls
                    break

            if defining_class is not tornado.web.RequestHandler:
                supported.append(m.upper())

        self.set_header('Access-Control-Allow-Methods', ', '.join(supported))
        self.set_header('Access-Control-Allow-Headers', 'Accept, Content-Type')
        self.set_status(204)
        self.finish()


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
        if 'tag' not in self.json_body:
            return self.send_error(400, message="bad request")
        if 'name' not in self.json_body['tag']:
            return self.send_error(400, message="bad request")

        try:
            with db.transaction():
                t = Tag.create(name=self.json_body['tag']['name'])
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
        self.write_json(True)


class DocumentsHandler(BaseHandler):
    """
    Handles the routes:
        GET     /api/documents
        POST    /api/documents
    """
    def get(self):
        self.serialize('document', Document.select())

    def post(self):
        try:
            # TODO: real parameters here
            with db.transaction():
                t = Document.create(
                    path=self.json_body['path'],
                    amount=self.json_body['amount'],
                )
        except peewee.IntegrityError:
            return self.send_error(409, message="document already exists")

        self.set_status(201)
        self.serialize('document', t)


class DocumentHandler(BaseHandler):
    """
    Handles the routes:
        GET     /api/documents/<id>
        DELETE  /api/documents/<id>
    """
    def get(self, document_id):
        try:
            document = Document.get(Document.id == document_id)
        except Document.DoesNotExist:
            return self.send_error(404, message="document does not exist")

        self.serialize('document', document)

    def delete(self, document_id):
        try:
            document = Document.get(Document.id == document_id)
        except Document.DoesNotExist:
            return self.send_error(404, message="document does not exist")

        document.delete_instance()
        self.write_json(True)


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
        except OSError:
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


class UploadImageHandler(BaseHandler):
    def post(self):
        log.info("files = %r", self.request.files)

        if 'image' not in self.request.files:
            return self.send_error(400, message="no image given")
