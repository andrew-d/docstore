import os

from flask import Flask
from flask.ext.sqlalchemy import SQLAlchemy

from . import scanner


# Set up and configure app
app = Flask(__name__)
app.config.from_object('config')

# Set up and configure everything else
db = SQLAlchemy(app)
app.config['scanners'] = scanner.Config


# This needs to go at the end.
app.logger.info("Configuration finished")
from . import models, views, util
