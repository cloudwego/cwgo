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
	successMsgSyncRepository = "sync repository successfully"
)

type SyncRepositoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSyncRepositoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SyncRepositoryLogic {
	return &SyncRepositoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SyncRepositoryLogic) SyncRepository(req *repository.SyncRepositoryByIdReq) (res *repository.SyncRepositoryByIdRes) {
	// TODO: get certain agent by repo id
	client, err := l.svcCtx.Manager.GetAgentClient()
	if err != nil {
		logger.Logger.Error(consts.ErrMsgRpcGetClient, zap.Error(err))
		return &repository.SyncRepositoryByIdRes{
			Code: consts.ErrNumRpcGetClient,
			Msg:  consts.ErrMsgRpcGetClient,
		}
	}

	rpcRes, err := client.SyncRepositoryById(l.ctx, &agent.SyncRepositoryByIdReq{
		Ids: req.Ids,
	})
	if err != nil {
		logger.Logger.Error(consts.ErrMsgRpcConnectClient, zap.Error(err))
		return &repository.SyncRepositoryByIdRes{
			Code: consts.ErrNumRpcConnectClient,
			Msg:  consts.ErrMsgRpcConnectClient,
		}
	}
	if rpcRes.Code != 0 {
		return &repository.SyncRepositoryByIdRes{
			Code: rpcRes.Code,
			Msg:  rpcRes.Msg,
		}
	}

	return &repository.SyncRepositoryByIdRes{
		Code: 0,
		Msg:  successMsgSyncRepository,
	}
}
