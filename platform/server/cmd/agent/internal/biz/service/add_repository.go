/*
*
*  * Copyright 2022 CloudWeGo Authors
*  *
*  * Licensed under the Apache License, Version 2.0 (the "License");
*  * you may not use this file except in compliance with the License.
*  * You may obtain a copy of the License at
*  *
*  *     http://www.apache.org/licenses/LICENSE-2.0
*  *
*  * Unless required by applicable law or agreed to in writing, software
*  * distributed under the License is distributed on an "AS IS" BASIS,
*  * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
*  * See the License for the specific language governing permissions and
*  * limitations under the License.
*
 */

package service

import (
	"context"

	"github.com/cloudwego/cwgo/platform/server/cmd/agent/internal/svc"
	"github.com/cloudwego/cwgo/platform/server/shared/consts"
	"github.com/cloudwego/cwgo/platform/server/shared/errx"
	agent "github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/agent"
	"github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/model"
)

type AddRepositoryService struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
} // NewAddRepositoryService new AddRepositoryService
func NewAddRepositoryService(ctx context.Context, svcCtx *svc.ServiceContext) *AddRepositoryService {
	return &AddRepositoryService{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Run create note info
func (s *AddRepositoryService) Run(req *agent.AddRepositoryReq) (resp *agent.AddRepositoryRes, err error) {
	repo := model.Repository{
		RepositoryType:   req.RepositoryType,
		RepositoryDomain: req.RepositoryDomain,
		RepositoryOwner:  req.RepositoryOwner,
		RepositoryName:   req.RepositoryName,
		StoreType:        req.StoreType,
		RepositoryBranch: req.Branch,
	}

	// validate repo info add repo to memory
	err = s.svcCtx.RepoManager.AddClient(&repo)
	if err != nil {
		if errx.GetCode(err) == consts.ErrNumTokenInvalid {
			return &agent.AddRepositoryRes{
				Code: consts.ErrNumTokenInvalid,
				Msg:  err.Error(),
			}, nil
		}

		return &agent.AddRepositoryRes{
			Code: -1,
			Msg:  err.Error(),
		}, nil
	}

	// save repo info to db
	_, err = s.svcCtx.DaoManager.Repository.AddRepository(s.ctx, repo)
	if err != nil {
		if errx.GetCode(err) == consts.ErrNumDatabaseDuplicateRecord {
			return &agent.AddRepositoryRes{
				Code: consts.ErrNumDatabaseDuplicateRecord,
				Msg:  "repository is already exist",
			}, nil
		}

		return &agent.AddRepositoryRes{
			Code: consts.ErrNumDatabase,
			Msg:  consts.ErrMsgDatabase,
		}, nil
	}

	return &agent.AddRepositoryRes{
		Code: 0,
		Msg:  "add repository successfully",
	}, nil
}
