from flask import render_template

from . import app
from . import models as m


@app.route("/")
@app.route("/index")
def index():
    return render_template('index.html')


@app.route("/stats")
def stats():
    stats = {
        'num_documents': m.Document.query.count(),
    }

    return render_template('stats.html', stats=stats)
