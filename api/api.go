package api

import (
	"gitbeam/service"
	"github.com/go-chi/chi/v5"
)

type API struct {
	service service.GitBeamService
	logger  *logrus.Logger
}

func New(service service.GitBeamService) *API {
	return &API{service: service}
}

func (a API) Routes(mux chi.Mux) {
	// Mount all route paths here.
}
