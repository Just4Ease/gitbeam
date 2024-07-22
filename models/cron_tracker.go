package models

import "time"

type CronTask struct {
	RepoName  string     `json:"repoName"`
	OwnerName string     `json:"ownerName"`
	FromDate  *time.Time `json:"fromDate"`
	ToDate    *time.Time `json:"toDate"`
}
