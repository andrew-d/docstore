import falcon

from .hooks import check_media_type
#from .routes.item import ItemResource, ItemsResource
from .routes.file import FileResource, FilesResource, FileContentResource


api = falcon.API(before=[check_media_type])
#api.add_route('/api/items',                     ItemsResource())
#api.add_route('/api/items/{item_id}',           ItemResource())
api.add_route('/api/files',                     FilesResource())
api.add_route('/api/files/{file_id}',           FileResource())
api.add_route('/api/files/{file_id}/content',   FileContentResource())
