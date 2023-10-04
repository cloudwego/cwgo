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

package registry

import (
	"context"
	"github.com/cloudwego/cwgo/platform/server/cmd/api/internal/svc"
	registry "github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/registry"
	"github.com/cloudwego/cwgo/platform/server/shared/logger"
	"go.uber.org/zap"
	"net/http"
)

const (
	successMsgDeregister = "" // TODO: to be filled...
)

type DeregisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeregisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeregisterLogic {
	return &DeregisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeregisterLogic) Deregister(req *registry.DeregisterReq) (res *registry.DeRegisterRes) {
	err := l.svcCtx.BuiltinRegistry.Deregister(req.ServiceId)
	if err != nil {
		logger.Logger.Error("deregister service failed", zap.Error(err))
		return &registry.DeRegisterRes{
			Code: http.StatusBadRequest,
			Msg:  "internal err",
		}
	}

	return &registry.DeRegisterRes{
		Code: 0,
		Msg:  successMsgDeregister,
	}
}
