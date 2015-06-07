package main

import (
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/tylerb/graceful"

	"github.com/andrew-d/docstore/api"
	"github.com/andrew-d/docstore/datastore"
)

const Timeout = 10 * time.Second

var (
	// Closed when we're shutting down.
	shutdown = make(chan struct{})

	// Waitgroup for all goroutines.
	wg sync.WaitGroup

	// Logger
	log = logrus.New()
)

func init() {
	log.Out = os.Stderr
	log.Formatter = &logrus.TextFormatter{}
}

func main() {
	// Set up signal handling
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	if err := startAPI(); err != nil {
		log.Fatal(err)
	}
	if err := startFrontend(); err != nil {
		log.Fatal(err)
	}

	// Wait for a signal, then stop all servers.
	<-interrupt
	close(shutdown)
	log.Info("shutting down...")

	// Wait for the servers to finish.
	wg.Wait()
	log.Info("finished")
}

func startAPI() error {
	// Set up dependencies.
	// TODO
	var _ = datastore.Create

	// Create the API router
	router := api.Make(&api.APIServices{
		Documents: nil, // TODO
	})

	// Start serving it
	addr := "localhost:3002"
	serveGracefully(addr, router)
	log.WithField("addr", addr).Info("started API server")
	return nil
}

func startFrontend() error {
	// Set up dependencies.
	// TODO
	var _ = datastore.Create

	// Create the frontend router
	// TODO

	// Start serving it
	addr := "localhost:3001"
	serveGracefully(addr, nil)
	log.WithField("addr", addr).Info("started frontend server")
	return nil
}

func serveGracefully(addr string, handler http.Handler) {
	// Create a new graceful HTTP server that does not handle signals.
	srv := &graceful.Server{
		Timeout:          Timeout,
		NoSignalHandling: true,
		Server: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
	}

	// Shutdown the server when our close channel is signalled.
	wg.Add(1)
	go func() {
		<-shutdown
		log.WithField("addr", addr).Info("telling server to stop...")
		srv.Stop(Timeout)
		wg.Done()
	}()

	// Run the server in (yet another) goroutine.
	wg.Add(1)
	go func() {
		// TODO: handle this error somehow?
		if err := srv.ListenAndServe(); err != nil {
			log.WithField("err", err).Error("error starting graceful listener")
		}
		log.WithField("addr", addr).Info("server finished")
		wg.Done()
	}()
}
