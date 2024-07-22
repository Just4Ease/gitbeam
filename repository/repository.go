package repository

import (
	"context"
	"gitbeam/models"
	"time"
)

//go:generate mockgen -source=repository.go -destination=../mocks/data_store_mock.go -package=mocks
type DataStore interface {
	StoreRepository(ctx context.Context, payload *models.Repo) error
	ListRepos(context.Context) ([]*models.Repo, error)
	GetRepoByOwner(ctx context.Context, owner *models.OwnerAndRepoName) (*models.Repo, error)

	SaveCommit(ctx context.Context, payload *models.Commit) error
	ListCommits(ctx context.Context, filter models.CommitFilters) ([]*models.Commit, error)
	GetLastCommit(ctx context.Context, owner *models.OwnerAndRepoName, startTime *time.Time) (*models.Commit, error)
	GetCommitBySHA(ctx context.Context, owner models.OwnerAndRepoName, sha string) (*models.Commit, error)
	GetTopCommitAuthors(ctx context.Context, filter models.CommitFilters) ([]*models.TopCommitAuthor, error)
}

type CronServiceStore interface {
	SaveCronTask(ctx context.Context, task models.CronTask) error
	GetCronTask(ctx context.Context, owner models.OwnerAndRepoName) (*models.CronTask, error)
	DeleteCronTask(ctx context.Context, owner models.OwnerAndRepoName) error
	ListCronTask(ctx context.Context) ([]*models.CronTask, error)
}
