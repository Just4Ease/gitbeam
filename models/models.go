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
	Id          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	URL         string         `json:"url"`
	ForkCounts  map[string]int `json:"forkCounts"`
	Language    string         `json:"language"`
	Attributes  map[string]any `json:"attributes,omitempty"`
}

func (r Repo) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required),
		validation.Field(&r.URL, validation.Required),
	)
}
