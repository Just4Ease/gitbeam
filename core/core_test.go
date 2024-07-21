package core

import (
	"context"
	"gitbeam/mocks"
	"gitbeam/models"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetRepoByOwnerAndRepoName(t *testing.T) {
	logger := logrus.New()

	service := NewGitBeamService(logger, nil, nil, nil)

	ctx := context.Background()
	tests := []struct {
		TestName      string
		OwnerName     string
		RepoName      string
		ExpectedError error
	}{
		{
			TestName:      "Valid Owner and Repo Name",
			OwnerName:     "Just4Ease",
			RepoName:      "gitbeam",
			ExpectedError: nil,
		},
		{
			TestName:      "Empty Owner and Valid Repo Name",
			OwnerName:     "",
			RepoName:      "gitbeam",
			ExpectedError: ErrOwnerAndRepoNameRequired,
		},
		{
			TestName:      "Valid Owner and Empty Repo Name",
			OwnerName:     "Just4Ease",
			RepoName:      "",
			ExpectedError: ErrOwnerAndRepoNameRequired,
		},
		{
			TestName:      "Valid Owner and Invalid Repo Name",
			OwnerName:     "Just4Ease",
			RepoName:      "gitbeamxxxx",
			ExpectedError: ErrGithubRepoNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.TestName, func(t *testing.T) {
			repo, err := service.GetByOwnerAndRepoName(ctx, test.OwnerName, test.RepoName)
			if test.ExpectedError != nil {
				assert.Equal(t, test.ExpectedError, err)
				assert.Nil(t, repo)
			} else {
				assert.Nil(t, err)
				assert.NotNil(t, repo)
				assert.Equal(t, test.RepoName, repo.Name)
				assert.Equal(t, test.OwnerName, repo.Owner)
			}
		})
	}
}

func TestStartBeamingRepoRepositoryCommits(t *testing.T) {
	logger := logrus.New()

	controller := gomock.NewController(t)
	dataStore := mocks.NewMockDataStore(controller)
	eventStore := mocks.NewMockEventStore(controller)

	service := NewGitBeamService(logger, eventStore, dataStore, nil)

	ctx := context.Background()

	tests := []struct {
		TestName      string
		OwnerName     string
		RepoName      string
		ExpectedError error
		ExpectedRepo  *models.Repo
	}{
		{
			TestName:      "Beam Valid Owner and Repo Name",
			OwnerName:     "Just4Ease",
			RepoName:      "MongoLeon",
			ExpectedError: nil,
			ExpectedRepo:  &models.Repo{},
		},
		{
			TestName:      "Beam Already Beamed Repo",
			OwnerName:     "",
			RepoName:      "MongoLeon",
			ExpectedError: ErrOwnerAndRepoNameRequired,
		},
		{
			TestName:      "Valid Owner and Empty Repo Name",
			OwnerName:     "Just4Ease",
			RepoName:      "",
			ExpectedError: ErrOwnerAndRepoNameRequired,
		},
		{
			TestName:      "Valid Owner and Invalid Repo Name",
			OwnerName:     "Just4Ease",
			RepoName:      "mongoleonnnn",
			ExpectedError: ErrGithubRepoNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.TestName, func(t *testing.T) {
			repo, err := service.StartBeamingCommits(ctx, models.OwnerAndRepoName{
				OwnerName: "",
				RepoName:  "",
			})

			if test.ExpectedError != nil {
				assert.Equal(t, test.ExpectedError, err)
				assert.Nil(t, repo)
			} else {
				assert.Nil(t, err)
				assert.NotNil(t, repo)
				assert.Equal(t, test.RepoName, repo.Name)
				assert.Equal(t, test.OwnerName, repo.Owner)
			}
		})
	}
}
