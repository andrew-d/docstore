from flask import Blueprint, current_app, request
from flask.ext import restful
from flask.ext.restful import reqparse

from .models import db
from . import models as m


api_app = Blueprint('api_app', __name__, url_prefix='/api')
api = restful.Api(api_app)


class DocumentListAPI(restful.Resource):
    def get(self):
        """Get list of all documents"""
        parser = reqparse.RequestParser()
        parser.add_argument('page', location='args', type=int, default=1)
        parser.add_argument('per_page', location='args', type=int, default=10)
        args = parser.parse_args()

        paginated = m.Document.query.paginate(args['page'],
                                              per_page=args['per_page'],
                                              error_out=False)
        documents = [document.as_json() for document in paginated.items]
        tags = [t.as_json() for doc in paginated.items for t in doc.tags]

        return {
            'documents': documents,
            'tags': tags,
            'meta': {
                'total_pages': paginated.pages,
            },
        }

    def post(self):
        """Create a new document"""
        restful.abort(501,
                      message='Creation not supported - '
                              'use the explicit scan/upload endpoints')


class DocumentAPI(restful.Resource):
    def get(self, id):
        """Get a specific document"""
        doc = m.Document.query.get(id)
        if not doc:
            restful.abort(404, message='Document not found')

        return {
            'document': doc.as_json(),
            'tags': [t.as_json() for t in doc.tags],
        }

    def put(self, id):
        """Update a specific document"""
        doc = m.Document.query.get(id)
        if not doc:
            restful.abort(404, message='Document not found')

        parser = reqparse.RequestParser()
        parser.add_argument('document', location='json', type=dict, required=True)
        args = parser.parse_args()

        try:
            doc.name = args['document']['name']
            doc.meta = args['document']['meta']

            tags = []
            for t in args['document']['tags']:
                tag = m.Tag.query.get(t)
                if not tag:
                    restful.abort(404, message='Tag "%d" not found' % (t,))

                tags.append(tag)

            doc.tags = tags
        except KeyError:
            restful.abort(400)

        db.session.add(doc)
        db.session.commit()
        return {
            'document': doc.as_json(),
            'tags': [t.as_json() for t in doc.tags],
        }


    def delete(self, id):
        """Delete a specific document"""
        doc = m.Document.query.get(id)
        if not doc:
            restful.abort(404, message='Document not found')

        db.session.delete(doc)
        db.session.commit()
        return {}


class TagListAPI(restful.Resource):
    def get(self):
        """Get list of all tags"""
        parser = reqparse.RequestParser()
        parser.add_argument('page', type=int, default=1)
        parser.add_argument('per_page', type=int, default=10)
        args = parser.parse_args()

        paginated = m.Tag.query.paginate(args['page'],
                                         per_page=args['per_page'],
                                         error_out=False)
        tags = [tag.as_json() for tag in paginated.items]
        return {
            'tags': tags,
            'meta': {
                'total_pages': paginated.pages,
            },
        }

    def post(self):
        """Create a new tag"""
        # TODO: implement me
        restful.abort(501)


class TagAPI(restful.Resource):
    def get(self, id):
        """Get a specific tag"""
        tag = m.Tag.query.get(id)
        if not tag:
            restful.abort(404, message='Tag not found')

        return {'tag': tag.as_json()}

    def put(self, id):
        """Update a specific tag"""
        tag = m.Tag.query.get(id)
        if not tag:
            restful.abort(404, message='Tag not found')

        parser = reqparse.RequestParser()
        parser.add_argument('tag', location='json', type=dict, required=True)
        args = parser.parse_args()

        try:
            tag.name = args['tag']['name']
        except KeyError:
            restful.abort(400)

        db.session.add(tag)
        db.session.commit()
        return {'tag': tag.as_json()}

    def delete(self, id):
        """Delete a specific tag"""
        tag = m.Tag.query.get(id)
        if not tag:
            restful.abort(404, message='Tag not found')

        db.session.delete(tag)
        db.session.commit()
        return {}


class CapabilitiesAPI(restful.Resource):
    def get(self):
        return {
            'scanning': current_app.config['ENABLE_SCANNING'],
            'ocr':      current_app.config['ENABLE_OCR'],
            # TODO: need to get this from system, not hard-coded
            'scanners': [
                {'id': 1, 'name': 'foo'},
                {'id': 2, 'name': 'bar'},
            ],
        }



api.add_resource(DocumentListAPI, "/documents")
api.add_resource(DocumentAPI, "/documents/<int:id>")
api.add_resource(TagListAPI, "/tags")
api.add_resource(TagAPI, "/tags/<int:id>")
api.add_resource(CapabilitiesAPI, "/capabilities")
