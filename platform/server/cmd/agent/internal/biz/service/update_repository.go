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

package service

import (
	"context"

	"github.com/cloudwego/cwgo/platform/server/cmd/agent/internal/svc"
	"github.com/cloudwego/cwgo/platform/server/shared/consts"
	"github.com/cloudwego/cwgo/platform/server/shared/errx"
	agent "github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/agent"
	"github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/model"
)

type UpdateRepositoryService struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
} // NewUpdateRepositoryService new UpdateRepositoryService
func NewUpdateRepositoryService(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRepositoryService {
	return &UpdateRepositoryService{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Run create note info
func (s *UpdateRepositoryService) Run(req *agent.UpdateRepositoryReq) (resp *agent.UpdateRepositoryResp, err error) {
	// validate repo info
	repoModel, err := s.svcCtx.DaoManager.Repository.GetRepository(s.ctx, req.Id)
	if err != nil {
		if errx.GetCode(err) == consts.ErrNumDatabaseRecordNotFound {
			return &agent.UpdateRepositoryResp{
				Code: consts.ErrNumDatabaseRecordNotFound,
				Msg:  "repo id not exist",
			}, nil
		}
	}

	if req.Branch == repoModel.RepositoryBranch {
		return &agent.UpdateRepositoryResp{
			Code: consts.ErrNumParamRepositoryBranch,
			Msg:  "repo branch already switched",
		}, nil
	}

	if req.Branch != "" {
		repoClient, err := s.svcCtx.RepoManager.GetClient(req.Id)
		if err != nil {
			return &agent.UpdateRepositoryResp{
				Code: errx.GetCode(err),
				Msg:  err.Error(),
			}, nil
		}

		isValid, err := repoClient.ValidateRepoBranch(req.Branch)
		if err != nil {
			return &agent.UpdateRepositoryResp{
				Code: consts.ErrNumRepoValidateBranch,
				Msg:  consts.ErrMsgRepoValidateBranch,
			}, nil
		}

		if !isValid {
			return &agent.UpdateRepositoryResp{
				Code: consts.ErrNumParamRepositoryBranch,
				Msg:  consts.ErrMsgParamRepositoryBranch,
			}, nil
		}
	}

	// update repo info
	err = s.svcCtx.DaoManager.Repository.UpdateRepository(s.ctx, model.Repository{
		Id:               req.Id,
		RepositoryBranch: req.Branch,
		Status:           req.Status,
	})
	if err != nil {
		if errx.GetCode(err) == consts.ErrNumDatabaseRecordNotFound {
			return &agent.UpdateRepositoryResp{
				Code: consts.ErrNumDatabaseRecordNotFound,
				Msg:  "repo id not exist",
			}, nil
		}
		return &agent.UpdateRepositoryResp{
			Code: consts.ErrNumDatabase,
			Msg:  consts.ErrMsgDatabase,
		}, nil
	}

	if req.Status == consts.RepositoryStatusNumInactive {
		s.svcCtx.RepoManager.DelClient(req.Id)
	}
	if req.Branch != "" {
		client, err := s.svcCtx.RepoManager.GetClient(req.Id)
		if err != nil {
			return &agent.UpdateRepositoryResp{
				Code: errx.GetCode(err),
				Msg:  err.Error(),
			}, nil
		}

		client.UpdateBranch(req.Branch)
	}

	return &agent.UpdateRepositoryResp{
		Code: 0,
		Msg:  "update repository successfully",
	}, nil
}
