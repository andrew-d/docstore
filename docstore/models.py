from peewee import *
from datetime import datetime


database_proxy = Proxy()


class BaseModel(Model):
    class Meta(object):
        database = database_proxy


class Tag(BaseModel):
    id = PrimaryKeyField()
    name = CharField(null=False, unique=True)

    def as_json(self):
        return {
            'id':   self.id,
            'name': self.name,
        }


class Item(BaseModel):
    id = PrimaryKeyField()
    created_at = DateTimeField(null=False, default=datetime.utcnow)

    def as_json(self):
        return {
            'id':         self.id,
            'created_at': self.created_at.isoformat(),
        }


class ItemTags(BaseModel):
    item = ForeignKeyField(Item)
    tag = ForeignKeyField(Tag)


# TODO: configure me
database_proxy.initialize(SqliteDatabase(":memory:"))
database_proxy.create_tables([
    Tag, Item, ItemTags
], safe=True)
