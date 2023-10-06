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
	"github.com/cloudwego/cwgo/platform/server/shared/utils"
)

const (
	successMsgUpdateIDL = "" // TODO: to be filled...
)

type UpdateIDLLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateIDLLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateIDLLogic {
	return &UpdateIDLLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateIDLLogic) UpdateIDL(req *idl.UpdateIDLReq) (res *idl.UpdateIDLRes) {
	if !utils.ValidStrings(req.MainIdlPath, req.ServiceName) {
		return &idl.UpdateIDLRes{
			Code: 400,
			Msg:  "err: The input field contains an empty string",
		}
	}

	err := l.svcCtx.DaoManager.Idl.UpdateIDL(req.ID, req.RepositoryID, req.MainIdlPath, req.ServiceName)
	if err != nil {
		return &idl.UpdateIDLRes{
			Code: 400,
			Msg:  err.Error(),
		}
	}

	return &idl.UpdateIDLRes{
		Code: 0,
		Msg:  successMsgUpdateIDL,
	}
}
