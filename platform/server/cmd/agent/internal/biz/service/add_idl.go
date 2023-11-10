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
	"github.com/cloudwego/cwgo/platform/server/shared/parser"
	"github.com/cloudwego/cwgo/platform/server/shared/utils"
	"net/http"
	"path/filepath"
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

	idlPid, owner, repoName, err := repoClient.ParseIdlUrl(req.MainIdlPath)

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

	// obtain dependent file paths
	var importPaths []string
	switch idlType {
	case consts.IdlTypeNumThrift:
		thriftFile := &parser.ThriftFile{}
		importPaths, err = thriftFile.GetDependentFilePaths(req.MainIdlPath)
		if err != nil {
			return &agent.AddIDLRes{
				Code: http.StatusBadRequest,
				Msg:  "get dependent file paths error",
			}, nil
		}
	case consts.IdlTypeNumProto:
		protoFile := &parser.ProtoFile{}
		importPaths, err = protoFile.GetDependentFilePaths(req.MainIdlPath)
		return &agent.AddIDLRes{
			Code: http.StatusBadRequest,
			Msg:  "get dependent file paths error",
		}, nil
	}
	importIDLs := make([]*model.ImportIDL, len(importPaths))

	// calculate the hash value and add it to the importIDLs slice
	for i, importPath := range importPaths {
		calculatedPath := filepath.Join(idlPid, importPath)
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
	repo, err := s.svcCtx.DaoManager.Repository.GetRepository(s.ctx, req.RepositoryId)

	if err != nil {
		return &agent.AddIDLRes{
			Code: http.StatusBadRequest,
			Msg:  "internal err",
		}, nil
	}

	if req.ServiceRepositoryName == "" {
		req.ServiceRepositoryName = "cwgo_" + repoName
	}
	serviceRepoURL, err := repoClient.AutoCreateRepository(owner, req.ServiceRepositoryName)
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
		return &agent.AddIDLRes{
			Code: http.StatusBadRequest,
			Msg:  "internal err",
		}, nil
	}

	// add idl
	err = s.svcCtx.DaoManager.Idl.AddIDL(s.ctx, model.IDL{
		RepositoryId:        req.RepositoryId,
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
