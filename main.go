package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/schema"
	"github.com/jmoiron/sqlx"
	"github.com/lann/squirrel"
	"github.com/lidashuang/goji-gzip"
	flag "github.com/ogier/pflag"
	"github.com/zenazn/goji/graceful"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"

	"github.com/andrew-d/docstore/controllers"
	"github.com/andrew-d/docstore/models"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

var (
	flagHost          string
	flagPort          uint16
	flagDebug         bool
	flagDataDirectory string
	flagDbType        string
	flagDbConn        string

	log = logrus.New()
)

func init() {
	flag.StringVar(&flagHost, "host", "localhost", "host to listen on")
	flag.Uint16VarP(&flagPort, "port", "p", 8080, "port to listen on")
	flag.BoolVarP(&flagDebug, "debug", "d", false, "run in debug mode")
	flag.StringVar(&flagDataDirectory, "data", filepath.Join(".", "data"),
		"path to store data in")
	flag.StringVar(&flagDbType, "dbtype", "sqlite3", "type of database to use")
	flag.StringVar(&flagDbConn, "dbconn", filepath.Join(".", "data", "docstore.db"),
		"database connection string")
}

func main() {
	flag.Parse()

	// Create data directory
	filesDir, err := filepath.Abs(filepath.Join(flagDataDirectory, "files"))
	if err != nil {
		log.WithField("err", err).Error("Error while getting absolute files directory")
	}

	err = os.MkdirAll(filesDir, 0700)
	if err != nil {
		log.WithField("err", err).Fatal("Could not create files directory")
	}

	// Connect to the database
	db, err := sqlx.Connect(flagDbType, flagDbConn)
	if err != nil {
		log.WithFields(logrus.Fields{
			"dbtype": flagDbType,
			"dbconn": flagDbConn,
			"err":    err,
		}).Fatal("Could not open db")
	}
	defer db.Close()

	var sq squirrel.StatementBuilderType
	if flagDbType == "postgres" {
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
		Decoder: schema.NewDecoder(),
		Debug:   flagDebug,
		Logger:  log.WithField("package", "controllers"),
	}

	m := web.New()
	m.Use(middleware.EnvInit)
	m.Use(middleware.RequestID)
	m.Use(logMiddleware)
	m.Use(recoverMiddleware)
	m.Use(middleware.AutomaticOptions)
	m.Use(gzip.GzipHandler)

	// Create API mux
	api := web.New()
	api.Use(jsonMiddleware)
	api.Use(corsMiddleware)

	// Set up controllers/routes
	tagController := controllers.TagController{AppController: appController}
	api.Get("/api/tags", tagController.Action(tagController.GetAll))
	api.Get("/api/tags/:tag_id", tagController.Action(tagController.GetOne))
	api.Delete("/api/tags/:tag_id", tagController.Action(tagController.Delete))
	api.Post("/api/tags", tagController.Action(tagController.Create))

	documentController := controllers.DocumentController{AppController: appController}
	api.Get("/api/documents", documentController.Action(documentController.GetAll))
	api.Get("/api/documents/:document_id", documentController.Action(documentController.GetOne))
	api.Delete("/api/documents/:document_id", documentController.Action(documentController.Delete))
	api.Post("/api/documents", documentController.Action(documentController.Create))

	fileController := controllers.FileController{
		AppController: appController,
		FilePath:      filesDir,
	}
	api.Get("/api/documents/:document_id/files", fileController.Action(fileController.GetAll))
	api.Post("/api/documents/:document_id/files/upload", fileController.Action(fileController.Upload))
	api.Post("/api/documents/:document_id/files/scan", fileController.Action(fileController.Scan))
	api.Get("/api/documents/:document_id/files/:file_id", fileController.Action(fileController.GetOne))
	api.Get("/api/documents/:document_id/files/:file_id/content", fileController.Action(fileController.Content))

	m.Handle("/api/*", api)

	// Set up graceful listener
	graceful.HandleSignals()
	graceful.PreHook(func() { log.Info("Received signal, gracefully stopping...") })
	graceful.PostHook(func() { log.Info("Stopped") })

	// Start listening
	addr := fmt.Sprintf("%s:%d", flagHost, flagPort)
	log.WithField("addr", addr).Info("Starting server...")
	if err := graceful.ListenAndServe(addr, m); err != nil {
		log.WithField("err", err).Error("Error while listening")
		return
	}
	graceful.Wait()
}
