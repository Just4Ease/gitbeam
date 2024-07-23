package api

import (
	gitRepos "gitbeam/api/pb/repos"
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
	repo, err := a.reposRPC.ListGitRepositories(r.Context(), &gitRepos.Void{})
	if err != nil {
		a.logger.WithError(err).Error("failed to fetch list of repositories")
		statusCode := http.StatusBadRequest
		//if errors.Is(err, core.ErrGithubRepoNotFound) {
		//	statusCode = http.StatusNotFound
		//}

		utils.WriteHTTPError(w, statusCode, err)
		return
	}

	utils.WriteHTTPSuccess(w, "Successfully retrieved list of repositories", repo.Repos)
}

func (a API) getRepoByOwnerAndRepoName(w http.ResponseWriter, r *http.Request) {
	repo, err := a.reposRPC.GetGitRepo(r.Context(), &gitRepos.GetGitRepoRequest{
		OwnerName: chi.URLParam(r, "ownerName"),
		RepoName:  chi.URLParam(r, "repoName"),
	})
	if err != nil {
		a.logger.WithError(err).Error("error getting repo by owner/repo")
		statusNotFound := http.StatusNotFound

		utils.WriteHTTPError(w, statusNotFound, err)
		return
	}

	utils.WriteHTTPSuccess(w, "Successfully retrieved repo", repo)

}
