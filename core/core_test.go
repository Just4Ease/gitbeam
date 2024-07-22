package core

import (
	"context"
	"encoding/json"
	"gitbeam/mocks"
	"gitbeam/models"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGetRepoByOwnerAndRepoName(t *testing.T) {
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
	}{
		{
			TestName:      "Valid Owner and Repo Name",
			OwnerName:     "brave",
			RepoName:      "brave-browser",
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

			owner := &models.OwnerAndRepoName{
				OwnerName: test.OwnerName,
				RepoName:  test.RepoName,
			}

			if test.ExpectedError == nil {
				f, err := os.Open("../seeds/brave.brave-browser.json")
				assert.NoError(t, err)
				var repo models.Repo
				assert.Nil(t, json.NewDecoder(f).Decode(&repo))
				dataStore.EXPECT().GetRepoByOwner(ctx, owner).MaxTimes(1).Return(&repo, nil)
			} else {
				dataStore.EXPECT().GetRepoByOwner(ctx, owner).MaxTimes(1).Return(nil, nil)
			}

			repo, err := service.GetByOwnerAndRepoName(ctx, owner)
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
