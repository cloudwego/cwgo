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
	"github.com/cloudwego/cwgo/platform/server/shared/utils"
)

const (
	successMsgAddRepository = "" // TODO: to be filled...
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
	if !utils.ValidStrings(req.RepositoryURL, req.Token) {
		return &repository.AddRepositoryRes{
			Code: 400,
			Msg:  "err: The input field contains an empty string",
		}
	}

	if !utils.ValidRepoType(req.RepositoryType) {
		return &repository.AddRepositoryRes{
			Code: 400,
			Msg:  "err: Incorrect repository type",
		}
	}

	err := l.svcCtx.DaoManager.Repository.AddRepository(req.RepositoryURL, req.Token, consts.RepositoryStatusActive, req.RepositoryType)
	if err != nil {
		return &repository.AddRepositoryRes{
			Code: 400,
			Msg:  err.Error(),
		}
	}

	return &repository.AddRepositoryRes{
		Code: 0,
		Msg:  successMsgAddRepository,
	}
}
