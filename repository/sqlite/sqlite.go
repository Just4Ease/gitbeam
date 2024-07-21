package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"gitbeam/models"
	"gitbeam/repository"
	"time"
)

// In a real world application, I would use https://entgo.io/ for MySQL/SQLite/Postgresql ( RMDBs ) or mongodb directly.
// But for this exercise, without too many dependencies I'm using the native go sql driver on sqlite db.
type sqliteRepo struct {
	dataStore *sql.DB
}

func (s sqliteRepo) GetRepoByOwner(ctx context.Context, owner *models.OwnerAndRepoName) (*models.Repo, error) {
	row := s.dataStore.QueryRowContext(ctx,
		`SELECT * from repos WHERE owner_name = ? AND repo_name = ? LIMIT 1`, owner.OwnerName, owner.RepoName)
	return scanRepoRow(row)
}

func (s sqliteRepo) GetLastCommit(ctx context.Context, owner *models.OwnerAndRepoName, startTime *time.Time) (*models.Commit, error) {
	clause := `SELECT * from commits WHERE owner_name = ? AND repo_name = ?`
	if startTime != nil {
		clause = fmt.Sprintf("%s AND commit_date >= '%s'", clause, startTime.Format(time.RFC3339))
	}

	query := fmt.Sprintf(`%s ORDER BY commit_date DESC LIMIT 1`, clause)
	row := s.dataStore.QueryRowContext(ctx,
		query, owner.OwnerName, owner.RepoName)
	return scanCommitRow(row)
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

func (s sqliteRepo) CountSavedCommits(ctx context.Context, owner models.OwnerAndRepoName) (int64, error) {
	query := `SELECT COUNT(*) FROM commits WHERE owner_name = ? AND repo_name = ?`
	var count int64
	err := s.dataStore.QueryRowContext(ctx, query, owner.OwnerName, owner).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s sqliteRepo) ListCommits(ctx context.Context, filter models.ListCommitFilter) ([]*models.Commit, error) {

	if filter.Limit >= 100 || filter.Limit <= 0 {
		filter.Limit = 100
	}

	querySQL := `
        SELECT * FROM commits WHERE owner_name = ? AND repo_name =  ? 
        LIMIT ? OFFSET ?`

	rows, err := s.dataStore.QueryContext(ctx, querySQL,
		filter.Owner.OwnerName,
		filter.Owner.RepoName,
		filter.Limit,
		filter.Page,
	)

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

func (s sqliteRepo) GetCommitBySHA(ctx context.Context, owner models.OwnerAndRepoName, sha string) (*models.Commit, error) {
	row := s.dataStore.QueryRowContext(ctx, "SELECT * from commits WHERE owner_name = ? AND repo_name = ? AND sha = ? LIMIT 1", owner.OwnerName, owner.RepoName, sha)
	return scanCommitRow(row)
}

func (s sqliteRepo) SaveCommit(ctx context.Context, commit *models.Commit) error {
	if existingCommit, _ := s.GetCommitBySHA(ctx, models.OwnerAndRepoName{
		OwnerName: commit.OwnerName,
		RepoName:  commit.RepoName,
	}, commit.SHA); existingCommit != nil {
		// TODO: carry out commit update on disk this is in the case a commit message was updated via the git append command from it's source.
		return nil
	}

	insertSQL := `
        INSERT INTO commits (
            sha,
			message,
			author,
			repo_name,
			owner_name,
			url,
			parent_commit_ids,
			commit_date           
		)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	serializedParentCommitIds, err := json.Marshal(commit.ParentCommitIDs)
	if err != nil {
		fmt.Println("error serializing parent commit ids", err)
		return err
	}

	fmt.Println("commit.sha and parent: ", commit.SHA, string(serializedParentCommitIds))

	_, err = s.dataStore.ExecContext(ctx, insertSQL,
		commit.SHA,
		commit.Message,
		commit.Author,
		commit.RepoName,
		commit.OwnerName,
		commit.URL,
		string(serializedParentCommitIds),
		commit.Date.Format(time.RFC3339),
	)

	if err != nil {
		fmt.Println(err, " err writing.")
	}

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
