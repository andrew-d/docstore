package api

import (
	"net/http"

	"github.com/andrew-d/docstore/models"
	"github.com/andrew-d/docstore/services"
)

func handleDocumentList(w http.ResponseWriter, r *http.Request, s *APIServices) {
	var opts services.DocumentsListOptions

	if err := s.Decoder.Decode(&opts, r.URL.Query()); err != nil {
		renderFailure(w, http.StatusBadRequest, M{
			"err": err,
		})
		return
	}

	docs, err := s.Documents.List(&opts)
	if err != nil {
		renderError(w, "error listing documents", M{
			"err": err,
		})
		return
	}

	if docs == nil {
		docs = []*models.Document{}
	}

	renderSuccess(w, docs)
}

func handleDocumentCreate(w http.ResponseWriter, r *http.Request, s *APIServices) {
	panic("unimplemented")
}

func handleDocumentGet(w http.ResponseWriter, r *http.Request, s *APIServices) {
	panic("unimplemented")
}

func handleDocumentDelete(w http.ResponseWriter, r *http.Request, s *APIServices) {
	panic("unimplemented")
}
