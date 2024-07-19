package api

import (
	"errors"
	"gitbeam/core"
	"gitbeam/utils"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (a API) listRepositories(w http.ResponseWriter, r *http.Request) {

}

func (a API) getRepoById(w http.ResponseWriter, r *http.Request) {

}

func (a API) getRepoByOwnerAndRepoName(w http.ResponseWriter, r *http.Request) {
	ownerName := r.URL.Query().Get("ownerName")
	repoName := r.URL.Query().Get("repoName")

	repo, err := a.service.GetByOwnerAndRepoName(r.Context(), ownerName, repoName)
	if err != nil {
		a.logger.WithError(err).Error("error getting repo by owner/repo")
		statusNotFound := http.StatusBadRequest
		if errors.Is(err, core.ErrRepositoryNotFound) {
			statusNotFound = http.StatusNotFound
		}

		utils.WriteHTTPError(w, statusNotFound, err)
		return
	}

	utils.WriteHTTPSuccess(w, "Successfully retrieved repo", repo)

}

func (a API) beamRepoCommits(w http.ResponseWriter, r *http.Request) {

}

func (a API) newReposRoute() chi.Router {
	router := chi.NewRouter()

	router.Get("/", a.getRepoByOwnerAndRepoName)
	router.Get("/{id}", a.getRepoById)
	router.Post("/beam", a.beamRepoCommits)

	return router
}
