package core

import (
	"context"
	"encoding/json"
	"errors"
	"gitbeam/events/topics"
	"gitbeam/models"
	"gitbeam/repository"
	"gitbeam/store"
	"gitbeam/utils"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/go-github/v63/github"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
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

	repo := &models.Repo{
		Id:            gitRepo.GetID(),
		Name:          gitRepo.GetName(),
		Owner:         gitRepo.GetOwner().GetLogin(),
		Description:   gitRepo.GetDescription(),
		URL:           gitRepo.GetHTMLURL(),
		Languages:     gitRepo.GetLanguage(),
		ForkCount:     gitRepo.GetForksCount(),
		StarCount:     gitRepo.GetStargazersCount(),
		OpenIssues:    gitRepo.GetOpenIssues(),
		WatchersCount: gitRepo.GetWatchersCount(),
		TimeCreated:   gitRepo.GetCreatedAt().Time,
		TimeUpdated:   gitRepo.GetUpdatedAt().Time,
		Meta:          make(map[string]any),
	}

	// Take the raw git repo response:, gitRepo -> ([]bytes||string) -> map[string]any
	// Intentionally ignoring this error message.
	// Note, this field can be removed totally... It serves no purpose at the moment.
	_ = utils.UnPack(gitRepo, &repo.Meta)

	return repo, nil
}

func (g GitBeamService) ListCommits(ctx context.Context, filters models.ListCommitFilter) ([]*models.Commit, error) {
	useLogger := g.logger.WithContext(ctx).WithField("methodName", "ListCommits")

	var commits []*models.Commit
	var err error
	hasAttemptedRetry := false

retry:
	commits, err = g.dataStore.ListCommits(ctx, filters)
	if err != nil {
		useLogger.WithError(err).Errorln("failed to list commits from database")
		return make([]*models.Commit, 0), nil
	}

	if len(commits) == 0 && !hasAttemptedRetry {
		commit, err := g.dataStore.GetLastCommit(ctx, filters.Owner)
		if err != nil {
			return make([]*models.Commit, 0), nil
		}

		_ = g.FetchAndSaveCommits(ctx, &filters.Owner, commit.Date)
		hasAttemptedRetry = true
		goto retry
	}

	return commits, err
}

func (g GitBeamService) FetchAndSaveCommits(ctx context.Context, owner *models.OwnerAndRepoName, startTimeCursor time.Time) error {
	useLogger := g.logger.WithContext(ctx).WithField("methodName", "GetByOwnerAndRepoName")

	repo, err := g.GetByOwnerAndRepoName(ctx, owner.OwnerName, owner.RepoName)
	if err != nil {
		useLogger.WithError(err).Errorln("xxxGetByOwnerAndRepoName")
		return err
	}

	pageNumber := 1

repeat:
	gitCommits, _, err := g.githubClient.Repositories.ListCommits(ctx, repo.Owner, repo.Name, &github.CommitsListOptions{
		Since: startTimeCursor,
		Until: time.Now(),
		ListOptions: github.ListOptions{
			Page:    pageNumber,
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

	// if the previous/above attempt to list commits from github had data, then let's check if a new page will have data,
	// else we exit until FetchAndSaveCommits is called by a cron.
	if len(gitCommits) > 0 {
		pageNumber += 1
		goto repeat
	}

	return nil
}

func (g GitBeamService) StartBeamingCommits(ctx context.Context, payload models.BeamRepoCommitsRequest) (*models.Repo, error) {
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

	if payload.StartTime != nil {
		t, err := time.Parse(time.DateTime, *payload.StartTime)
		if err != nil {
			useLogger.WithError(err).Error("time.Parse: invalid start time provided, must confirm to YYYY-MM-DD HH:MM:SS")
			return nil, err
		}

		repo.Meta["startTime"] = t // To be consumed by the commit background activity.
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
