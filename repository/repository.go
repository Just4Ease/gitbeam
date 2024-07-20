package repository

import (
	"context"
	"gitbeam/models"
)

type DataStore interface {
	StoreRepository(ctx context.Context, payload *models.Repo) error
	ListRepos(context.Context) ([]*models.Repo, error)

	// TODO: Define commit filters.
	ListCommits(ctx context.Context) ([]*models.Commit, error)
	GetCommitById(ctx context.Context, id string) (*models.Commit, error)
	SaveCommit(ctx context.Context, payload *models.Commit) error
}
