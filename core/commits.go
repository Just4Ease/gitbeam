package core

import (
	"context"
	"errors"
	"gitbeam/models"
	"github.com/google/go-github/v63/github"
)

var (
	ErrCommitNotFound = errors.New("commit not found")
)

func (g GitBeamService) ListCommits(ctx context.Context, filters models.CommitFilters) ([]*models.Commit, error) {
	useLogger := g.logger.WithContext(ctx).WithField("methodName", "ListCommits")

	var commits []*models.Commit
	var err error
	hasAttemptedRetry := false

retry:
	commits, err = g.dataStore.ListCommits(ctx, filters)
	if err != nil {
		useLogger.WithError(err).Errorln("failed to list commits from database")
		return make([]*models.Commit, 0), nil
	}

	if len(commits) == 0 && !hasAttemptedRetry {
		_ = g.FetchAndSaveCommits(ctx, filters)
		hasAttemptedRetry = true
		goto retry
	}

	return commits, err
}

func (g GitBeamService) GetTopCommitAuthors(ctx context.Context, filters models.CommitFilters) ([]*models.TopCommitAuthor, error) {
	useLogger := g.logger.WithContext(ctx).WithField("methodName", "GetTopCommitAuthors")

	list, err := g.dataStore.GetTopCommitAuthors(ctx, filters)
	if err != nil {
		useLogger.WithError(err).Errorln("failed to list top commit author from database")
		return make([]*models.TopCommitAuthor, 0), nil
	}

	return list, nil
}

func (g GitBeamService) GetCommitsBySha(ctx context.Context, owner models.OwnerAndRepoName, sha string) (*models.Commit, error) {
	useLogger := g.logger.WithContext(ctx).WithField("methodName", "GetCommitsBySha")
	commit, err := g.dataStore.GetCommitBySHA(ctx, owner, sha)
	if err != nil {
		useLogger.WithError(err).Errorln("failed to fetch commit by owner and sha details from the dataStore.")
		return nil, ErrCommitNotFound
	}

	return commit, nil
}

func (g GitBeamService) FetchAndSaveCommits(ctx context.Context, filters models.CommitFilters) error {
	useLogger := g.logger.WithContext(ctx).WithField("methodName", "FetchAndSaveCommits")
	pageNumber := 1

	// TODO: use the list commit thing as a way to internally skip existing records before pulling them from github.
	//list, _ := g.dataStore.ListCommits(ctx, filters)

	//lastCommit, _ := g.dataStore.GetLastCommit(ctx, &filters.Owner, filters.FromDate)
	//if lastCommit != nil {
	//	*filters.FromDate = lastCommit.Date.Add(time.Millisecond)
	//	useLogger.WithFields(logrus.Fields{
	//		"repo_name":  owner.RepoName,
	//		"owner_name": owner.OwnerName,
	//		"start_time": startTimeCursor,
	//	}).Infoln("updated start time to be 1ms greater than the last record in the database so as not to repeat commits or waste rate limits.")
	//}

repeat:
	gitCommits, response, err := g.githubClient.Repositories.ListCommits(ctx, filters.OwnerName, filters.OwnerName, &github.CommitsListOptions{
		Since: filters.FromDate.Time,
		Until: filters.ToDate.Time,
		ListOptions: github.ListOptions{
			Page:    pageNumber,
			PerPage: 1000,
		},
	})

	if err != nil {
		useLogger.WithError(err).Error("failed to list commits from github")
		return err
	}

	useLogger.WithField("rate_limits", response.Rate).Infoln("Raw API response from GitHub")

	for _, gitCommit := range gitCommits {
		c := gitCommit.GetCommit()

		commit := &models.Commit{
			SHA:             gitCommit.GetSHA(),
			Message:         c.GetMessage(),
			Author:          gitCommit.GetCommitter().GetLogin(),
			Date:            c.Committer.GetDate().Time,
			URL:             gitCommit.GetHTMLURL(),
			OwnerName:       filters.OwnerName,
			RepoName:        filters.RepoName,
			ParentCommitIDs: make([]string, 0),
		}

		parents := gitCommit.Parents
		for _, parent := range parents {
			commit.ParentCommitIDs = append(commit.ParentCommitIDs, parent.GetSHA())
		}

		if err := g.dataStore.SaveCommit(ctx, commit); err != nil {
			useLogger.WithError(err).Errorln("error saving commit to storage.")
			return err
		}
	}

	// if the previous/above attempt to list commits from github had data, then let's check if a new page will have data,
	// else we exit until FetchAndSaveCommits is called by a cron.
	if len(gitCommits) > 0 && response.Rate.Remaining > 0 { // TODO: Apply rate limiting rules to respect github's rate limit flow.
		pageNumber += 1
		goto repeat
	}

	return nil
}
