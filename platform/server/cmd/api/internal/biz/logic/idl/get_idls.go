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
	"github.com/cloudwego/cwgo/platform/server/cmd/api/internal/svc"
	"github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/idl"
	"github.com/cloudwego/cwgo/platform/server/shared/utils"
)

const (
	successMsgGetIDLs = "" // TODO: to be filled...
)

type GetIDLsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetIDLsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetIDLsLogic {
	return &GetIDLsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetIDLsLogic) GetIDLs(req *idl.GetIDLsReq) (res *idl.GetIDLsRes) {
	if !utils.ValidOrder(req.Order) || !utils.ValidOrderBy(req.OrderBy) {
		return &idl.GetIDLsRes{
			Code: 400,
			Msg:  "err: invalid field",
			Data: nil,
		}
	}
	idls, err := l.svcCtx.DaoManager.Idl.GetIDLs(req.Page, req.Limit, req.Order, req.OrderBy)
	if err != nil {
		return &idl.GetIDLsRes{
			Code: 400,
			Msg:  err.Error(),
			Data: nil,
		}
	}

	return &idl.GetIDLsRes{
		Code: 0,
		Msg:  successMsgGetIDLs,
		Data: &idl.GetIDLsResData{Idls: idls},
	}
}
