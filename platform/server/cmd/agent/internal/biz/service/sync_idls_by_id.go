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
	"github.com/cloudwego/cwgo/platform/server/shared/log"
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
	for _, syncMainIdlId := range req.Ids {
		idlEntityWithRepoInfo, err := s.svcCtx.DaoManager.Idl.GetIDL(s.ctx, syncMainIdlId)
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

		idlRepoModel, err := s.svcCtx.DaoManager.Repository.GetRepository(s.ctx, idlEntityWithRepoInfo.IdlRepositoryId)
		if err != nil {
			return &agent.SyncIDLsByIdRes{
				Code: consts.ErrNumDatabase,
				Msg:  consts.ErrMsgDatabase,
			}, nil
		}

		repoClient, err := s.svcCtx.RepoManager.GetClient(idlRepoModel.Id)
		if err != nil {
			if errx.GetCode(err) == consts.ErrNumTokenInvalid {
				// repo token is invalid or expired
				return &agent.SyncIDLsByIdRes{
					Code: consts.ErrNumTokenInvalid,
					Msg:  consts.ErrMsgTokenInvalid,
				}, nil
			}

			log.Error(consts.ErrMsgRepoGetClient, zap.Error(err), zap.Int64("repo_id", idlRepoModel.Id))
			return &agent.SyncIDLsByIdRes{
				Code: consts.ErrNumRepoGetClient,
				Msg:  consts.ErrMsgRepoGetClient,
			}, nil
		}

		idlPid, owner, repoName, err := repoClient.ParseFileUrl(
			utils.GetRepoFullUrl(
				idlRepoModel.RepositoryType,
				fmt.Sprintf("https://%s/%s/%s",
					idlRepoModel.RepositoryDomain,
					idlRepoModel.RepositoryOwner,
					idlRepoModel.RepositoryName,
				),
				idlRepoModel.RepositoryBranch,
				idlEntityWithRepoInfo.MainIdlPath,
			),
		)
		if err != nil {
			log.Error(consts.ErrMsgParamRepositoryUrl, zap.Error(err))
			return &agent.SyncIDLsByIdRes{
				Code: consts.ErrNumParamRepositoryUrl,
				Msg:  consts.ErrMsgParamRepositoryUrl,
			}, nil
		}

		// get the entire repository archive
		archiveData, err := repoClient.GetRepositoryArchive(owner, repoName, idlRepoModel.RepositoryBranch)
		if err != nil {
			log.Error(consts.ErrMsgRepoGetArchive, zap.Error(err))
			return &agent.SyncIDLsByIdRes{
				Code: consts.ErrNumRepoGetArchive,
				Msg:  consts.ErrMsgRepoGetArchive,
			}, nil
		}

		// create temp dir
		tempDir, err := os.MkdirTemp(consts.TempDir, strconv.FormatInt(idlRepoModel.Id, 10))
		if err != nil {
			if os.IsNotExist(err) {
				err = os.Mkdir(consts.TempDir, 0o700)
				if err != nil {
					log.Error(consts.ErrMsgCommonCreateTempDir, zap.Error(err))
					return &agent.SyncIDLsByIdRes{
						Code: consts.ErrNumCommonCreateTempDir,
						Msg:  consts.ErrMsgCommonCreateTempDir,
					}, nil
				}

				tempDir, err = os.MkdirTemp(consts.TempDir, strconv.FormatInt(idlRepoModel.Id, 10))
				if err != nil {
					log.Error(consts.ErrMsgCommonCreateTempDir, zap.Error(err))
					return &agent.SyncIDLsByIdRes{
						Code: consts.ErrNumCommonCreateTempDir,
						Msg:  consts.ErrMsgCommonCreateTempDir,
					}, nil
				}
			} else {
				log.Error(consts.ErrMsgCommonCreateTempDir, zap.Error(err))
				return &agent.SyncIDLsByIdRes{
					Code: consts.ErrNumCommonCreateTempDir,
					Msg:  consts.ErrMsgCommonCreateTempDir,
				}, nil
			}
		}
		defer os.RemoveAll(tempDir)

		tempDirRepo := tempDir + "/" + consts.TempDirRepo

		// the archive type of GitHub is tarball instead of tar
		isTarBall := false
		if idlRepoModel.RepositoryType == consts.RepositoryTypeNumGithub {
			isTarBall = true
		}

		// extract the tar package and persist it to a temporary file
		archiveName, err := utils.UnTar(archiveData, tempDirRepo, isTarBall)
		if err != nil {
			log.Error(consts.ErrMsgRepoParseArchive, zap.Error(err))
			return &agent.SyncIDLsByIdRes{
				Code: consts.ErrNumRepoParseArchive,
				Msg:  consts.ErrMsgRepoParseArchive,
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

		idlParser := parser.NewParser(idlType)
		if idlParser == nil {
			return &agent.SyncIDLsByIdRes{
				Code: consts.ErrNumIdlFileExtension,
				Msg:  consts.ErrMsgIdlFileExtension,
			}, nil
		}
		var importPaths []string
		var importBaseDirPath string
		importBaseDirPath, importPaths, err = idlParser.GetDependentFilePaths(tempDirRepo+"/"+archiveName, idlPid)
		if err != nil {
			return &agent.SyncIDLsByIdRes{
				Code: consts.ErrNumIdlGetDependentFilePath,
				Msg:  consts.ErrMsgIdlGetDependentFilePath,
			}, nil
		}

		needToSync := false
		importIDLs := make([]*model.ImportIDL, 0)

		// compare main idl
		mainIdlHash, err := repoClient.GetLatestCommitHash(owner, repoName, idlPid, idlRepoModel.RepositoryBranch)
		if err != nil {
			log.Error("get latest commit hash failed", zap.Error(err))
			return &agent.SyncIDLsByIdRes{
				Code: consts.ErrNumRepoGetCommitHash,
				Msg:  consts.ErrMsgRepoGetCommitHash,
			}, nil
		}

		if mainIdlHash != idlEntityWithRepoInfo.CommitHash {
			needToSync = true
		} else {
			// if mail idl is not changed
			// then compare imported idl files

			// calculate the mainIdlHa value and add it to the importIDLs slice
			for _, importPath := range importPaths {
				calculatedPath := filepath.ToSlash(filepath.Join(importBaseDirPath, importPath))
				commitHash, err := repoClient.GetLatestCommitHash(owner, repoName, calculatedPath, idlRepoModel.RepositoryBranch)
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
			if len(importIDLs) == len(idlEntityWithRepoInfo.ImportIdls) {
				// create a map to find imports
				existingImportIDLsMap := make(map[string]struct{})
				for _, importIDL := range importIDLs {
					// use IdlPath as key
					existingImportIDLsMap[importIDL.CommitHash] = struct{}{}
				}

				// compare import idl
				for _, dbImportIDL := range idlEntityWithRepoInfo.ImportIdls {
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
		}

		if !needToSync {
			continue
		}

		idlEntityWithRepoInfo.ImportIdls = importIDLs

		if idlType != consts.IdlTypeNumProto {
			importBaseDirPath = ""
		}
		err = s.svcCtx.GenerateCode(s.ctx, repoClient,
			tempDir, importBaseDirPath, idlEntityWithRepoInfo, idlRepoModel, archiveName)
		if err != nil {
			return &agent.SyncIDLsByIdRes{
				Code: errx.GetCode(err),
				Msg:  err.Error(),
			}, nil
		}

		err = s.svcCtx.DaoManager.Idl.Sync(s.ctx, model.IDL{
			Id:         idlEntityWithRepoInfo.Id,
			CommitHash: mainIdlHash,
			ImportIdls: importIDLs,
		})
		if err != nil {
			log.Error("sync idl content to dao failed", zap.Error(err))
			return &agent.SyncIDLsByIdRes{
				Code: consts.ErrNumDatabase,
				Msg:  consts.ErrMsgDatabase,
			}, nil
		}

	}

	return &agent.SyncIDLsByIdRes{
		Code: 0,
		Msg:  "sync idls successfully",
	}, nil
}
