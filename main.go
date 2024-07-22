package main

import (
	"context"
	"fmt"
	"gitbeam/api"
	"gitbeam/core"
	"gitbeam/cron"
	"gitbeam/events"
	"gitbeam/repository"
	"gitbeam/repository/sqlite"
	"gitbeam/store"
	"github.com/go-chi/chi/v5"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var eventStore store.EventStore
	var dataStore repository.DataStore
	var cronStore repository.CronServiceStore
	var err error

	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)
	router := chi.NewRouter()

	//Using SQLite as the mini persistent storage.
	//( in a real world system, this would be any production level or vendor managed db )
	if dataStore, err = sqlite.NewSqliteRepo("data.db"); err != nil {
		logger.WithError(err).Fatal("failed to initialize sqlite database repository for cron store.")
	}

	// A channel based pub/sub messaging system.
	//( in a real world system, this would be apache-pulsar, kafka, nats.io or rabbitmq )
	eventStore = store.NewEventStore(logger)

	// If the dependencies were more than 3, I would use a variadic function to inject them.
	//Clarity is better here for this exercise.
	coreService := core.NewGitBeamService(logger, eventStore, dataStore, nil)

	// To handle event-based background activities. ( in a real world system, this would be apache-pulsar, kafka, nats.io or rabbitmq )
	go events.NewEventHandler(eventStore, logger, coreService).Listen()

	//Using SQLite as the mini persistent storage.
	//( in a real world system, this would be any production level or vendor managed db )
	if cronStore, err = sqlite.NewSqliteCronStore("cron_store.db"); err != nil {
		logger.WithError(err).Fatal("failed to initialize sqlite database repository for cron store.")
	}

	cronService := cron.NewCronService(cronStore, coreService, logger)
	go cronService.Start()

	api.New(coreService, cronService, logger).Routes(router)

	startAndManageHTTPServer(router, logger)
}

func startAndManageHTTPServer(router *chi.Mux, logger *logrus.Logger) {
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Channel to listen for signals
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	// Channel to notify the server has been stopped
	shutdownChan := make(chan bool)

	// Start server in a goroutine
	go func() {
		logger.Info("Started Server")
		if err := server.ListenAndServe(); err != nil {
			logger.WithError(err).Error("failed to start server")
			fmt.Printf("ListenAndServe(): %s\n", err)
		}
	}()

	// Listen for shutdown signal
	go func() {
		<-signalChan
		logger.Info("Shutting down server...")

		// Create a deadline to wait for.
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Attempt to gracefully shut down the server
		if err := server.Shutdown(ctx); err != nil {
			logger.WithError(err).Error("Server forced to shutdown")
		}

		close(shutdownChan)
	}()

	// Wait for shutdown signal
	<-shutdownChan
	logger.Info("Server gracefully stopped...")
}
