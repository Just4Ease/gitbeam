package core

import (
	"context"
	"encoding/json"
	"errors"
	"gitbeam/events/topics"
	"gitbeam/models"
	"github.com/google/go-github/v63/github"
	"time"
)

var (
	ErrCommitNotFound = errors.New("commit not found")
)

func (g GitBeamService) ListCommits(ctx context.Context, filters models.ListCommitFilter) ([]*models.Commit, error) {
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

func (g GitBeamService) GetCommitsBySha(ctx context.Context, owner models.OwnerAndRepoName, sha string) (*models.Commit, error) {
	useLogger := g.logger.WithContext(ctx).WithField("methodName", "GetCommitsBySha")
	commit, err := g.dataStore.GetCommitBySHA(ctx, owner, sha)
	if err != nil {
		useLogger.WithError(err).Errorln("failed to fetch commit by owner and sha details from the dataStore.")
		return nil, ErrCommitNotFound
	}

	return commit, nil
}

func (g GitBeamService) FetchAndSaveCommits(ctx context.Context, filters models.ListCommitFilter) error {
	useLogger := g.logger.WithContext(ctx).WithField("methodName", "FetchAndSaveCommits")
	pageNumber := 1

	//lastCommit, _ := g.dataStore.GetLastCommit(ctx, &filters.Owner, filters.StartTime)
	//if lastCommit != nil {
	//	*filters.StartTime = lastCommit.Date.Add(time.Millisecond)
	//	useLogger.WithFields(logrus.Fields{
	//		"repo_name":  owner.RepoName,
	//		"owner_name": owner.OwnerName,
	//		"start_time": startTimeCursor,
	//	}).Infoln("updated start time to be 1ms greater than the last record in the database so as not to repeat commits or waste rate limits.")
	//}

repeat:
	gitCommits, response, err := g.githubClient.Repositories.ListCommits(ctx, filters.OwnerName, filters.OwnerName, &github.CommitsListOptions{
		//Since: startTimeCursor,
		Until: time.Now(),
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

func (g GitBeamService) StartBeamingCommits(ctx context.Context, payload models.BeamRepoCommitsRequest) (*models.Repo, error) {
	useLogger := g.logger.WithContext(ctx).WithField("methodName", "StartBeamingCommits")
	repo, err := g.GetByOwnerAndRepoName(ctx, &models.OwnerAndRepoName{
		OwnerName: payload.OwnerName,
		RepoName:  payload.RepoName,
	})
	if err != nil {
		useLogger.WithError(err).Errorln("GetByOwnerAndRepoName")
		return nil, err
	}

	if !repo.IsSaved {
		if err := g.dataStore.StoreRepository(ctx, repo); err != nil {
			useLogger.WithError(err).Errorln("StoreRepository")
			return nil, err
		}
	}

	if payload.StartTime != nil {
		t, err := time.Parse(time.DateTime, *payload.StartTime)
		if err != nil {
			useLogger.WithError(err).Error("time.Parse: invalid start time provided, must confirm to YYYY-MM-DD HH:MM:SS")
			return nil, err
		}

		repo.Meta["startTime"] = t.Format(time.RFC3339) // To be consumed by the commit background activity.
	}

	data, err := json.Marshal(repo)
	if err != nil {
		useLogger.WithError(err).Errorln("json.Marshal: failed to marshal repo before publishing to event store.")
		return nil, err
	}

	// This is a channel-based event store, so checking of errors aren't needed here.
	_ = g.eventStore.Publish(topics.RepoCreated, data)

	return repo, nil
}
