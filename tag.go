package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/schema"
	"github.com/jmoiron/sqlx"
	"github.com/lann/squirrel"
	"github.com/zenazn/goji/web"
)

var _ = schema.NewDecoder

type Tag struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

func routeTagsGetAll(c web.C, w http.ResponseWriter, r *http.Request) {
	db := c.Env["db"].(*sqlx.DB)

	allTags := []Tag{}
	err := db.Select(&allTags, `SELECT * FROM tags ORDER BY name ASC`)
	if err != nil {
		log.WithField("err", err).Error("Error getting tags")
		render.JSON(w, http.StatusInternalServerError, map[string]interface{}{
			"error":   err.Error(),
			"message": "error getting tags",
		})
		return
	}

	render.JSON(w, http.StatusOK, allTags)
}

func routeTagsCreate(c web.C, w http.ResponseWriter, r *http.Request) {
	db := c.Env["db"].(*sqlx.DB)

	var createParams struct {
		Name string `json:"name" validate:"nonzero"`
	}

	if err := json.NewDecoder(r.Body).Decode(&createParams); err != nil {
		render.JSON(w, http.StatusBadRequest, map[string]interface{}{
			"error":   err.Error(),
			"message": "invalid decoding body JSON",
		})
		return
	}
	if errs := v.ValidateAndTag(&createParams, "json"); errs != nil {
		render.JSON(w, http.StatusBadRequest, map[string]interface{}{
			"error":   errorsToStrings(errs),
			"message": "invalid input",
		})
		return
	}

	ret, err := db.NamedExec(iQuery(`INSERT INTO tags (name) VALUES (:name)`),
		map[string]interface{}{
			"name": createParams.Name,
		})
	if err != nil {
		log.WithField("err", err).Error("Error saving tag")
		render.JSON(w, http.StatusInternalServerError, map[string]interface{}{
			"error":   err.Error(),
			"message": "error saving tag",
		})
		return
	}

	id, _ := ret.LastInsertId()
	render.JSON(w, http.StatusOK, Tag{
		Id:   id,
		Name: createParams.Name,
	})
}

func routeTagsGetOne(c web.C, w http.ResponseWriter, r *http.Request) {
	db := c.Env["db"].(*sqlx.DB)

	id, err := parseIntParam(c, "tag_id")
	if err != nil {
		render.JSON(w, http.StatusBadRequest, map[string]interface{}{
			"error":   err.Error(),
			"message": "invalid tag id",
		})
		return
	}

	var t Tag
	sql, args, _ := sq.Select("*").From("tags").Where(squirrel.Eq{"id": id}).ToSql()
	err = db.Get(&t, sql, args...)
	if t.Id == 0 {
		render.JSON(w, http.StatusNotFound, map[string]interface{}{
			"error":   err.Error(),
			"message": "tag not found",
		})
		return
	}

	// TODO: Load this tag's documents

	render.JSON(w, http.StatusOK, t)
}
