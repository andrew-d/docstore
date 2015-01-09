import falcon

from .hooks import check_media_type
from .routes.item import ItemResource, ItemsResource


api = falcon.API(before=[check_media_type])
api.add_route('/api/items',                     ItemsResource())
api.add_route('/api/items/{item_id}',           ItemResource())
