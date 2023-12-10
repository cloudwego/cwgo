package repository

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/cloudwego/cwgo/platform/server/shared/consts"
	"github.com/cloudwego/cwgo/platform/server/shared/dao"
	"github.com/cloudwego/cwgo/platform/server/shared/errx"
	"github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/model"
	"github.com/cloudwego/cwgo/platform/server/shared/logger"
	"github.com/patrickmn/go-cache"
	"go.uber.org/zap"
)

type Manager struct {
	daoManager *dao.Manager

	repositoryClientsCache    *cache.Cache
	repositoryApiClientsCache *cache.Cache

	sync.RWMutex
}

const (
	repositoryClientDefaultExpiration    = 24 * time.Hour
	repositoryApiClientDefaultExpiration = 24 * time.Hour
)

func NewRepoManager(daoManager *dao.Manager) (*Manager, error) {
	manager := &Manager{
		daoManager:                daoManager,
		repositoryClientsCache:    cache.New(repositoryClientDefaultExpiration, 1*time.Minute),
		repositoryApiClientsCache: cache.New(repositoryApiClientDefaultExpiration, 1*time.Minute),
		RWMutex:                   sync.RWMutex{},
	}

	return manager, nil
}

func (rm *Manager) AddClient(repositoryModel *model.Repository) (err error) {
	var repositoryClient IRepository

	if repositoryModel.TokenId != 0 {
		// if repo has existed token, then get the token info
		tokenModel, err := rm.daoManager.Token.GetTokenById(context.TODO(), repositoryModel.TokenId)
		if err == nil {
			switch repositoryModel.RepositoryType {
			case consts.RepositoryTypeNumGitLab:
				repositoryClient, err = NewGitLabApi(
					tokenModel.RepositoryDomain,
					tokenModel.Token,
					repositoryModel.RepositoryOwner,
					repositoryModel.RepositoryName,
				)
			case consts.RepositoryTypeNumGithub:
				repositoryClient, err = NewGitHubApi(
					tokenModel.Token,
					repositoryModel.RepositoryOwner,
					repositoryModel.RepositoryName,
				)
			default:
				return consts.ErrParamRepositoryType
			}
			if err != nil {
				if errx.GetCode(err) != consts.ErrNumTokenInvalid {
					return err
				}
			}
		}
	} else {
		// if repo's token is invalid or has no token
		// then search valid token in database
		tokenModels, err := rm.daoManager.Token.GetActiveTokenForDomain(context.TODO(), repositoryModel.RepositoryDomain)
		if err != nil {
			return err
		}

		waitChan := make(chan struct{})
		exitChan := make(chan struct {
			RepositoryClient IRepository
			TokenId          int64
		})

		for _, tokenModel := range tokenModels {
			go func(tokenModel *model.Token) {
				defer func() {
					waitChan <- struct{}{}
				}()

				internalRepositoryClient, err := NewGitLabApi(
					tokenModel.RepositoryDomain,
					tokenModel.Token,
					repositoryModel.RepositoryOwner,
					repositoryModel.RepositoryName,
				)
				if err != nil {
					logger.Logger.Error("init repo client failed", zap.Error(err))
					return
				}

				exitChan <- struct {
					RepositoryClient IRepository
					TokenId          int64
				}{RepositoryClient: internalRepositoryClient, TokenId: tokenModel.Id}
			}(tokenModel)
		}

		for i := 0; i < len(tokenModels); i++ {
			select {
			case <-waitChan:

			case chanRes := <-exitChan:
				repositoryClient = chanRes.RepositoryClient
				repositoryModel.TokenId = chanRes.TokenId
				break
			}
		}
	}

	if repositoryClient == nil {
		return consts.ErrTokenInvalid
	}

	if repositoryModel.RepositoryBranch == "" {
		// if branch is empty, then switch to default branch
		defaultBranch, err := repositoryClient.GetRepoDefaultBranch()
		if err != nil {
			return err
		}
		repositoryModel.RepositoryBranch = defaultBranch
	}

	rm.Lock()
	rm.repositoryClientsCache.SetDefault(strconv.FormatInt(repositoryModel.Id, 10), repositoryClient)
	rm.Unlock()

	return nil
}

func (rm *Manager) GetClient(repoId int64) (IRepository, error) {
	rm.RLock()
	if clientIface, ok := rm.repositoryClientsCache.Get(strconv.FormatInt(repoId, 10)); !ok {
		rm.RUnlock()
		repoModel, err := rm.daoManager.Repository.GetRepository(context.Background(), repoId)
		if err != nil {
			return nil, err
		}

		err = rm.AddClient(repoModel)
		if err != nil {
			if err == consts.ErrTokenInvalid {
				// exist token is invalid (expired maybe)
				// change repo status to inactive
				err = rm.daoManager.Repository.ChangeRepositoryStatus(context.Background(), repoModel.Id, consts.RepositoryStatusNumInactive)
				if err != nil {
					return nil, err
				}

				return nil, consts.ErrTokenInvalid
			}
			return nil, err
		}

		return rm.GetClient(repoId)
	} else {
		rm.RUnlock()
		return clientIface.(IRepository), nil
	}
}
