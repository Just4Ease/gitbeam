package api

import (
	"errors"
	"gitbeam/core"
	"gitbeam/models"
	"gitbeam/utils"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (a API) newReposRoute() chi.Router {
	router := chi.NewRouter()

	router.Get("/{ownerName}/{repoName}", a.getRepoByOwnerAndRepoName)
	router.Get("/", a.listRepositories)

	return router
}

func (a API) listRepositories(w http.ResponseWriter, r *http.Request) {
	repo, err := a.coreService.ListRepos(r.Context())
	if err != nil {
		a.logger.WithError(err).Error("failed to fetch list of repositories")
		statusNotFound := http.StatusBadRequest
		if errors.Is(err, core.ErrGithubRepoNotFound) {
			statusNotFound = http.StatusNotFound
		}

		utils.WriteHTTPError(w, statusNotFound, err)
		return
	}

	utils.WriteHTTPSuccess(w, "Successfully retrieved list of repositories", repo)
}

func (a API) getRepoByOwnerAndRepoName(w http.ResponseWriter, r *http.Request) {
	owner := models.OwnerAndRepoName{
		OwnerName: chi.URLParam(r, "ownerName"),
		RepoName:  chi.URLParam(r, "repoName"),
	}

	repo, err := a.coreService.GetByOwnerAndRepoName(r.Context(), &owner)
	if err != nil {
		a.logger.WithError(err).Error("error getting repo by owner/repo")
		statusNotFound := http.StatusBadRequest
		if errors.Is(err, core.ErrGithubRepoNotFound) {
			statusNotFound = http.StatusNotFound
		}

		utils.WriteHTTPError(w, statusNotFound, err)
		return
	}

	utils.WriteHTTPSuccess(w, "Successfully retrieved repo", repo)

}
