package repository

import (
	"context"
	"errors"
	"sync"

	"github.com/cloudwego/cwgo/platform/server/shared/consts"
	"github.com/cloudwego/cwgo/platform/server/shared/dao"
	"github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/repository"
	"github.com/cloudwego/cwgo/platform/server/shared/utils"
	"github.com/google/go-github/v53/github"
	"github.com/xanzy/go-gitlab"
	"golang.org/x/oauth2"
)

type Manager struct {
	daoManager *dao.Manager

	repositoryClients map[int64]IRepository

	sync.RWMutex
}

func NewRepoManager(daoManager *dao.Manager) (*Manager, error) {
	repositoryClients := make(map[int64]IRepository)

	manager := &Manager{
		daoManager:        daoManager,
		repositoryClients: repositoryClients,
		RWMutex:           sync.RWMutex{},
	}

	return manager, nil
}

func (rm *Manager) ClearClient() {
	rm.repositoryClients = make(map[int64]IRepository)
}

func (rm *Manager) AddClient(repository *repository.Repository) error {
	rm.Lock()
	defer rm.Unlock()

	switch repository.RepositoryType {
	case consts.RepositoryTypeNumGitLab:
		gitlabClient, err := NewGitlabClient(repository.Token)
		if err != nil {
			if utils.IsTokenError(err) {
				err = rm.daoManager.Repository.ChangeRepositoryStatus(repository.Id, consts.InvalidToken)
				if err != nil {
					return err
				}
			} else if utils.IsNetworkError(err) {
				return errors.New("client initialization request timeout")
			} else {
				return errors.New("client initialization request unknown error")
			}
		}

		rm.repositoryClients[repository.Id] = NewGitLabApi(gitlabClient)
	case consts.RepositoryTypeNumGithub:
		githubClient := NewGithubClient(repository.Token)

		rm.repositoryClients[repository.Id] = NewGitHubApi(githubClient)
	default:
		return errors.New("invalid repository type")
	}

	return nil
}

func (rm *Manager) DelClient(repository *repository.Repository) {
	rm.Lock()
	defer rm.Unlock()

	delete(rm.repositoryClients, repository.Id)
}

func (rm *Manager) GetClient(repoId int64) (IRepository, error) {
	rm.RLock()
	defer rm.RUnlock()

	if client, ok := rm.repositoryClients[repoId]; !ok {
		repo, err := rm.daoManager.Repository.GetRepository(repoId)
		if err != nil {
			return nil, err
		}
		err = rm.AddClient(repo)
		if err != nil {
			return nil, err
		}
		return rm.GetClient(repoId)
	} else {
		return client, nil
	}
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
