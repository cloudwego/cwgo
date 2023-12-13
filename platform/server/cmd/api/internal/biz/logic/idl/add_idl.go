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

package idl

import (
	"context"
	"net/url"
	"strings"

	"github.com/cloudwego/cwgo/platform/server/cmd/api/internal/biz/model/idl"
	"github.com/cloudwego/cwgo/platform/server/cmd/api/internal/svc"
	"github.com/cloudwego/cwgo/platform/server/shared/consts"
	"github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/agent"
	"github.com/cloudwego/cwgo/platform/server/shared/logger"
	"go.uber.org/zap"
)

const (
	successMsgAddIDL = "add idl successfully"
)

type AddIDLLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddIDLLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddIDLLogic {
	return &AddIDLLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddIDLLogic) AddIDL(req *idl.AddIDLReq) (res *idl.AddIDLRes) {
	urlParsed, err := url.Parse(req.MainIdlPath)
	if err != nil {
		return &idl.AddIDLRes{
			Code: consts.ErrNumParamMainIdlPath,
			Msg:  consts.ErrMsgParamMainIdlPath,
		}
	}
	if urlParsed.Scheme != "http" && urlParsed.Scheme != "https" {
		return &idl.AddIDLRes{
			Code: consts.ErrNumParamMainIdlPath,
			Msg:  consts.ErrMsgParamMainIdlPath,
		}
	}

	req.ServiceName = strings.Replace(req.ServiceName, "/", "_", -1)
	req.ServiceName = strings.Replace(req.ServiceName, ".", "_", -1)

	if req.ServiceRepositoryName == "" {
		req.ServiceName = req.ServiceRepositoryName
	}

	client, err := l.svcCtx.Manager.GetAgentClient()
	if err != nil {
		logger.Logger.Error(consts.ErrMsgRpcGetClient, zap.Error(err))
		return &idl.AddIDLRes{
			Code: consts.ErrNumRpcGetClient,
			Msg:  consts.ErrMsgRpcGetClient,
		}
	}

	rpcRes, err := client.AddIDL(l.ctx, &agent.AddIDLReq{
		RepositoryId:          req.RepositoryID,
		MainIdlPath:           req.MainIdlPath,
		ServiceName:           req.ServiceName,
		ServiceRepositoryName: req.ServiceRepositoryName,
	})
	if err != nil {
		logger.Logger.Error(consts.ErrMsgRpcConnectClient, zap.Error(err))
		return &idl.AddIDLRes{
			Code: consts.ErrNumRpcConnectClient,
			Msg:  consts.ErrMsgRpcConnectClient,
		}
	}
	if rpcRes.Code != 0 {
		return &idl.AddIDLRes{
			Code: rpcRes.Code,
			Msg:  rpcRes.Msg,
		}
	}

	return &idl.AddIDLRes{
		Code: 0,
		Msg:  successMsgAddIDL,
	}
}
