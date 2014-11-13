from flask.ext.wtf import Form
from flask.ext.wtf.file import FileField, FileRequired, FileAllowed
from wtforms import StringField, SubmitField, SelectField
from wtforms.validators import Required, DataRequired


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


class AddTagsForm(Form):
    tags = StringField('Tags to apply', validators=[DataRequired()])
    add = SubmitField('Add Tags')


class UploadFileForm(Form):
    file = FileField('file', validators=[
        FileAllowed(ALLOWED_EXTS, 'Invalid document type'),
    ])
    upload = SubmitField('Upload')


class ScanFileForm(Form):
    scanner_name = SelectField("Scanner Name", validators=[
        RequiredIf("scan"),
    ])
    scan = SubmitField('Scan New')


class UploadDocumentForm(UploadFileForm):
    """
    Extension of uploading a file that also requires us to specify a document
    name, and allows applying tags.
    """
    name = StringField('Document name', validators=[
        DataRequired(),
    ])
    tags = StringField('Tags to apply', validators=[])


class ScanDocumentForm(ScanFileForm):
    name = StringField('Document name', validators=[
        DataRequired(),
    ])
    tags = StringField('Tags to apply', validators=[])
