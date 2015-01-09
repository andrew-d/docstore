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
            'files':      [x.id for x in self.files],
        }


class ItemTag(BaseModel):
    item = ForeignKeyField(Item)
    tag = ForeignKeyField(Tag)


class File(BaseModel):
    id = PrimaryKeyField()
    name = CharField(null=False, unique=True)
    size = IntegerField(null=False)
    created_at = DateTimeField(null=False, default=datetime.utcnow)
    item = ForeignKeyField(Item, related_name='files')

    def as_json(self):
        return {
            'id':         self.id,
            'name':       self.name,
            'size':       self.size,
            'created_at': self.created_at.isoformat(),
            'item':       self.item.id,
        }


# TODO: configure me
database_proxy.initialize(SqliteDatabase(":memory:"))
database_proxy.create_tables([
    File,
    Item,
    ItemTag,
    Tag,
], safe=True)
