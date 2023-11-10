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
	agent "github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/agent"
	"net/http"
)

type GetRepositoriesService struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
} // NewGetRepositoriesService new GetRepositoriesService
func NewGetRepositoriesService(ctx context.Context, svcCtx *svc.ServiceContext) *GetRepositoriesService {
	return &GetRepositoriesService{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Run create note info
func (s *GetRepositoriesService) Run(req *agent.GetRepositoriesReq) (resp *agent.GetRepositoriesRes, err error) {
	repos, err := s.svcCtx.DaoManager.Repository.GetRepositoryList(s.ctx, req.Page, req.Limit, req.Order, req.OrderBy)
	if err != nil {
		return &agent.GetRepositoriesRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Data: nil,
		}, nil
	}

	return &agent.GetRepositoriesRes{
		Code: 0,
		Msg:  "get repositories successfully",
		Data: &agent.GetRepositoriesResData{
			Repositories: repos,
		},
	}, nil
}
