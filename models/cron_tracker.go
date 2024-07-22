package models

import "time"

type CronTracker struct {
	RepoName  string     `json:"repoName"`
	OwnerName string     `json:"ownerName"`
	NextTick  time.Time  `json:"nextTick"`
	FromDate  *time.Time `json:"fromDate"`
	ToDate    *time.Time `json:"toDate"`
}
