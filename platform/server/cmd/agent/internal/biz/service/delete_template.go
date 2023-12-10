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

package service

import (
	"context"

	"github.com/cloudwego/cwgo/platform/server/cmd/agent/internal/svc"
	"github.com/cloudwego/cwgo/platform/server/shared/consts"
	agent "github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/agent"
)

type DeleteTemplateService struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
} // NewDeleteTemplateService new DeleteTemplateService
func NewDeleteTemplateService(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteTemplateService {
	return &DeleteTemplateService{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Run create note info
func (s *DeleteTemplateService) Run(req *agent.DeleteTemplateReq) (resp *agent.DeleteTemplateRes, err error) {
	err = s.svcCtx.DaoManager.Template.DeleteTemplate(s.ctx, req.Ids)
	if err != nil {
		return &agent.DeleteTemplateRes{
			Code: consts.ErrNumDatabase,
			Msg:  consts.ErrMsgDatabase,
		}, nil
	}

	return &agent.DeleteTemplateRes{
		Code: 0,
		Msg:  "delete template successfully",
	}, nil
}
