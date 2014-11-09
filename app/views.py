import hashlib
import datetime

from flask import render_template, request, redirect, url_for, flash
from flask.ext.uploads import UploadNotAllowed

from . import app, db, uploads
from . import models as m
from . import forms as f


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
    try:
        return int(request.args.get('page') or 1)
    except ValueError:
        # TODO
        raise


@app.route("/")
@app.route("/index")
def index():
    return render_template('index.html')


@app.route("/stats")
def stats():
    stats = {
        'num_documents': m.Document.query.count(),
        'num_tags':      m.Tag.query.count(),
    }

    return render_template('stats.html', stats=stats)


@app.route("/documents", methods=['GET', 'POST'])
def documents():
    upload_form = f.UploadDocument()

    if upload_form.validate_on_submit():
        uploaded = upload_form.file.data

        # Hash the file to get the filename we're using
        fhash = hashlib.sha256()
        fsize = 0
        try:
            for chunk in read_in_chunks(uploaded):
                fhash.update(chunk)
                fsize += len(chunk)
        finally:
            uploaded.seek(0)

        app.logger.info("Got uploaded file '%s' with hash: %s",
                        uploaded.filename, fhash.hexdigest())

        # Check if this document exists already.
        _, ext = uploaded.filename.rsplit('.', 1)
        existing = (m.Document.query.
                    filter_by(filename=fhash.hexdigest()+'.'+ext).
                    first())
        if existing is not None:
            flash('This document already exists', 'warning')
            # TODO: redirect to this document directly
            return redirect(url_for('documents'))

        # Save with the hash as a filename.  Appending the '.' means that
        # Flask-Uploads will automatically fill in the extension.
        try:
            fname = uploads.save(uploaded, name=fhash.hexdigest() + '.')
        except UploadNotAllowed:
            flash('Upload was not allowed', 'error')
            return redirect(url_for('documents'))

        newdoc = m.Document(name=upload_form.name.data,
                            filename=fname,
                            file_size=fsize)

        if upload_form.tags.data:
            # TODO: proper shell splitting with quotes
            new_tags = upload_form.tags.data.split(' ')
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

        flash('Uploaded new document', 'success')
        return redirect(url_for('documents'))

    page = get_page()
    documents = (m.Document.query.
                    order_by(m.Document.created).
                    paginate(page, error_out=False)
                 )

    return render_template('documents.html',
                           upload_form=upload_form,
                           documents=documents)


@app.route("/documents/<int:id>", methods=['GET', 'POST'])
def single_document(id):
    doc = m.Document.query.get_or_404(id)
    return render_template('single_document.html',
                           document=doc)


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
