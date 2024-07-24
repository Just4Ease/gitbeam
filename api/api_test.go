package api

import (
	"context"
	"encoding/json"
	"fmt"
	"gitbeam/api/pb/commits"
	gitRepos "gitbeam/api/pb/repos"
	"gitbeam/mocks"
	"gitbeam/models"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestListRepositories(t *testing.T) {
	logger := logrus.New()
	controller := gomock.NewController(t)
	defer controller.Finish()

	ctx := context.Background()
	repoRPCMock := mocks.NewMockGitBeamRepositoryServiceClient(controller)
	repoRPCMock.EXPECT().ListGitRepositories(gomock.Any(), &gitRepos.Void{}).MaxTimes(1).Return(
		&gitRepos.ListGitRepositoriesResponse{
			Repos: []*gitRepos.Repo{
				{Name: "chromium"},
				{Name: "brave"},
			},
		}, nil)

	router := chi.NewMux()
	New(nil, repoRPCMock, logger).Routes(router)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/repos", nil)

	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "chromium")
	assert.Contains(t, rr.Body.String(), "brave")
}

func TestGetRepo(t *testing.T) {
	logger := logrus.New()
	controller := gomock.NewController(t)
	defer controller.Finish()

	ownerName := "chromium"
	repoName := "chromium"

	repoRPCMock := mocks.NewMockGitBeamRepositoryServiceClient(controller)
	repoRPCMock.EXPECT().GetGitRepo(gomock.Any(), &gitRepos.GetGitRepoRequest{
		OwnerName: ownerName,
		RepoName:  repoName,
	}).MaxTimes(1).Return(
		&gitRepos.Repo{
			Name:  repoName,
			Owner: ownerName,
		}, nil)

	router := chi.NewMux()
	New(nil, repoRPCMock, logger).Routes(router)

	req, err := http.NewRequest(http.MethodGet, "/repos/chromium/chromium", nil)
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	assert.Contains(t, rr.Body.String(), repoName)
	assert.Contains(t, rr.Body.String(), ownerName)
	var result models.Result
	assert.Nil(t, json.NewDecoder(rr.Body).Decode(&result))
	assert.Equal(t, true, result.Success)
}

func TestListCommits(t *testing.T) {
	logger := logrus.New()
	controller := gomock.NewController(t)
	defer controller.Finish()

	ownerName := "chromium"
	repoName := "chromium"

	//ctx := context.Background()
	mockRepoRPC := mocks.NewMockGitBeamRepositoryServiceClient(controller)
	mockCommitsRPC := mocks.NewMockGitBeamCommitsServiceClient(controller)

	file, err := os.Open("../seeds/commit.json")
	assert.Nil(t, err)

	var commit *commits.Commit
	assert.Nil(t, json.NewDecoder(file).Decode(&commit))

	mockRepoRPC.EXPECT().GetGitRepo(gomock.Any(), &gitRepos.GetGitRepoRequest{
		OwnerName: ownerName,
		RepoName:  repoName,
	}).MaxTimes(1).Return(
		&gitRepos.Repo{
			Name:  repoName,
			Owner: ownerName,
		}, nil)

	mockCommitsRPC.EXPECT().ListCommits(gomock.Any(), &commits.CommitFilterParams{
		Page:      1,
		Limit:     1,
		OwnerName: ownerName,
		RepoName:  repoName,
	}).MaxTimes(1).Return(
		&commits.ListCommitResponse{
			Data: []*commits.Commit{
				commit,
			},
		},
		nil,
	)
	router := chi.NewMux()
	New(mockCommitsRPC, mockRepoRPC, logger).Routes(router)

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/commits?ownerName=%s&repoName=%s&page=1&limit=1&", ownerName, repoName), nil)

	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), commit.Author)
	assert.Contains(t, rr.Body.String(), commit.Sha)
}
