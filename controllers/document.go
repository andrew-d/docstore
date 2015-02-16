package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/lann/squirrel"
	"github.com/zenazn/goji/web"

	"github.com/andrew-d/docstore/models"
)

type DocumentController struct {
	AppController
}

func (c *DocumentController) GetAll(ctx web.C, w http.ResponseWriter, r *http.Request) error {
	allDocuments := []models.Document{}
	err := c.DB.Select(&allDocuments, `SELECT * FROM documents ORDER BY id ASC`)
	if err != nil {
		return VError{err, "error getting documents", http.StatusInternalServerError}
	}

	c.JSON(w, http.StatusOK, allDocuments)
	return nil
}

func (c *DocumentController) GetOne(ctx web.C, w http.ResponseWriter, r *http.Request) error {
	id, err := c.ParseIntParam(ctx, "document_id")
	if err != nil {
		return err
	}

	var t models.Document
	sql, args, _ := c.Builder.Select("*").From("documents").Where(squirrel.Eq{"id": id}).ToSql()
	err = c.DB.Get(&t, sql, args...)
	if t.Id == 0 {
		return VError{err, "document not found", 404}
	}

	// TODO: Load this document's tags

	c.JSON(w, http.StatusOK, t)
	return nil
}

func (c *DocumentController) Create(ctx web.C, w http.ResponseWriter, r *http.Request) error {
	var createParams struct {
		Tags []int64 `json:"tags"`
	}

	if err := json.NewDecoder(r.Body).Decode(&createParams); err != nil {
		return VError{err, "invalid body JSON", http.StatusBadRequest}
	}

	createdAt := time.Now().UTC().Unix()

	// Insert everything in a transaction
	tx := c.DB.MustBegin()
	defer tx.Rollback()

	// TODO: collection
	ret, err := tx.NamedExec(c.iQuery(`INSERT INTO documents (created_at, collection_id) VALUES (:created_at, :collection_id)`),
		map[string]interface{}{
			"created_at":    createdAt,
			"collection_id": 0,
		})
	if err != nil {
		return VError{err, "error saving document", http.StatusInternalServerError}
	}

	id, _ := ret.LastInsertId()

	// Insert all tags
	for _, tag_id := range createParams.Tags {
		var tag models.Tag

		// Check the tag exists
		sql, args, _ := c.Builder.Select("*").From("tags").Where(squirrel.Eq{"id": id}).ToSql()
		err = tx.Get(&tag, sql, args...)
		if tag.Id == 0 {
			return VError{err, fmt.Sprintf("tag %d not found", tag_id), http.StatusNotFound}
		}

		// Insert it
		_, err := tx.NamedExec(`INSERT INTO document_tags (document_id, tag_id) VALUES (:document_id, :tag_id)`,
			map[string]interface{}{
				"document_id": id,
				"tag_id":      tag_id,
			})
		if err != nil {
			return VError{err, "error saving document tag", http.StatusInternalServerError}
		}
	}

	// If we get here, our transaction succeeded
	err = tx.Commit()
	if err != nil {
		return VError{err, "error committing transaction", http.StatusInternalServerError}
	}

	// TODO: render document tags too
	c.JSON(w, http.StatusOK, models.Document{
		Id:           id,
		CreatedAt:    createdAt,
		CollectionId: 0, // TODO
	})
	return nil
}
