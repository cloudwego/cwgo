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

package token

import (
	"context"

	"github.com/cloudwego/cwgo/platform/server/cmd/api/internal/svc"
	"github.com/cloudwego/cwgo/platform/server/shared/consts"
	"github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/agent"
	token "github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/token"
	"github.com/cloudwego/cwgo/platform/server/shared/logger"
	"go.uber.org/zap"
)

const (
	successMsgGetToken = "get tokens successfully"
)

type GetTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTokenLogic {
	return &GetTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTokenLogic) GetToken(req *token.GetTokenReq) (res *token.GetTokenRes) {
	if req.Order != consts.OrderNumInc && req.Order != consts.OrderNumDec {
		return &token.GetTokenRes{
			Code: consts.ErrNumParamOrderNum,
			Msg:  consts.ErrMsgParamOrderNum,
			Data: nil,
		}
	}

	switch req.OrderBy {
	case "last_update_time", "last_sync_time", "create_time", "update_time", "":

	default:
		return &token.GetTokenRes{
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
		return &token.GetTokenRes{
			Code: consts.ErrNumParamRepositoryType,
			Msg:  consts.ErrMsgParamRepositoryType,
			Data: nil,
		}
	}

	client, err := l.svcCtx.Manager.GetAgentClient()
	if err != nil {
		logger.Logger.Error(consts.ErrMsgRpcGetClient, zap.Error(err))
		return &token.GetTokenRes{
			Code: consts.ErrNumRpcGetClient,
			Msg:  consts.ErrMsgRpcGetClient,
		}
	}

	rpcRes, err := client.GetToken(l.ctx, &agent.GetTokenReq{
		Page:             req.Page,
		Limit:            req.Limit,
		Order:            req.Order,
		OrderBy:          req.OrderBy,
		RepositoryType:   req.RepositoryType,
		RepositoryDomain: req.RepositoryDomain,
		Owner:            req.Owner,
	})
	if err != nil {
		logger.Logger.Error(consts.ErrMsgRpcConnectClient, zap.Error(err))
		return &token.GetTokenRes{
			Code: consts.ErrNumRpcConnectClient,
			Msg:  consts.ErrMsgRpcConnectClient,
		}
	}
	if rpcRes.Code != 0 {
		return &token.GetTokenRes{
			Code: rpcRes.Code,
			Msg:  rpcRes.Msg,
		}
	}

	return &token.GetTokenRes{
		Code: 0,
		Msg:  successMsgGetToken,
		Data: &token.GetTokenResData{
			Tokens: rpcRes.Data.Tokens,
			Total:  rpcRes.Data.Total,
		},
	}
}
