import logging

from bottle import abort, request, response

from ..app import app
from ..models import Tag


LOG = logging.getLogger(__name__)


@app.get('/api/tags')
def tags_get_many(db):
    # If we have a 'name' parameter, we need to lookup the individual tag, or
    # return a 404 error.
    if request.query.name:
        tag = db.query(Tag).filter_by(name=request.query.name).first()
        if not tag:
            abort(404, '')

        # Return just the single tag
        return {
            'tags': [tag.as_json()],
            'meta': {
                'total': 1,
            },
        }

    # Pagination
    try:
        offset = int(request.query.offset or 0)
        limit = int(request.query.limit or 20)
    except ValueError:
        offset = 0
        limit = 20

    # Get all tags, offset by the limit
    query = (db.query(Tag)
             .order_by(Tag.id)
             .offset(offset)
             .limit(limit)
             )

    tags = [x.as_json() for x in query]

    return {
        'tags': tags,
        'meta': {
            'total': db.query(Tag).count(),
        },
    }

@app.post('/api/tags')
def tags_post(db):
    if (not request.json or
            'tag' not in request.json or
            'name' not in request.json['tag']):
        abort(400, 'No name given')

    tag = Tag(name=request.json['tag']['name'])
    db.add(tag)
    db.commit()

    return {
        'tag': tag.as_json(),
    }


@app.get('/api/tags/<tag_id:int>')
def tags_get_one(tag_id, db):
    tag = db.query(Tag).filter_by(id=tag_id).first()
    if not tag:
        abort(404, '')

    return {
        'tag': tag.as_json(),
    }


@app.delete('/api/tags/<tag_id:int>')
def tags_delete_one(tag_id, db):
    tag = db.query(Tag).filter_by(id=tag_id).first()
    if not tag:
        abort(404, '')

    tag.delete_instance()
    response.status = 204
