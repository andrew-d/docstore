import os
import sys
import logging

import tornado
import tornado.web
import tornado.ioloop
from tornado.options import options

from playhouse.sqlite_ext import SqliteExtDatabase

import pyinsane.abstract as pyinsane

from .app import APP_ROOT, db
from . import routes as r
from . import models
from .util import Serializer


log = logging.getLogger(__name__)


def get_device():
    if options.debug:
        log.debug("scanning for devices")
        devices = pyinsane.get_devices()
        for i, dev in enumerate(devices):
            log.info("found device %d: %s/%s/%s",
                     i, dev.name, dev.vendor, dev.model)

    if options.device is None:
        log.warn("no device specified, scanning will be unavailable")
        return None

    device = pyinsane.Scanner(name=options.device)
    log.debug("using device: name = %s, vendor = %s, model = %s",
              device.name, device.vendor, device.model)
    return device


def main():
    tornado.options.parse_command_line()

    if options.storage is None:
        log.error("no storage directory given")
        return
    storage = os.path.realpath(os.path.abspath(options.storage))
    log.info("using storage path: %s", storage)

    if options.database is None:
        log.error("no database path given")
        return

    # Create the database and init. the proxy
    real_db = SqliteExtDatabase(options.database, threadlocals=True)
    db.initialize(real_db)

    #check_dependencies()

    settings = {
        'debug': options.debug,
        'device': get_device(),
        'storage': storage,
        'serializer': Serializer(),
    }

    application = tornado.web.Application([
        (r"/",                      r.MainHandler),
        (r"/api/tags",              r.TagsHandler),
        (r"/api/tags/(\d+)",        r.TagHandler),
        (r"/api/receipts",          r.ReceiptsHandler),
        (r"/api/receipts/(\d+)",    r.ReceiptHandler),
        (r"/api/images",            r.ImagesHandler),
        (r"/api/images/(\d+)",      r.ImageHandler),
        (r"/api/images/(\d+)/data", r.ImageDataHandler),
        (r"/api/images/scan",       r.ScanImageHandler),
    ], **settings)

    # Create database tables
    db.create_tables([
        models.Image,
        models.Receipt,
        models.Tag,
        models.ReceiptToTag,
        models.FTSReceipt
    ], safe=True)

    # Start app
    log.info("Listening on port %d", 8888)
    application.listen(8888)

    try:
        tornado.ioloop.IOLoop.instance().start()
    except KeyboardInterrupt:
        log.info("Finished")
    finally:
        sys.stdout.flush()
        sys.stderr.flush()
