import os
from whoosh.index import create_in, open_dir
from whoosh.fields import DATETIME, ID, KEYWORD, Schema, TEXT

from . import app


schema = Schema(
    name=TEXT(stored=True),
    created=DATETIME,
    tags=KEYWORD,
    data=TEXT
)


def create_index():
    pth = app.config['INDEX_PATH']
    os.makedirs(pth)
    ix = create_in(pth, schema)


def open_index():
    return open_dir(app.config['INDEX_PATH'])


def add_to_index(documents):
    ix = app.config['index']
    writer = ix.writer()

    for doc in documents:
        # TODO: add 'data' field
        writer.add_document(
            name=unicode(doc.name),
            created=doc.created,
            tags=u' '.join(t.name for t in doc.tags),
            data=u''
        )

    writer.commit()
