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

package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/cloudwego/cwgo/platform/server/cmd/agent/internal/svc"
	"github.com/cloudwego/cwgo/platform/server/shared/consts"
	"github.com/cloudwego/cwgo/platform/server/shared/errx"
	agent "github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/agent"
	"github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/model"
	"github.com/cloudwego/cwgo/platform/server/shared/logger"
	"github.com/cloudwego/cwgo/platform/server/shared/parser"
	"github.com/cloudwego/cwgo/platform/server/shared/utils"
	"go.uber.org/zap"
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
			if errx.GetCode(err) == consts.ErrNumDatabaseRecordNotFound {
				return &agent.SyncIDLsByIdRes{
					Code: consts.ErrNumDatabaseRecordNotFound,
					Msg:  "idl not exist",
				}, nil
			}

			return &agent.SyncIDLsByIdRes{
				Code: consts.ErrNumDatabase,
				Msg:  consts.ErrMsgDatabase,
			}, nil
		}

		repoModel, err := s.svcCtx.DaoManager.Repository.GetRepository(s.ctx, idlModel.IdlRepositoryId)
		if err != nil {
			return &agent.SyncIDLsByIdRes{
				Code: consts.ErrNumDatabase,
				Msg:  consts.ErrMsgDatabase,
			}, nil
		}

		repoClient, err := s.svcCtx.RepoManager.GetClient(repoModel.Id)
		if err != nil {
			if errx.GetCode(err) == consts.ErrNumTokenInvalid {
				// repo token is invalid or expired
				return &agent.SyncIDLsByIdRes{
					Code: consts.ErrNumTokenInvalid,
					Msg:  consts.ErrMsgTokenInvalid,
				}, nil
			}

			logger.Logger.Error(consts.ErrMsgRepoGetClient, zap.Error(err), zap.Int64("repo_id", repoModel.Id))
			return &agent.SyncIDLsByIdRes{
				Code: consts.ErrNumRepoGetClient,
				Msg:  consts.ErrMsgRepoGetClient,
			}, nil
		}

		idlPid, owner, repoName, err := repoClient.ParseFileUrl(
			utils.GetRepoFullUrl(
				repoModel.RepositoryType,
				fmt.Sprintf("https://%s/%s/%s",
					repoModel.RepositoryDomain,
					repoModel.RepositoryOwner,
					repoModel.RepositoryName,
				),
				consts.MainRef,
				idlModel.MainIdlPath,
			),
		)
		if err != nil {
			logger.Logger.Error(consts.ErrMsgParamRepositoryUrl, zap.Error(err))
			return &agent.SyncIDLsByIdRes{
				Code: consts.ErrNumParamRepositoryUrl,
				Msg:  consts.ErrMsgParamRepositoryUrl,
			}, nil
		}

		_, err = repoClient.GetFile(owner, repoName, idlPid, consts.MainRef)
		if err != nil {
			logger.Logger.Error(consts.ErrMsgRepoGetFile, zap.Error(err))
			return &agent.SyncIDLsByIdRes{
				Code: consts.ErrNumRepoGetFile,
				Msg:  consts.ErrMsgRepoGetFile,
			}, nil
		}

		// determine the idl type for subsequent calculations of different types
		idlType, err := utils.DetermineIdlType(idlPid)
		if err != nil {
			return &agent.SyncIDLsByIdRes{
				Code: consts.ErrNumIdlFileExtension,
				Msg:  consts.ErrMsgIdlFileExtension,
			}, nil
		}

		// create temp dir
		tempDir, err := os.MkdirTemp(consts.TempDir, strconv.FormatInt(repoModel.Id, 10))
		if err != nil {
			if os.IsNotExist(err) {
				err = os.Mkdir(consts.TempDir, 0o700)
				if err != nil {
					logger.Logger.Error(consts.ErrMsgCommonCreateTempDir, zap.Error(err))
					return &agent.SyncIDLsByIdRes{
						Code: consts.ErrNumCommonCreateTempDir,
						Msg:  consts.ErrMsgCommonCreateTempDir,
					}, nil
				}

				tempDir, err = os.MkdirTemp(consts.TempDir, strconv.FormatInt(repoModel.Id, 10))
				if err != nil {
					logger.Logger.Error(consts.ErrMsgCommonCreateTempDir, zap.Error(err))
					return &agent.SyncIDLsByIdRes{
						Code: consts.ErrNumCommonCreateTempDir,
						Msg:  consts.ErrMsgCommonCreateTempDir,
					}, nil
				}
			} else {
				logger.Logger.Error(consts.ErrMsgCommonCreateTempDir, zap.Error(err))
				return &agent.SyncIDLsByIdRes{
					Code: consts.ErrNumCommonCreateTempDir,
					Msg:  consts.ErrMsgCommonCreateTempDir,
				}, nil
			}
		}
		defer os.RemoveAll(tempDir)

		// get the entire repository archive
		archiveData, err := repoClient.GetRepositoryArchive(owner, repoName, consts.MainRef)
		if err != nil {
			logger.Logger.Error(consts.ErrMsgRepoGetArchive, zap.Error(err))
			return &agent.SyncIDLsByIdRes{
				Code: consts.ErrNumRepoGetArchive,
				Msg:  consts.ErrMsgRepoGetArchive,
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
			logger.Logger.Error(consts.ErrMsgRepoParseArchive, zap.Error(err))
			return &agent.SyncIDLsByIdRes{
				Code: consts.ErrNumRepoParseArchive,
				Msg:  consts.ErrMsgRepoParseArchive,
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
					Code: consts.ErrNumIdlGetDependentFilePath,
					Msg:  consts.ErrMsgIdlGetDependentFilePath,
				}, nil
			}
		case consts.IdlTypeNumProto:
			protoFile := &parser.ProtoFile{}
			importPaths, err = protoFile.GetDependentFilePaths(tempDir + "/" + archiveName + idlPid)
			return &agent.SyncIDLsByIdRes{
				Code: consts.ErrNumIdlGetDependentFilePath,
				Msg:  consts.ErrMsgIdlGetDependentFilePath,
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
					Code: consts.ErrNumRepoGetCommitHash,
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
				Code: consts.ErrNumRepoGetCommitHash,
				Msg:  consts.ErrMsgRepoGetCommitHash,
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
				Code: consts.ErrNumDatabase,
				Msg:  consts.ErrMsgDatabase,
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