package api

import (
	"encoding/json"
	"errors"
	"gitbeam/core"
	"gitbeam/models"
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

func (a API) startBeamingRepoCommits(w http.ResponseWriter, r *http.Request) {
	var payload models.OwnerAndRepoName
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		a.logger.WithError(err).Error("error decoding payload")
		utils.WriteHTTPError(w, http.StatusBadRequest, err)
		return
	}

	repo, err := a.service.StartBeamingCommits(r.Context(), payload)
	if err != nil {
		a.logger.WithError(err).Error("error attempting to beam repository commits.")
		statusNotFound := http.StatusBadRequest
		if errors.Is(err, core.ErrRepositoryNotFound) {
			statusNotFound = http.StatusNotFound
		}

		utils.WriteHTTPError(w, statusNotFound, err)
		return
	}

	utils.WriteHTTPSuccess(w, "Successfully started beaming repo", repo)
}

func (a API) stopBeamingRepoCommits(w http.ResponseWriter, r *http.Request) {
	return
}

func (a API) newReposRoute() chi.Router {
	router := chi.NewRouter()

	router.Get("/", a.getRepoByOwnerAndRepoName)
	router.Get("/{id}", a.getRepoById)
	router.Post("/start-beam", a.startBeamingRepoCommits)
	router.Post("/stop-beam", a.stopBeamingRepoCommits)

	return router
}
