package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"gitbeam/models"
	"time"
)

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

func (s sqliteRepo) ListCommits(ctx context.Context, filter models.ListCommitFilter) ([]*models.Commit, error) {

	if filter.Limit <= 0 {
		filter.Limit = 100
	}

	clause := `SELECT * from commits WHERE owner_name = ? AND repo_name = ?`
	if filter.FromDate != nil {
		clause = fmt.Sprintf("%s AND commit_date >= '%s'", clause, filter.FromDate.Format(time.RFC3339))
	}

	if filter.ToDate != nil {
		clause = fmt.Sprintf(`%s AND commit_date <= '%s'`, clause, filter.ToDate.Format(time.RFC3339))
	}

	query := fmt.Sprintf(`%s ORDER BY commit_date DESC LIMIT ? OFFSET ?`, clause)

	rows, err := s.dataStore.QueryContext(ctx, query,
		filter.OwnerName,
		filter.RepoName,
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
