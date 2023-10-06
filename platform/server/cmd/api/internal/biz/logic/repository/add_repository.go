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

package repository

import (
	"context"
	"github.com/cloudwego/cwgo/platform/server/cmd/api/internal/biz/model/repository"
	"github.com/cloudwego/cwgo/platform/server/cmd/api/internal/svc"
	"github.com/cloudwego/cwgo/platform/server/shared/consts"
	"github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/agent"
	"github.com/cloudwego/cwgo/platform/server/shared/logger"
	"go.uber.org/zap"
	"net/http"
	"net/url"
)

const (
	successMsgAddRepository = "add repository successfully"
)

type AddRepositoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddRepositoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddRepositoryLogic {
	return &AddRepositoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddRepositoryLogic) AddRepository(req *repository.AddRepositoryReq) (res *repository.AddRepositoryRes) {
	_, err := url.Parse(req.RepositoryURL)
	if err != nil {
		return &repository.AddRepositoryRes{
			Code: http.StatusBadRequest,
			Msg:  "invalid repository url path",
		}
	}

	if _, ok := consts.RepositoryTypeNumMap[int(req.RepositoryType)]; !ok {
		return &repository.AddRepositoryRes{
			Code: http.StatusBadRequest,
			Msg:  "invalid repository type",
		}
	}

	client, err := l.svcCtx.Manager.GetAgentClient()
	if err != nil {
		logger.Logger.Error("get rpc client failed", zap.Error(err))
		return &repository.AddRepositoryRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
		}
	}

	rpcRes, err := client.AddRepository(l.ctx, &agent.AddRepositoryReq{
		RepositoryType: req.RepositoryType,
		RepositoryUrl:  req.RepositoryURL,
		Token:          req.Token,
	})
	if err != nil {
		logger.Logger.Error("connect to rpc client failed", zap.Error(err))
		return &repository.AddRepositoryRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
		}
	}
	if rpcRes.Code != 0 {
		return &repository.AddRepositoryRes{
			Code: http.StatusBadRequest,
			Msg:  rpcRes.Msg,
		}
	}

	return &repository.AddRepositoryRes{
		Code: 0,
		Msg:  successMsgAddRepository,
	}
}
