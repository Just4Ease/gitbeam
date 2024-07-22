package sqlite

import (
	"database/sql"
	"gitbeam/repository"
)

// In a real world application, I would use https://entgo.io/ for MySQL/SQLite/Postgresql ( RMDBs ) or mongodb directly.
// But for this exercise, without too many dependencies I'm using the native go sql driver on sqlite db.
type sqliteRepo struct {
	dataStore *sql.DB
}

func NewSqliteRepo(db *sql.DB) (repository.DataStore, error) {
	if _, err := db.Exec(repoTableSetup); err != nil {
		return nil, err
	}
	if _, err := db.Exec(commitsTableSetup); err != nil {
		return nil, err
	}
	if _, err := db.Exec(cronTrackerTableSetup); err != nil {
		return nil, err
	}
	return &sqliteRepo{
		dataStore: db,
	}, nil
}
