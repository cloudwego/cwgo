package repository

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cloudwego/cwgo/platform/server/shared/consts"
	"github.com/cloudwego/cwgo/platform/server/shared/dao"
	"github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/model"
	"github.com/cloudwego/cwgo/platform/server/shared/utils"
	"github.com/google/go-github/v56/github"
	"github.com/patrickmn/go-cache"
	"github.com/xanzy/go-gitlab"
)

var (
	ErrTokenInvalid = errors.New("token is invalid")
)

type Manager struct {
	daoManager *dao.Manager

	repositoryClients      map[int64]IRepository
	repositoryClientsCache *cache.Cache

	sync.RWMutex
}

const (
	repositoryClientDefaultExpiration = 24 * time.Hour
)

func NewRepoManager(daoManager *dao.Manager) (*Manager, error) {
	repositoryClients := make(map[int64]IRepository)

	manager := &Manager{
		daoManager:             daoManager,
		repositoryClients:      repositoryClients,
		repositoryClientsCache: cache.New(repositoryClientDefaultExpiration, 1*time.Minute),
		RWMutex:                sync.RWMutex{},
	}

	return manager, nil
}

func (rm *Manager) AddClient(repository *model.Repository) error {

	switch repository.RepositoryType {
	case consts.RepositoryTypeNumGitLab:
		gitlabClient, err := NewGitlabClient(repository.Token)
		if err != nil {
			if utils.IsNetworkError(err) {
				return errors.New("client initialization request timeout")
			} else {
				return errors.New("client initialization request unknown error")
			}
		}

		rm.Lock()
		rm.repositoryClientsCache.SetDefault(strconv.FormatInt(repository.Id, 10), NewGitLabApi(gitlabClient))
		rm.Unlock()
	case consts.RepositoryTypeNumGithub:
		githubClient, err := NewGithubClient(repository.Token)
		if err != nil {
			return err
		}

		rm.Lock()
		rm.repositoryClientsCache.SetDefault(strconv.FormatInt(repository.Id, 10), NewGitHubApi(githubClient))
		rm.Unlock()
	default:
		return errors.New("invalid repository type")
	}

	return nil
}

func (rm *Manager) DelClient(repository *model.Repository) {
	rm.Lock()
	defer rm.Unlock()

	delete(rm.repositoryClients, repository.Id)
}

func (rm *Manager) GetClient(repoId int64) (IRepository, error) {
	rm.RLock()
	if clientIface, ok := rm.repositoryClientsCache.Get(strconv.FormatInt(repoId, 10)); !ok {
		rm.RUnlock()
		repo, err := rm.daoManager.Repository.GetRepository(context.Background(), repoId)
		if err != nil {
			return nil, err
		}

		err = rm.AddClient(repo)
		if err != nil {
			if err == ErrTokenInvalid {
				// exist token is invalid (expired maybe)
				// change repo status to inactive
				err = rm.daoManager.Repository.ChangeRepositoryStatus(context.Background(), repo.Id, consts.RepositoryStatusNumInactive)
				if err != nil {
					return nil, err
				}
			}
			return nil, err
		}

		return rm.GetClient(repoId)
	} else {
		rm.RUnlock()
		return clientIface.(IRepository), nil
	}
}

func NewGitlabClient(token string) (*gitlab.Client, error) {
	var client *gitlab.Client
	var err error

	if consts.ProxyUrl != "" {
		proxyUrl, _ := url.Parse(consts.ProxyUrl)
		httpClient := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
		client, err = gitlab.NewClient(token, gitlab.WithHTTPClient(httpClient))
	} else {
		client, err = gitlab.NewClient(token)
	}

	if err != nil {
		if strings.Contains(err.Error(), "401 Unauthorized") {
			return nil, ErrTokenInvalid
		}

		return nil, err
	}

	return client, nil
}

func NewGithubClient(token string) (*github.Client, error) {
	var httpClient *http.Client

	if consts.ProxyUrl != "" {
		proxyUrl, _ := url.Parse(consts.ProxyUrl)
		httpClient = &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
	}

	client := github.NewClient(httpClient).WithAuthToken(token)

	_, _, err := client.Meta.Get(context.Background())
	if err != nil {
		if githubErr, ok := err.(*github.ErrorResponse); ok {
			if githubErr.Message == "Bad credentials" {
				return nil, ErrTokenInvalid
			}
		}

		return nil, err
	}

	return client, nil
}
