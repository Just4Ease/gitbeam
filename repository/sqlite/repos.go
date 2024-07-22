package sqlite

import (
	"context"
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

func (s sqliteRepo) GetRepoByOwner(ctx context.Context, owner *models.OwnerAndRepoName) (*models.Repo, error) {
	row := s.dataStore.QueryRowContext(ctx,
		`SELECT * from repos WHERE owner_name = ? AND repo_name = ? LIMIT 1`, owner.OwnerName, owner.RepoName)
	return scanRepoRow(row)
}

func (s sqliteRepo) StoreRepository(ctx context.Context, payload *models.Repo) error {
	insertSQL := `
        INSERT INTO repos (
			repo_name,
			owner_name,
			description,
			url,
			repo_languages,
			meta,
			forks_count,
			stars_count,
			watchers_count,
			open_issues_count,
			time_created,
			time_updated
		)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	meta, _ := json.Marshal(payload.Meta)
	_, err := s.dataStore.ExecContext(ctx, insertSQL,
		payload.Name,
		payload.Owner,
		payload.Description,
		payload.URL,
		payload.Languages,
		string(meta),
		payload.ForkCount,
		payload.StarCount,
		payload.WatchersCount,
		payload.OpenIssues,
		payload.TimeCreated.Format(time.RFC3339),
		payload.TimeUpdated.Format(time.RFC3339),
	)

	return err
}

func (s sqliteRepo) ListRepos(ctx context.Context) ([]*models.Repo, error) {
	// Define the SQL query for listing transaction history by source account ID.
	querySQL := `
        SELECT * FROM repos` // TODO: Apply repo filters.

	rows, err := s.dataStore.QueryContext(ctx, querySQL)
	if err != nil {
		return nil, err
	}
	var repos []*models.Repo
	defer rows.Close()
	for rows.Next() {
		repo, err := scanRepoRows(rows)
		if err != nil {
			return nil, err
		}

		repos = append(repos, repo)
	}

	return repos, nil
}
