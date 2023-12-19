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
)

type DeleteRepositoriesService struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
} // NewDeleteRepositoriesService new DeleteRepositoriesService
func NewDeleteRepositoriesService(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteRepositoriesService {
	return &DeleteRepositoriesService{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Run create note info
func (s *DeleteRepositoriesService) Run(req *agent.DeleteRepositoriesReq) (resp *agent.DeleteRepositoriesRes, err error) {
	err = s.svcCtx.DaoManager.Repository.DeleteRepository(s.ctx, req.Ids)
	if err != nil {
		if errx.GetCode(err) == consts.ErrNumDatabaseRecordNotFound {
			return &agent.DeleteRepositoriesRes{
				Code: consts.ErrNumDatabaseRecordNotFound,
				Msg:  "repo id not exist",
			}, nil
		}

		return &agent.DeleteRepositoriesRes{
			Code: consts.ErrNumDatabase,
			Msg:  consts.ErrMsgDatabase,
		}, nil
	}

	for _, id := range req.Ids {
		s.svcCtx.RepoManager.DelClient(id)
	}

	return &agent.DeleteRepositoriesRes{
		Code: 0,
		Msg:  "delete repositories successfully",
	}, nil
}
