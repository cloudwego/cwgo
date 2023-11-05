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

	idlPid, owner, repoName, err := repoClient.ParseUrl(req.MainIdlPath)

	_, err = repoClient.GetFile(owner, repoName, idlPid, consts.MainRef)
	if err != nil {
		return &agent.AddIDLRes{
			Code: http.StatusBadRequest,
			Msg:  "invalid main idl path",
		}, nil
	}

	// Obtain the commit hash for the main IDL
	mainIdlHash, err := repoClient.GetLatestCommitHash(owner, repoName, idlPid, consts.MainRef)
	if err != nil {
		return &agent.AddIDLRes{
			Code: http.StatusBadRequest,
			Msg:  "invalid main idl path",
		}, nil
	}

	// Determine the idl type for subsequent calculations of different types
	idlType, err := utils.DetermineIdlType(idlPid)
	if err != nil {
		return &agent.AddIDLRes{
			Code: http.StatusBadRequest,
			Msg:  "incorrect idl type",
		}, nil
	}

	// Obtain dependent file paths
	var importPath []string
	switch idlType {
	case consts.IdlTypeNumThrift:
		thriftFile := &parser.ThriftFile{}
		importPath, err = thriftFile.GetDependentFilePaths(req.MainIdlPath)
		if err != nil {
			return &agent.AddIDLRes{
				Code: http.StatusBadRequest,
				Msg:  "get dependent file paths error",
			}, nil
		}
	case consts.IdlTypeNumProto:
		protoFile := &parser.ProtoFile{}
		importPath, err = protoFile.GetDependentFilePaths(req.MainIdlPath)
		return &agent.AddIDLRes{
			Code: http.StatusBadRequest,
			Msg:  "get dependent file paths error",
		}, nil
	}
	importIDLs := make([]*model.ImportIDL, 0)

	// Calculate the hash value and add it to the importIDLs slice
	for _, path := range importPath {
		calculatedPath := filepath.Join(idlPid, path)
		commitHash, err := repoClient.GetLatestCommitHash(owner, repoName, calculatedPath, consts.MainRef)
		if err != nil {
			return &agent.AddIDLRes{
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

	// add idl
	err = s.svcCtx.DaoManager.Idl.AddIDL(s.ctx, model.IDL{
		RepositoryId: req.RepositoryId,
		MainIdlPath:  req.MainIdlPath,
		ServiceName:  req.ServiceName,
		ImportIdls:   importIDLs,
		CommitHash:   mainIdlHash,
	})
	if err != nil {
		return &agent.AddIDLRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
		}, nil
	}

	return
}
