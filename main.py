#!/usr/bin/env python

import os
import sys
import logging
sys.path.insert(0, os.path.join(os.path.dirname(__file__), "packages"))

import bottle

from docstore.api import app


# Configure the application
# TODO: get these configuration variables from somewhere
app.config['debug'] = True
app.config['docstore.dbconn'] = 'sqlite:///:memory:'
app.config['docstore.data_path'] = os.path.join(
    os.path.dirname(__file__), "data")


# Configure logging
logging.basicConfig(
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
    level=logging.DEBUG if app.config.get('debug') else logging.INFO
)



bottle.run(app=app,
           host='localhost',
           port=8000,
           server='bjoern')
