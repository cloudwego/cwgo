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
	"net"

	token "github.com/cloudwego/cwgo/platform/server/cmd/api/internal/biz/model/token"
	"github.com/cloudwego/cwgo/platform/server/cmd/api/internal/svc"
	"github.com/cloudwego/cwgo/platform/server/shared/consts"
	"github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/agent"
	"github.com/cloudwego/cwgo/platform/server/shared/logger"
	"go.uber.org/zap"
)

const (
	successMsgAddToken = "add token successfully"
)

type AddTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddTokenLogic {
	return &AddTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddTokenLogic) AddToken(req *token.AddTokenReq) (res *token.AddTokenRes) {
	if req.RepositoryType != consts.RepositoryTypeNumGithub {
		_, err := net.LookupHost(req.RepositoryDomain)
		if err != nil {
			return &token.AddTokenRes{
				Code: consts.ErrNumParamDomain,
				Msg:  consts.ErrMsgParamDomain,
				Data: nil,
			}
		}
	} else {
		req.RepositoryDomain = "github.com"
	}

	client, err := l.svcCtx.Manager.GetAgentClient()
	if err != nil {
		logger.Logger.Error(consts.ErrMsgRpcGetClient, zap.Error(err))
		return &token.AddTokenRes{
			Code: consts.ErrNumRpcGetClient,
			Msg:  consts.ErrMsgRpcGetClient,
		}
	}

	rpcRes, err := client.AddToken(l.ctx, &agent.AddTokenReq{
		RepositoryType:   req.RepositoryType,
		RepositoryDomain: req.RepositoryDomain,
		Token:            req.Token,
	})
	if err != nil {
		logger.Logger.Error(consts.ErrMsgRpcConnectClient, zap.Error(err))
		return &token.AddTokenRes{
			Code: consts.ErrNumRpcConnectClient,
			Msg:  consts.ErrMsgRpcConnectClient,
		}
	}
	if rpcRes.Code != 0 {
		return &token.AddTokenRes{
			Code: rpcRes.Code,
			Msg:  rpcRes.Msg,
		}
	}

	return &token.AddTokenRes{
		Code: 0,
		Msg:  successMsgAddToken,
		Data: &token.AddTokenResData{
			Owner:          rpcRes.Data.Owner,
			ExpirationTime: rpcRes.Data.ExpirationTime,
		},
	}
}
