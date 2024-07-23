package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"go/build"
	"os"
	"path/filepath"
)

const ServiceName = "gitbeam"

type Secrets struct {
	CommitsMonitorURL string `json:"COMMITS_MONITOR_URL"`
	RepoManagerURL    string `json:"REPO_MANAGER_URL"`
	Port              string
}

var ss Secrets

func init() {
	importPath := fmt.Sprintf("%s/config", ServiceName)
	p, err := build.Default.Import(importPath, "", build.FindOnly)
	if err == nil {
		env := filepath.Join(p.Dir, "../.env")
		_ = godotenv.Load(env)
	}

	ss = Secrets{}
	ss.CommitsMonitorURL = os.Getenv("COMMITS_MONITOR_URL")
	ss.RepoManagerURL = os.Getenv("REPO_MANAGER_URL")
	if ss.Port = os.Getenv("PORT"); ss.Port == "" {
		ss.Port = "80"
	}
}

// GetSecrets is used to get value from the Secrets runtime.
func GetSecrets() Secrets {
	return ss
}
