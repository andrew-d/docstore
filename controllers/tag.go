package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/lann/squirrel"
	"github.com/zenazn/goji/web"

	"github.com/andrew-d/docstore/models"
)

type TagController struct {
	AppController
}

func (c *TagController) GetAll(ctx web.C, w http.ResponseWriter, r *http.Request) error {
	allTags := []models.Tag{}
	err := c.DB.Select(&allTags, `SELECT * FROM tags ORDER BY name ASC`)
	if err != nil {
		return VError{err, "error getting tags", http.StatusInternalServerError}
	}

	c.JSON(w, http.StatusOK, allTags)
	return nil
}

func (c *TagController) GetOne(ctx web.C, w http.ResponseWriter, r *http.Request) error {
	id, err := c.parseIntParam(ctx, "tag_id")
	if err != nil {
		return err
	}

	var t models.Tag
	sql, args, _ := (c.Builder.
		Select("*").
		From("tags").
		Where(squirrel.Eq{"id": id}).
		ToSql())
	err = c.DB.Get(&t, sql, args...)
	if t.Id == 0 {
		return VError{err, "tag not found", 404}
	}

	// TODO: Load this tag's documents

	c.JSON(w, http.StatusOK, t)
	return nil
}

func (c *TagController) Create(ctx web.C, w http.ResponseWriter, r *http.Request) error {
	var createParams struct {
		Name string `json:"name" validate:"nonzero"`
	}

	if err := json.NewDecoder(r.Body).Decode(&createParams); err != nil {
		return VError{err, "invalid body JSON", http.StatusBadRequest}
	}
	if errs := v.ValidateAndTag(&createParams, "json"); errs != nil {
		return VError{errs, "invalid input", http.StatusBadRequest}
	}

	sql, args, _ := (c.Builder.
		Insert("tags").
		Columns("name").
		Values(createParams.Name).
		ToSql())
	ret, err := c.DB.Exec(c.iQuery(sql), args...)
	if err != nil {
		return VError{err, "error saving tag", http.StatusInternalServerError}
	}

	id, _ := ret.LastInsertId()
	c.JSON(w, http.StatusOK, models.Tag{
		Id:   id,
		Name: createParams.Name,
	})
	return nil
}
