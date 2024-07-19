package sqlite

import (
	"context"
	"database/sql"
	"gitbeam/models"
	"gitbeam/repository"
)

type sqliteRepo struct {
	dataStore *sql.DB
}

func (s sqliteRepo) ListCommits(ctx context.Context) ([]*models.Commit, error) {
	//TODO implement me
	panic("implement me")
}

func (s sqliteRepo) GetCommitById(ctx context.Context, id string) (*models.Commit, error) {
	//TODO implement me
	panic("implement me")
}

func (s sqliteRepo) SaveCommit(ctx context.Context, payload *models.Commit) error {
	//TODO implement me
	panic("implement me")
}

func NewSqliteRepo() repository.Repository {
	return &sqliteRepo{}
}
