package repository

import (
	"context"
	"gitbeam/models"
)

type DataStore interface {
	ListCommits(ctx context.Context) ([]*models.Commit, error)
	GetCommitById(ctx context.Context, id string) (*models.Commit, error)
	SaveCommit(ctx context.Context, payload *models.Commit) error
}
