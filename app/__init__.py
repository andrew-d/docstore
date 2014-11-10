import os

from flask import Flask
from flask.ext.sqlalchemy import SQLAlchemy


# Set up and configure app
app = Flask(__name__)
app.config.from_object('config')

# Set up and configure everything else
db = SQLAlchemy(app)


# This needs to go at the end.
from . import models, views, util
