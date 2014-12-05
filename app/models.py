import datetime

from flask.ext.sqlalchemy import SQLAlchemy


db = SQLAlchemy()


tags_rel = db.Table("tags_rel",
    db.Column("tag_id", db.Integer, db.ForeignKey("tags.id")),
    db.Column("document_id", db.Integer, db.ForeignKey("documents.id")),
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
    document_id = db.Column(db.Integer, db.ForeignKey('documents.id'))

    def __repr__(self):
        return '<File %r>' % (self.name,)

    def as_json(self):
        return {
            'id':       self.id,
            'size':     self.size,
            'name':     self.name,
            'created':  self.created.isoformat(),
            'document': self.document_id,
        }


class Document(db.Model):
    __tablename__ = "documents"

    id = db.Column(db.Integer, primary_key=True)
    name = db.Column(db.String, nullable=False, index=True, unique=True)
    created = db.Column(db.DateTime, nullable=False,
                        default=datetime.datetime.utcnow)

    # Associated files
    files = db.relationship('File', backref='document', lazy='dynamic')

    # General metadata
    meta = db.Column(db.PickleType, nullable=False, default={})

    tags = db.relationship('Tag', secondary=tags_rel,
                           backref=db.backref('documents', lazy='dynamic'))

    def __repr__(self):
        return "<Document %r>" % (self.name,)

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
        For each tag in tags, get or create it and apply it to this document.
        Note: this does not remove tags from the current document.
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

    def __repr__(self):
        return "<Tag %r>" % (self.name,)

    def as_json(self):
        return {
            'id':        self.id,
            'name':      self.name,
            'documents': [x.id for x in self.documents],
        }

    @classmethod
    def get_or_create(klass, name):
        instance = klass.query.filter(klass.name == name).first()
        if instance:
            return instance

        instance = klass(name=name)
        db.session.add(instance)
        return instance
