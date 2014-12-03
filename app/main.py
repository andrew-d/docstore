import os

from flask import Flask


app = Flask('docstore')


def setup_app():
    # Generated configuration
    app.config['UPLOAD_FOLDER'] = os.path.join(app.config['DATA_DIRECTORY'],
                                               'uploads')
    app.config['INDEX_PATH'] = os.path.join(app.config['DATA_DIRECTORY'], 'search')

    # Other configuration
    from .models import db
    db.init_app(app)

    from .api import api_app
    app.register_blueprint(api_app)

    ###app.config['scanners'] = scanner.Config

    ###from .search import open_index
    ###app.config['index'] = open_index()
