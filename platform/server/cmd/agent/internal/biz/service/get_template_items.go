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

package service

import (
	"context"

	"github.com/cloudwego/cwgo/platform/server/cmd/agent/internal/svc"
	"github.com/cloudwego/cwgo/platform/server/shared/consts"
	agent "github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/agent"
)

type GetTemplateItemsService struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
} // NewGetTemplateItemsService new GetTemplateItemsService
func NewGetTemplateItemsService(ctx context.Context, svcCtx *svc.ServiceContext) *GetTemplateItemsService {
	return &GetTemplateItemsService{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Run create note infoaw
func (s *GetTemplateItemsService) Run(req *agent.GetTemplateItemsReq) (resp *agent.GetTemplateItemsRes, err error) {
	templateItems, err := s.svcCtx.DaoManager.Template.GetTemplateItemList(s.ctx, req.TemplateId, req.Page, req.Limit, req.Order, req.OrderBy)
	if err != nil {
		return &agent.GetTemplateItemsRes{
			Code: consts.ErrNumDatabase,
			Msg:  consts.ErrMsgDatabase,
			Data: nil,
		}, nil
	}

	return &agent.GetTemplateItemsRes{
		Code: 0,
		Msg:  "get template items successfully",
		Data: &agent.GetTemplateItemsResData{
			TemplateItems: templateItems,
		},
	}, nil
}
