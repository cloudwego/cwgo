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
	"github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/model"
)

type GetIDLsService struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
} // NewGetIDLsService new GetIDLsService
func NewGetIDLsService(ctx context.Context, svcCtx *svc.ServiceContext) *GetIDLsService {
	return &GetIDLsService{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Run create note info
func (s *GetIDLsService) Run(req *agent.GetIDLsReq) (resp *agent.GetIDLsRes, err error) {
	idls, total, err := s.svcCtx.DaoManager.Idl.GetIDLList(s.ctx,
		model.IDL{
			ServiceName: req.ServiceName,
		},
		req.Page, req.Limit, req.Order, req.OrderBy,
	)
	if err != nil {
		return &agent.GetIDLsRes{
			Code: consts.ErrNumDatabase,
			Msg:  consts.ErrMsgDatabase,
			Data: nil,
		}, err
	}

	return &agent.GetIDLsRes{
		Code: 0,
		Msg:  "get idls successfully",
		Data: &agent.GetIDLsResData{
			Idls:  idls,
			Total: int32(total),
		},
	}, nil
}
