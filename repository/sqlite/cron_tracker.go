package sqlite

import (
	"context"
	"database/sql"
	"gitbeam/models"
	"gitbeam/repository"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

const cronTrackerTableSetup = `
CREATE TABLE IF NOT EXISTS cron_tracker (
		repo_name TEXT,
		owner_name TEXT,
		from_date DATETIME,
		to_date DATETIME,
		UNIQUE (repo_name, owner_name)
)
`

func scanCronTrackerRow(row *sql.Row) (*models.CronTracker, error) {
	var fromDateString string
	var toDateString string
	var cronTracker models.CronTracker
	var err error
	if err = row.Scan(
		&cronTracker.RepoName,
		&cronTracker.OwnerName,
		&fromDateString,
		&toDateString,
	); err != nil {
		return nil, err
	}

	if fromDateString != "" {
		t, _ := time.Parse(time.RFC3339, fromDateString)
		cronTracker.FromDate = &t
	}

	if toDateString != "" {
		t, _ := time.Parse(time.RFC3339, toDateString)
		cronTracker.ToDate = &t
	}

	return &cronTracker, nil
}

func scanCronTrackerRows(rows *sql.Rows) (*models.CronTracker, error) {
	var fromDateString string
	var toDateString string
	var cronTracker models.CronTracker
	var err error
	if err = rows.Scan(
		&cronTracker.RepoName,
		&cronTracker.OwnerName,
		&fromDateString,
		&toDateString,
	); err != nil {
		return nil, err
	}

	if fromDateString != "" {
		t, _ := time.Parse(time.RFC3339, fromDateString)
		cronTracker.FromDate = &t
	}

	if toDateString != "" {
		t, _ := time.Parse(time.RFC3339, toDateString)
		cronTracker.ToDate = &t
	}

	return &cronTracker, nil
}

func (s sqliteRepo) ListCronTrackers(ctx context.Context) ([]*models.CronTracker, error) {
	querySQL := `SELECT * FROM cron_tracker`

	rows, err := s.dataStore.QueryContext(ctx, querySQL)
	if err != nil {
		return nil, err
	}

	var list []*models.CronTracker
	defer rows.Close()
	for rows.Next() {
		item, err := scanCronTrackerRows(rows)
		if err != nil {
			return nil, err
		}

		list = append(list, item)
	}

	return list, nil
}

func (s sqliteRepo) SaveCronTracker(ctx context.Context, payload models.CronTracker) error {
	insertSQL := `
        INSERT INTO cron_tracker (
			repo_name,
			owner_name,
			from_date,
			to_date
		)
        VALUES (?, ?, ?, ?)`

	_, err := s.dataStore.ExecContext(ctx, insertSQL,
		payload.RepoName,
		payload.OwnerName,
		payload.FromDate,
		payload.ToDate,
	)
	return err
}

func (s sqliteRepo) GetCronTracker(ctx context.Context, owner models.OwnerAndRepoName) (*models.CronTracker, error) {
	row := s.dataStore.QueryRowContext(ctx,
		`SELECT * from cron_tracker WHERE owner_name = ? AND repo_name = ? LIMIT 1`, owner.OwnerName, owner.RepoName)
	return scanCronTrackerRow(row)
}

func (s sqliteRepo) DeleteCronTracker(ctx context.Context, owner models.OwnerAndRepoName) error {
	_, err := s.dataStore.ExecContext(ctx,
		`DELETE from cron_tracker WHERE owner_name = ? AND repo_name = ?`, owner.OwnerName, owner.RepoName)

	return err
}

func NewSqliteCronStore(dbName string) (repository.CronServiceStore, error) {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(cronTrackerTableSetup); err != nil {
		return nil, err
	}
	return &sqliteRepo{
		dataStore: db,
	}, nil
}
