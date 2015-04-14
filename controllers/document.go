package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/lann/squirrel"
	"github.com/zenazn/goji/web"

	"github.com/andrew-d/docstore/daos"
)

type DocumentController struct {
	AppController
}

func (c *DocumentController) GetAll(ctx web.C, w http.ResponseWriter, r *http.Request) error {
	var queryParams struct {
		Limit  uint64 `schema:"limit"`
		Offset uint64 `schema:"offset"`
	}

	// Default value
	queryParams.Limit = 20

	if err := c.Decoder.Decode(&queryParams, r.URL.Query()); err != nil {
		return VError{err, "could not decode query parameters", 400}
	}

	// TODO: embed in struct?
	dao := daos.DocumentDAO{
		DB:      c.DB,
		Builder: c.Builder,
	}

	sql, args, _ := (c.Builder.
		Select("*").
		From("documents").
		Offset(queryParams.Offset).
		Limit(queryParams.Limit).
		ToSql())
	allDocuments, allTags, err := dao.LoadDocuments(sql, args...)
	if err != nil {
		return VError{err, "error getting documents", http.StatusInternalServerError}
	}

	// TODO: render document collections
	c.JSON(w, http.StatusOK, M{
		"documents": allDocuments,
		"tags":      allTags,
	})
	return nil
}

func (c *DocumentController) GetOne(ctx web.C, w http.ResponseWriter, r *http.Request) error {
	id, err := c.parseIntParam(ctx, "document_id")
	if err != nil {
		return err
	}

	// TODO: embed in struct?
	dao := daos.DocumentDAO{
		DB:      c.DB,
		Builder: c.Builder,
	}

	// Load document
	sql, args, _ := (c.Builder.
		Select("*").
		From("documents").
		Where(squirrel.Eq{"id": id}).
		ToSql())
	doc, tags, err := dao.LoadDocument(sql, args...)
	if err != nil {
		return VError{err, "error loading document", 500}
	}
	if doc.Id == 0 {
		return VError{err, "document not found", 404}
	}

	// TODO: render collection

	// Return document
	c.JSON(w, http.StatusOK, M{
		"document": &doc, // TODO: this is a hack to get MarshalJSON working
		"tags":     tags,
	})
	return nil
}

func (c *DocumentController) Create(ctx web.C, w http.ResponseWriter, r *http.Request) error {
	var createParams struct {
		Name         string  `json:"name"`
		Tags         []int64 `json:"tags"`
		CollectionId int64   `json:"collection_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&createParams); err != nil {
		return VError{err, "invalid body JSON", http.StatusBadRequest}
	}

	if len(createParams.Name) < 5 {
		return VError{nil, "document name must be longer than 5 characters", 400}
	}

	// TODO: embed in struct?
	dao := daos.DocumentDAO{
		DB:      c.DB,
		Builder: c.Builder,
	}

	doc, tags, err := dao.CreateDocument(createParams.Name,
		createParams.Tags,
		createParams.CollectionId)
	if err != nil {
		return err
	}

	c.JSON(w, http.StatusOK, M{
		"document": &doc, // TODO: this is a hack to get MarshalJSON working
		"tags":     tags,
	})
	return nil
}

func (c *DocumentController) Delete(ctx web.C, w http.ResponseWriter, r *http.Request) error {
	id, err := c.parseIntParam(ctx, "document_id")
	if err != nil {
		return err
	}

	sql, args, _ := (c.Builder.
		Delete("documents").
		Where(squirrel.Eq{"id": id}).
		ToSql())
	ret, err := c.DB.Exec(sql, args...)
	if err != nil {
		return VError{err, "error deleting document", 500}
	}
	if rows, _ := ret.RowsAffected(); rows == 0 {
		return VError{err, "document not found", 404}
	}

	w.WriteHeader(204)
	return nil
}
