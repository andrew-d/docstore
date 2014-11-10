from flask.ext.wtf import Form
from flask.ext.wtf.file import FileField, FileRequired, FileAllowed
from wtforms import StringField, SubmitField, SelectField
from wtforms.validators import Required, DataRequired, Optional


class RequiredIf(Required):
    """
    Validate that a field exists only if another field is set and has
    a truthy value.  Adapted from::
        https://stackoverflow.com/questions/8463209/how-to-make-a-field-conditionally-optional-in-wtforms
    """

    def __init__(self, other_field, *args, **kwargs):
        self.other_field = other_field
        super(RequiredIf, self).__init__(*args, **kwargs)

    def __call__(self, form, field):
        other_field = form._fields.get(self.other_field)
        if other_field is None:
            raise Exception('no field named "%s" in form' % self.other_field)

        if bool(other_field.data):
            super(RequiredIf, self).__call__(form, field)


class FileRequiredIf(FileRequired):
    """
    As RequiredIf, but for FileFields instead.
    """

    def __init__(self, other_field, *args, **kwargs):
        self.other_field = other_field
        super(FileRequiredIf, self).__init__(*args, **kwargs)

    def __call__(self, form, field):
        other_field = form._fields.get(self.other_field)
        if other_field is None:
            raise Exception('no field named "%s" in form' % self.other_field)

        if bool(other_field.data):
            super(FileRequiredIf, self).__call__(form, field)


ALLOWED_EXTS = (
    'jpg gpeg png gif svg bmp'.split() +
    'rtf doc docx xls xlsx'.split() +
    'csv json xml yaml yml'.split() +
    'tar tgz gz bz2 xz zip 7z'.split()
)


class UploadDocument(Form):
    name = StringField('Document name', validators=[
        DataRequired(),
    ])
    tags = StringField('Tags to apply', validators=[])
    upload = SubmitField('Upload')
    scan = SubmitField('Scan New')
    file = FileField('file', validators=[
        FileRequiredIf(other_field='upload'),
        FileAllowed(ALLOWED_EXTS, 'Invalid document type'),
    ])
    scanner_name = SelectField("Scanner Name", validators=[
        RequiredIf("scan"),
    ])
