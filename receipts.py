#!/usr/bin/env python

import os
import json
import datetime

import peewee
from peewee import *
from playhouse.sqlite_ext import FTSModel, SqliteExtDatabase

import tornado.ioloop
import tornado.web
from tornado.options import define, options


################################################################################
## Configuration
APP_ROOT = os.path.dirname(os.path.realpath(__file__))
DATABASE = os.path.join(APP_ROOT, 'receipts.db')

db = SqliteExtDatabase(DATABASE, threadlocals=True)


################################################################################
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


################################################################################
## Models
class BaseModel(Model):
    class Meta:
        database = db


class Receipt(BaseModel):
    """The model representing a single receipt."""
    id       = PrimaryKeyField()
    amount   = DecimalField(null=False)
    ocr_data = TextField(null=False)
    path     = TextField(null=False)
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
    receipt = ForeignKeyField(Receipt)
    tag = ForeignKeyField(Tag)

    class Meta: # Fix highlighting ):
        primary_key = CompositeKey('receipt', 'tag')


################################################################################
## Handlers
class BaseHandler(tornado.web.RequestHandler):
    def write_error(self, status_code, **kwargs):
        if 'message' not in kwargs:
            if status_code == 405:
                kwargs['message'] = 'Invalid HTTP method'
            else:
                kwargs['message'] = 'Unknown error'

        kwargs["error"] = True
        self.set_status(status_code)
        self.json(kwargs)

    def serialize(self, obj, *, name=None, force_serialize=False):
        if isinstance(obj, Model):
            obj = SERIALIZE.serialize_object(obj)
        elif (isinstance(obj, peewee.Query) or
              isinstance(obj, (list, tuple)) and force_serialize):
            if name is None:
                raise RuntimeError("must provide a name when serializing multiple models")

            obj = {name: [SERIALIZE.serialize_object(x) for x in obj]}
        else:
            raise RuntimeError("unknown type to serialize: %r" % (obj,))

        return self.json(obj)

    def json(self, obj):
        self.set_header('Content-Type', 'application/json; charset=utf-8')
        self.write(json.dumps(obj))


class MainHandler(BaseHandler):
    def get(self):
        self.write("Hello, world")


class TagsHandler(BaseHandler):
    def get(self):
        self.serialize(Tag.select(), name='tags')

    def post(self):
        params = json.loads(self.request.body.decode('latin1'))
        try:
            with db.transaction():
                t = Tag.create(name=params['name'])
        except IntegrityError:
            return self.send_error(409, message="tag already exists")

        self.set_status(201)
        self.serialize(t)


class TagHandler(BaseHandler):
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
        except IntegrityError:
            return self.send_error(409, message="receipt already exists")

        self.set_status(201)
        self.serialize(t)


class ReceiptHandler(BaseHandler):
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


# TODO:
#   - ScannedImage model
#       - Endpoint to scan these
#       - Clean up after 24 hours have passed (sched module?)
#       - OCR the input receipt


################################################################################
## App setup
application = tornado.web.Application([
    (r"/", MainHandler),
    (r"/api/tags", TagsHandler),
    (r"/api/tags/(\d+)", TagHandler),
    (r"/api/receipts", ReceiptsHandler),
    (r"/api/receipts/(\d+)", ReceiptHandler),
])

if __name__ == "__main__":
    tornado.options.parse_command_line()

    # Create database tables
    db.create_tables([Receipt, Tag, ReceiptToTag, FTSReceipt], safe=True)

    # Start app
    application.listen(8888)
    tornado.ioloop.IOLoop.instance().start()
