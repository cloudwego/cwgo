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
	"github.com/cloudwego/cwgo/platform/server/shared/logger"
	"go.uber.org/zap"
	"net/http"
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
		idl, err := s.svcCtx.DaoManager.Idl.GetIDL(v)
		if err != nil {
			resp.Code = 400
			resp.Msg = err.Error()
			return resp, err
		}

		repo, err := s.svcCtx.DaoManager.Repository.GetRepository(idl.RepositoryId)
		if err != nil {
			resp.Code = 400
			resp.Msg = err.Error()
			return resp, err
		}

		client, err := s.svcCtx.RepoManager.GetClient(repo.Id)
		if err != nil {
			logger.Logger.Error("get repo client failed", zap.Error(err), zap.Int64("repo_id", repo.Id))
			return &agent.SyncIDLsByIdRes{
				Code: http.StatusInternalServerError,
				Msg:  "internal err",
			}, nil
		}

		owner, repoName, idlPath, err := client.ParseUrl(idl.MainIdlPath)
		if err != nil {
			logger.Logger.Error("parse repo url failed", zap.Error(err))
			return &agent.SyncIDLsByIdRes{
				Code: http.StatusInternalServerError,
				Msg:  "internal err",
			}, nil
		}

		file, err := client.GetFile(owner, repoName, idlPath, consts.MainRef)
		if err != nil {
			logger.Logger.Error("get repo file failed", zap.Error(err))
			return &agent.SyncIDLsByIdRes{
				Code: http.StatusInternalServerError,
				Msg:  "internal err",
			}, nil
		}

		err = s.svcCtx.DaoManager.Idl.SyncIDLContent(idl.Id, string(file.Content))
		if err != nil {
			logger.Logger.Error("sync idl content to dao failed", zap.Error(err))
			return &agent.SyncIDLsByIdRes{
				Code: http.StatusInternalServerError,
				Msg:  "internal err",
			}, nil
		}
	}

	resp.Code = 0
	resp.Msg = "sync IDLs success"

	return resp, nil
}
