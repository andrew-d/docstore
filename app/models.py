from . import db


tags_rel = db.Table("tags_rel",
    db.Column("tag_id", db.Integer, db.ForeignKey("tags.id")),
    db.Column("document_id", db.Integer, db.ForeignKey("documents.id")),
)


class Document(db.Model):
    __tablename__ = "documents"

    id = db.Column(db.Integer, primary_key=True)
    name = db.Column(db.String, index=True, unique=True)
    tags = db.relationship('Tag', secondary=tags_rel,
                           backref=db.backref('documents', lazy='dynamic'))

    def __repr__(self):
        return "<Document %r>" % (self.name,)


class Tag(db.Model):
    __tablename__ = "tags"

    id = db.Column(db.Integer, primary_key=True)
    name = db.Column(db.String, index=True, unique=True)
