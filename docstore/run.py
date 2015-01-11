from .app import app

# Importing this registers our routes
from . import api as _dummy


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


if __name__ == "__main__":
    # Configure the application
    configure_app()

    # Run the application
    app.run(host='localhost', port=8000, server='bjoern')
