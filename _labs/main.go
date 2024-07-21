package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/go-github/v63/github"
	"time"
)

func main() {
	ghClient := github.NewClient(nil)

	commits, _, _ := ghClient.Repositories.ListCommits(context.Background(), "brave", "brave-browser", &github.CommitsListOptions{
		Since: time.Time{},
		Until: time.Time{},
		ListOptions: github.ListOptions{
			PerPage: 2,
			Page:    1,
		},
	})
	prettyJson(commits)
}

const (
	empty = ""
	tab   = "\t"
)

func prettyJson(data interface{}) {
	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)
	encoder.SetIndent(empty, tab)

	err := encoder.Encode(data)
	if err != nil {
		return
	}
	fmt.Print(buffer.String())
}
