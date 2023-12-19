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
	successMsgGetRepositories = "get repositories successfully"
)

type GetRepositoriesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetRepositoriesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRepositoriesLogic {
	return &GetRepositoriesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRepositoriesLogic) GetRepositories(req *repository.GetRepositoriesReq) (res *repository.GetRepositoriesRes) {
	if req.Order != consts.OrderNumInc && req.Order != consts.OrderNumDec {
		return &repository.GetRepositoriesRes{
			Code: consts.ErrNumParamOrderNum,
			Msg:  consts.ErrMsgParamOrderNum,
			Data: nil,
		}
	}

	switch req.OrderBy {
	case "last_update_time", "last_sync_time", "create_time", "update_time", "":

	default:
		return &repository.GetRepositoriesRes{
			Code: consts.ErrNumParamOrderBy,
			Msg:  consts.ErrMsgParamOrderBy,
			Data: nil,
		}
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = consts.DefaultLimit
	}

	if req.RepositoryType < 0 || req.RepositoryType > consts.RepositoryTypeNum {
		return &repository.GetRepositoriesRes{
			Code: consts.ErrNumParamRepositoryType,
			Msg:  consts.ErrMsgParamRepositoryType,
			Data: nil,
		}
	}

	if req.StoreType < 0 || req.StoreType > consts.RepositoryStoreTypeNum {
		return &repository.GetRepositoriesRes{
			Code: consts.ErrNumParamStoreType,
			Msg:  consts.ErrMsgParamStoreType,
			Data: nil,
		}
	}

	client, err := l.svcCtx.Manager.GetAgentClient()
	if err != nil {
		logger.Logger.Error(consts.ErrMsgRpcGetClient, zap.Error(err))
		return &repository.GetRepositoriesRes{
			Code: consts.ErrNumRpcGetClient,
			Msg:  consts.ErrMsgRpcGetClient,
		}
	}

	rpcRes, err := client.GetRepositories(l.ctx, &agent.GetRepositoriesReq{
		Page:             req.Page,
		Limit:            req.Limit,
		Order:            req.Order,
		OrderBy:          req.OrderBy,
		RepositoryType:   req.RepositoryType,
		StoreType:        req.StoreType,
		RepositoryDomain: req.RepositoryDomain,
		RepositoryOwner:  req.RepositoryOwner,
		RepositoryName:   req.RepositoryName,
	})
	if err != nil {
		logger.Logger.Error(consts.ErrMsgRpcConnectClient, zap.Error(err))
		return &repository.GetRepositoriesRes{
			Code: consts.ErrNumRpcConnectClient,
			Msg:  consts.ErrMsgRpcConnectClient,
		}
	}
	if rpcRes.Code != 0 {
		return &repository.GetRepositoriesRes{
			Code: rpcRes.Code,
			Msg:  rpcRes.Msg,
		}
	}

	return &repository.GetRepositoriesRes{
		Code: 0,
		Msg:  successMsgGetRepositories,
		Data: &repository.GetRepositoriesResData{
			Repositories: rpcRes.Data.Repositories,
			Total:        rpcRes.Data.Total,
		},
	}
}
