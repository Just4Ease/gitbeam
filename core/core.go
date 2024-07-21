package core

import (
	"context"
	"errors"
	"gitbeam/models"
	"gitbeam/repository"
	"gitbeam/store"
	"gitbeam/utils"
	"github.com/google/go-github/v63/github"
	"github.com/sirupsen/logrus"
	"net/http"
)

var (
	ErrGithubRepoNotFound       = errors.New("github repo not found")
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

func (g GitBeamService) GetByOwnerAndRepoName(ctx context.Context, owner *models.OwnerAndRepoName) (*models.Repo, error) {
	useLogger := g.logger.WithContext(ctx).WithField("methodName", "GetByOwnerAndRepoName")

	if err := owner.Validate(); err != nil {
		useLogger.WithError(err).Error("owner is invalid. please provide a valid ownerName and repoName")
		return nil, ErrOwnerAndRepoNameRequired
	}

	existingRepo, err := g.dataStore.GetRepoByOwner(ctx, owner)
	if err == nil && existingRepo != nil {
		existingRepo.IsSaved = true
		return existingRepo, nil
	}

	gitRepo, _, err := g.githubClient.Repositories.Get(ctx, owner.OwnerName, owner.RepoName)
	if err != nil {
		useLogger.WithError(err).Errorln("GetByOwnerAndRepoName")
		return nil, ErrGithubRepoNotFound
	}

	repo := &models.Repo{
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

func (g GitBeamService) ListRepos(ctx context.Context) ([]*models.Repo, error) {
	useLogger := g.logger.WithContext(ctx).WithField("methodName", "ListRepos")
	list, err := g.dataStore.ListRepos(ctx)
	if err != nil {
		useLogger.WithError(err).Errorln("failed to list repositories")
		return make([]*models.Repo, 0), nil
	}
	return list, nil
}
