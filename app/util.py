import humanize
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
    A Jinja2 global function that uses moment.js to make pretty times.
    Note that, while we *could* do this with humanize, moment.js allows us to
    ensure that the times are local to the user's browser (e.g. if we wanted
    to run docstore on a server not in the same time zone).

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

@app.template_filter('humansize')
def humansize(s, binary=True):
    return humanize.naturalsize(s, binary=binary)
