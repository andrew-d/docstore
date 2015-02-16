package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/zenazn/goji/web"
)

type Document struct {
	Id           int64     `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	Tags         []Tag     `json:"tags,omitempty" gorm:"many2many:document_tags"`
	CollectionId int64     `json:"collection_id"`

	// TODO: return collection vs. collection id?
}

func (d *Document) loadTags(db gorm.DB) {
	var documentTags []Tag
	db.Model(d).Association("Tags").Find(&documentTags)
	d.Tags = documentTags
}

// Return an array of all documents, without their tags
func routeDocumentsGetAll(c web.C, w http.ResponseWriter, r *http.Request) {
	db := c.Env["db"].(gorm.DB)

	var allDocuments []Document
	db.Find(&allDocuments)
	render.JSON(w, http.StatusOK, allDocuments)
}

// Create a single document.  Properly handles nested tags.
func routeDocumentsCreate(c web.C, w http.ResponseWriter, r *http.Request) {
	db := c.Env["db"].(gorm.DB)
	db.LogMode(true)
	defer db.LogMode(false)

	type createParams struct {
		Tags []int64 `json:"tags"`
	}

	var params createParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		render.JSON(w, http.StatusBadRequest, map[string]interface{}{
			"error":   err.Error(),
			"message": "invalid decoding body JSON",
		})
		return
	}

	// First, create the document
	d := Document{}
	db.Create(&d)

	// Append tags
	var tags []Tag
	for _, tag_id := range params.Tags {
		var tag Tag

		if db.First(&tag, tag_id).RecordNotFound() {
			render.JSON(w, http.StatusNotFound, map[string]interface{}{
				"error":   nil,
				"message": fmt.Sprintf("tag %d not found", tag_id),
			})
			return
		}

		tags = append(tags, tag)
	}
	db.Model(&d).Association("Tags").Append(tags)

	render.JSON(w, http.StatusOK, d)
}

// Get a single document, returning it along with its tags.
func routeDocumentsGetOne(c web.C, w http.ResponseWriter, r *http.Request) {
	db := c.Env["db"].(gorm.DB)

	id, err := parseIntParam(c, "document_id")
	if err != nil {
		render.JSON(w, http.StatusBadRequest, map[string]interface{}{
			"error":   err,
			"message": "invalid document",
		})
		return
	}

	var d Document
	if db.First(&d, id).RecordNotFound() {
		render.JSON(w, http.StatusNotFound, map[string]interface{}{
			"error":   nil,
			"message": "document not found",
		})
		return
	}

	d.loadTags(db)

	render.JSON(w, http.StatusOK, d)
}
