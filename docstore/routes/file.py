import os
import hashlib
import logging

from bottle import abort, request, response

from ..app import app
from ..models import File, Tag


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

    # Hash contents
    data_hash = hashlib.sha256(request.forms['data']).hexdigest()

    LOG.debug("Uploading file '%s' with %d bytes of data and hash %s",
              request.forms['filename'],
              len(request.forms['data']),
              data_hash
              )

    # Create a new file with this name
    ff = File(
        name=request.forms.filename,
        hash=data_hash,
        size=len(request.forms.data)
    )
    # TODO: tags or collection

    # Write the file to disk if it doesn't already exist.
    fpath = os.path.join(app.config['docstore.data_path'], 'files', data_hash)
    if not os.path.exists(fpath):
        with open(fpath, 'wb') as f:
            f.write(request.forms['data'])

        # TODO: generate thumbnails

    db.add(ff)
    db.commit()


@app.get('/api/files/<file_id:int>')
def files_get_one(file_id, db):
    ff = db.query(File).filter_by(id=file_id).first()
    if not ff:
        abort(404)

    return {
        'file': ff.as_json(),
    }


@app.put('/api/files/<file_id:int>')
def files_modify_one(file_id, db):
    if not request.json or 'file' not in request.json:
        abort(400)

    ff = db.query(File).filter_by(id=file_id).first()
    if not ff:
        abort(404)

    new_data = request.json['file']
    ff.name = new_data['name']

    # Reset list of tags
    new_tags = []
    for tag_id in new_data['tags']:
        if isinstance(tag_id, str):
            tag_id = int(tag_id)

        tag = db.query(Tag).filter_by(id=tag_id).first()
        if not tag:
            abort(400)

        new_tags.append(tag)

    ff.tags = new_tags

    db.add(ff)
    db.commit()

    return {
        'file': ff.as_json(),
    }


@app.delete('/api/files/<file_id:int>')
def files_delete_one(file_id, db):
    ff = db.query(File).filter_by(id=file_id).first()
    if not ff:
        abort(404)

    ff.delete_instance()
    response.status = 204


@app.get('/api/files/<file_id:int>/content')
def files_get_content(file_id, db):
    ff = db.query(File).filter_by(id=file_id).first()
    if not ff:
        abort(404)

    # TODO: serve the contents of the file
    # TODO: be sure to try out the WSGI sendfile, if it exists
    abort(501, 'Not Implemented')
