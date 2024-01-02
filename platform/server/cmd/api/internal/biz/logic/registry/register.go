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

package registry

import (
	"context"
	"net/http"

	"github.com/cloudwego/cwgo/platform/server/cmd/api/internal/svc"
	registry "github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/registry"
	"github.com/cloudwego/cwgo/platform/server/shared/log"
	"go.uber.org/zap"
)

const (
	successMsgRegister = "register service successfully"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *registry.RegisterReq) (res *registry.RegisterRes) {
	err := l.svcCtx.BuiltinRegistry.Register(req.ServiceId, req.Host, int(req.Port))
	if err != nil {
		log.Error("register service failed", zap.Error(err))
		return &registry.RegisterRes{
			Code: http.StatusBadRequest,
			Msg:  "internal err",
		}
	}

	return &registry.RegisterRes{
		Code: 0,
		Msg:  successMsgRegister,
	}
}
