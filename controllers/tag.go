package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/lann/squirrel"
	"github.com/zenazn/goji/web"

	"github.com/andrew-d/docstore/models"
)

type TagController struct {
	AppController
}

func (c *TagController) GetAll(ctx web.C, w http.ResponseWriter, r *http.Request) error {
	var sql string
	var args []interface{}

	// If we have a 'ids[]' parameter, it limits what tags we load.
	query := r.URL.Query()
	if ids, ok := query["ids[]"]; ok {
		var tagIds []int64
		for _, id := range ids {
			i, err := strconv.ParseUint(id, 10, 64)
			if err != nil {
				continue
			}

			tagIds = append(tagIds, int64(i))
		}

		sql, args, _ = (c.Builder.
			Select("*").
			From("tags").
			Where(squirrel.Eq{"id": tagIds}).
			OrderBy("name ASC").
			ToSql())
	} else {
		sql = `SELECT * FROM tags ORDER BY name ASC`
		args = []interface{}{}
	}

	allTags := []models.Tag{}
	err := c.DB.Select(&allTags, sql, args...)
	if err != nil {
		return VError{err, "error getting tags", http.StatusInternalServerError}
	}

	empty := []int64{}
	for i := 0; i < len(allTags); i++ {
		allTags[i].Documents = empty
	}

	c.JSON(w, http.StatusOK, M{
		"tags": allTags,
	})
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
	t.Documents = []int64{}

	c.JSON(w, http.StatusOK, M{
		"tag": t,
	})
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
	c.JSON(w, http.StatusOK, M{
		"tag": models.Tag{
			Id:        id,
			Name:      createParams.Name,
			Documents: []int64{},
		},
	})
	return nil
}

func (c *TagController) Delete(ctx web.C, w http.ResponseWriter, r *http.Request) error {
	id, err := c.parseIntParam(ctx, "tag_id")
	if err != nil {
		return err
	}

	sql, args, _ := (c.Builder.
		Delete("tags").
		Where(squirrel.Eq{"id": id}).
		ToSql())
	ret, err := c.DB.Exec(sql, args...)
	if err != nil {
		return VError{err, "error deleting tag", 500}
	}
	if rows, _ := ret.RowsAffected(); rows == 0 {
		return VError{err, "tag not found", 404}
	}

	w.WriteHeader(204)
	return nil
}
