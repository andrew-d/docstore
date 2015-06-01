package router

import (
	"github.com/gorilla/mux"
)

func NewAPIRouter(base *mux.Router) *mux.Router {
	if base == nil {
		base = mux.NewRouter()
	}

	base.StrictSlash(true)

	base.Path("/documents").Methods("GET").Name(ListDocuments)
	base.Path("/documents").Methods("POST").Name(CreateDocument)
	base.Path("/documents/{id:[0-9]+}").Methods("GET").Name(GetDocument)
	base.Path("/documents/{id:[0-9]+}").Methods("DELETE").Name(DeleteDocument)

	return base
}
