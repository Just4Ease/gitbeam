package core

import (
	"context"
	"errors"
	"gitbeam/models"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/go-github/v63/github"
	"github.com/sirupsen/logrus"
)

var (
	ErrRepositoryNotFound       = errors.New("repository not found")
	ErrOwnerAndRepoNameRequired = errors.New("owner and repo name required")
)

type GitBeamService struct {
	githubClient *github.Client
	logger       *logrus.Logger
}

func NewGitBeamService(logger *logrus.Logger) *GitBeamService {
	client := github.NewClient(nil) // Didn't need to pass this as a top level dependency into the git beam service.
	return &GitBeamService{
		githubClient: client,
		logger:       logger.WithField("serviceName", "GitBeamService").Logger,
	}
}

func (g GitBeamService) getByOwnerAndRepoName(ctx context.Context, ownerName, repoName string) (*models.Repo, error) {
	useLogger := g.logger.WithContext(ctx).WithField("methodName", "getByOwnerAndRepoName")

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
		useLogger.WithError(err).Errorln("getByOwnerAndRepoName")
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
