from jinja2 import Markup

from . import app


@app.template_filter('pluralize')
def pluralize(number, singular = '', plural = 's'):
    if number == 1:
        return singular
    else:
        return plural


@app.template_global('momentjs')
class momentjs(object):
    """
    A global function that uses moment.js to make pretty times.

    Taken from:
        http://blog.miguelgrinberg.com/post/the-flask-mega-tutorial-part-xiii-dates-and-times
    """
    def __init__(self, timestamp):
        self.timestamp = timestamp

    def render(self, fmt):
        s = "<script>\ndocument.write(moment(\"%s\").%s);\n</script>" % (
            self.timestamp.strftime("%Y-%m-%dT%H:%M:%S Z"),
            fmt
        )
        return Markup(s)

    def format(self, fmt):
        return self.render("format(\"%s\")" % fmt)

    def calendar(self):
        return self.render("calendar()")

    def fromNow(self):
        return self.render("fromNow()")
