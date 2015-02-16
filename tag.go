package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/schema"
	"github.com/jinzhu/gorm"
	"github.com/zenazn/goji/web"
)

var _ = schema.NewDecoder

type Tag struct {
	Id        int64      `json:"id"`
	Name      string     `json:"name" validate:"nonzero" sql:"not null;unique"`
	Documents []Document `json:"documents,omitempty" gorm:"many2many:document_tags"`
}

func routeTagsGetAll(c web.C, w http.ResponseWriter, r *http.Request) {
	db := c.Env["db"].(gorm.DB)

	allTags := []Tag{}
	db.Find(&allTags)

	render.JSON(w, http.StatusOK, allTags)
}

func routeTagsCreate(c web.C, w http.ResponseWriter, r *http.Request) {
	db := c.Env["db"].(gorm.DB)

	t := &Tag{}
	if err := json.NewDecoder(r.Body).Decode(t); err != nil {
		render.JSON(w, http.StatusBadRequest, map[string]interface{}{
			"error":   err.Error(),
			"message": "invalid decoding body JSON",
		})
		return
	}
	if errs := v.ValidateAndTag(t, "json"); errs != nil {
		render.JSON(w, http.StatusBadRequest, map[string]interface{}{
			"error":   errorsToStrings(errs),
			"message": "invalid input",
		})
		return
	}

	db.Create(t)
	render.JSON(w, http.StatusOK, t)
}

func routeTagsGetOne(c web.C, w http.ResponseWriter, r *http.Request) {
	db := c.Env["db"].(gorm.DB)

	id, err := parseIntParam(c, "tag_id")
	if err != nil {
		render.JSON(w, http.StatusBadRequest, map[string]interface{}{
			"error":   err,
			"message": "invalid tag",
		})
		return
	}

	var t Tag
	if db.First(&t, id).RecordNotFound() {
		render.JSON(w, http.StatusNotFound, map[string]interface{}{
			"error":   nil,
			"message": "tag not found",
		})
		return
	}

	// Load this tag's documents
	var tagDocuments []Document
	db.Model(&t).Association("Documents").Find(&tagDocuments)

	// Save on model to be rendered
	t.Documents = tagDocuments

	render.JSON(w, http.StatusOK, t)
}
