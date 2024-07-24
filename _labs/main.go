package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/go-github/v63/github"
)

func main() {
	ghClient := github.NewClient(nil)

	//repo, _, _ := ghClient.Repositories.Get(context.Background(), "brave", "brave-browser")
	//
	//prettyJson(repo)

	// Page 1: f9bbab669b6f33b228c6bbda8477d643ad58388e and d5dffa29e74bebe402cbf23c2b1c2ecf33e84971
	commit, _, _ := ghClient.Repositories.GetCommit(context.Background(), "chromium", "chromium", "a70fc91846eaa0da2db1de18b8f344b485eb7996", nil)

	author := commit.GetCommit().GetAuthor().GetName()
	parents := commit.Parents
	url := commit.GetHTMLURL()
	commitURL := commit.GetCommit().GetHTMLURL()

	fmt.Println("author: ", author)
	fmt.Println("html1: ", url)
	fmt.Println("html2: ", commitURL)
	for _, parent := range parents {
		fmt.Println("parent sha: ", parent.GetSHA())
	}
	prettyJson(commit)
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
