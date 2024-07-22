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

func (g GitBeamService) GetLastCommit(ctx context.Context, owner models.OwnerAndRepoName) (*models.Commit, error) {
	useLogger := g.logger.WithContext(ctx).WithField("methodName", "GetCommitsBySha")
	commit, err := g.dataStore.GetLastCommit(ctx, &owner, nil)
	if err != nil {
		useLogger.WithError(err).Errorln("failed to fetch last commit from the dataStore")
		return nil, ErrCommitNotFound
	}

	return commit, nil
}

func (g GitBeamService) FetchAndSaveCommits(ctx context.Context, filters models.CommitFilters) error {
	useLogger := g.logger.WithContext(ctx).WithField("methodName", "FetchAndSaveCommits")
	pageNumber := 1

	//list, _ := g.dataStore.ListCommits(ctx, filters)
	//if len(list) > 1 {
	//	// TODO: use the list commit thing as a way to internally skip existing records before pulling new ones from from github.
	//	//filters.ToDate, _ = models.Parse(list[0].Date.Format(time.DateOnly))
	//	//if len(list) > 1 {
	//	//	filters.FromDate, _ = models.Parse(list[len(list)-1].Date.Format(time.DateOnly))
	//	//}
	//	return nil
	//}

	ghOptions := github.CommitsListOptions{
		ListOptions: github.ListOptions{
			Page:    pageNumber,
			PerPage: 1000,
		},
	}

	if filters.FromDate != nil {
		ghOptions.Since = filters.FromDate.Time
	}

	if filters.ToDate != nil {
		ghOptions.Until = filters.ToDate.Time
	}
run:
	gitCommits, response, err := g.githubClient.Repositories.ListCommits(ctx, filters.OwnerName, filters.OwnerName, &ghOptions)

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
		goto run
	}

	return nil
}
