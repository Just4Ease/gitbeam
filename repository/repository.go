package repository

import (
	"context"
	"gitbeam/models"
)

//go:generate mockgen -source=repository.go -destination=../mocks/data_store_mock.go -package=mocks
type DataStore interface {
	StoreRepository(ctx context.Context, payload *models.Repo) error
	ListRepos(context.Context) ([]*models.Repo, error)
	GetRepoByOwner(ctx context.Context, owner *models.OwnerAndRepoName) (*models.Repo, error)
	ListCommits(ctx context.Context, filter models.ListCommitFilter) ([]*models.Commit, error)
	CountSavedCommits(ctx context.Context, owner models.OwnerAndRepoName) (int64, error)
	GetLastCommit(ctx context.Context, owner models.OwnerAndRepoName) (*models.Commit, error)
	GetCommitBySHA(ctx context.Context, sha string) (*models.Commit, error)
	SaveCommit(ctx context.Context, payload *models.Commit) error
}
