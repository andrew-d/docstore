import bottle

app = bottle.Bottle()


@app.hook('startup')
def install_db():
    """
    Install the bottle-sqlalchemy plugin when the application starts.
    """
    from bottle.ext import sqlalchemy
    from sqlalchemy import create_engine

    from .models import Base

    engine = create_engine(
        app.config.get('docstore.dbconn', 'sqlite:///:memory:'),
        echo=app.config.get('debug', False)
    )

    plugin = sqlalchemy.Plugin(
        engine,
        Base.metadata,
        keyword='db',
        create=True, # Execute `metadata.create_all(engine)`
        commit=True, # Plugin commits changes after route is executed
        use_kwargs=False
    )
    app.install(plugin)


@app.hook('startup')
def create_data_dir():
    """
    Create the data directories on startup.  Will create:

        /the/data/dir/files/
        /the/data/dir/thumbnails/
    """
    import os

    data_path = app.config['docstore.data_path']
    for fname in ['files', 'thumbnails']:
        os.makedirs(os.path.join(data_path, fname))
