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

package template

import (
	"context"
	"github.com/cloudwego/cwgo/platform/server/cmd/api/internal/svc"
	"github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/template"
	"github.com/cloudwego/cwgo/platform/server/shared/utils"
)

const (
	successMsgAddTemplateItem = "" // TODO: to be filled...
)

type AddTemplateItemLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddTemplateItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddTemplateItemLogic {
	return &AddTemplateItemLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddTemplateItemLogic) AddTemplateItem(req *template.AddTemplateItemReq) (res *template.AddTemplateItemRes) {
	if !utils.ValidStrings(req.Name, req.Content) {
		return &template.AddTemplateItemRes{
			Code: 400,
			Msg:  "err: The input field contains an empty string",
		}
	}
	err := l.svcCtx.DaoManager.Template.AddTemplateItem(req.TemplateId, req.Name, req.Content)
	if err != nil {
		return &template.AddTemplateItemRes{
			Code: 400,
			Msg:  err.Error(),
		}
	}

	return &template.AddTemplateItemRes{
		Code: 0,
		Msg:  successMsgAddTemplateItem,
	}
}
