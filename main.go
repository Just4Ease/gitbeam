package main

import (
	"gitbeam/api"
	"gitbeam/core"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"os"
)

func main() {

	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)

	app := core.NewGitBeamService(logger)

	router := chi.NewRouter()

	api.New(app, logger).Routes(router)

}
