# -*- coding: utf8 -*-

import os


class Config(object):
    DEBUG = False
    TESTING = False
    SECRET_KEY = ('cb04a4d1284a33e7d2a6cebc864c70acf5112992' +
                  'f50c3a9caf561cd7d013576f1b0e36b129c29b6f')

    # Data is stored in subfolders of this directory
    DATA_DIRECTORY = os.path.abspath(os.path.dirname(__file__))

    # Control whether certain features are enabled
    ENABLE_SCANNING = True
    ENABLE_OCR = True


class Production(Config):
    pass


class Development(Config):
    DEBUG = True
    DATA_DIRECTORY = Config.DATA_DIRECTORY
    SQLALCHEMY_DATABASE_URI = 'sqlite:///%s?check_same_thread=False' % (
        os.path.join(DATA_DIRECTORY, 'app.db'))
