import datetime


from . import db


tags_rel = db.Table("tags_rel",
    db.Column("tag_id", db.Integer, db.ForeignKey("tags.id")),
    db.Column("document_id", db.Integer, db.ForeignKey("documents.id")),
)


class File(db.Model):
    __tablename__ = "files"

    id = db.Column(db.Integer, primary_key=True)
    size = db.Column(db.Integer, nullable=False)
    name = db.Column(db.String, nullable=False, index=True, unique=True)
    created = db.Column(db.DateTime, nullable=False,
                        default=datetime.datetime.utcnow)
    document_id = db.Column(db.Integer, db.ForeignKey('documents.id'))


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

    def apply_tags(self, tags):
        if not tags:
            return self

        # TODO: proper shell splitting with quotes
        new_tags = tags.split(' ')
        for tagname in new_tags:
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

    @classmethod
    def get_or_create(klass, name):
        instance = klass.query.filter(klass.name == name).first()
        if instance:
            return instance

        instance = klass(name=name)
        db.session.add(instance)
        return instance
