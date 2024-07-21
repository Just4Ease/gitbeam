package models

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"time"
)

type Commit struct {
	Message         string    `json:"message"`
	Author          string    `json:"author"`
	Date            time.Time `json:"date"`
	URL             string    `json:"url"`
	SHA             string    `json:"id"`
	RepoId          string    `json:"repoId"`
	ParentCommitIDs []string  `json:"parentCommitIDs"`
	Branch          string    `json:"branch"`
}

func (commit Commit) Validate() error {
	return validation.ValidateStruct(&commit,
		validation.Field(&commit.Message, validation.Required),
		validation.Field(&commit.Author, validation.Required),
		validation.Field(&commit.Date, validation.Required),
		validation.Field(&commit.RepoId, validation.Required),
		validation.Field(&commit.SHA, validation.Required),
		validation.Field(&commit.RepoId, validation.Required),
		validation.Field(&commit.Branch, validation.Required),
		validation.Field(&commit.ParentCommitIDs, validation.Required),
	)
}

type ListCommitFilter struct {
	Limit int64            `json:"limit"`
	Page  int64            `json:"page"`
	Owner OwnerAndRepoName `json:"owner"`
}

type Repo struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Owner       string `json:"owner"`
	Description string `json:"description"`
	URL         string `json:"url"`
	ForkCount   int    `json:"forkCounts"`
	Language    string `json:"language"`
}

func (r Repo) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required),
		validation.Field(&r.URL, validation.Required),
	)
}
