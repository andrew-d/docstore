import logging

from bottle import abort, request, response

from ..app import app
from ..models import Collection


LOG = logging.getLogger(__name__)


@app.get('/api/collections')
def collections_get_many(db):
    # Pagination
    try:
        offset = int(request.query.offset or 0)
        limit = int(request.query.limit or 20)
    except ValueError:
        offset = 0
        limit = 20

    # Get all collections, offset by the limit
    query = (db.query(Collection)
             .order_by(Collection.id)
             .offset(offset)
             .limit(limit)
             )

    collections = [x.as_json() for x in query]

    return {
        'collections': collections,
        'meta': {
            'total': db.query(Collection).count(),
        },
    }

@app.post('/api/collections')
def collections_post(db):
    abort(501, 'Not Implemented')


@app.get('/api/collections/<collection_id:int>')
def collections_get_one(collection_id, db):
    collection = db.query(Collection).filter_by(id=collection_id).first()
    if not collection:
        abort(404, '')

    return {
        'collection': collection.as_json(),
    }


@app.delete('/api/collections/<collection_id:int>')
def collections_delete_one(collection_id, db):
    collection = db.query(Collection).filter_by(id=collection_id).first()
    if not collection:
        abort(404, '')

    collection.delete_instance()
    response.status = 204
