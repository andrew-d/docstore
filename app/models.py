import datetime

from flask.ext.sqlalchemy import SQLAlchemy


db = SQLAlchemy()


# General overview
#   - The base model in this application is an `Item`
#   - An Item can have multiple `File`s
#   - Each File represents a single file on-disk, along with its size, when it
#     was created, and the "type" of the file.  The "type" is used to determine
#     how to display the file in the browser.
#   - Files are stored in content-addressible storage - i.e. each file is
#     stored in a directory named by its hash.
#   - An Item can have multple `Tag`s
#   - Each Tag is a single name
#   - Tags can be aliased to another Tag


tags_rel = db.Table("tags_rel",
    db.Column("tag_id", db.Integer, db.ForeignKey("tags.id")),
    db.Column("item_id", db.Integer, db.ForeignKey("items.id")),
)


class File(db.Model):
    __tablename__ = "files"

    # Note: we don't use the name as a primary key because we want to be able
    # to edit the file (e.g. cropping images) without changing the model's ID.
    id = db.Column(db.Integer, primary_key=True)
    size = db.Column(db.Integer, nullable=False)
    name = db.Column(db.String, nullable=False, index=True, unique=True)
    created = db.Column(db.DateTime, nullable=False,
                        default=datetime.datetime.utcnow)
    item_id = db.Column(db.Integer, db.ForeignKey('items.id'))

    def __repr__(self):
        return '<File %r>' % (self.name,)

    def as_json(self):
        return {
            'id':       self.id,
            'size':     self.size,
            'name':     self.name,
            'created':  self.created.isoformat(),
            'item':     self.item_id,
        }


class Item(db.Model):
    __tablename__ = "items"

    id = db.Column(db.Integer, primary_key=True)
    name = db.Column(db.String, nullable=False, index=True, unique=True)
    created = db.Column(db.DateTime, nullable=False,
                        default=datetime.datetime.utcnow)

    # Associated files
    files = db.relationship('File', backref='item', lazy='dynamic')

    # General metadata
    meta = db.Column(db.PickleType, nullable=False, default={})

    tags = db.relationship('Tag', secondary=tags_rel,
                           backref=db.backref('items', lazy='dynamic'))

    def __repr__(self):
        return "<Item %r>" % (self.name,)

    def as_json(self):
        return {
            'id':       self.id,
            'name':     self.name,
            'created':  self.created.isoformat(),
            'files':    [f.id for f in self.files],
            'meta':     self.meta,
            'tags':     [t.id for t in self.tags],
        }

    def apply_tags(self, tags):
        """
        For each tag in tags, get or create it and apply it to this item.
        Note: this does not remove tags from the current item.
        """
        if not tags:
            return self

        for tagname in tags:
            t = Tag.get_or_create(tagname)

            if t not in self.tags:
                self.tags.append(t)

        return self


class Tag(db.Model):
    __tablename__ = "tags"

    id = db.Column(db.Integer, primary_key=True)
    name = db.Column(db.String, nullable=False, index=True, unique=True)

    # Indexed since, when listing tags, we look for .where(Tag.alias_for_id == None)
    alias_for_id = db.Column(db.Integer,
                             db.ForeignKey("tags.id"),
                             nullable=True, index=True)
    alias_for = db.relationship(lambda: Tag, remote_side=id, backref='aliases')

    def __repr__(self):
        return "<Tag %r>" % (self.name,)

    def as_json(self):
        return {
            'id':        self.id,
            'name':      self.name,
            'items': [x.id for x in self.items],
        }

    @classmethod
    def get_or_create(klass, name, resolve_aliases=True):
        instance = klass.query.filter(klass.name == name).first()
        if instance:
            # If this is an alias for another tag, return it
            if instance.alias_for is not None and resolve_aliases:
                return instance.alias_for

            return instance

        instance = klass(name=name)
        db.session.add(instance)
        return instance
