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

package idl

import (
	"context"
	"github.com/cloudwego/cwgo/platform/server/cmd/api/internal/svc"
	"github.com/cloudwego/cwgo/platform/server/shared/consts"
	"github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/idl"
	"github.com/cloudwego/cwgo/platform/server/shared/utils"
)

const (
	successMsgSyncIDLs = "" // TODO: to be filled...
)

type SyncIDLsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSyncIDLsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SyncIDLsLogic {
	return &SyncIDLsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SyncIDLsLogic) SyncIDLs(req *idl.SyncIDLsByIdReq) (res *idl.SyncIDLsByIdRes) {
	for _, v := range req.Ids {
		Idl, err := l.svcCtx.DaoManager.Idl.GetIDL(v)
		if err != nil {
			return &idl.SyncIDLsByIdRes{
				Code: 400,
				Msg:  err.Error(),
			}
		}

		repoType, err := l.svcCtx.DaoManager.Repository.GetRepoTypeByID(Idl.RepositoryId)
		if err != nil {
			return &idl.SyncIDLsByIdRes{
				Code: 400,
				Msg:  err.Error(),
			}
		}

		switch repoType {
		case consts.GitLab:
			ref := consts.MainRef
			owner, repoName, idlPid, err := utils.ParseGitlabIdlURL(Idl.MainIdlPath)
			if err != nil {
				return &idl.SyncIDLsByIdRes{
					Code: 400,
					Msg:  err.Error(),
				}
			}
			file, err := l.svcCtx.RepoManager.Gitlab.GetFile(owner, repoName, idlPid, ref)
			err = l.svcCtx.DaoManager.Idl.SyncIDLContent(Idl.Id, string(file.Content))
			if err != nil {
				return nil
			}
		}
	}

	return &idl.SyncIDLsByIdRes{
		Code: 0,
		Msg:  successMsgSyncIDLs,
	}
}
