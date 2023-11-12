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
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type AddIDLService struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
} // NewAddIDLService new AddIDLService
func NewAddIDLService(ctx context.Context, svcCtx *svc.ServiceContext) *AddIDLService {
	return &AddIDLService{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Run create note info
func (s *AddIDLService) Run(req *agent.AddIDLReq) (resp *agent.AddIDLRes, err error) {
	// check main idl path
	repoClient, err := s.svcCtx.RepoManager.GetClient(req.RepositoryId)
	if err != nil {
		if err == repository.ErrTokenInvalid {
			// repo token is invalid or expired
			return &agent.AddIDLRes{
				Code: http.StatusBadRequest,
				Msg:  err.Error(),
			}, nil
		}

		return &agent.AddIDLRes{
			Code: http.StatusBadRequest,
			Msg:  "can not get the client",
		}, nil
	}

	idlPid, owner, repoName, err := repoClient.ParseIdlUrl(req.MainIdlPath)
	if err != nil {
		return &agent.AddIDLRes{
			Code: http.StatusBadRequest,
			Msg:  "can not parse the IDL url",
		}, nil
	}

	_, err = repoClient.GetFile(owner, repoName, idlPid, consts.MainRef)
	if err != nil {
		return &agent.AddIDLRes{
			Code: http.StatusBadRequest,
			Msg:  "invalid main idl path",
		}, nil
	}

	// obtain the commit hash for the main IDL
	mainIdlHash, err := repoClient.GetLatestCommitHash(owner, repoName, idlPid, consts.MainRef)
	if err != nil {
		return &agent.AddIDLRes{
			Code: http.StatusBadRequest,
			Msg:  "invalid main idl path",
		}, nil
	}

	// determine the idl type for subsequent calculations of different types
	idlType, err := utils.DetermineIdlType(idlPid)
	if err != nil {
		return &agent.AddIDLRes{
			Code: http.StatusBadRequest,
			Msg:  "incorrect idl type",
		}, nil
	}

	repo, err := s.svcCtx.DaoManager.Repository.GetRepository(s.ctx, req.RepositoryId)
	if err != nil {
		logger.Logger.Error("get repository failed", zap.Error(err))
		return &agent.AddIDLRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
		}, nil
	}

	// create temp dir
	tempDir, err := ioutil.TempDir("", strconv.FormatInt(repo.Id, 10))
	if err != nil {
		logger.Logger.Error("create temp dir failed", zap.Error(err))
		return &agent.AddIDLRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
		}, nil
	}
	defer os.RemoveAll(tempDir)

	// get the entire repository archive
	archiveData, err := repoClient.GetRepositoryArchive(owner, repoName, consts.MainRef)
	if err != nil {
		logger.Logger.Error("get archive failed", zap.Error(err))
		return &agent.AddIDLRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
		}, nil
	}

	// the archive type of GitHub is tarball instead of tar
	isTarBall := false
	if repo.RepositoryType == consts.RepositoryTypeNumGithub {
		isTarBall = true
	}

	// extract the tar package and persist it to a temporary file
	archiveName, err := utils.UnTar(archiveData, tempDir, isTarBall)
	if err != nil {
		logger.Logger.Error("parse archive failed", zap.Error(err))
		return &agent.AddIDLRes{
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
			return &agent.AddIDLRes{
				Code: http.StatusBadRequest,
				Msg:  "get dependent file paths error",
			}, nil
		}
	case consts.IdlTypeNumProto:
		protoFile := &parser.ProtoFile{}
		importPaths, err = protoFile.GetDependentFilePaths(tempDir + "/" + archiveName + idlPid)
		return &agent.AddIDLRes{
			Code: http.StatusBadRequest,
			Msg:  "get dependent file paths error",
		}, nil
	}
	importIDLs := make([]*model.ImportIDL, len(importPaths))

	mainIdlDir := filepath.Dir(idlPid)
	// calculate the hash value and add it to the importIDLs slice
	for i, importPath := range importPaths {
		calculatedPath := filepath.ToSlash(filepath.Join(mainIdlDir, importPath))
		commitHash, err := repoClient.GetLatestCommitHash(owner, repoName, calculatedPath, consts.MainRef)
		if err != nil {
			return &agent.AddIDLRes{
				Code: http.StatusBadRequest,
				Msg:  "cannot get depended idl latest commit hash",
			}, nil
		}

		importIDLs[i] = &model.ImportIDL{
			IdlPath:    importPath,
			CommitHash: commitHash,
		}
	}

	if req.ServiceRepositoryName == "" {
		req.ServiceRepositoryName = "cwgo_" + repoName
	}

	isPrivacy, err := repoClient.GetRepositoryPrivacy(owner, repoName)
	if err != nil {
		return &agent.AddIDLRes{
			Code: http.StatusBadRequest,
			Msg:  "internal err",
		}, nil
	}

	serviceRepoURL, err := repoClient.AutoCreateRepository(owner, req.ServiceRepositoryName, isPrivacy)
	if err != nil {
		return &agent.AddIDLRes{
			Code: http.StatusBadRequest,
			Msg:  "internal err",
		}, nil
	}

	serviceRepoId, err := s.svcCtx.DaoManager.Repository.AddRepository(s.ctx, model.Repository{
		RepositoryType: repo.RepositoryType,
		StoreType:      consts.RepositoryStoreTypeNumService,
		RepositoryUrl:  serviceRepoURL,
		Token:          repo.Token,
	})
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return &agent.AddIDLRes{
				Code: http.StatusBadRequest,
				Msg:  "repository is already exist",
			}, nil
		}
		return &agent.AddIDLRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
		}, nil
	}

	// add idl
	err = s.svcCtx.DaoManager.Idl.AddIDL(s.ctx, model.IDL{
		IdlRepositoryId:     req.RepositoryId,
		ServiceRepositoryId: serviceRepoId,
		MainIdlPath:         req.MainIdlPath,
		ServiceName:         req.ServiceName,
		ImportIdls:          importIDLs,
		CommitHash:          mainIdlHash,
	})
	if err != nil {
		return &agent.AddIDLRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
		}, nil
	}

	return
}
