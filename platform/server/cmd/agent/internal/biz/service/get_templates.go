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

type GetTemplatesService struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
} // NewGetTemplatesService new GetTemplatesService
func NewGetTemplatesService(ctx context.Context, svcCtx *svc.ServiceContext) *GetTemplatesService {
	return &GetTemplatesService{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Run create note info
func (s *GetTemplatesService) Run(req *agent.GetTemplatesReq) (resp *agent.GetTemplatesRes, err error) {
	templates, err := s.svcCtx.DaoManager.Template.GetTemplateList(s.ctx, req.Page, req.Limit, req.Order, req.OrderBy)
	if err != nil {
		return &agent.GetTemplatesRes{
			Code: consts.ErrNumDatabase,
			Msg:  consts.ErrMsgDatabase,
			Data: nil,
		}, nil
	}

	return &agent.GetTemplatesRes{
		Code: 0,
		Msg:  "get templates successfully",
		Data: &agent.GetTemplatesResData{
			Templates: templates,
		},
	}, nil
}
