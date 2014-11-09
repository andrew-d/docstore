import os

from flask import Flask
from flask.ext.sqlalchemy import SQLAlchemy
from flask.ext.uploads import (
    AllExcept,
    EXECUTABLES,
    UploadSet,
    configure_uploads
)


# Set up and configure app
app = Flask(__name__)
app.config.from_object('config')

# Set up and configure everything else
db = SQLAlchemy(app)
uploads = UploadSet('documents',
                    AllExcept(EXECUTABLES),
                    default_dest=lambda app: app.config['UPLOAD_FOLDER'])
configure_uploads(app, (uploads,))


# This needs to go at the end.
from . import models, views, util
