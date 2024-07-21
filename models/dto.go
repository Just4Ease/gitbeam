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
	OwnerName string `json:"ownerName"`
	RepoName  string `json:"repoName"`
}

type BeamRepoCommitsRequest struct {
	OwnerAndRepoName `json:",inline"`
	StartTime        *string `json:"startTime,omitempty"`
}

func (s OwnerAndRepoName) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.OwnerName, validation.Required),
		validation.Field(&s.RepoName, validation.Required))
}
