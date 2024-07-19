package service

import "context"

type GitBeamService struct {
}

func NewGitBeamService() *GitBeamService {
	return &GitBeamService{}
}

func (g GitBeamService) getRepoByName(ctx context.Context, name string) error {
	panic("implement me")
}
