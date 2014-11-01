import logging
import subprocess

log = logging.getLogger(__name__)

# Tesseract wrapper - tries OCRing an image.  Returns None if something failed,
# the OCR-d text otherwise.
def ocr_image(image_file, language='eng'):
    proc = subprocess.Popen(['tesseract', 'stdin', 'stdout', '-l', language],
                            stdin=image_file,
                            stdout=subprocess.PIPE,
                            stderr=subprocess.PIPE
                            )

    try:
        stdout, stderr = proc.communicate(timeout=15)
    except subprocess.TimeoutExpired:
        log.warn("OCR process timed out - killing...")
        proc.kill()
        stdout, stderr = proc.communicate()

    if proc.returncode != 0:
        log.warn("OCR process exited with return code %d", proc.returncode)
        return None

    # TODO: convert from bytes to string?
    return stdout

