/*
*
 * Copyright 2023 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
*
*/

package repository

import (
	"context"
	"strconv"
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

	repositoryClientsCache *cache.Cache
}

const (
	repositoryClientDefaultExpiration = 24 * time.Hour
)

func NewRepoManager(daoManager *dao.Manager) (*Manager, error) {
	manager := &Manager{
		daoManager:             daoManager,
		repositoryClientsCache: cache.New(repositoryClientDefaultExpiration, 1*time.Minute),
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
					repositoryModel.RepositoryBranch,
				)
			case consts.RepositoryTypeNumGithub:
				repositoryClient, err = NewGitHubApi(
					tokenModel.Token,
					repositoryModel.RepositoryOwner,
					repositoryModel.RepositoryName,
					repositoryModel.RepositoryBranch,
				)
			default:
				return consts.ErrParamRepositoryType
			}
			if err != nil {
				if errx.GetCode(err) != consts.ErrNumTokenInvalid {
					return err
				}
				repositoryClient = nil
			}
		}
	}
	if repositoryClient == nil {
		// if repo's token is invalid or has no token
		// then search valid token in database
		tokenModels, err := rm.daoManager.Token.GetActiveTokenForDomain(context.TODO(), repositoryModel.RepositoryDomain)
		if err != nil {
			return err
		}

		// use token which it's owner is same with repo owner
		for _, tokenModel := range tokenModels {
			if tokenModel.Owner == repositoryModel.RepositoryOwner {
				switch repositoryModel.RepositoryType {
				case consts.RepositoryTypeNumGitLab:
					repositoryClient, err = NewGitLabApi(
						tokenModel.RepositoryDomain,
						tokenModel.Token,
						repositoryModel.RepositoryOwner,
						repositoryModel.RepositoryName,
						repositoryModel.RepositoryBranch,
					)
				case consts.RepositoryTypeNumGithub:
					repositoryClient, err = NewGitHubApi(
						tokenModel.Token,
						repositoryModel.RepositoryOwner,
						repositoryModel.RepositoryName,
						repositoryModel.RepositoryBranch,
					)
				default:
					return consts.ErrParamRepositoryType
				}
				if err != nil {
					if errx.GetCode(err) != consts.ErrNumTokenInvalid {
						return err
					}
					repositoryClient = nil
				} else {
					repositoryModel.TokenId = tokenModel.Id
					break
				}
			}
		}

		if repositoryClient == nil {
			// if there is no token corresponding to the repo owner
			// then check all token corresponding to the repo domain
			waitChan := make(chan struct{})
			exitChan := make(chan struct {
				RepositoryClient IRepository
				TokenId          int64
			})

			goNum := 0
			for _, tokenModel := range tokenModels {
				if tokenModel.Owner == repositoryModel.RepositoryOwner {
					continue
				}

				goNum++
				go func(tokenModel *model.Token) {
					defer func() {
						waitChan <- struct{}{}
					}()

					internalRepositoryClient, err := NewGitLabApi(
						tokenModel.RepositoryDomain,
						tokenModel.Token,
						repositoryModel.RepositoryOwner,
						repositoryModel.RepositoryName,
						repositoryModel.RepositoryBranch,
					)
					if err != nil {
						return
					}

					logger.Logger.Debug("get token for repo",
						zap.Int64("repo_id", repositoryModel.Id),
						zap.Int64("token_id", tokenModel.Id),
					)

					exitChan <- struct {
						RepositoryClient IRepository
						TokenId          int64
					}{RepositoryClient: internalRepositoryClient, TokenId: tokenModel.Id}
				}(tokenModel)
			}

			for i := 0; i < goNum; i++ {
				select {
				case <-waitChan:

				case chanRes := <-exitChan:
					repositoryClient = chanRes.RepositoryClient
					repositoryModel.TokenId = chanRes.TokenId
					break
				}
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

	rm.repositoryClientsCache.SetDefault(strconv.FormatInt(repositoryModel.Id, 10), repositoryClient)

	return nil
}

func (rm *Manager) DelClient(repoId int64) {
	rm.repositoryClientsCache.Delete(strconv.FormatInt(repoId, 10))
}

func (rm *Manager) GetClient(repoId int64) (IRepository, error) {
	if clientIface, ok := rm.repositoryClientsCache.Get(strconv.FormatInt(repoId, 10)); !ok {
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
		return clientIface.(IRepository), nil
	}
}
