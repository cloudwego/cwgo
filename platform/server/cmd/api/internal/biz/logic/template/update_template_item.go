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
	"github.com/cloudwego/cwgo/platform/server/shared/log"
	"go.uber.org/zap"
)

const (
	successMsgUpdateTemplateItem = "update template item successfully"
)

type UpdateTemplateItemLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateTemplateItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTemplateItemLogic {
	return &UpdateTemplateItemLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateTemplateItemLogic) UpdateTemplateItem(req *template.UpdateTemplateItemReq) (res *template.UpdateTemplateItemRes) {
	client, err := l.svcCtx.Manager.GetAgentClient()
	if err != nil {
		log.Error(consts.ErrMsgRpcGetClient, zap.Error(err))
		return &template.UpdateTemplateItemRes{
			Code: consts.ErrNumRpcGetClient,
			Msg:  consts.ErrMsgRpcGetClient,
		}
	}

	rpcRes, err := client.UpdateTemplateItem(l.ctx, &agent.UpdateTemplateItemReq{
		Id:      req.Id,
		Name:    req.Name,
		Content: req.Content,
	})
	if err != nil {
		log.Error(consts.ErrMsgRpcConnectClient, zap.Error(err))
		return &template.UpdateTemplateItemRes{
			Code: consts.ErrNumRpcConnectClient,
			Msg:  consts.ErrMsgRpcConnectClient,
		}
	}
	if rpcRes.Code != 0 {
		return &template.UpdateTemplateItemRes{
			Code: rpcRes.Code,
			Msg:  rpcRes.Msg,
		}
	}

	return &template.UpdateTemplateItemRes{
		Code: 0,
		Msg:  successMsgUpdateTemplateItem,
	}
}
