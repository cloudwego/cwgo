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
	"github.com/cloudwego/cwgo/platform/server/cmd/api/internal/biz/model/template"
	"github.com/cloudwego/cwgo/platform/server/cmd/api/internal/svc"
	"github.com/cloudwego/cwgo/platform/server/shared/utils"
)

const (
	successMsgAddTemplate = "" // TODO: to be filled...
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
	if !utils.ValidStrings(req.Name) {
		return &template.AddTemplateRes{
			Code: 400,
			Msg:  "err: The input field contains an empty string",
		}
	}
	//TODO: valid template type

	if !utils.ValidStrings(req.Name) {
		return &template.AddTemplateRes{
			Code: 400,
			Msg:  "err: The input field contains an empty string",
		}
	}
	err := l.svcCtx.DaoManager.Template.AddTemplate(req.Name, req.Type)
	if err != nil {
		return &template.AddTemplateRes{
			Code: 400,
			Msg:  err.Error(),
		}
	}

	return &template.AddTemplateRes{
		Code: 0,
		Msg:  successMsgAddTemplate,
	}
}
