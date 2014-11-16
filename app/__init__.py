from flask import Flask
from flask.ext.sqlalchemy import SQLAlchemy
from flask.ext.wtf import CsrfProtect

from . import scanner, middleware


# Set up and configure app
app = Flask(__name__)
app.wsgi_app = middleware.MethodRewriteMiddleware(app.wsgi_app,
                                                  query_param="_methodov")
app.config.from_object('config')

# Set up and configure everything else
db = SQLAlchemy(app)
csrf = CsrfProtect(app)
app.config['scanners'] = scanner.Config
app.config['WTF_CSRF_METHODS'] = ['POST', 'PUT', 'PATCH', 'DELETE']


# This needs to go at the end.
app.logger.info("Configuration finished")
from . import models, views, util
