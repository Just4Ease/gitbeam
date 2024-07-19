package sqlite

import (
	"context"
	"database/sql"
	"gitbeam/models"
	"gitbeam/repository"
)

const repoTableSetup = `
CREATE TABLE IF NOT EXISTS repos (
    	id INT PRIMARY KEY UNIQUE,
		repo_name TEXT,
		owner_name TEXT,
		description TEXT,
		url TEXT,
		fork_count INT,
        repo_language TEXT
)
`

const commitsTableSetup = `
CREATE TABLE IF NOT EXISTS commits (
    	id INT PRIMARY KEY UNIQUE,
		message FLOAT,
		author TEXT,
		commit_timestamp DATETIME,
		parent_commit_id TEXT
)
`

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

func NewSqliteRepo(db *sql.DB) (repository.DataStore, error) {
	if _, err := db.Exec(repoTableSetup); err != nil {
		return nil, err
	}
	if _, err := db.Exec(commitsTableSetup); err != nil {
		return nil, err
	}
	return &sqliteRepo{
		dataStore: db,
	}, nil
}
