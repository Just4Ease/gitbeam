package models

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type Result struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

type OwnerAndRepoName struct {
	OwnerName string `json:"ownerName" schema:"ownerName"`
	RepoName  string `json:"repoName" schema:"repoName"`
}

type MirrorRepoCommitsRequest struct {
	OwnerAndRepoName `json:",inline"`
	FromDate         *Date `json:"fromDate,omitempty"`
	ToDate           *Date `json:"toDate,omitempty"`
}

func (s OwnerAndRepoName) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.OwnerName, validation.Required),
		validation.Field(&s.RepoName, validation.Required))
}
