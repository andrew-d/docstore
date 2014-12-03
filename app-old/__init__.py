import os

from flask import Flask
from flask.ext.sqlalchemy import SQLAlchemy
from flask.ext.wtf import CsrfProtect

from . import middleware, scanner


# Set up and configure app
app = Flask(__name__)
app.wsgi_app = middleware.MethodRewriteMiddleware(app.wsgi_app,
                                                  query_param="_methodov")
app.config.from_object('config')

# Generated configuration
app.config['UPLOAD_FOLDER'] = os.path.join(app.config['DATA_DIRECTORY'],
                                           'uploads')
app.config['INDEX_PATH'] = os.path.join(app.config['DATA_DIRECTORY'], 'search')


# Set up and configure everything else
## Database
db = SQLAlchemy(app)

## CSRF protection
csrf = CsrfProtect(app)
app.config['WTF_CSRF_METHODS'] = ['POST', 'PUT', 'PATCH', 'DELETE']

## Scanner
app.config['scanners'] = scanner.Config

## Searching
from .search import open_index
app.config['index'] = open_index()


# This needs to go at the end.
app.logger.info("Configuration finished")
from . import models, search, util, views
