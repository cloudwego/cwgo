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
	"github.com/cloudwego/cwgo/platform/server/shared/consts"
	agent "github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/agent"
	"github.com/cloudwego/cwgo/platform/server/shared/utils"
)

type SyncIDLsByIdService struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
} // NewSyncIDLsByIdService new SyncIDLsByIdService
func NewSyncIDLsByIdService(ctx context.Context, svcCtx *svc.ServiceContext) *SyncIDLsByIdService {
	return &SyncIDLsByIdService{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Run create note info
func (s *SyncIDLsByIdService) Run(req *agent.SyncIDLsByIdReq) (resp *agent.SyncIDLsByIdRes, err error) {
	for _, v := range req.Ids {
		Idl, err := s.svcCtx.DaoManager.Idl.GetIDL(v)
		if err != nil {
			resp.Code = 400
			resp.Msg = err.Error()
			return resp, err
		}

		repo, err := s.svcCtx.DaoManager.Repository.GetRepository(Idl.RepositoryId)
		if err != nil {
			resp.Code = 400
			resp.Msg = err.Error()
			return resp, err
		}

		switch repo.Type {
		case consts.GitLab:
			ref := consts.MainRef
			owner, repoName, idlPath, err := utils.ParseGitlabIdlURL(Idl.MainIdlPath)
			if err != nil {
				resp.Code = 400
				resp.Msg = err.Error()
				return resp, err
			}
			file, err := s.svcCtx.RepoManager.GitLab.GetFile(repo.Id, owner, repoName, idlPath, ref)
			err = s.svcCtx.DaoManager.Idl.SyncIDLContent(Idl.Id, string(file.Content))
			if err != nil {
				resp.Code = 400
				resp.Msg = err.Error()
				return resp, err
			}
		}
	}

	resp.Code = 0
	resp.Msg = "sync IDLs success"

	return resp, nil
}
