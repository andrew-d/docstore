import logging
logging.basicConfig(
    format='DOCSTORE: [%(asctime)s] %(levelname)s: %(message)s',
    datefmt='%m/%d/%Y %I:%M:%S %p',
    level=logging.DEBUG
)

from .api import api
