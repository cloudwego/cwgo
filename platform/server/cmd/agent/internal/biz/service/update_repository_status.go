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
	"github.com/cloudwego/cwgo/platform/server/shared/utils"
)

type UpdateRepositoryStatusService struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
} // NewUpdateRepositoryStatusService new UpdateRepositoryStatusService
func NewUpdateRepositoryStatusService(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRepositoryStatusService {
	return &UpdateRepositoryStatusService{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Run create note info
func (s *UpdateRepositoryStatusService) Run(req *agent.UpdateRepositoryStatusReq) (resp *agent.UpdateRepositoryStatusRes, err error) {
	utils.ValidStatus(req.Status)
	repo, err := s.svcCtx.DaoManager.Repository.GetRepository(req.Id)
	if err != nil {
		resp.Code = 400
		resp.Msg = err.Error()
		return resp, err
	}

	if req.Status == consts.Active {
		err = s.svcCtx.RepoManager.AddClient(repo)
		if err != nil {
			resp.Code = 400
			resp.Msg = err.Error()
			return resp, err
		}
	} else if req.Status == consts.DisActive {
		s.svcCtx.RepoManager.DelClient(repo)
	}

	resp.Code = 0
	resp.Msg = "update status success"

	return resp, nil
}
