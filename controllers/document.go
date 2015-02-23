package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lann/squirrel"
	"github.com/zenazn/goji/web"

	"github.com/andrew-d/docstore/models"
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

	allDocuments := []models.Document{}
	sql, args, _ := (c.Builder.
		Select("*").
		From("documents").
		Offset(queryParams.Offset).
		Limit(queryParams.Limit).
		ToSql())
	err := c.DB.Select(&allDocuments, sql, args...)
	if err != nil {
		return VError{err, "error getting documents", http.StatusInternalServerError}
	}

	// TODO: should we load associated tags for these documents?
	c.JSON(w, http.StatusOK, M{
		"documents": allDocuments,
	})
	return nil
}

func (c *DocumentController) GetOne(ctx web.C, w http.ResponseWriter, r *http.Request) error {
	id, err := c.parseIntParam(ctx, "document_id")
	if err != nil {
		return err
	}

	// Load document
	var doc models.Document
	sql, args, _ := (c.Builder.
		Select("*").
		From("documents").
		Where(squirrel.Eq{"id": id}).
		ToSql())
	err = c.DB.Get(&doc, sql, args...)
	if doc.Id == 0 {
		return VError{err, "document not found", 404}
	}

	// Load all tags for document
	var tags []models.Tag
	sql, args, _ = (c.Builder.
		Select("tags.id AS id", "name").
		From("tags").
		Join("document_tags ON document_tags.tag_id == tags.id").
		Join("documents ON document_tags.document_id == documents.id").
		Where(squirrel.Eq{"documents.id": doc.Id}).
		ToSql())
	fmt.Println("sql", sql)
	err = c.DB.Select(&tags, sql, args...)
	if err != nil {
		return VError{err, "error getting document's tags", http.StatusInternalServerError}
	}

	// TODO: load collection

	// Return document
	c.JSON(w, http.StatusOK, M{
		"document": doc.WithTags(tags),
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

	// Create document model
	ret := models.Document{
		Name:      createParams.Name,
		CreatedAt: time.Now().UTC().Unix(),
	}

	// If given, add document to collection
	if createParams.CollectionId > 0 {
		var coll models.Collection
		sql, args, _ := (c.Builder.
			Select("*").
			From("collections").
			Where(squirrel.Eq{"id": createParams.CollectionId}).
			ToSql())
		err := c.DB.Get(&coll, sql, args...)
		if coll.Id == 0 {
			return VError{err, fmt.Sprintf("collection %d not found", createParams.CollectionId),
				http.StatusNotFound}
		}

		ret.CollectionId.Int64 = coll.Id
		ret.CollectionId.Valid = true
	}

	// Insert everything in a transaction
	var tags []models.Tag
	err := c.inTransaction(func(tx *sqlx.Tx) error {
		// Insert the document
		sql, args, _ := (c.Builder.
			Insert("documents").
			Columns("name", "created_at", "collection_id").
			Values(ret.Name, ret.CreatedAt, ret.CollectionId).
			ToSql())
		sqlRes, err := tx.Exec(c.iQuery(sql), args...)
		if err != nil {
			return VError{err, "error saving document", http.StatusInternalServerError}
		}

		// Save the returned ID
		id, _ := sqlRes.LastInsertId()
		ret.Id = id

		// Insert all tags
		for _, tag_id := range createParams.Tags {
			var tag models.Tag

			// Check the tag exists
			sql, args, _ := (c.Builder.
				Select("*").
				From("tags").
				Where(squirrel.Eq{"id": id}).
				ToSql())
			err = tx.Get(&tag, sql, args...)
			if tag.Id == 0 {
				return VError{err, fmt.Sprintf("tag %d not found", tag_id),
					http.StatusNotFound}
			}

			// Save for rendering
			tags = append(tags, tag)

			// Insert it
			sql, args, _ = (c.Builder.
				Insert("document_tags").
				Columns("document_id", "tag_id").
				Values(id, tag_id).
				ToSql())
			_, err := tx.Exec(sql, args...)
			if err != nil {
				return VError{err, "error saving document tag",
					http.StatusInternalServerError}
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	c.JSON(w, http.StatusOK, M{
		"document": ret.WithTags(tags), // Attach tag IDs to document
		"tags":     tags,
	})
	return nil
}
