# -*- coding: utf8 -*-

import os


CSRF_ENABLED = True
SECRET_KEY   = ('cb04a4d1284a33e7d2a6cebc864c70acf5112992' +
                'f50c3a9caf561cd7d013576f1b0e36b129c29b6f')
DATA_DIRECTORY = os.path.abspath(os.path.dirname(__file__))

if os.environ.get('DATABASE_URL') is None:
    SQLALCHEMY_DATABASE_URI = ('sqlite:///' + os.path.join(DATA_DIRECTORY, 'app.db') +
                               '?check_same_thread=False')
else:
    SQLALCHEMY_DATABASE_URI = os.environ['DATABASE_URL']
