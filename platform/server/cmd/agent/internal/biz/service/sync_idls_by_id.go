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
	"github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/model"
	"github.com/cloudwego/cwgo/platform/server/shared/logger"
	"github.com/cloudwego/cwgo/platform/server/shared/parser"
	"github.com/cloudwego/cwgo/platform/server/shared/utils"
	"go.uber.org/zap"
	"net/http"
	"path/filepath"
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
		idl, err := s.svcCtx.DaoManager.Idl.GetIDL(s.ctx, v)
		if err != nil {
			resp.Code = 400
			resp.Msg = err.Error()
			return resp, err
		}

		repo, err := s.svcCtx.DaoManager.Repository.GetRepository(s.ctx, idl.RepositoryId)
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

		owner, repoName, idlPid, err := client.ParseUrl(idl.MainIdlPath)
		if err != nil {
			logger.Logger.Error("parse repo url failed", zap.Error(err))
			return &agent.SyncIDLsByIdRes{
				Code: http.StatusInternalServerError,
				Msg:  "internal err",
			}, nil
		}

		_, err = client.GetFile(owner, repoName, idlPid, consts.MainRef)
		if err != nil {
			logger.Logger.Error("get repo file failed", zap.Error(err))
			return &agent.SyncIDLsByIdRes{
				Code: http.StatusInternalServerError,
				Msg:  "internal err",
			}, nil
		}

		// determine the idl type for subsequent calculations of different types
		idlType, err := utils.DetermineIdlType(idlPid)
		if err != nil {
			return &agent.SyncIDLsByIdRes{
				Code: http.StatusBadRequest,
				Msg:  "incorrect idl type",
			}, nil
		}

		// obtain dependent file paths
		var importPath []string
		switch idlType {
		case consts.IdlTypeNumThrift:
			thriftFile := &parser.ThriftFile{}
			importPath, err = thriftFile.GetDependentFilePaths(idl.MainIdlPath)
			if err != nil {
				return &agent.SyncIDLsByIdRes{
					Code: http.StatusBadRequest,
					Msg:  "get dependent file paths error",
				}, nil
			}
		case consts.IdlTypeNumProto:
			protoFile := &parser.ProtoFile{}
			importPath, err = protoFile.GetDependentFilePaths(idl.MainIdlPath)
			return &agent.SyncIDLsByIdRes{
				Code: http.StatusBadRequest,
				Msg:  "get dependent file paths error",
			}, nil
		}

		importIDLs := make([]*model.ImportIDL, 0)

		// calculate the hash value and add it to the importIDLs slice
		for _, path := range importPath {
			calculatedPath := filepath.Join(idlPid, path)
			commitHash, err := client.GetLatestCommitHash(owner, repoName, calculatedPath, consts.MainRef)
			if err != nil {
				return &agent.SyncIDLsByIdRes{
					Code: http.StatusBadRequest,
					Msg:  "cannot get depended idl latest commit hash",
				}, nil
			}

			importIDL := &model.ImportIDL{
				IdlPath:    path,
				CommitHash: commitHash,
			}

			importIDLs = append(importIDLs, importIDL)
		}

		// use a bool value to judge whether to sync
		needToSync := false
		// create a map to find imports
		existingImportIDLsMap := make(map[string]bool)
		for _, importIDL := range importIDLs {
			// use IdlPath as key
			existingImportIDLsMap[importIDL.CommitHash] = true
		}

		// compare import idl
		for _, dbImportIDL := range idl.ImportIdls {
			if existingImportIDLsMap[dbImportIDL.CommitHash] {
				// importIDL exist in importIDLs then continue
				continue
			} else {
				needToSync = true
				break
			}
		}

		hash, err := client.GetLatestCommitHash(owner, repoName, idlPid, consts.MainRef)
		if err != nil {
			logger.Logger.Error("get latest commit hash failed", zap.Error(err))
			return &agent.SyncIDLsByIdRes{
				Code: http.StatusInternalServerError,
				Msg:  "internal err",
			}, nil
		}

		if hash != idl.CommitHash {
			needToSync = true
		}

		if !needToSync {
			continue
		}

		err = s.svcCtx.DaoManager.Idl.Sync(s.ctx, model.IDL{
			Id:         idl.Id,
			CommitHash: hash,
			ImportIdls: importIDLs,
		})
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
