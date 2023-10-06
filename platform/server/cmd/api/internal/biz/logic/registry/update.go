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
	registry "github.com/cloudwego/cwgo/platform/server/cmd/api/internal/biz/model/registry"
	"github.com/cloudwego/cwgo/platform/server/cmd/api/internal/svc"
	"github.com/cloudwego/cwgo/platform/server/shared/logger"
	"go.uber.org/zap"
	"net/http"
)

const (
	successMsgUpdate = "" // TODO: to be filled...
)

type UpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLogic) Update(req *registry.UpdateReq) (res *registry.UpdateRes) {
	err := l.svcCtx.BuiltinRegistry.Update(req.ServiceID)
	if err != nil {
		logger.Logger.Error("update service failed", zap.Error(err))
		return &registry.UpdateRes{
			Code: http.StatusBadRequest,
			Msg:  "internal err",
		}
	}

	return &registry.UpdateRes{
		Code: 0,
		Msg:  successMsgUpdate,
	}
}
