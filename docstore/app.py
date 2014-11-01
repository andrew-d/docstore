import os
from concurrent.futures import ThreadPoolExecutor

from tornado.options import define
from peewee import Proxy


APP_ROOT = os.path.dirname(os.path.realpath(__file__))

#db = SqliteExtDatabase(DATABASE, threadlocals=True)
scan_pool = ThreadPoolExecutor(max_workers=1)
ocr_pool  = ThreadPoolExecutor(max_workers=10)

# Lazily initialize the database
db = Proxy()

define("debug", default=False, help="run the application in debug mode")
define("device", default=None, help="name of the device to use as a scanner")
define("database", default=None, help="location of the database")
define("storage", default=None, help="location to store scanned receipts in")
