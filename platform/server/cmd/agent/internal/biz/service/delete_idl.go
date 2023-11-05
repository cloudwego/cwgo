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
	agent "github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/agent"
	"net/http"
)

type DeleteIDLService struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
} // NewDeleteIDLService new DeleteIDLService
func NewDeleteIDLService(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteIDLService {
	return &DeleteIDLService{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Run create note info
func (s *DeleteIDLService) Run(req *agent.DeleteIDLsReq) (resp *agent.DeleteIDLsRes, err error) {
	err = s.svcCtx.DaoManager.Idl.DeleteIDLs(s.ctx, req.Ids)
	if err != nil {
		return &agent.DeleteIDLsRes{
			Code: http.StatusBadRequest,
			Msg:  "internal error",
		}, err
	}

	return &agent.DeleteIDLsRes{
		Code: 0,
		Msg:  "delete idl successfully",
	}, nil
}
