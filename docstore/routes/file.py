import logging

import falcon

from ..base import BaseResource
from ..models import File
from ..hooks import deserialize, serialize


LOG = logging.getLogger(__name__)


class FilesResource(BaseResource):
    """
    /api/files
    """
    @falcon.after(serialize)
    def on_get(self, req, resp):
        # Pagination
        try:
            offset = int(req.get_param('offset') or 0)
            limit = int(req.get_param('limit') or 20)
        except ValueError:
            offset = 0
            limit = 20

        # Get all files, offset by the limit
        query = (File
                 .select()
                 .offset(offset)
                 .limit(limit)
                 .order_by(File.id)
                 )

        files = [x.as_json() for x in query]

        req.context['doc'] = {
            'files': files,
            'meta': {
                'total': File.select().count(),
            },
        }

    def on_post(self, req, resp):
        raise falcon.HTTPError(falcon.HTTP_501,
                               'Not Implemented',
                               'TODO implement me')


class FileResource(BaseResource):
    """
    /api/files/{file_id}
    """
    @falcon.after(serialize)
    def on_get(self, req, resp, file_id):
        try:
            file_id = int(file_id)
        except ValueError:
            raise falcon.HTTPBadRequest('Invalid File ID',
                                        'The file ID is not a valid integer.')

        try:
            ff = File.get(File.id == file_id)
        except File.DoesNotExist:
            raise falcon.HTTPNotFound()

        req.context['doc'] = {'file': ff.as_json()}

    def on_delete(self, req, resp, file_id):
        try:
            file_id = int(file_id)
        except ValueError:
            raise falcon.HTTPBadRequest('Invalid File ID',
                                        'The file ID is not a valid integer.')

        try:
            ff = File.get(File.id == file_id)
        except File.DoesNotExist:
            raise falcon.HTTPNotFound()

        ff.delete_instance()
        resp.status = falcon.HTTP_204


class FileContentResource(BaseResource):
    def on_get(self, req, resp, file_id):
        try:
            file_id = int(file_id)
        except ValueError:
            raise falcon.HTTPBadRequest('Invalid File ID',
                                        'The file ID is not a valid integer.')

        try:
            ff = File.get(File.id == file_id)
        except File.DoesNotExist:
            raise falcon.HTTPNotFound()

        # TODO: serve the contents of the file
        # TODO: be sure to try out the WSGI sendfile, if it exists
        raise falcon.HTTPError(falcon.HTTP_501,
                               'Not Implemented',
                               'TODO implement me')
