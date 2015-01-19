import logging

from bottle import abort, request, response

from ..app import app
from ..models import File


LOG = logging.getLogger(__name__)


@app.get('/api/files')
def files_get_many(db):
    # Pagination
    try:
        offset = int(request.query.offset or 0)
        limit = int(request.query.limit or 20)
    except ValueError:
        offset = 0
        limit = 20

    # Get all files, offset by the limit
    query = (db.query(File)
             .order_by(File.id)
             .offset(offset)
             .limit(limit)
             )

    files = [x.as_json() for x in query]

    return {
        'files': files,
        'meta': {
            'total': db.query(File).count(),
        },
    }


@app.post('/api/files')
def files_post(db):
    abort(501, 'Use the /upload endpoint instead')


@app.post('/api/files/upload')
def files_upload(db):
    # Ensure we have the appropriate fields
    if 'data' not in request.forms or 'filename' not in request.forms:
        abort(400, 'Invalid Request')

    LOG.debug("Uploading file '%s' with %d bytes of data",
              request.forms['filename'],
              len(request.forms['data'])
              )

    # Create a new file with this name
    ff = File(name=request.forms.filename,
              size=len(request.forms.data))
    # TODO: save contents somewhere
    # TODO: tags or collection

    db.add(ff)
    db.commit()


@app.get('/api/files/<file_id:int>')
def files_get_one(file_id, db):
    ff = db.query(File).filter_by(id=file_id).first()
    if not ff:
        abort(404, '')

    return {
        'file': ff.as_json(),
    }


@app.delete('/api/files/<file_id:int>')
def files_delete_one(file_id, db):
    ff = db.query(File).filter_by(id=file_id).first()
    if not ff:
        abort(404, '')

    ff.delete_instance()
    response.status = 204


@app.get('/api/files/<file_id:int>/content')
def files_get_content(file_id, db):
    ff = db.query(File).filter_by(id=file_id).first()
    if not ff:
        abort(404, '')

    # TODO: serve the contents of the file
    # TODO: be sure to try out the WSGI sendfile, if it exists
    abort(501, 'Not Implemented')
