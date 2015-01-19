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


file_tags = Table('file_tags', Base.metadata,
    Column('file_id', Integer, ForeignKey('files.id'), nullable=False),
    Column('tag_id', Integer, ForeignKey('tags.id'), nullable=False)
)


collection_files = Table('collection_files', Base.metadata,
    Column('collection_id', Integer, ForeignKey('collections.id'), nullable=False),
    Column('file_id', Integer, ForeignKey('files.id'), nullable=False)
)


class Tag(Base):
    __tablename__ = "tags"

    id = Column(Integer, primary_key=True)
    name = Column(String, nullable=False, unique=True)

    def as_json(self):
        return {
            'id':    self.id,
            'name':  self.name,
            'files': [x.id for x in self.files],
        }


class File(Base):
    __tablename__ = "files"

    id = Column(Integer, primary_key=True)

    name = Column(String, nullable=False, unique=True)  # Name as uploaded
    hash = Column(String(32), nullable=False)           # SHA256 hash of data
    size = Column(Integer, nullable=False)              # Size of file data
    created_at = Column(DateTime, nullable=False, default=datetime.utcnow)

    tags = relationship('Tag',
                        secondary=file_tags,
                        backref='files')
    collections = relationship('Collection',
                               secondary=collection_files,
                               backref='files')

    def as_json(self):
        return {
            'id':          self.id,
            'name':        self.name,
            'size':        self.size,
            'created_at':  self.created_at.isoformat(),
            'tags':        [x.id for x in self.tags],
            'collections': [x.id for x in self.collections],
        }


class Collection(Base):
    __tablename__ = "collections"

    id = Column(Integer, primary_key=True)
    name = Column(String, nullable=False, unique=True)

    # A collection can optionally be a member of another collection
    parent_id = Column(Integer, ForeignKey('collections.id'), nullable=True)
    parent = relationship('Collection',
                          remote_side='Collection.id',
                          backref='children')


    def as_json(self):
        return {
            'id':       self.id,
            'name':     self.name,
            'parent':   self.parent,
            'children': [x.id for x in self.children],
            'files':    [x.id for x in self.files],
        }
