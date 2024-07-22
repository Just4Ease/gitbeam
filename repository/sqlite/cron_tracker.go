package sqlite

import (
	"context"
	"database/sql"
	"gitbeam/models"
	"time"
)

const cronTrackerTableSetup = `
CREATE TABLE IF NOT EXISTS commits (
		repo_name TEXT,
		owner_name TEXT,
		next_tick DATETIME,
		from_date DATETIME,
		to_date DATETIME,
		UNIQUE (repo_name, owner_name)
)
`

func scanCronTrackerRow(row *sql.Row) (*models.CronTracker, error) {
	var nextTickString string
	var fromDateString string
	var toDateString string
	var cronTracker models.CronTracker
	var err error
	if err = row.Scan(
		&cronTracker.RepoName,
		&cronTracker.OwnerName,
		&nextTickString,
		&fromDateString,
		&toDateString,
	); err != nil {
		return nil, err
	}

	cronTracker.NextTick, err = time.Parse(time.RFC3339, nextTickString)
	if err != nil {
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

func (s sqliteRepo) SaveCronTracker(ctx context.Context, tracker models.CronTracker) error {
	//TODO implement me
	panic("implement me")
}

func (s sqliteRepo) GetCronTracker(ctx context.Context, owner models.OwnerAndRepoName) (*models.CronTracker, error) {
	//TODO implement me
	panic("implement me")
}

func (s sqliteRepo) DeleteCronTracker(ctx context.Context, owner models.OwnerAndRepoName) error {
	//TODO implement me
	panic("implement me")
}
