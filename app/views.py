import os
import hashlib
import datetime
import mimetypes

from flask import (
    abort,
    flash,
    redirect,
    render_template,
    request,
    send_from_directory,
    url_for
)
from werkzeug import secure_filename

from . import app, db
from . import models as m
from . import forms as f
from . import scanner


def read_in_chunks(infile, chunk_size=1*1024*1024):
    while True:
        chunk = infile.read(chunk_size)
        if chunk:
            yield chunk

        else:
            # The chunk was empty, which means we're at the end
            # of the file
            return


def get_page():
    """
    Helper function to get the current page # from a request.
    """
    try:
        return int(request.args.get('page') or 1)
    except ValueError:
        abort(400)


@app.route("/")
@app.route("/index")
def index():
    return render_template('index.html',)


@app.errorhandler(400)
def bad_request(error):
    return render_template('error.html', error=error), 400


@app.errorhandler(404)
def not_found(error):
    return render_template('error.html', error=error), 404


@app.errorhandler(500)
def internal_server_error(error):
    return render_template('error.html', error=error), 500


@app.route("/stats")
def stats():
    stats = {
        'num_documents': m.Document.query.count(),
        'num_tags':      m.Tag.query.count(),
    }

    return render_template('stats.html',
                           stats=stats,
                           scanners=app.config['scanners'].scanner_info)



def handle_uploaded_file(name, f, tags=None):
    """
    Handle an uploaded file (and an optional list of tags to apply).
    Returns (document, exists), where 'document' is the new Document
    object, and 'exists' indicates whether or not it already existed.
    'document' will be None on error.
    """
    # Hash the file to get the filename we're using
    fhash = hashlib.sha256()
    fsize = 0
    try:
        for chunk in read_in_chunks(f):
            fhash.update(chunk)
            fsize += len(chunk)
    finally:
        f.seek(0)

    app.logger.info("Got uploaded file '%s' with hash: %s",
                    f.filename, fhash.hexdigest())

    # Generate needed filenames
    sfname = secure_filename(f.filename)
    _, ext = os.path.splitext(sfname)
    fname = fhash.hexdigest() + ext

    # Check if this document exists already.
    existing = m.Document.query.filter_by(filename=fname).first()
    if existing is not None:
        return existing, True

    # Save the file in our file storage, create database object
    f.save(os.path.join(app.config['UPLOAD_FOLDER'], fname))
    newdoc = m.Document(name=name,
                        filename=fname,
                        file_size=fsize)

    if tags:
        # TODO: proper shell splitting with quotes
        new_tags = tags.split(' ')
        for tagname in new_tags:
            t = m.Tag.get_or_create(tagname)
            newdoc.tags.append(t)

    # Add a tag for the current year
    now = datetime.datetime.utcnow().year
    year_tag = m.Tag.get_or_create('year:' + str(now))
    if year_tag not in newdoc.tags:
        newdoc.tags.append(year_tag)

    # All set - commit everything
    db.session.commit()
    return newdoc, False


def scan_document(name, scanner, tags=None):
    """
    As per handle_uploaded_file, except that we retrieve an image
    from an attached scanner.
    """
    data, fhash, fsize = scanner.scan_image()
    fname = fhash + '.png'

    # Check if this document exists already.
    existing = m.Document.query.filter_by(filename=fname).first()
    if existing is not None:
        return existing, True

    # Save the file in our file storage, create database object
    scanned_image.save(
        os.path.join(app.config['UPLOAD_FOLDER'], fname),
        'PNG'
    )
    newdoc = m.Document(name=name,
                        filename=fname,
                        file_size=fsize)

    if tags:
        # TODO: proper shell splitting with quotes
        new_tags = tags.split(' ')
        for tagname in new_tags:
            t = m.Tag.get_or_create(tagname)
            newdoc.tags.append(t)

    # Add a tag for the current year
    now = datetime.datetime.utcnow().year
    year_tag = m.Tag.get_or_create('year:' + str(now))
    if year_tag not in newdoc.tags:
        newdoc.tags.append(year_tag)

    # All set - commit everything
    db.session.commit()
    return newdoc, False


@app.route("/documents", methods=['GET', 'POST'])
def documents():
    upload_form = f.UploadDocument()

    # Fill in the choices for scanner names
    upload_form.scanner_name.choices = [
        (d.name, d.vendor + ' - ' + d.model)
        for d in app.config['scanners'].scanner_info
    ]

    if upload_form.validate_on_submit():
        if upload_form.upload.data:
            doc, exists = handle_uploaded_file(upload_form.name.data,
                                               upload_form.file.data,
                                               tags=upload_form.tags.data)
            if exists is True:
                flash('This document already exists', 'warning')
                return redirect(url_for('single_document', id=doc.id))

            if doc is None:
                flash('Upload was not allowed', 'error')
                return redirect(url_for('documents'))

            flash('Uploaded new document', 'success')
            return redirect(url_for('single_document', id=doc.id))

        elif upload_form.scan.data:
            scanner_name = upload_form.scanner_name.data
            app.logger.info("Scanning document with scanner %s", scanner_name)

            doc, exists = scan_document(upload_form.name.data,
                                        upload_form.scanner_name.data,
                                        tags=upload_form.tags.data)
            if exists is True:
                flash('This document already exists', 'warning')
                return redirect(url_for('single_document', id=doc.id))

            if doc is None:
                flash('Could not scan the document', 'error')
                return redirect(url_for('documents'))

            flash('Scanned new document', 'success')
            return redirect(url_for('single_document', id=doc.id))



    page = get_page()
    documents = (m.Document.query.
                    order_by(m.Document.created).
                    paginate(page, error_out=False)
                 )

    return render_template('documents.html',
                           upload_form=upload_form,
                           documents=documents,
                           have_scanner=app.config['scanners'].have_scanner)


@app.route("/documents/<int:id>", methods=['GET', 'POST'])
def single_document(id):
    doc = m.Document.query.get_or_404(id)
    return render_template('single_document.html',
                           document=doc)


@app.route("/documents/<int:id>/content")
def single_document_data(id):
    doc = m.Document.query.get_or_404(id)
    return send_from_directory(app.config['UPLOAD_FOLDER'], doc.filename)


@app.route("/tags")
def tags():
    page = get_page()
    tags = (m.Tag.query.
            order_by(m.Tag.name).
            paginate(page, error_out=False)
           )
    return render_template('tags.html',
                           tags=tags)


@app.route("/tags/<int:id>")
def single_tag(id):
    page = get_page()
    tag = m.Tag.query.get_or_404(id)
    tag_documents = (tag.documents.
                     order_by(m.Document.id).
                     paginate(page, error_out=False)
                    )
    return render_template('single_tag.html',
                           tag=tag,
                           tag_documents=tag_documents)