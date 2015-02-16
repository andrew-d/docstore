package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lann/squirrel"
	"github.com/zenazn/goji/web"
)

type Document struct {
	Id           int64 `json:"id"`
	CreatedAt    int64 `json:"created_at" db:"created_at"`
	CollectionId int64 `json:"collection_id" db:"collection_id"`
}

// Return an array of all documents, without their tags
func routeDocumentsGetAll(c web.C, w http.ResponseWriter, r *http.Request) {
	db := c.Env["db"].(*sqlx.DB)

	allDocuments := []Document{}
	err := db.Select(&allDocuments, `SELECT * FROM documents ORDER BY id ASC`)
	if err != nil {
		log.WithField("err", err).Error("Error getting documents")
		render.JSON(w, http.StatusInternalServerError, map[string]interface{}{
			"error":   err.Error(),
			"message": "error getting documents",
		})
		return
	}

	render.JSON(w, http.StatusOK, allDocuments)
}

// Create a single document.  Properly handles nested tags.
func routeDocumentsCreate(c web.C, w http.ResponseWriter, r *http.Request) {
	db := c.Env["db"].(*sqlx.DB)

	var createParams struct {
		Tags []int64 `json:"tags"`
	}

	if err := json.NewDecoder(r.Body).Decode(&createParams); err != nil {
		render.JSON(w, http.StatusBadRequest, map[string]interface{}{
			"error":   err.Error(),
			"message": "invalid decoding body JSON",
		})
		return
	}

	createdAt := time.Now().UTC().Unix()

	// Insert everything in a transaction
	tx := db.MustBegin()
	defer tx.Rollback()

	// TODO: collection
	ret, err := tx.NamedExec(iQuery(`INSERT INTO documents (created_at, collection_id) VALUES (:created_at, :collection_id)`),
		map[string]interface{}{
			"created_at":    createdAt,
			"collection_id": 0,
		})
	if err != nil {
		log.WithField("err", err).Error("Error saving document")
		render.JSON(w, http.StatusInternalServerError, map[string]interface{}{
			"error":   err.Error(),
			"message": "error saving document",
		})
		return
	}

	id, _ := ret.LastInsertId()

	// Insert all tags
	for _, tag_id := range createParams.Tags {
		var tag Tag

		// Check the tag exists
		sql, args, _ := sq.Select("*").From("tags").Where(squirrel.Eq{"id": id}).ToSql()
		err = tx.Get(&tag, sql, args...)
		if tag.Id == 0 {
			render.JSON(w, http.StatusNotFound, map[string]interface{}{
				"error":   err.Error(),
				"message": fmt.Sprintf("tag %d not found", tag_id),
			})
			return
		}

		// Insert it
		_, err := tx.NamedExec(`INSERT INTO document_tags (document_id, tag_id) VALUES (:document_id, :tag_id)`,
			map[string]interface{}{
				"document_id": id,
				"tag_id":      tag_id,
			})
		if err != nil {
			log.WithField("err", err).Error("Error saving document tag")
			render.JSON(w, http.StatusInternalServerError, map[string]interface{}{
				"error":   err.Error(),
				"message": "error saving document tag",
			})
			return
		}
	}

	// If we get here, our transaction succeeded
	err = tx.Commit()
	if err != nil {
		log.WithField("err", err).Error("Error committing transaction")
		render.JSON(w, http.StatusInternalServerError, map[string]interface{}{
			"error":   err.Error(),
			"message": "error committing transaction",
		})
		return
	}

	// TODO: render document tags too
	render.JSON(w, http.StatusOK, Document{
		Id:           id,
		CreatedAt:    createdAt,
		CollectionId: 0, // TODO
	})
}

// Get a single document, returning it along with its tags.
func routeDocumentsGetOne(c web.C, w http.ResponseWriter, r *http.Request) {
	db := c.Env["db"].(*sqlx.DB)

	id, err := parseIntParam(c, "document_id")
	if err != nil {
		render.JSON(w, http.StatusBadRequest, map[string]interface{}{
			"error":   err.Error(),
			"message": "invalid document id",
		})
		return
	}

	var doc Document
	sql, args, _ := sq.Select("*").From("documents").Where(squirrel.Eq{"id": id}).ToSql()
	err = db.Get(&doc, sql, args...)
	if doc.Id == 0 {
		render.JSON(w, http.StatusNotFound, map[string]interface{}{
			"error":   err.Error(),
			"message": "document not found",
		})
		return
	}

	// TODO: Load this document's tags

	render.JSON(w, http.StatusOK, doc)
}
