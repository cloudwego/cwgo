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
	"github.com/cloudwego/cwgo/platform/server/shared/repository"
	"github.com/cloudwego/cwgo/platform/server/shared/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

type SyncIDLsByIdService struct {
	ctx          context.Context
	svcCtx       *svc.ServiceContext
	agentService agent.AgentService
} // NewSyncIDLsByIdService new SyncIDLsByIdService
func NewSyncIDLsByIdService(ctx context.Context, svcCtx *svc.ServiceContext, agentService agent.AgentService) *SyncIDLsByIdService {
	return &SyncIDLsByIdService{
		ctx:          ctx,
		svcCtx:       svcCtx,
		agentService: agentService,
	}
}

// Run create note info
func (s *SyncIDLsByIdService) Run(req *agent.SyncIDLsByIdReq) (resp *agent.SyncIDLsByIdRes, err error) {
	for _, v := range req.Ids {
		idlModel, err := s.svcCtx.DaoManager.Idl.GetIDL(s.ctx, v)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return &agent.SyncIDLsByIdRes{
					Code: http.StatusBadRequest,
					Msg:  "idl not exist",
				}, nil
			}
			return &agent.SyncIDLsByIdRes{
				Code: http.StatusInternalServerError,
				Msg:  "internal err",
			}, nil
		}

		repoModel, err := s.svcCtx.DaoManager.Repository.GetRepository(s.ctx, idlModel.IdlRepositoryId)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return &agent.SyncIDLsByIdRes{
					Code: http.StatusBadRequest,
					Msg:  "idl not exist",
				}, nil
			}
			return &agent.SyncIDLsByIdRes{
				Code: http.StatusInternalServerError,
				Msg:  "internal err",
			}, nil
		}

		repoClient, err := s.svcCtx.RepoManager.GetClient(repoModel.Id)
		if err != nil {
			if err == repository.ErrTokenInvalid {
				// repo token is invalid or expired
				return &agent.SyncIDLsByIdRes{
					Code: http.StatusBadRequest,
					Msg:  err.Error(),
				}, nil
			}
			logger.Logger.Error("get repo client failed", zap.Error(err), zap.Int64("repo_id", repoModel.Id))
			return &agent.SyncIDLsByIdRes{
				Code: http.StatusInternalServerError,
				Msg:  "internal err",
			}, nil
		}

		idlPid, owner, repoName, err := repoClient.ParseIdlUrl(
			utils.GetRepoFullUrl(
				repoModel.RepositoryType,
				repoModel.RepositoryUrl,
				consts.MainRef,
				idlModel.MainIdlPath,
			),
		)
		if err != nil {
			logger.Logger.Error("parse repo url failed", zap.Error(err))
			return &agent.SyncIDLsByIdRes{
				Code: http.StatusInternalServerError,
				Msg:  "internal err",
			}, nil
		}

		_, err = repoClient.GetFile(owner, repoName, idlPid, consts.MainRef)
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

		// create temp dir
		tempDir, err := os.MkdirTemp(consts.TempDir, strconv.FormatInt(repoModel.Id, 10))
		if err != nil {
			if os.IsNotExist(err) {
				err = os.Mkdir(consts.TempDir, 0700)
				if err != nil {
					logger.Logger.Error("create temp dir failed", zap.Error(err))
					return &agent.SyncIDLsByIdRes{
						Code: http.StatusInternalServerError,
						Msg:  "internal err",
					}, nil
				}

				tempDir, err = os.MkdirTemp(consts.TempDir, strconv.FormatInt(repoModel.Id, 10))
				if err != nil {
					logger.Logger.Error("create temp dir failed", zap.Error(err))
					return &agent.SyncIDLsByIdRes{
						Code: http.StatusInternalServerError,
						Msg:  "internal err",
					}, nil
				}
			} else {
				logger.Logger.Error("create temp dir failed", zap.Error(err))
				return &agent.SyncIDLsByIdRes{
					Code: http.StatusInternalServerError,
					Msg:  "internal err",
				}, nil
			}
		}
		defer os.RemoveAll(tempDir)

		// get the entire repository archive
		archiveData, err := repoClient.GetRepositoryArchive(owner, repoName, consts.MainRef)
		if err != nil {
			logger.Logger.Error("get archive failed", zap.Error(err))
			return &agent.SyncIDLsByIdRes{
				Code: http.StatusInternalServerError,
				Msg:  "internal err",
			}, nil
		}

		// the archive type of GitHub is tarball instead of tar
		isTarBall := false
		if repoModel.RepositoryType == consts.RepositoryTypeNumGithub {
			isTarBall = true
		}

		// extract the tar package and persist it to a temporary file
		archiveName, err := utils.UnTar(archiveData, tempDir, isTarBall)
		if err != nil {
			logger.Logger.Error("parse archive failed", zap.Error(err))
			return &agent.SyncIDLsByIdRes{
				Code: http.StatusInternalServerError,
				Msg:  "internal err",
			}, nil
		}

		// obtain dependent file paths
		var importPaths []string
		switch idlType {
		case consts.IdlTypeNumThrift:
			thriftFile := &parser.ThriftFile{}
			importPaths, err = thriftFile.GetDependentFilePaths(tempDir + "/" + archiveName + idlPid)
			if err != nil {
				return &agent.SyncIDLsByIdRes{
					Code: http.StatusBadRequest,
					Msg:  "get dependent file paths error",
				}, nil
			}
		case consts.IdlTypeNumProto:
			protoFile := &parser.ProtoFile{}
			importPaths, err = protoFile.GetDependentFilePaths(tempDir + "/" + archiveName + idlPid)
			return &agent.SyncIDLsByIdRes{
				Code: http.StatusBadRequest,
				Msg:  "get dependent file paths error",
			}, nil
		}

		importIDLs := make([]*model.ImportIDL, 0)

		mainIdlDir := filepath.Dir(idlPid)
		// calculate the hash value and add it to the importIDLs slice
		for _, importPath := range importPaths {
			calculatedPath := filepath.ToSlash(filepath.Join(mainIdlDir, importPath))
			commitHash, err := repoClient.GetLatestCommitHash(owner, repoName, calculatedPath, consts.MainRef)
			if err != nil {
				return &agent.SyncIDLsByIdRes{
					Code: http.StatusBadRequest,
					Msg:  "cannot get depended idl latest commit hash",
				}, nil
			}

			importIDL := &model.ImportIDL{
				IdlPath:    calculatedPath,
				CommitHash: commitHash,
			}

			importIDLs = append(importIDLs, importIDL)
		}

		// use a bool value to judge whether to sync
		needToSync := false
		if len(importIDLs) == len(idlModel.ImportIdls) {
			// create a map to find imports
			existingImportIDLsMap := make(map[string]struct{})
			for _, importIDL := range importIDLs {
				// use IdlPath as key
				existingImportIDLsMap[importIDL.CommitHash] = struct{}{}
			}

			// compare import idl
			for _, dbImportIDL := range idlModel.ImportIdls {
				if _, ok := existingImportIDLsMap[dbImportIDL.CommitHash]; ok {
					// importIDL exist in importIDLs then continue
					continue
				} else {
					needToSync = true
					break
				}
			}
		} else {
			needToSync = true
		}

		// compare main idl
		hash, err := repoClient.GetLatestCommitHash(owner, repoName, idlPid, consts.MainRef)
		if err != nil {
			logger.Logger.Error("get latest commit hash failed", zap.Error(err))
			return &agent.SyncIDLsByIdRes{
				Code: http.StatusInternalServerError,
				Msg:  "internal err",
			}, nil
		}

		if hash != idlModel.CommitHash {
			needToSync = true
		}

		needToSync = true // TODO: delete

		if !needToSync {
			continue
		}

		err = s.svcCtx.DaoManager.Idl.Sync(s.ctx, model.IDL{
			Id:         idlModel.Id,
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

		res, err := s.agentService.GenerateCode(s.ctx, &agent.GenerateCodeReq{
			IdlId: v,
		})
		if res.Code != 0 {
			return &agent.SyncIDLsByIdRes{
				Code: res.Code,
				Msg:  res.Msg,
			}, nil
		}

	}

	return &agent.SyncIDLsByIdRes{
		Code: 0,
		Msg:  "sync idls successfully",
	}, nil
}
