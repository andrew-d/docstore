from werkzeug import url_decode

class MethodRewriteMiddleware(object):
    """
    A middleware that allows overriding the HTTP method via a configurable
    query-string parameter.
    """
    VALID_METHODS = ['DELETE', 'GET', 'PATCH', 'POST', 'PUT']

    def __init__(self, app, query_param='_method'):
        self.app = app
        self.query_param = query_param

    def __call__(self, environ, start_response):
        if self.query_param in environ.get('QUERY_STRING', ''):
            args = url_decode(environ['QUERY_STRING'])
            method = args.get(self.query_param)

            if method and method in self.VALID_METHODS:
                method = method.encode('ascii', 'replace')
                environ['REQUEST_METHOD'] = method

        return self.app(environ, start_response)

