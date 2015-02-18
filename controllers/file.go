package controllers

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Sirupsen/logrus"
	"github.com/lann/squirrel"
	"github.com/zenazn/goji/web"

	"github.com/andrew-d/docstore/models"
)

type FileController struct {
	AppController

	// Location to save files in
	FilePath string
}

func (c *FileController) getDocument(ctx web.C) (*models.Document, error) {
	id, err := c.parseIntParam(ctx, "document_id")
	if err != nil {
		return nil, err
	}

	var doc models.Document
	sql, args, _ := (c.Builder.
		Select("*").
		From("documents").
		Where(squirrel.Eq{"id": id}).
		ToSql())
	err = c.DB.Get(&doc, sql, args...)
	if doc.Id == 0 {
		return nil, VError{err, "document not found", 404}
	}

	return &doc, nil
}

func (c *FileController) GetAll(ctx web.C, w http.ResponseWriter, r *http.Request) error {
	var doc *models.Document
	var err error

	if doc, err = c.getDocument(ctx); err != nil {
		return err
	}

	allFiles := []models.File{}
	sql, args, _ := (c.Builder.
		Select("*").
		From("files").
		Where(squirrel.Eq{"document_id": doc.Id}).
		ToSql())
	err = c.DB.Select(&allFiles, sql, args...)
	if err != nil {
		return VError{err, "error getting files", http.StatusInternalServerError}
	}

	c.JSON(w, http.StatusOK, allFiles)
	return nil
}

func (c *FileController) GetOne(ctx web.C, w http.ResponseWriter, r *http.Request) error {
	id, err := c.parseIntParam(ctx, "file_id")
	if err != nil {
		return err
	}

	var doc *models.Document
	if doc, err = c.getDocument(ctx); err != nil {
		return err
	}

	var f models.File
	sql, args, _ := (c.Builder.
		Select("*").
		From("files").
		Where(squirrel.Eq{"document_id": doc.Id, "id": id}).
		ToSql())
	err = c.DB.Get(&f, sql, args...)
	if f.Id == 0 {
		return VError{err, "file not found", 404}
	}

	c.JSON(w, http.StatusOK, f)
	return nil
}

func (c *FileController) Content(ctx web.C, w http.ResponseWriter, r *http.Request) error {
	id, err := c.parseIntParam(ctx, "file_id")
	if err != nil {
		return err
	}

	var doc *models.Document
	if doc, err = c.getDocument(ctx); err != nil {
		return err
	}

	var f models.File
	sql, args, _ := (c.Builder.
		Select("*").
		From("files").
		Where(squirrel.Eq{"document_id": doc.Id, "id": id}).
		ToSql())
	err = c.DB.Get(&f, sql, args...)
	if f.Id == 0 {
		return VError{err, "file not found", 404}
	}

	// Open the file on-disk
	file, err := os.Open(filepath.Join(c.FilePath, f.Hash))
	if err != nil {
		if err == os.ErrNotExist {
			return VError{err, "file not found", 404}
		}

		return VError{err, "error opening file", 500}
	}
	defer file.Close()

	// TODO: set real mime type before writing
	w.Header().Set("Content-Type", "application/octet-stream")

	// Copy the file to the client
	n, err := io.Copy(w, file)
	if err != nil {
		// Note: we can't return an error here, since we've (possibly) already
		// started writing to the client.  The best we can do is log and move on.
		c.Logger.WithFields(logrus.Fields{
			"file_id":     id,
			"document_id": doc.Id,
			"err":         err,
			"bytesSent":   n,
		}).Error("Error sending file to client")
		return nil
	}

	return nil
}

func (c *FileController) Upload(ctx web.C, w http.ResponseWriter, r *http.Request) error {
	var doc *models.Document
	var err error

	if doc, err = c.getDocument(ctx); err != nil {
		return err
	}

	var createParams struct {
		Filename string `schema:"filename"`
		Data     string `schema:"data"`
	}

	if err = r.ParseForm(); err != nil {
		return VError{err, "could not parse form body", 400}
	}

	if err = c.Decoder.Decode(&createParams, r.PostForm); err != nil {
		return VError{err, "could not decode form contents", 400}
	}

	hashBytes := sha256.Sum256([]byte(createParams.Data))
	hash := hex.EncodeToString(hashBytes[:])

	writePath := filepath.Join(c.FilePath, hash)

	// Note: the following flags mean:
	//	- Read-Write access
	//	- Create the file if it doesn't exist
	//	- Fail if it *does* exist
	//
	// Since we store files by their hash, we can avoid writing the file to
	// disk if it's already there.
	f, err := os.OpenFile(writePath, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil && err != os.ErrExist {
		return err
	}
	defer f.Close()

	if _, err = io.WriteString(f, createParams.Data); err != nil {
		return err
	}

	// File exists - great.  Create the file entry.
	sql, args, _ := (c.Builder.
		Insert("files").
		Columns("name", "hash", "document_id").
		Values(createParams.Filename, hash, doc.Id).
		ToSql())
	ret, err := c.DB.Exec(c.iQuery(sql), args...)
	if err != nil {
		return VError{err, "error saving document", http.StatusInternalServerError}
	}

	id, _ := ret.LastInsertId()

	// Set the location of this file, return an empty response
	w.Header().Set("Location", fmt.Sprintf("/api/documents/%d/files/%d", doc.Id, id))
	w.WriteHeader(201)
	return nil
}
