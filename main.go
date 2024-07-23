package main

import (
	"context"
	"errors"
	"fmt"
	"gitbeam/api"
	"gitbeam/api/pb/commits"
	gitRepos "gitbeam/api/pb/repos"
	"gitbeam/config"
	"github.com/go-chi/chi/v5"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)
	router := chi.NewRouter()
	secrets := config.GetSecrets()

	var repoServiceRPC gitRepos.GitBeamRepositoryServiceClient
	var commitsServiceRPC commits.GitBeamCommitsServiceClient
	var err error

	if repoServiceRPC, err = connectRPC[gitRepos.GitBeamRepositoryServiceClient](
		secrets.RepoManagerURL,
		func(connection grpc.ClientConnInterface) any {
			return gitRepos.NewGitBeamRepositoryServiceClient(connection)
		}); err != nil {
		logger.WithError(err).Fatal("failed to connect to repos RPC server")
	}

	if commitsServiceRPC, err = connectRPC[commits.GitBeamCommitsServiceClient](
		secrets.CommitsMonitorURL,
		func(connection grpc.ClientConnInterface) any {
			return commits.NewGitBeamCommitsServiceClient(connection)
		}); err != nil {
		logger.WithError(err).Fatal("failed to connect to commits RPC server")
	}

	api.New(commitsServiceRPC, repoServiceRPC, logger).Routes(router)
	startAndManageHTTPServer(router, secrets.Port, logger)
}

func startAndManageHTTPServer(router *chi.Mux, port string, logger *logrus.Logger) {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
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

func connectRPC[T any](address string, fn connectRPCFunc) (T, error) {
	a := new(T) // As nil.
	connection, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return *a, err
	}
	out, ok := fn(connection).(T)
	if !ok {
		return *a, errors.New("invalid RPC type expectation")
	}
	return out, nil
}

type connectRPCFunc func(connection grpc.ClientConnInterface) any
