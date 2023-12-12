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
	"github.com/cloudwego/cwgo/platform/server/cmd/api/internal/biz/model/idl"
	"github.com/cloudwego/cwgo/platform/server/cmd/api/internal/svc"
	"github.com/cloudwego/cwgo/platform/server/shared/consts"
	"github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/agent"
	"github.com/cloudwego/cwgo/platform/server/shared/logger"
	"go.uber.org/zap"
)

const (
	successMsgUpdateIDL = "update idl successfully"
)

type UpdateIDLLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateIDLLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateIDLLogic {
	return &UpdateIDLLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateIDLLogic) UpdateIDL(req *idl.UpdateIDLReq) (res *idl.UpdateIDLRes) {
	if req.Status != 0 {
		if _, ok := consts.IdlStatusNumMap[int(req.Status)]; !ok {
			return &idl.UpdateIDLRes{
				Code: consts.ErrNumParamIdlStatus,
				Msg:  consts.ErrMsgParamIdlStatus,
			}
		}
	}

	client, err := l.svcCtx.Manager.GetAgentClient()
	if err != nil {
		logger.Logger.Error(consts.ErrMsgRpcGetClient, zap.Error(err))
		return &idl.UpdateIDLRes{
			Code: consts.ErrNumRpcGetClient,
			Msg:  consts.ErrMsgRpcGetClient,
		}
	}

	rpcRes, err := client.UpdateIDL(l.ctx, &agent.UpdateIDLReq{
		RepositoryId: req.ID,
		ServiceName:  req.ServiceName,
		Status:       req.Status,
	})
	if err != nil {
		logger.Logger.Error(consts.ErrMsgRpcConnectClient, zap.Error(err))
		return &idl.UpdateIDLRes{
			Code: consts.ErrNumRpcConnectClient,
			Msg:  consts.ErrMsgRpcConnectClient,
		}
	}
	if rpcRes.Code != 0 {
		return &idl.UpdateIDLRes{
			Code: rpcRes.Code,
			Msg:  rpcRes.Msg,
		}
	}

	return &idl.UpdateIDLRes{
		Code: 0,
		Msg:  successMsgUpdateIDL,
	}
}
