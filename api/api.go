package api

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/justinas/alice"

	"github.com/andrew-d/docstore/router"
	"github.com/andrew-d/docstore/services"
)

var (
	log = logrus.New().WithField("package", "github.com/andrew-d/docstore/api")
)

type APIServices struct {
	Documents services.DocumentsService

	Decoder *schema.Decoder
}

func Make(s *APIServices) *mux.Router {
	// Default services.
	if s.Decoder != nil {
		s.Decoder = schema.NewDecoder()
	}

	// TODO: more middleware?
	mw := alice.New(jsonMiddleware)

	wrap := func(fn func(http.ResponseWriter, *http.Request, *APIServices)) http.Handler {
		// Make a handler that injects the services
		handler := func(w http.ResponseWriter, r *http.Request) {
			fn(w, r, s)
		}

		// Return the function wrapped in our middleware.
		return mw.ThenFunc(handler)
	}

	r := router.NewAPIRouter(nil)

	// Set the handlers for our named routes.
	r.Get(router.ListDocuments).Handler(wrap(handleDocumentList))
	r.Get(router.CreateDocument).Handler(wrap(handleDocumentCreate))
	r.Get(router.GetDocument).Handler(wrap(handleDocumentGet))
	r.Get(router.DeleteDocument).Handler(wrap(handleDocumentDelete))

	return r
}

func jsonMiddleware(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
