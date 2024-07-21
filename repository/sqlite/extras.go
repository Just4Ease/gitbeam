package sqlite

import (
	"database/sql"
	"encoding/json"
	"gitbeam/models"
	"time"
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
    	sha TEXT PRIMARY KEY UNIQUE,
		message FLOAT,
		author TEXT,
		url TEXT,
		commit_timestamp DATETIME,
		parent_commit_ids TEXT,
		branch TEXT
)
`

func scanCommitRows(rows *sql.Rows) (*models.Commit, error) {
	var dateString string
	var serializedParentCommitIDs string
	var commit models.Commit
	var err error
	if err = rows.Scan(
		&commit.SHA,
		&commit.Message,
		&commit.Author,
		&commit.URL,
		&dateString,
		&serializedParentCommitIDs,
		&commit.Branch,
	); err != nil {
		return nil, err
	}

	commit.Date, err = time.Parse(time.RFC3339, dateString)
	if err != nil {
		return nil, err
	}

	commit.ParentCommitIDs, err = deserializeParentCommitIds(serializedParentCommitIDs)
	if err != nil {
		return nil, err
	}

	return &commit, nil
}

func scanCommitRow(row *sql.Row) (*models.Commit, error) {
	var dateString string
	var serializedParentCommitIDs string
	var commit models.Commit
	var err error
	if err = row.Scan(
		&commit.SHA,
		&commit.Message,
		&commit.Author,
		&commit.URL,
		&dateString,
		&serializedParentCommitIDs,
		&commit.Branch,
	); err != nil {
		return nil, err
	}

	commit.Date, err = time.Parse(time.RFC3339, dateString)
	if err != nil {
		return nil, err
	}

	commit.ParentCommitIDs, err = deserializeParentCommitIds(serializedParentCommitIDs)
	if err != nil {
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

func deserializeParentCommitIds(data string) ([]string, error) {
	var ids []string
	err := json.Unmarshal([]byte(data), &ids)
	if err != nil {
		return nil, err
	}
	return ids, nil
}
