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
	"github.com/cloudwego/cwgo/platform/server/cmd/api/internal/svc"
	"github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/repository"
	"github.com/cloudwego/cwgo/platform/server/shared/utils"
)

const (
	successMsgUpdateRepository = "" // TODO: to be filled...
)

type UpdateRepositoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateRepositoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRepositoryLogic {
	return &UpdateRepositoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateRepositoryLogic) UpdateRepository(req *repository.UpdateRepositoryReq) (res *repository.UpdateRepositoryRes) {
	if !utils.ValidStrings(req.Id, req.Token, req.Status) {
		return &repository.UpdateRepositoryRes{
			Code: 400,
			Msg:  "err: The input field contains an empty string",
		}
	}

	err := l.svcCtx.DaoManager.Repository.UpdateRepository(req.Id, req.Token, req.Status)
	if err != nil {
		return &repository.UpdateRepositoryRes{
			Code: 400,
			Msg:  err.Error(),
		}
	}
	//TODO: 给agent发送status变更

	return &repository.UpdateRepositoryRes{
		Code: 0,
		Msg:  successMsgUpdateRepository,
	}
}
