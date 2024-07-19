package models

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"time"
)

type Commit struct {
	Message string    `json:"message"`
	Author  string    `json:"author"`
	Date    time.Time `json:"date"`
	URL     string    `json:"url"`
	Id      string    `json:"id"`
	RepoId  string    `json:"repoId"`
}

func (commit Commit) Validate() error {
	return validation.ValidateStruct(&commit,
		validation.Field(&commit.Message, validation.Required),
		validation.Field(&commit.Author, validation.Required),
		validation.Field(&commit.Date, validation.Required),
		validation.Field(&commit.RepoId, validation.Required),
		validation.Field(&commit.Id, validation.Required),
		validation.Field(&commit.RepoId, validation.Required),
	)
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
