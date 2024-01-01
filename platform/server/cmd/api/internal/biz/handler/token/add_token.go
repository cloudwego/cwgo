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

package token

import (
	"context"

	tokenlogic "github.com/cloudwego/cwgo/platform/server/cmd/api/internal/biz/logic/token"
	"github.com/cloudwego/cwgo/platform/server/cmd/api/internal/svc"
	"github.com/cloudwego/cwgo/platform/server/cmd/api/pkg/model/response"
	tokenmodel "github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/token"
	"github.com/cloudwego/cwgo/platform/server/shared/logger"
	"github.com/cloudwego/hertz/pkg/app"
	hertzconsts "github.com/cloudwego/hertz/pkg/protocol/consts"
	"go.uber.org/zap"
)

// AddToken .
// @router /api/token [POST]
func AddToken(ctx context.Context, rCtx *app.RequestContext) {
	var err error
	var req tokenmodel.AddTokenReq
	err = rCtx.BindAndValidate(&req)
	if err != nil {
		logger.Logger.Debug("parse http request failed", zap.Error(err), zap.Reflect("http request", req))
		response.Fail(
			rCtx,
			hertzconsts.StatusBadRequest,
			hertzconsts.StatusBadRequest,
			err.Error(),
		)
		return
	}

	logger.Logger.Debug("http request args", zap.Reflect("args", req))

	l := tokenlogic.NewAddTokenLogic(ctx, svc.Svc)

	res := l.AddToken(&req)

	logger.Logger.Debug("http response args", zap.Reflect("args", res))

	if res.Code != 0 {
		response.Fail(rCtx, hertzconsts.StatusBadRequest, res.Code, res.Msg)
		return
	}

	response.OkWithData(rCtx, res.Msg, res.Data)
}
