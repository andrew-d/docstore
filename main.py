#!/usr/bin/env python

import bottle
from bottle.ext import sqlalchemy
from sqlalchemy import create_engine

from docstore.api import app
from docstore.models import Base


TESTING = True


# Configure the application
# TODO: get the connection string from elsewhere
engine = create_engine('sqlite:///:memory:', echo=True)
plugin = sqlalchemy.Plugin(
    engine,
    Base.metadata,
    keyword='db',
    create=True, # Execute `metadata.create_all(engine)`
    commit=True, # Plugin commits changes after route is executed
    use_kwargs=False
)
app.install(plugin)


# For testing
if TESTING:
    plugin.metadata.create_all(engine)
    db = plugin.create_session(bind=engine)

    from docstore.models import File, Collection, Tag

    f1 = File(name='File 1', size=1234)
    f2 = File(name='File 2', size=456)

    t1 = Tag(name='tag1')
    t2 = Tag(name='t2')

    c1 = Collection(name='Collection')

    f1.tags.append(t1)
    f2.tags.append(t2)
    c1.files.append(f1)
    c1.files.append(f2)

    db.add_all([f1, f2, t1, t2, c1])
    db.commit()


bottle.run(app=app,
           host='localhost',
           port=8000,
           server='bjoern')
