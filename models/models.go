package models

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"time"
)

type Commit struct {
	Date            time.Time `json:"date"`
	Message         string    `json:"message"`
	Author          string    `json:"author"`
	URL             string    `json:"url"`
	SHA             string    `json:"id"`
	RepoId          string    `json:"repoId"`
	ParentCommitIDs []string  `json:"parentCommitIDs"`
}

func (commit Commit) Validate() error {
	return validation.ValidateStruct(&commit,
		validation.Field(&commit.Message, validation.Required),
		validation.Field(&commit.Author, validation.Required),
		validation.Field(&commit.Date, validation.Required),
		validation.Field(&commit.RepoId, validation.Required),
		validation.Field(&commit.SHA, validation.Required),
		validation.Field(&commit.RepoId, validation.Required),
		validation.Field(&commit.ParentCommitIDs, validation.Required),
	)
}

type ListCommitFilter struct {
	Limit int64            `json:"limit"`
	Page  int64            `json:"page"`
	Owner OwnerAndRepoName `json:"owner"`
}

func (l ListCommitFilter) Validate() error {
	return validation.ValidateStruct(&l,
		validation.Field(&l.Owner, validation.Required),
		validation.Field(&l.Owner.OwnerName, validation.Required),
		validation.Field(&l.Owner.RepoName, validation.Required),
	)
}

type Repo struct {
	Id            int64          `json:"id"`
	Name          string         `json:"name"`
	Owner         string         `json:"owner"`
	Description   string         `json:"description"`
	URL           string         `json:"url"`
	Languages     string         `json:"language"`
	Meta          map[string]any `json:"meta"`
	ForkCount     int            `json:"forkCounts"`
	StarCount     int            `json:"starCounts"`
	OpenIssues    int            `json:"openIssues"`
	WatchersCount int            `json:"watchersCount"`
	TimeCreated   time.Time      `json:"timeCreated"`
	TimeUpdated   time.Time      `json:"timeUpdated"`
}

func (r Repo) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required),
		validation.Field(&r.URL, validation.Required),
	)
}
