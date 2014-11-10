import signal
import multiprocessing
from io import BytesIO
from collections import namedtuple


ScannerInfo = namedtuple('ScannerInfo', 'name vendor model')


def _scan_image(scanner_name):
    import pyinsane.abstract_th as pyinsane

    device = pyinsane.Scanner(name=scanner_name)

    # TODO: multiple page support
    scan_session = device.scan(multiple=False)
    try:
        while True:
            scan_session.scan.read()
    except EOFError:
        pass

    scanned_image = scan_session.images[0]

    # Serialize the image as PNG, and hash it.
    with BytesIO() as memf:
        scanned_image.save(memf, 'PNG')
        data = memf.getvalue()

    fhash = hashlib.sha256(data).hexdigest()
    fsize = len(data)

    return data, fhash, len(data)


def _get_scanner_info():
    import pyinsane.abstract_th as pyinsane
    devices = pyinsane.get_devices()
    return [ScannerInfo(name=d.name,
                        vendor=d.vendor,
                        model=d.model) for d in devices]


def _init_worker():
    # This allows us to Ctrl-C the main process and not have it hang.
    signal.signal(signal.SIGINT, signal.SIG_IGN)


# Note: this must be created before the various internal functions above, or
# it won't be able to "see" those functions.
POOL = multiprocessing.Pool(1, _init_worker)


def scan_image(scanner_name):
    return POOL.apply(_scan_image, [scanner_name])
    #future = EXECUTOR.submit(_scan_image, scanner_name)
    #return future.result()


class _Config(object):
    """
    Class-level caching configuration for scanner information.

    Note: the reason we have this as a (singleton) object is because,
    for some reason, getting the information during startup causes
    Flask to hang.  So, I made it a memoized property instead.
    """
    def __init__(self):
        self.__scanner_info = None

    @property
    def scanner_info(self):
        if self.__scanner_info is None:
            self.__scanner_info = POOL.apply(_get_scanner_info)

        return self.__scanner_info

    @property
    def have_scanner(self):
        return len(self.scanner_info) > 0


Config = _Config()
