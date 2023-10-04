package repository

import (
	"context"
	"errors"
	"sync"

	"github.com/cloudwego/cwgo/platform/server/shared/consts"
	"github.com/cloudwego/cwgo/platform/server/shared/dao"
	"github.com/cloudwego/cwgo/platform/server/shared/utils"
	"github.com/google/go-github/v53/github"
	"github.com/xanzy/go-gitlab"
	"golang.org/x/oauth2"
)

type Manager struct {
	GitHub   *GitHubApi
	GitLab   *GitLabApi
	MuGitHub sync.Mutex
	MuGitlab sync.Mutex
}

func NewRepoManager(DaoManager *dao.Manager) (*Manager, error) {
	manager := &Manager{
		GitHub: &GitHubApi{make(map[int64]*github.Client)},
		GitLab: &GitLabApi{make(map[int64]*gitlab.Client)},
	}

	repos, err := DaoManager.Repository.GetAllRepositories()
	if err != nil {
		return nil, err
	}

	for _, r := range repos {
		switch r.Type {
		case consts.GitLab:
			manager.GitLab.Client[r.Id], err = NewGitlabClient(r.Token)
			if err != nil {
				if utils.IsTokenError(err) {
					err = DaoManager.Repository.ChangeRepositoryStatus(r.Id, consts.InvalidToken)
					if err != nil {
						return nil, err
					}
				} else if utils.IsNetworkError(err) {
					return nil, errors.New("client initialization request timeout")
				} else {
					return nil, errors.New("client initialization request unknown error")
				}
			}

		case consts.Github:
			manager.GitHub.Client[r.Id] = NewGithubClient(r.Token)
		}
	}

	return manager, nil
}

func NewGitlabClient(token string) (*gitlab.Client, error) {
	client, err := gitlab.NewClient(token)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func NewGithubClient(token string) *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc)
}
