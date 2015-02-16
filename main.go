package main

import (
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/jinzhu/gorm"
	"github.com/mccoyst/validate"
	flag "github.com/ogier/pflag"
	renderpkg "github.com/unrolled/render"
	"github.com/zenazn/goji/graceful"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"

	_ "github.com/mattn/go-sqlite3"
)

var (
	flagPort uint16

	log    = logrus.New()
	render *renderpkg.Render
	v      validate.V
)

func init() {
	flag.Uint16VarP(&flagPort, "port", "p", 8080, "port to listen on")

	render = renderpkg.New(renderpkg.Options{
		IndentJSON: true,
	})

	v = make(validate.V)
	v["nonzero"] = func(i interface{}) error {
		var ok bool

		switch v := i.(type) {
		case int:
			ok = v != 0
		case uint:
			ok = v != 0
		case string:
			ok = len(v) > 0
		}

		if !ok {
			return fmt.Errorf("should be nonzero")
		}

		return nil
	}
	v["nonempty"] = func(i interface{}) error {
		var ok bool

		switch v := i.(type) {
		case []int64:
			ok = len(v) > 0
		}

		if !ok {
			return fmt.Errorf("should be nonempty")
		}

		return nil
	}
	v["notnil"] = func(i interface{}) error {
		if i == nil {
			return fmt.Errorf("should not be nil")
		}
		return nil
	}
}

func main() {
	db, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		log.WithField("err", err).Fatal("Could not open db")
	}

	// Database migration
	db.AutoMigrate(&Document{}, &Tag{}, &File{})

	m := web.New()
	m.Use(middleware.EnvInit)
	m.Use(middleware.RequestID)
	m.Use(logMiddleware)
	m.Use(recoverMiddleware)
	m.Use(middleware.AutomaticOptions)

	// Create API mux
	api := web.New()
	api.Use(jsonMiddleware)
	api.Use(corsMiddleware)
	api.Use(func(c *web.C, h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			c.Env["db"] = db
			h.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	})

	api.Get("/api/tags", routeTagsGetAll)
	api.Post("/api/tags", routeTagsCreate)
	api.Get("/api/tags/:tag_id", routeTagsGetOne)

	api.Get("/api/documents", routeDocumentsGetAll)
	api.Post("/api/documents", routeDocumentsCreate)
	api.Get("/api/documents/:document_id", routeDocumentsGetOne)

	m.Handle("/api/*", api)

	// Set up graceful listener
	graceful.HandleSignals()
	graceful.PreHook(func() { log.Info("Received signal, gracefully stopping...") })
	graceful.PostHook(func() { log.Info("Stopped") })

	// Start listening
	addr := fmt.Sprintf("localhost:%d", flagPort)
	log.WithField("addr", addr).Info("Starting server...")
	if err := graceful.ListenAndServe(addr, m); err != nil {
		log.WithField("err", err).Error("Error while listening")
		return
	}
	graceful.Wait()
}
