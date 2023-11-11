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
	agent "github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/agent"
	"github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/model"
	"github.com/cloudwego/cwgo/platform/server/shared/repository"
	"net/http"
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
		RepositoryType: req.RepositoryType,
		StoreType:      req.StoreType,
		RepositoryUrl:  req.RepositoryUrl,
		Token:          req.Token,
	}

	// validate repo info add repo to memory
	err = s.svcCtx.RepoManager.AddClient(&repo)
	if err != nil {
		if err == repository.ErrTokenInvalid {
			return &agent.AddRepositoryRes{
				Code: http.StatusBadRequest,
				Msg:  err.Error(),
			}, nil
		}

		return &agent.AddRepositoryRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
		}, nil
	}

	// save repo info to db
	_, err = s.svcCtx.DaoManager.Repository.AddRepository(s.ctx, repo)
	if err != nil {
		if err == consts.ErrDuplicateRecord {
			return &agent.AddRepositoryRes{
				Code: http.StatusBadRequest,
				Msg:  "repository is already exist",
			}, nil
		}

		return &agent.AddRepositoryRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
		}, nil
	}

	return &agent.AddRepositoryRes{
		Code: 0,
		Msg:  "add repository successfully",
	}, nil
}
