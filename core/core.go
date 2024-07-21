package core

import (
	"context"
	"encoding/json"
	"errors"
	"gitbeam/events/topics"
	"gitbeam/models"
	"gitbeam/repository"
	"gitbeam/store"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/go-github/v63/github"
	"github.com/sirupsen/logrus"
	"net/http"
)

var (
	ErrRepositoryNotFound       = errors.New("dataStore not found")
	ErrOwnerAndRepoNameRequired = errors.New("owner and repo name required")
)

type GitBeamService struct {
	githubClient *github.Client
	logger       *logrus.Logger
	dataStore    repository.DataStore
	eventStore   store.EventStore
}

func NewGitBeamService(
	logger *logrus.Logger,
	eventStore store.EventStore,
	dataStore repository.DataStore,
	httpClient *http.Client, // Nullable.
) *GitBeamService {
	client := github.NewClient(httpClient) // Didn't need to pass this as a top level dependency into the git beam service.
	return &GitBeamService{
		githubClient: client,
		dataStore:    dataStore,
		eventStore:   eventStore,
		logger:       logger.WithField("serviceName", "GitBeamService").Logger,
	}
}

func (g GitBeamService) GetByOwnerAndRepoName(ctx context.Context, ownerName, repoName string) (*models.Repo, error) {
	useLogger := g.logger.WithContext(ctx).WithField("methodName", "GetByOwnerAndRepoName")

	if err := validation.Validate(ownerName, validation.Required); err != nil {
		useLogger.WithError(err).Error("failed to validate owner name")
		return nil, ErrOwnerAndRepoNameRequired
	}

	if err := validation.Validate(repoName, validation.Required); err != nil {
		useLogger.WithError(err).Error("failed to validate repo name")
		return nil, ErrOwnerAndRepoNameRequired
	}

	gitRepo, _, err := g.githubClient.Repositories.Get(ctx, ownerName, repoName)
	if err != nil {
		useLogger.WithError(err).Errorln("GetByOwnerAndRepoName")
		return nil, ErrRepositoryNotFound
	}

	//repo.
	repo := &models.Repo{
		Id:          gitRepo.GetID(),
		Name:        gitRepo.GetName(),
		Owner:       gitRepo.GetOwner().GetLogin(),
		Description: gitRepo.GetDescription(),
		URL:         gitRepo.GetHTMLURL(),
		ForkCount:   gitRepo.GetForksCount(),
		Language:    gitRepo.GetLanguage(),
	}

	return repo, nil
}

func (g GitBeamService) ListCommits(ctx context.Context, ownerName, repoName string) ([]*models.Commit, error) {
	useLogger := g.logger.WithContext(ctx).WithField("methodName", "ListCommits")

	var commits []*models.Commit
	var err error
	hasAttemptedRetry := false

retry:
	commits, err = g.dataStore.ListCommits(ctx)
	if err != nil {
		useLogger.WithError(err).Errorln("failed to list commits from database")
		return make([]*models.Commit, 0), nil
	}

	if len(commits) == 0 && !hasAttemptedRetry {
		_ = g.FetchAndSaveCommits(ctx, ownerName, repoName)
		hasAttemptedRetry = true
		goto retry
	}

	return commits, err
}

func (g GitBeamService) FetchAndSaveCommits(ctx context.Context, ownerName, repoName string) error {
	useLogger := g.logger.WithContext(ctx).WithField("methodName", "GetByOwnerAndRepoName")

	repo, err := g.GetByOwnerAndRepoName(ctx, ownerName, repoName)
	if err != nil {
		useLogger.WithError(err).Errorln("GetByOwnerAndRepoName")
		return err
	}

	gitCommits, _, err := g.githubClient.Repositories.ListCommits(ctx, repo.Owner, repo.Name, &github.CommitsListOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	})

	for _, gitCommit := range gitCommits {
		commit := &models.Commit{
			SHA:             gitCommit.GetSHA(),
			Message:         gitCommit.GetCommit().GetMessage(),
			Author:          gitCommit.GetCommit().Committer.GetLogin(),
			Date:            gitCommit.GetCommit().Committer.GetDate().Time,
			URL:             gitCommit.GetCommit().GetURL(),
			ParentCommitIDs: make([]string, 0),
		}

		parents := gitCommit.GetCommit().Parents
		for _, parent := range parents {
			commit.ParentCommitIDs = append(commit.ParentCommitIDs, parent.GetSHA())
		}

		if err := g.dataStore.SaveCommit(ctx, commit); err != nil {
			useLogger.WithError(err).Errorln("SaveCommit")
			return err
		}
	}

	return nil
}

func (g GitBeamService) StartBeamingCommits(ctx context.Context, payload models.OwnerAndRepoName) (*models.Repo, error) {
	useLogger := g.logger.WithContext(ctx).WithField("methodName", "StartBeamingCommits")
	repo, err := g.GetByOwnerAndRepoName(ctx, payload.OwnerName, payload.RepoName)
	if err != nil {
		useLogger.WithError(err).Errorln("GetByOwnerAndRepoName")
		return nil, err
	}

	if err := g.dataStore.StoreRepository(ctx, repo); err != nil {
		useLogger.WithError(err).Errorln("StoreRepository")
		return nil, err
	}

	data, err := json.Marshal(repo)
	if err != nil {
		useLogger.WithError(err).Errorln("json.Marshal: failed to marshal repo before publishing to event store.")
		return nil, err
	}

	// This is a channel-based event store, so checking of errors aren't needed here.
	_ = g.eventStore.Publish(topics.RepoCreated, data)

	return repo, nil
}
