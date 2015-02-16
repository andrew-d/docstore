package controllers

import (
	"net/http"

	"github.com/lann/squirrel"
	"github.com/zenazn/goji/web"

	"github.com/andrew-d/docstore/models"
)

type FileController struct {
	AppController
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

/*
func (c *FileController) Create(ctx web.C, w http.ResponseWriter, r *http.Request) error {
	var doc *models.Document
	var err error

	if doc, err = c.getDocument(ctx); err != nil {
		return err
	}

	var createParams struct {
		// TODO
	}

	if err := json.NewDecoder(r.Body).Decode(&createParams); err != nil {
		return VError{err, "invalid body JSON", http.StatusBadRequest}
	}

	return nil
}
*/
