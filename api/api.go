package api

import (
	"gitbeam/core"
	"gitbeam/cron"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type API struct {
	coreService *core.GitBeamService
	cronService *cron.Service
	logger      *logrus.Logger
}

func New(coreService *core.GitBeamService, cronService *cron.Service, logger *logrus.Logger) *API {
	return &API{
		coreService: coreService,
		cronService: cronService,
		logger:      logger.WithField("serviceName", "apiRouter").Logger,
	}
}

func (a API) Routes(router *chi.Mux) {
	// Mount all route paths here.
	router.Mount("/repos", a.newReposRoute())
	router.Mount("/commits", a.newCommitsRoute())
}
