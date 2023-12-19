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

	"github.com/cloudwego/cwgo/platform/server/cmd/api/internal/biz/model/repository"
	"github.com/cloudwego/cwgo/platform/server/cmd/api/internal/svc"
	"github.com/cloudwego/cwgo/platform/server/shared/consts"
	"github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/agent"
	"github.com/cloudwego/cwgo/platform/server/shared/logger"
	"go.uber.org/zap"
)

const (
	successMsgDeleteRepository = "delete repository successfully"
)

type DeleteRepositoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteRepositoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteRepositoryLogic {
	return &DeleteRepositoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteRepositoryLogic) DeleteRepository(req *repository.DeleteRepositoriesReq) (res *repository.DeleteRepositoriesRes) {
	client, err := l.svcCtx.Manager.GetAgentClient()
	if err != nil {
		logger.Logger.Error(consts.ErrMsgRpcGetClient, zap.Error(err))
		return &repository.DeleteRepositoriesRes{
			Code: consts.ErrNumRpcGetClient,
			Msg:  consts.ErrMsgRpcGetClient,
		}
	}

	rpcRes, err := client.DeleteRepositories(l.ctx, &agent.DeleteRepositoriesReq{
		Ids: req.Ids,
	})
	if err != nil {
		logger.Logger.Error(consts.ErrMsgRpcConnectClient, zap.Error(err))
		return &repository.DeleteRepositoriesRes{
			Code: consts.ErrNumRpcConnectClient,
			Msg:  consts.ErrMsgRpcConnectClient,
		}
	}
	if rpcRes.Code != 0 {
		return &repository.DeleteRepositoriesRes{
			Code: rpcRes.Code,
			Msg:  rpcRes.Msg,
		}
	}

	return &repository.DeleteRepositoriesRes{
		Code: 0,
		Msg:  successMsgDeleteRepository,
	}
}
