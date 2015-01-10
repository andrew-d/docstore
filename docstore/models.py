from datetime import datetime

from sqlalchemy import (
    Column,
    DateTime,
    ForeignKey,
    Integer,
    String,
    Table,
)
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import relationship, backref


Base = declarative_base()


item_tags = Table('item_tags', Base.metadata,
    Column('item_id', Integer, ForeignKey('items.id')),
    Column('tag_id', Integer, ForeignKey('tags.id'))
)


collection_items = Table('collection_items', Base.metadata,
    Column('collection_id', Integer, ForeignKey('collections.id')),
    Column('item_id', Integer, ForeignKey('items.id'))
)


class Tag(Base):
    __tablename__ = "tags"

    id = Column(Integer, primary_key=True)
    name = Column(String, nullable=False, unique=True)

    def as_json(self):
        return {
            'id':   self.id,
            'name': self.name,
        }


class Item(Base):
    __tablename__ = "items"

    id = Column(Integer, primary_key=True)

    discriminator = Column('type', String(50))
    __mapper_args__ = {'polymorphic_on': discriminator}

    # Each item can have tags and possibly be part of collections
    # TODO: is this polymorphic?
    tags = relationship('Tag',
                        secondary=item_tags,
                        backref='items')
    collections = relationship('Collection',
                               secondary=collection_items,
                               backref='children')


class File(Item):
    __tablename__ = "files"
    __mapper_args__ = {'polymorphic_identity': 'file'}

    file_id = Column('id', Integer, ForeignKey('items.id'), primary_key=True)

    name = Column(String, nullable=False, unique=True)
    size = Column(Integer, nullable=False)
    created_at = Column(DateTime, nullable=False, default=datetime.utcnow)

    def as_json(self):
        return {
            'id':         self.id,
            'name':       self.name,
            'size':       self.size,
            'created_at': self.created_at.isoformat(),
            #'item':       self.item.id,
        }


class Collection(Item):
    __tablename__ = "collections"
    __mapper_args__ = {'polymorphic_identity': 'collection'}

    coll_id = Column('id', Integer, ForeignKey('items.id'), primary_key=True)
    name = Column(String, nullable=False, unique=True)



# TODO: configure me
engine = create_engine('sqlite://:memory:')
Base.metadata.create_all(engine)
