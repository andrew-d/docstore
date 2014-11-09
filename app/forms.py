from flask.ext.wtf import Form
from flask.ext.wtf.file import FileField, FileRequired, FileAllowed
from wtforms import StringField
from wtforms.validators import DataRequired

from . import uploads


class UploadDocument(Form):
    name = StringField('Document name', validators=[
        DataRequired(),
    ])
    file = FileField('file', validators=[
        FileRequired(),
        FileAllowed(uploads, 'Invalid document type'),
    ])
    tags = StringField('Tags to apply', validators=[])
