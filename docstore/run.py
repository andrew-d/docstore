import bottle
from .app import app

# Importing this registers our routes
from . import api as _dummy


def add_dummy_data(db):
    from .models import File, Collection, Tag

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



def configure_app():
    """Set up application configuration"""
    from sqlalchemy import create_engine
    from bottle.ext import sqlalchemy
    from .models import Base

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
    plugin.metadata.create_all(engine)
    add_dummy_data(plugin.create_session(bind=engine))


if __name__ == "__main__":
    # Configure the application
    configure_app()

    # Add middleware.
    ##from .middleware import StripPathMiddleware
    ##myapp = StripPathMiddleware(app)
    myapp = app

    # Run the application
    bottle.run(app=myapp, host='localhost', port=8000, server='bjoern')
