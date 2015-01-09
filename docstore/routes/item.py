import logging

import falcon

from ..base import BaseResource
from ..models import Item
from ..hooks import deserialize, serialize


LOG = logging.getLogger(__name__)


class ItemsResource(BaseResource):
    """
    /api/items
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

        # Get all items and all tags, offset by the limit
        query = (Item
                 .select()
                 .offset(offset)
                 .limit(limit)
                 .order_by(Item.id)
                 )

        items = [item.as_json() for item in query]

        req.context['doc'] = {
            'items': items,
            'meta': {
                'total': Item.select().count(),
            },
        }

    @falcon.before(deserialize)
    @falcon.after(serialize)
    def on_post(self, req, resp, doc):
        # Validate incoming fields
        # TODO

        # Create the new item
        item = Item.create()

        # Return the new item
        req.context['doc'] = {
            'item': item.as_json(),
        }


class ItemResource(BaseResource):
    """
    /api/items/{item_id}
    """
    @falcon.after(serialize)
    def on_get(self, req, resp, item_id):
        try:
            item_id = int(item_id)
        except ValueError:
            raise falcon.HTTPBadRequest('Invalid Item ID',
                                        'The item ID is not a valid integer.')

        try:
            item = Item.get(Item.id == item_id)
        except Item.DoesNotExist:
            raise falcon.HTTPNotFound()

        req.context['doc'] = {'item': item.as_json()}

    @falcon.before(deserialize)
    def on_put(self, req, resp, item_id, doc):
        try:
            item_id = int(item_id)
        except ValueError:
            raise falcon.HTTPBadRequest('Invalid Item ID',
                                        'The item ID is not a valid integer.')

        # TODO: implement
        raise falcon.HTTPError(falcon.HTTP_501,
                               'Not Implemented',
                               'TODO implement me')

    def on_delete(self, req, resp, item_id):
        try:
            item_id = int(item_id)
        except ValueError:
            raise falcon.HTTPBadRequest('Invalid Item ID',
                                        'The item ID is not a valid integer.')

        try:
            item = Item.get(Item.id == item_id)
        except Item.DoesNotExist:
            raise falcon.HTTPNotFound()

        item.delete_instance()
        resp.status = falcon.HTTP_204
