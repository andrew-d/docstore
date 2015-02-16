package main

import (
	"fmt"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/jmoiron/sqlx"
	"github.com/lann/squirrel"
	flag "github.com/ogier/pflag"
	"github.com/zenazn/goji/graceful"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"

	"github.com/andrew-d/docstore/controllers"
	"github.com/andrew-d/docstore/models"

	_ "github.com/mattn/go-sqlite3"
)

var (
	flagPort  uint16
	flagDebug bool

	log = logrus.New()
)

func init() {
	flag.Uint16VarP(&flagPort, "port", "p", 8080, "port to listen on")
	flag.BoolVarP(&flagDebug, "debug", "d", false, "run in debug mode")
}

func main() {
	db, err := sqlx.Connect("sqlite3", ":memory:")
	if err != nil {
		log.WithField("err", err).Fatal("Could not open db")
	}

	// TODO: properly configured
	var sq squirrel.StatementBuilderType
	if false {
		sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	} else {
		sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question)
	}

	// Set up schema
	tx := db.MustBegin()
	for i, stmt := range models.Schema() {
		log.WithField("index", i).Debug("Executing schema statement")
		tx.MustExec(strings.TrimSpace(stmt))
	}
	err = tx.Commit()
	if err != nil {
		log.WithField("err", err).Fatal("Error committing schema transaction")
	}

	// Create global controller
	appController := controllers.AppController{
		DB:      db,
		Builder: sq,
		Debug:   flagDebug,
	}

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

	// Set up controllers/routes
	tagController := controllers.TagController{AppController: appController}
	api.Get("/api/tags", tagController.Action(tagController.GetAll))
	api.Get("/api/tags/:tag_id", tagController.Action(tagController.GetOne))
	api.Post("/api/tags", tagController.Action(tagController.Create))

	documentController := controllers.DocumentController{AppController: appController}
	api.Get("/api/documents", documentController.Action(documentController.GetAll))
	api.Get("/api/documents/:document_id", documentController.Action(documentController.GetOne))
	api.Post("/api/documents", documentController.Action(documentController.Create))

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
