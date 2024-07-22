package models

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"time"
)

type Commit struct {
	Date            time.Time `json:"date"`
	Message         string    `json:"message"`
	Author          string    `json:"author"`
	RepoName        string    `json:"repoName"`
	OwnerName       string    `json:"ownerName"`
	URL             string    `json:"url"`
	SHA             string    `json:"sha"`
	ParentCommitIDs []string  `json:"parentCommitIDs"`
}

func (commit Commit) Validate() error {
	return validation.ValidateStruct(&commit,
		validation.Field(&commit.SHA, validation.Required, is.Hexadecimal),
		validation.Field(&commit.Message, validation.Required),
		validation.Field(&commit.Author, validation.Required),
		validation.Field(&commit.URL, validation.Required, is.URL),
		validation.Field(&commit.OwnerName, validation.Required),
		validation.Field(&commit.RepoName, validation.Required),
		validation.Field(&commit.Date, validation.Required),
		validation.Field(&commit.ParentCommitIDs, validation.Required),
	)
}

type ListCommitFilter struct {
	Limit            int64 `json:"limit" schema:"limit,omitempty"`
	Page             int64 `json:"page" schema:"page,omitempty"`
	OwnerAndRepoName `json:",inline" schema:",inline"`
	FromDate         *Date `json:"fromDate" schema:"fromDate,omitempty"`
	ToDate           *Date `json:"toDate" schema:"toDate,omitempty"`
}

type Repo struct {
	TimeCreated   time.Time      `json:"timeCreated"`
	TimeUpdated   time.Time      `json:"timeUpdated"`
	Name          string         `json:"name"`
	Owner         string         `json:"owner"`
	Description   string         `json:"description"`
	URL           string         `json:"url"`
	Languages     string         `json:"language"`
	ForkCount     int            `json:"forkCounts"`
	StarCount     int            `json:"starCounts"`
	OpenIssues    int            `json:"openIssues"`
	WatchersCount int            `json:"watchersCount"`
	IsSaved       bool           `json:"isSaved"`
	Meta          map[string]any `json:"meta"`
}

func (r Repo) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required),
		validation.Field(&r.URL, validation.Required),
	)
}
