from bottle import error

# Import all the routes.  As each route is imported, they will register their
# handlers on the app instance.
from .routes import file
#from .routes import collection
from .routes import tag


@error(404)
def error404(error):
    return 'Not Found'
