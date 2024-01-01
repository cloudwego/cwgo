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
	successMsgAddTemplate = "add template successfully"
)

type AddTemplateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddTemplateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddTemplateLogic {
	return &AddTemplateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddTemplateLogic) AddTemplate(req *template.AddTemplateReq) (res *template.AddTemplateRes) {
	if _, ok := consts.TemplateTypeNumMap[int(req.Type)]; !ok {
		return &template.AddTemplateRes{
			Code: consts.ErrNumParamTemplateType,
			Msg:  consts.ErrMsgParamTemplateType,
		}
	}

	client, err := l.svcCtx.Manager.GetAgentClient()
	if err != nil {
		logger.Logger.Error(consts.ErrMsgRpcGetClient, zap.Error(err))
		return &template.AddTemplateRes{
			Code: consts.ErrNumRpcGetClient,
			Msg:  consts.ErrMsgRpcGetClient,
		}
	}

	rpcRes, err := client.AddTemplate(l.ctx, &agent.AddTemplateReq{
		Name: req.Name,
		Type: req.Type,
	})
	if err != nil {
		logger.Logger.Error(consts.ErrMsgRpcConnectClient, zap.Error(err))
		return &template.AddTemplateRes{
			Code: consts.ErrNumRpcConnectClient,
			Msg:  consts.ErrMsgRpcConnectClient,
		}
	}
	if rpcRes.Code != 0 {
		return &template.AddTemplateRes{
			Code: rpcRes.Code,
			Msg:  rpcRes.Msg,
		}
	}

	return &template.AddTemplateRes{
		Code: 0,
		Msg:  successMsgAddTemplate,
	}
}
