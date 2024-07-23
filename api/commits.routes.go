package api

import (
	"encoding/json"
	"errors"
	"gitbeam/api/pb/commits"
	gitRepos "gitbeam/api/pb/repos"
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
	router.Post("/start-monitoring", a.startMonitoringRepoCommits)
	router.Post("/stop-monitoring", a.stopMonitoringRepoCommits)

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
	rpcFilter := &commits.CommitFilterParams{
		Page:      filter.Page,
		Limit:     filter.Limit,
		OwnerName: filter.OwnerName,
		RepoName:  filter.RepoName,
		FromDate:  "",
		ToDate:    "",
	}

	if filter.FromDate != nil {
		rpcFilter.FromDate = filter.FromDate.String()
	}

	if filter.ToDate != nil {
		rpcFilter.ToDate = filter.ToDate.String()
	}

	list, err := a.commitsRPC.ListCommits(r.Context(), rpcFilter)
	if err != nil {
		utils.WriteHTTPError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteHTTPSuccess(w, "Success", list.Data)
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
	rpcFilter := &commits.CommitFilterParams{
		Page:      filter.Page,
		Limit:     filter.Limit,
		OwnerName: filter.OwnerName,
		RepoName:  filter.RepoName,
		FromDate:  "",
		ToDate:    "",
	}

	if filter.FromDate != nil {
		rpcFilter.FromDate = filter.FromDate.String()
	}

	if filter.ToDate != nil {
		rpcFilter.ToDate = filter.ToDate.String()
	}

	list, err := a.commitsRPC.ListTopCommitAuthor(r.Context(), rpcFilter)
	if err != nil {
		utils.WriteHTTPError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteHTTPSuccess(w, "Success", list.Data)
}

func (a API) getCommitBySha(w http.ResponseWriter, r *http.Request) {
	owner := commits.CommitByOwnerAndShaParams{
		OwnerName: chi.URLParam(r, "ownerName"),
		RepoName:  chi.URLParam(r, "repoName"),
		Sha:       chi.URLParam(r, "sha"),
	}

	commit, err := a.commitsRPC.GetCommitByOwnerAndSHA(r.Context(), &owner)
	if err != nil {
		statusNotFound := http.StatusBadRequest

		//switch {
		//case errors.Is(err, core.ErrGithubRepoNotFound), errors.Is(err, core.ErrCommitNotFound):
		//	statusNotFound = http.StatusNotFound
		//}

		utils.WriteHTTPError(w, statusNotFound, err)
		return
	}

	utils.WriteHTTPSuccess(w, "Successfully retrieved commit", commit)
}

func (a API) startMonitoringRepoCommits(w http.ResponseWriter, r *http.Request) {
	useLogger := a.logger.WithContext(r.Context()).WithField("endpointName", "startMonitoringRepoCommits").Logger

	var payload commits.MonitorRepositoryCommitsConfigParams
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		useLogger.WithError(err).Error("failed to decode query params into filters.")
		utils.WriteHTTPError(w, http.StatusBadRequest, errors.New("Bad/Invalid Query Parameters"))
		return
	}

	_, err := a.reposRPC.GetGitRepo(r.Context(), &gitRepos.GetGitRepoRequest{
		OwnerName: payload.OwnerName,
		RepoName:  payload.RepoName,
	})
	if err != nil {
		useLogger.WithError(err).Error("failed to get repo from repo rpc service.")
		utils.WriteHTTPError(w, http.StatusBadRequest, err)
		return
	}

	_, err = a.commitsRPC.StartMonitoringRepositoryCommits(r.Context(), &payload)
	if err != nil {
		useLogger.WithError(err).Error("failed to start monitoring repository commits")
		statusCode := http.StatusBadRequest
		utils.WriteHTTPError(w, statusCode, err)
		return
	}

	utils.WriteHTTPSuccess(w, "Successfully started monitoring repo commits.", nil)
}

func (a API) stopMonitoringRepoCommits(w http.ResponseWriter, r *http.Request) {
	useLogger := a.logger.WithContext(r.Context()).WithField("endpointName", "stopMonitoringRepoCommits").Logger
	var payload models.OwnerAndRepoName
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		useLogger.WithError(err).Error("error decoding payload")
		utils.WriteHTTPError(w, http.StatusBadRequest, err)
		return
	}

	_, err := a.reposRPC.GetGitRepo(r.Context(), &gitRepos.GetGitRepoRequest{
		OwnerName: payload.OwnerName,
		RepoName:  payload.RepoName,
	})
	if err != nil {
		useLogger.WithError(err).WithField("payload", payload).Error("failed to get repo from repo rpc service.")
		utils.WriteHTTPError(w, http.StatusBadRequest, err)
		return
	}

	_, err = a.commitsRPC.StopMonitoringRepositoryCommits(r.Context(), &commits.StopMonitoringRepositoryCommitParams{
		OwnerName: payload.OwnerName,
		RepoName:  payload.RepoName,
	})
	if err != nil {
		useLogger.WithError(err).Error("failed to stop monitoring repository commits")
		statusCode := http.StatusBadRequest
		utils.WriteHTTPError(w, statusCode, err)
		return
	}
	//
	utils.WriteHTTPSuccess(w, "Successfully stopped monitoring repo commits.", nil)
}
