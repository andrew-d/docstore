package controllers

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Sirupsen/logrus"
	"github.com/lann/squirrel"
	"github.com/tjgq/sane"
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

	c.JSON(w, http.StatusOK, M{
		"files": allFiles,
	})
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

	c.JSON(w, http.StatusOK, M{"file": f})
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

	// Actually save the data
	id, err := c.newFromBytes(doc, createParams.Filename, []byte(createParams.Data))
	if err != nil {
		return err
	}

	// Set the location of this file, return an empty response
	w.Header().Set("Location", fmt.Sprintf("/api/documents/%d/files/%d", doc.Id, id))
	w.WriteHeader(201)
	return nil
}

func (c *FileController) Scan(ctx web.C, w http.ResponseWriter, r *http.Request) error {
	var doc *models.Document
	var err error

	if doc, err = c.getDocument(ctx); err != nil {
		return err
	}

	// TODO: do we even need this filename
	var createParams struct {
		Filename string `schema:"filename"`
	}

	if err = r.ParseForm(); err != nil {
		return VError{err, "could not parse form body", 400}
	}

	if err = c.Decoder.Decode(&createParams, r.PostForm); err != nil {
		return VError{err, "could not decode form contents", 400}
	}

	// Scan an image
	conn, err := sane.Open("DEVICE")
	if err != nil {
		return VError{err, "could not open scanner", 500}
	}
	defer conn.Close()

	img, err := conn.ReadImage()
	if err != nil {
		return VError{err, "could not scan image", 500}
	}

	// Save the image to a byte array
	var buf bytes.Buffer
	err = png.Encode(&buf, img)
	if err != nil {
		return VError{err, "error encoding image", 500}
	}

	// Actually save the image.
	id, err := c.newFromBytes(doc, createParams.Filename, buf.Bytes())
	if err != nil {
		return err
	}

	// Set the location of this file, return an empty response
	w.Header().Set("Location", fmt.Sprintf("/api/documents/%d/files/%d", doc.Id, id))
	w.WriteHeader(201)
	return nil
}

func (c *FileController) newFromBytes(doc *models.Document, filename string, data []byte) (int64, error) {
	// Get the file hash
	hashBytes := sha256.Sum256(data)
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
		// TODO: handle ErrExist
		return -1, err
	}
	defer f.Close()

	if _, err = f.Write(data); err != nil {
		// TODO: should remove partially-written file
		return -1, err
	}

	// File exists.  Create the file entry.
	sql, args, _ := (c.Builder.
		Insert("files").
		Columns("name", "hash", "document_id").
		Values(filename, hash, doc.Id).
		ToSql())
	ret, err := c.DB.Exec(c.iQuery(sql), args...)
	if err != nil {
		return -1, VError{err, "error saving document", http.StatusInternalServerError}
	}

	id, _ := ret.LastInsertId()
	return id, nil

}
