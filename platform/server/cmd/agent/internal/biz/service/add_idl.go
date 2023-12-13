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

type AddIDLService struct {
	ctx          context.Context
	svcCtx       *svc.ServiceContext
	agentService agent.AgentService
} // NewAddIDLService new AddIDLService
func NewAddIDLService(ctx context.Context, svcCtx *svc.ServiceContext, agentService agent.AgentService) *AddIDLService {
	return &AddIDLService{
		ctx:          ctx,
		svcCtx:       svcCtx,
		agentService: agentService,
	}
}

// Run create note info
func (s *AddIDLService) Run(req *agent.AddIDLReq) (resp *agent.AddIDLRes, err error) {
	// check main idl path
	repoClient, err := s.svcCtx.RepoManager.GetClient(req.RepositoryId)
	if err != nil {
		return &agent.AddIDLRes{
			Code: errx.GetCode(err),
			Msg:  err.Error(),
		}, nil
	}

	idlPid, owner, repoName, err := repoClient.ParseFileUrl(req.MainIdlPath)
	if err != nil {
		return &agent.AddIDLRes{
			Code: consts.ErrNumParamMainIdlPath,
			Msg:  consts.ErrMsgParamMainIdlPath,
		}, nil
	}

	isExist, err := s.svcCtx.DaoManager.Idl.CheckMainIdlIfExist(s.ctx, req.RepositoryId, idlPid)
	if err != nil {
		return &agent.AddIDLRes{
			Code: errx.GetCode(err),
			Msg:  err.Error(),
		}, nil
	}
	if isExist {
		return &agent.AddIDLRes{
			Code: consts.ErrNumIdlAlreadyExist,
			Msg:  consts.ErrMsgIdlAlreadyExist,
		}, nil
	}

	_, err = repoClient.GetFile(owner, repoName, idlPid, repoClient.GetBranch())
	if err != nil {
		return &agent.AddIDLRes{
			Code: consts.ErrNumRepoGetFile,
			Msg:  consts.ErrMsgRepoGetFile,
		}, nil
	}

	// obtain the commit hash for the main IDL
	mainIdlHash, err := repoClient.GetLatestCommitHash(owner, repoName, idlPid, repoClient.GetBranch())
	if err != nil {
		return &agent.AddIDLRes{
			Code: consts.ErrNumRepoGetCommitHash,
			Msg:  consts.ErrMsgRepoGetCommitHash,
		}, nil
	}

	// determine the idl type for subsequent calculations of different types
	idlType, err := utils.DetermineIdlType(idlPid)
	if err != nil {
		return &agent.AddIDLRes{
			Code: consts.ErrNumIdlFileExtension,
			Msg:  consts.ErrMsgIdlFileExtension,
		}, nil
	}

	idlRepoModel, err := s.svcCtx.DaoManager.Repository.GetRepository(s.ctx, req.RepositoryId)
	if err != nil {
		logger.Logger.Error("get repository failed", zap.Error(err))
		return &agent.AddIDLRes{
			Code: consts.ErrNumDatabase,
			Msg:  consts.ErrMsgDatabase,
		}, nil
	}

	// create temp dir
	tempDir, err := os.MkdirTemp(consts.TempDir, strconv.FormatInt(idlRepoModel.Id, 10))
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(consts.TempDir, 0o700)
			if err != nil {
				logger.Logger.Error(consts.ErrMsgCommonCreateTempDir, zap.Error(err))
				return &agent.AddIDLRes{
					Code: consts.ErrNumCommonCreateTempDir,
					Msg:  consts.ErrMsgCommonCreateTempDir,
				}, nil
			}
			tempDir, err = os.MkdirTemp(consts.TempDir, strconv.FormatInt(idlRepoModel.Id, 10))
			if err != nil {
				logger.Logger.Error(consts.ErrMsgCommonCreateTempDir, zap.Error(err))
				return &agent.AddIDLRes{
					Code: consts.ErrNumCommonCreateTempDir,
					Msg:  consts.ErrMsgCommonCreateTempDir,
				}, nil
			}
		} else {
			logger.Logger.Error(consts.ErrMsgCommonCreateTempDir, zap.Error(err))
			return &agent.AddIDLRes{
				Code: consts.ErrNumCommonCreateTempDir,
				Msg:  consts.ErrMsgCommonCreateTempDir,
			}, nil
		}
	}
	defer os.RemoveAll(tempDir)

	// get the entire repository archive
	archiveData, err := repoClient.GetRepositoryArchive(owner, repoName, repoClient.GetBranch())
	if err != nil {
		logger.Logger.Error(consts.ErrMsgRepoGetArchive, zap.Error(err))
		return &agent.AddIDLRes{
			Code: consts.ErrNumRepoGetArchive,
			Msg:  consts.ErrMsgRepoGetArchive,
		}, nil
	}

	// the archive type of GitHub is tarball instead of tar
	isTarBall := false
	if idlRepoModel.RepositoryType == consts.RepositoryTypeNumGithub {
		isTarBall = true
	}

	// extract the tar package and persist it to a temporary file
	archiveName, err := utils.UnTar(archiveData, tempDir, isTarBall)
	if err != nil {
		logger.Logger.Error(consts.ErrMsgRepoParseArchive, zap.Error(err))
		return &agent.AddIDLRes{
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
			return &agent.AddIDLRes{
				Code: consts.ErrNumIdlGetDependentFilePath,
				Msg:  consts.ErrMsgIdlGetDependentFilePath,
			}, nil
		}
	case consts.IdlTypeNumProto:
		protoFile := &parser.ProtoFile{}
		importPaths, err = protoFile.GetDependentFilePaths(tempDir + "/" + archiveName + idlPid)
		return &agent.AddIDLRes{
			Code: consts.ErrNumIdlGetDependentFilePath,
			Msg:  consts.ErrMsgIdlGetDependentFilePath,
		}, nil
	}
	importIDLs := make([]*model.ImportIDL, len(importPaths))

	mainIdlDir := filepath.Dir(idlPid)
	// calculate the hash value and add it to the importIDLs slice
	for i, importPath := range importPaths {
		calculatedPath := filepath.ToSlash(filepath.Join(mainIdlDir, importPath))
		commitHash, err := repoClient.GetLatestCommitHash(owner, repoName, calculatedPath, repoClient.GetBranch())
		if err != nil {
			logger.Logger.Error(consts.ErrMsgRepoGetCommitHash, zap.Error(err))
			return &agent.AddIDLRes{
				Code: consts.ErrNumRepoGetCommitHash,
				Msg:  consts.ErrMsgRepoGetCommitHash,
			}, nil
		}

		importIDLs[i] = &model.ImportIDL{
			IdlPath:    calculatedPath,
			CommitHash: commitHash,
		}
	}

	if req.ServiceRepositoryName == "" {
		req.ServiceRepositoryName = "cwgo_" + repoName
	}

	isPrivacy, err := repoClient.GetRepositoryPrivacy(owner, repoName)
	if err != nil {
		logger.Logger.Error(consts.ErrMsgRepoCreate, zap.Error(err))
		return &agent.AddIDLRes{
			Code: consts.ErrNumRepoCreate,
			Msg:  consts.ErrMsgRepoCreate,
		}, nil
	}

	serviceRepoURL, err := repoClient.AutoCreateRepository(owner, req.ServiceRepositoryName, isPrivacy)
	if err != nil {
		logger.Logger.Error(consts.ErrMsgRepoCreate, zap.Error(err))
		return &agent.AddIDLRes{
			Code: consts.ErrNumRepoCreate,
			Msg:  consts.ErrMsgRepoCreate,
		}, nil
	}

	domain, owner, repoName, err := utils.ParseRepoUrl(serviceRepoURL)
	if err != nil {
		return &agent.AddIDLRes{
			Code: consts.ErrNumParamRepositoryUrl,
			Msg:  consts.ErrMsgParamRepositoryUrl,
		}, nil
	}

	serviceRepoId, err := s.svcCtx.DaoManager.Repository.AddRepository(s.ctx, model.Repository{
		RepositoryType:   idlRepoModel.RepositoryType,
		StoreType:        consts.RepositoryStoreTypeNumService,
		RepositoryDomain: domain,
		RepositoryOwner:  owner,
		RepositoryName:   repoName,
		RepositoryBranch: consts.MainRef,
	})
	if err != nil {
		if errx.GetCode(err) == consts.ErrNumDatabaseDuplicateRecord {
			return &agent.AddIDLRes{
				Code: consts.ErrNumDatabaseDuplicateRecord,
				Msg:  consts.ErrMsgDatabaseDuplicateRecord,
			}, nil
		}
	}

	// add idl
	mainIdlId, err := s.svcCtx.DaoManager.Idl.AddIDL(s.ctx, model.IDL{
		IdlRepositoryId:     req.RepositoryId,
		ServiceRepositoryId: serviceRepoId,
		MainIdlPath:         idlPid,
		ServiceName:         req.ServiceName,
		ImportIdls:          importIDLs,
		CommitHash:          mainIdlHash,
	})
	if err != nil {
		return &agent.AddIDLRes{
			Code: consts.ErrNumDatabase,
			Msg:  consts.ErrMsgDatabase,
		}, nil
	}

	// TODO: is async?
	res, err := s.agentService.GenerateCode(s.ctx, &agent.GenerateCodeReq{
		IdlId: mainIdlId,
	})
	if res.Code != 0 {
		return &agent.AddIDLRes{
			Code: res.Code,
			Msg:  res.Msg,
		}, nil
	}

	return &agent.AddIDLRes{
		Code: 0,
		Msg:  "add idl successfully",
	}, nil
}
