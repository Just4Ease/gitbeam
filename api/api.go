package api

import (
	"gitbeam/core"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type API struct {
	service *core.GitBeamService
	logger  *logrus.Logger
}

func New(service *core.GitBeamService, logger *logrus.Logger) *API {
	return &API{service: service, logger: logger.WithField("serviceName", "api").Logger}
}

func (a API) Routes(router *chi.Mux) {
	// Mount all route paths here.
	router.Mount("/repos", a.newReposRoute())
	router.Mount("/commits", a.newCommitsRoute())
}
