package repository

import (
	"context"
	"github.com/google/go-github/v53/github"
	"github.com/xanzy/go-gitlab"
	"golang.org/x/oauth2"
)

type Manager struct {
	GitHub GitHubApi
	Gitlab GitLabApi
}

func NewRepoManager() (*Manager, error) {
	//查找
	token := ""
	NewGitlabClient(token)
	return &Manager{
		GitHub: GitHubApi{},
		Gitlab: GitLabApi{},
	}, nil
}

func NewGitlabClient(token string) (*GitLabApi, error) {
	client, err := gitlab.NewClient(token)
	if err != nil {
		return nil, err
	}

	return &GitLabApi{
		Client: client,
	}, nil
}

func NewGithubClient(token string) *GitHubApi {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	return &GitHubApi{
		Client: github.NewClient(tc),
	}
}
