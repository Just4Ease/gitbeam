package core

import (
	"context"
	"errors"
	"gitbeam/models"
	"gitbeam/repository"
	"gitbeam/store"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/go-github/v63/github"
	"github.com/sirupsen/logrus"
	"time"
)

var (
	ErrRepositoryNotFound       = errors.New("repository not found")
	ErrOwnerAndRepoNameRequired = errors.New("owner and repo name required")
)

type GitBeamService struct {
	githubClient *github.Client
	logger       *logrus.Logger
	repository   repository.DataStore
}

func NewGitBeamService(
	logger *logrus.Logger,
	eventStore store.EventStore,
	dataStore repository.DataStore,
) *GitBeamService {
	client := github.NewClient(nil) // Didn't need to pass this as a top level dependency into the git beam service.
	return &GitBeamService{
		githubClient: client,
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

	commits, err := g.repository.ListCommits(ctx)
	if err != nil {
		useLogger.WithError(err).Errorln("failed to list commits from database")
		return make([]*models.Commit, 0), nil
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
		Since: time.Time{},
		Until: time.Time{},
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	})

	for _, gitCommit := range gitCommits {
		commit := &models.Commit{
			Message: gitCommit.GetCommit().GetMessage(),
			Author:  gitCommit.GetCommit().Committer.GetLogin(),
			Date:    gitCommit.GetCommit().Committer.GetDate().Time,
			URL:     gitCommit.GetCommit().GetURL(),
			Id:      gitCommit.GetNodeID(),
			//RepoId:  gitCommit.R,
		}

		if err := g.repository.SaveCommit(ctx, commit); err != nil {
			useLogger.WithError(err).Errorln("SaveCommit")
			return err
		}
	}

	return nil
}
