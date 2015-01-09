import json

import falcon


def deserialize(req, resp, resource, params):
    """
    Convert the body from JSON when receiving a request, adding an
    additional 'doc' parameter to the function's arguments that
    contains the decoded JSON body.
    """
    body = req.stream.read()
    if not body:
        raise falcon.HTTPBadRequest('Empty request body',
                                    'A valid JSON document is required.')

    try:
        params['doc'] = json.loads(body.decode('utf-8'))
    except (ValueError, UnicodeDecodeError):
        raise falcon.HTTPBadRequest('Malformed JSON',
                                    'Could not decode the request body. The '
                                    'JSON was incorrect or not encoded as '
                                    'UTF-8.')


def serialize(req, resp, resource):
    """
    Set the response body of the request to be the value of the 'doc'
    key in the request context.
    """
    # TODO: escaping
    # TODO: prettyprint in development mode
    # TODO: prefix for JSON hijacking
    resp.body = json.dumps(req.context['doc'])


def check_media_type(req, resp, params):
    if not req.client_accepts_json:
        raise falcon.HTTPNotAcceptable(
            'This API only supports responses encoded as JSON.')

    if req.method in ('POST', 'PUT'):
        if not req.content_type == 'application/json':
            raise falcon.HTTPUnsupportedMediaType(
                'This API only supports requests encoded as JSON.')
