package api

import (
	"errors"
	"gitbeam/core"
	"gitbeam/models"
	"gitbeam/utils"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/schema"
	"net/http"
)

func (a API) newCommitsRoute() chi.Router {
	router := chi.NewRouter()

	router.Get("/", a.listCommits)
	router.Get("/top-authors", a.listTopCommitAuthors)
	router.Get("/{ownerName}/{repoName}/{sha}", a.getCommitBySha)
	router.Post("/start-mirroring", a.startBeamingRepoCommits)
	router.Post("/stop-mirroring", a.stopBeamingRepoCommits)

	return router
}

func (a API) listCommits(w http.ResponseWriter, r *http.Request) {
	useLogger := a.logger.WithContext(r.Context()).WithField("endpointName", "listCommits").Logger
	decoder := schema.NewDecoder()
	var filter models.CommitFilters
	if err := decoder.Decode(&filter, r.URL.Query()); err != nil {
		useLogger.WithError(err).Error("failed to decode query params into filters.")
		utils.WriteHTTPError(w, http.StatusBadRequest, errors.New("Bad/Invalid Query Parameters"))
		return
	}

	useLogger.WithField("filter", filter).Info("filters")
	list, err := a.service.ListCommits(r.Context(), filter)
	if err != nil {
		utils.WriteHTTPError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteHTTPSuccess(w, "Success", list)
}

func (a API) listTopCommitAuthors(w http.ResponseWriter, r *http.Request) {
	useLogger := a.logger.WithContext(r.Context()).WithField("endpointName", "listTopCommitAuthors").Logger
	decoder := schema.NewDecoder()
	var filter models.CommitFilters
	if err := decoder.Decode(&filter, r.URL.Query()); err != nil {
		useLogger.WithError(err).Error("failed to decode query params into filters.")
		utils.WriteHTTPError(w, http.StatusBadRequest, errors.New("Bad/Invalid Query Parameters"))
		return
	}

	useLogger.WithField("filter", filter).Info("filters")

	list, err := a.service.GetTopCommitAuthors(r.Context(), filter)
	if err != nil {
		utils.WriteHTTPError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteHTTPSuccess(w, "Success", list)
}

func (a API) getCommitBySha(w http.ResponseWriter, r *http.Request) {
	owner := models.OwnerAndRepoName{
		OwnerName: chi.URLParam(r, "ownerName"),
		RepoName:  chi.URLParam(r, "repoName"),
	}

	sha := chi.URLParam(r, "sha")

	commit, err := a.service.GetCommitsBySha(r.Context(), owner, sha)
	if err != nil {
		statusNotFound := http.StatusBadRequest

		switch {
		case errors.Is(err, core.ErrGithubRepoNotFound), errors.Is(err, core.ErrCommitNotFound):
			statusNotFound = http.StatusNotFound
		}

		utils.WriteHTTPError(w, statusNotFound, err)
		return
	}

	utils.WriteHTTPSuccess(w, "Successfully retrieved commit", commit)
}
