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
	successMsgGetIDLs = "get idls successfully"
)

type GetIDLsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetIDLsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetIDLsLogic {
	return &GetIDLsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetIDLsLogic) GetIDLs(req *idl.GetIDLsReq) (res *idl.GetIDLsRes) {
	if req.Order != consts.OrderNumInc && req.Order != consts.OrderNumDec {
		return &idl.GetIDLsRes{
			Code: consts.ErrNumParamOrderNum,
			Msg:  consts.ErrMsgParamOrderNum,
			Data: nil,
		}
	}

	switch req.OrderBy {
	case "last_sync_time", "create_time", "update_time", "":

	default:
		return &idl.GetIDLsRes{
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

	client, err := l.svcCtx.Manager.GetAgentClient()
	if err != nil {
		logger.Logger.Error(consts.ErrMsgRpcGetClient, zap.Error(err))
		return &idl.GetIDLsRes{
			Code: consts.ErrNumRpcGetClient,
			Msg:  consts.ErrMsgRpcGetClient,
		}
	}

	rpcRes, err := client.GetIDLs(l.ctx, &agent.GetIDLsReq{
		Page:        req.Page,
		Limit:       req.Limit,
		Order:       req.Order,
		OrderBy:     req.OrderBy,
		ServiceName: req.ServiceName,
	})
	if err != nil {
		logger.Logger.Error(consts.ErrMsgRpcConnectClient, zap.Error(err))
		return &idl.GetIDLsRes{
			Code: consts.ErrNumRpcConnectClient,
			Msg:  consts.ErrMsgRpcConnectClient,
		}
	}
	if rpcRes.Code != 0 {
		return &idl.GetIDLsRes{
			Code: rpcRes.Code,
			Msg:  rpcRes.Msg,
		}
	}

	return &idl.GetIDLsRes{
		Code: 0,
		Msg:  successMsgGetIDLs,
		Data: &idl.GetIDLsResData{
			Idls:  rpcRes.Data.Idls,
			Total: rpcRes.Data.Total,
		},
	}
}
