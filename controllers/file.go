package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/lann/squirrel"
	"github.com/zenazn/goji/web"

	"github.com/andrew-d/docstore/models"
)

type FileController struct {
	AppController
}

func (c *FileController) GetAll(ctx web.C, w http.ResponseWriter, r *http.Request) error {
	// TODO
}

func (c *FileController) GetOne(ctx web.C, w http.ResponseWriter, r *http.Request) error {
	// TODO
}

func (c *FileController) Create(ctx web.C, w http.ResponseWriter, r *http.Request) error {
	// TODO
}
