package sqlite

import (
	"database/sql"
	"encoding/json"
	"gitbeam/models"
	"gitbeam/utils"
	"time"
)

const repoTableSetup = `
CREATE TABLE IF NOT EXISTS repos (
		repo_name TEXT,
		owner_name TEXT,
		description TEXT,
		url TEXT,
		repo_languages TEXT,
		meta TEXT,
		forks_count INT,
		stars_count INT,
		watchers_count INT,
		open_issues_count INT,
		time_created DATETIME,
		time_updated DATETIME,
		UNIQUE (repo_name, owner_name)
)
`

func scanRepoRow(row *sql.Row) (*models.Repo, error) {
	var repo models.Repo
	var meta string
	if err := row.Scan(
		&repo.Name,
		&repo.Owner,
		&repo.Description,
		&repo.URL,
		&repo.Languages,
		&meta,
		&repo.ForkCount,
		&repo.StarCount,
		&repo.WatchersCount,
		&repo.OpenIssues,
		&repo.TimeCreated,
		&repo.TimeUpdated,
	); err != nil {
		return nil, err
	}

	_ = utils.UnPack(meta, &repo.Meta)
	if repo.Meta == nil {
		repo.Meta = make(map[string]any)
	}
	return &repo, nil
}

func scanRepoRows(rows *sql.Rows) (*models.Repo, error) {
	var repo models.Repo
	var meta string
	if err := rows.Scan(
		&repo.Name,
		&repo.Owner,
		&repo.Description,
		&repo.URL,
		&repo.Languages,
		&meta,
		&repo.ForkCount,
		&repo.StarCount,
		&repo.WatchersCount,
		&repo.OpenIssues,
		&repo.TimeCreated,
		&repo.TimeUpdated,
	); err != nil {
		return nil, err
	}

	_ = utils.UnPack(meta, &repo.Meta)
	if repo.Meta == nil {
		repo.Meta = make(map[string]any)
	}
	return &repo, nil
}

const commitsTableSetup = `
CREATE TABLE IF NOT EXISTS commits (
    	sha TEXT PRIMARY KEY,
		message TEXT,
		author TEXT,
		repo_name TEXT,
		owner_name TEXT,
		url TEXT,
		parent_commit_ids TEXT,
		commit_date DATETIME,
		UNIQUE (repo_name, owner_name, sha)
)
`

func deserializeParentCommitIds(data string) ([]string, error) {
	var ids []string
	err := json.Unmarshal([]byte(data), &ids)
	if err != nil {
		return nil, err
	}
	return ids, nil
}

func scanCommitRows(rows *sql.Rows) (*models.Commit, error) {
	var dateString string
	var serializedParentCommitIDs string
	var commit models.Commit
	var err error
	if err = rows.Scan(
		&commit.SHA,
		&commit.Message,
		&commit.Author,
		&commit.RepoName,
		&commit.OwnerName,
		&commit.URL,
		&serializedParentCommitIDs,
		&dateString,
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
		&commit.RepoName,
		&commit.OwnerName,
		&commit.URL,
		&serializedParentCommitIDs,
		&dateString,
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
