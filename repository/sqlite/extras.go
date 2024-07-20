package sqlite

import (
	"database/sql"
	"gitbeam/models"
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
		url TEXT,
		commit_timestamp DATETIME,
		parent_commit_id TEXT,
		branch TEXT
)
`

func scanCommitRows(rows *sql.Rows) (*models.Commit, error) {
	var commit models.Commit
	if err := rows.Scan(
		&commit.Id,
		&commit.Message,
		&commit.Author,
		&commit.URL,
		&commit.Date,
		&commit.ParentCommitId,
		&commit.Branch,
	); err != nil {
		return nil, err
	}

	return &commit, nil
}

func scanCommitRow(row *sql.Row) (*models.Commit, error) {
	var commit models.Commit
	if err := row.Scan(
		&commit.Id,
		&commit.Message,
		&commit.Author,
		&commit.URL,
		&commit.Date,
		&commit.ParentCommitId,
		&commit.Branch,
	); err != nil {
		return nil, err
	}

	return &commit, nil
}

func scanRepoRows(rows *sql.Rows) (*models.Repo, error) {
	var repo models.Repo
	if err := rows.Scan(
		&repo.Id,
		&repo.Name,
		&repo.Owner,
		&repo.Description,
		&repo.URL,
		&repo.ForkCount,
		&repo.Language,
	); err != nil {
		return nil, err
	}

	return &repo, nil
}

func scanRepoRow(row *sql.Row) (*models.Repo, error) {
	var repo models.Repo
	if err := row.Scan(
		&repo.Id,
		&repo.Name,
		&repo.Owner,
		&repo.Description,
		&repo.URL,
		&repo.ForkCount,
		&repo.Language,
	); err != nil {
		return nil, err
	}

	return &repo, nil
}
