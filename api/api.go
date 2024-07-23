package api

import (
	"gitbeam/api/pb/commits"
	gitRepos "gitbeam/api/pb/repos"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type API struct {
	commitsRPC commits.GitBeamCommitsServiceClient
	reposRPC   gitRepos.GitBeamRepositoryServiceClient
	logger     *logrus.Logger
}

func New(commitsRPC commits.GitBeamCommitsServiceClient, reposRPC gitRepos.GitBeamRepositoryServiceClient, logger *logrus.Logger) *API {
	return &API{
		commitsRPC: commitsRPC,
		reposRPC:   reposRPC,
		logger:     logger.WithField("serviceName", "apiRouter").Logger,
	}
}

func (a API) Routes(router *chi.Mux) {
	// Mount all route paths here.
	router.Mount("/repos", a.newReposRoute())
	router.Mount("/commits", a.newCommitsRoute())
}
