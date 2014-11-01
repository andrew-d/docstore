import datetime

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
from playhouse.sqlite_ext import FTSModel

from .app import db


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
