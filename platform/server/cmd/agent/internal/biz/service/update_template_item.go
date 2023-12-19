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
	"github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/model"
)

type UpdateTemplateItemService struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
} // NewUpdateTemplateItemService new UpdateTemplateItemService
func NewUpdateTemplateItemService(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTemplateItemService {
	return &UpdateTemplateItemService{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Run create note info
func (s *UpdateTemplateItemService) Run(req *agent.UpdateTemplateItemReq) (resp *agent.UpdateTemplateItemRes, err error) {
	err = s.svcCtx.DaoManager.Template.UpdateTemplateItem(s.ctx, model.TemplateItem{
		Id:      req.Id,
		Name:    req.Name,
		Content: req.Content,
	})
	if err != nil {
		return &agent.UpdateTemplateItemRes{
			Code: consts.ErrNumDatabase,
			Msg:  consts.ErrMsgDatabase,
		}, nil
	}

	return &agent.UpdateTemplateItemRes{
		Code: 0,
		Msg:  "update template item successfully",
	}, nil
}
