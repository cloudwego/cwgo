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
)

const (
	successMsgDeleteRepository = "" // TODO: to be filled...
)

type DeleteRepositoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteRepositoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteRepositoryLogic {
	return &DeleteRepositoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteRepositoryLogic) DeleteRepository(req *repository.DeleteRepositoriesReq) (res *repository.DeleteRepositoriesRes) {
	err := l.svcCtx.DaoManager.Repository.DeleteRepository(req.Ids)
	if err != nil {
		return &repository.DeleteRepositoriesRes{
			Code: 400,
			Msg:  err.Error(),
		}
	}

	return &repository.DeleteRepositoriesRes{
		Code: 0,
		Msg:  successMsgDeleteRepository,
	}
}
