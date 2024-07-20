package sqlite

import (
	"context"
	"database/sql"
	"gitbeam/models"
	"gitbeam/repository"
)

// In a real world application, I would use https://entgo.io/ for MySQL/SQLite/Postgresql ( RMDBs ) or mongodb directly.
// But for this exercise, without too many dependencies I'm using the native go sql driver on sqlite db.
type sqliteRepo struct {
	dataStore *sql.DB
}

func (s sqliteRepo) StoreRepository(ctx context.Context, payload *models.Repo) error {
	insertSQL := `
        INSERT INTO repos (id, repo_name, owner_name, description, url, fork_count, repo_language)
        VALUES (?, ?, ?, ?, ?, ?, ?)`

	_, err := s.dataStore.ExecContext(ctx, insertSQL,
		payload.Id,
		payload.Name,
		payload.Owner,
		payload.Description,
		payload.URL,
		payload.ForkCount,
		payload.Language,
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

func (s sqliteRepo) ListCommits(ctx context.Context) ([]*models.Commit, error) {
	// Define the SQL query for listing transaction history by source account ID.
	querySQL := `
        SELECT * FROM commits` // TODO: Apply commits filters.

	rows, err := s.dataStore.QueryContext(ctx, querySQL)
	if err != nil {
		return nil, err
	}
	var commits []*models.Commit
	defer rows.Close()
	for rows.Next() {
		commit, err := scanCommitRows(rows)
		if err != nil {
			return nil, err
		}

		commits = append(commits, commit)
	}

	return commits, nil
}

func (s sqliteRepo) GetCommitById(ctx context.Context, id string) (*models.Commit, error) {
	row := s.dataStore.QueryRowContext(ctx, "SELECT * from commits WHERE id = ?", id)
	return scanCommitRow(row)
}

func (s sqliteRepo) SaveCommit(ctx context.Context, commit *models.Commit) error {
	// TODO: Get duplicate and update or overwrite it with the commit data.
	// ( this is in the case a commit message was updated via the git append command from it's source )
	//s.getCommitDuplicateById(ctx, commit.Id)

	insertSQL := `
        INSERT INTO commits (id, message, author, url, commit_timestamp, parent_commit_id, branch)
        VALUES (?, ?, ?, ?, ?, ?, ?)`

	_, err := s.dataStore.ExecContext(ctx, insertSQL,
		commit.Id,
		commit.Message,
		commit.Author,
		commit.URL,
		commit.Date,
		commit.ParentCommitId,
		commit.Branch,
	)

	return err
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
