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

package template

import (
	"context"

	"github.com/cloudwego/cwgo/platform/server/cmd/api/internal/svc"
	"github.com/cloudwego/cwgo/platform/server/shared/consts"
	"github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/agent"
	"github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/template"
	"github.com/cloudwego/cwgo/platform/server/shared/logger"
	"go.uber.org/zap"
)

const (
	successMsgGetTemplates = "get templates successfully"
)

type GetTemplatesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTemplatesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTemplatesLogic {
	return &GetTemplatesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTemplatesLogic) GetTemplates(req *template.GetTemplatesReq) (res *template.GetTemplatesRes) {
	if req.Order != consts.OrderNumInc && req.Order != consts.OrderNumDec {
		return &template.GetTemplatesRes{
			Code: consts.ErrNumParamOrderNum,
			Msg:  consts.ErrMsgParamOrderNum,
			Data: nil,
		}
	}

	switch req.OrderBy {
	case "create_time":

	case "update_time":

	default:
		return &template.GetTemplatesRes{
			Code: consts.ErrNumParamOrderBy,
			Msg:  consts.ErrMsgParamOrderBy,
			Data: nil,
		}
	}

	if req.Page == 0 {
		req.Page = 1
	}
	if req.Limit == 0 {
		req.Limit = consts.DefaultLimit
	}

	client, err := l.svcCtx.Manager.GetAgentClient()
	if err != nil {
		logger.Logger.Error(consts.ErrMsgRpcGetClient, zap.Error(err))
		return &template.GetTemplatesRes{
			Code: consts.ErrNumRpcGetClient,
			Msg:  consts.ErrMsgRpcGetClient,
		}
	}

	rpcRes, err := client.GetTemplates(l.ctx, &agent.GetTemplatesReq{
		Page:    req.Page,
		Limit:   req.Limit,
		Order:   req.Order,
		OrderBy: req.OrderBy,
	})
	if err != nil {
		logger.Logger.Error(consts.ErrMsgRpcConnectClient, zap.Error(err))
		return &template.GetTemplatesRes{
			Code: consts.ErrNumRpcConnectClient,
			Msg:  consts.ErrMsgRpcConnectClient,
		}
	}
	if rpcRes.Code != 0 {
		return &template.GetTemplatesRes{
			Code: rpcRes.Code,
			Msg:  rpcRes.Msg,
		}
	}

	return &template.GetTemplatesRes{
		Code: 0,
		Msg:  successMsgGetTemplates,
		Data: &template.GetTemplatesResData{Templates: rpcRes.Data.Templates},
	}
}
