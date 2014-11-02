import json
import datetime

from peewee import (
    CompositeKey,
    DateTimeField,
    Field,
    ForeignKeyField,
    JOIN_LEFT_OUTER,
    Model,
    PrimaryKeyField,
    TextField,
)
from playhouse.sqlite_ext import FTSModel, SqliteExtDatabase

from docstore.app import db


class JSONField(Field):
    db_field = 'json'

    def db_value(self, value):
        if not isinstance(value, (dict, list, tuple)):
            raise ValueError("unknown type for JSONField: %s",
                             type(value).__name__)

        return json.dumps(value)

    def python_value(self, value):
        return json.loads(value)


# Register the JSONField type with SQLite
SqliteExtDatabase.register_fields({"json": "TEXT"})


class BaseModel(Model):
    class Meta:
        database = db


class Image(BaseModel):
    """The model representing an uploaded or scanned image"""
    id = PrimaryKeyField()
    hash = TextField(null=False)
    created = DateTimeField(null=False, default=datetime.datetime.now)

    @classmethod
    def get_expired(self):
        """
        Return all images that are more than 24 hours old and are not
        associated with a document.
        """
        expiry_time = datetime.timedelta(hours=24)
        abs_expiry = datetime.datetime.now() - expiry_time
        return (Image.
                select().
                where(Image.created < abs_expiry).
                join(Document, JOIN_LEFT_OUTER).
                where(Document.image >> None))


class Document(BaseModel):
    """The model representing a single document"""
    id = PrimaryKeyField()
    ocr_data = TextField(null=False)
    image = ForeignKeyField(Image, null=False)
    created = DateTimeField(null=False, default=datetime.datetime.now)

    # Other arbitrary data, stored as JSON.
    metadata = JSONField(null=False)


class FTSDocument(FTSModel):
    """The full-text search model for our Document"""
    document = ForeignKeyField(Document, primary_key=True)
    content = TextField()

    class Meta:  # Fix highlighting ):
        database = db

    @classmethod
    def store_note(cls, document):
        try:
            fts_document = FTSDocument.get(FTSDocument.document == document)
        except FTSDocument.DoesNotExist:
            fts_document = FTSDocument(document=document)
            force_insert = True
        else:
            force_insert = False

        fts_document.content = document.ocr_data
        fts_document.save(force_insert=force_insert)


class Tag(BaseModel):
    """The model representing a tag."""
    id = PrimaryKeyField()
    name = TextField(null=False, unique=True)


class DocumentToTag(BaseModel):
    """A simple "through" table for many-to-many relationship."""
    document = ForeignKeyField(Document, null=False)
    tag = ForeignKeyField(Tag, null=False)

    class Meta:  # Fix highlighting ):
        primary_key = CompositeKey('document', 'tag')
