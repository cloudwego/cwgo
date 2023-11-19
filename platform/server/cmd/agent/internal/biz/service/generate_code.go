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
	"github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/agent"
	"github.com/cloudwego/cwgo/platform/server/shared/logger"
	"github.com/cloudwego/cwgo/platform/server/shared/repository"
	"github.com/cloudwego/cwgo/platform/server/shared/utils"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

const (
	successMsgGenerateCode = "generate code successfully"
)

type GenerateCodeService struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
} // NewGenerateCodeService new GenerateCodeService
func NewGenerateCodeService(ctx context.Context, svcCtx *svc.ServiceContext) *GenerateCodeService {
	return &GenerateCodeService{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Run create note info
func (s *GenerateCodeService) Run(req *agent.GenerateCodeReq) (resp *agent.GenerateCodeRes, err error) {
	// get idl info by idl id
	idlModel, err := s.svcCtx.DaoManager.Idl.GetIDL(s.ctx, req.IdlId)
	if err != nil {
		logger.Logger.Error("get idl info failed", zap.Error(err))
		return &agent.GenerateCodeRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
		}, nil
	}

	// get repository info by repository id
	repoModel, err := s.svcCtx.DaoManager.Repository.GetRepository(s.ctx, idlModel.IdlRepositoryId)
	if err != nil {
		logger.Logger.Error("get repo info failed", zap.Error(err))
		return &agent.GenerateCodeRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
		}, nil
	}

	// get repo client
	client, err := s.svcCtx.RepoManager.GetClient(repoModel.Id)
	if err != nil {
		if err == repository.ErrTokenInvalid {
			// repo token is invalid or expired
			return &agent.GenerateCodeRes{
				Code: http.StatusBadRequest,
				Msg:  err.Error(),
			}, nil
		}
		logger.Logger.Error("get repo client failed", zap.Error(err), zap.Int64("repo_id", repoModel.Id))
		return &agent.GenerateCodeRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
		}, nil
	}

	// parsing URLs to obtain information
	idlPid, owner, repoName, err := client.ParseIdlUrl(
		utils.GetRepoFullUrl(
			repoModel.RepositoryType,
			repoModel.RepositoryUrl,
			consts.MainRef,
			idlModel.MainIdlPath,
		),
	)
	if err != nil {
		logger.Logger.Error("parse repo url failed", zap.Error(err))
		return &agent.GenerateCodeRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
		}, nil
	}

	// create temp dir
	tempDir, err := ioutil.TempDir("", strconv.FormatInt(repoModel.Id, 10))
	if err != nil {
		logger.Logger.Error("create temp dir failed", zap.Error(err))
		return &agent.GenerateCodeRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
		}, nil
	}
	defer os.RemoveAll(tempDir)

	// get the entire repository archive
	archiveData, err := client.GetRepositoryArchive(owner, repoName, consts.MainRef)
	if err != nil {
		logger.Logger.Error("get archive failed", zap.Error(err))
		return &agent.GenerateCodeRes{
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
		return &agent.GenerateCodeRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
		}, nil
	}

	// generate code using cwgo
	err = s.svcCtx.Generator.Generate(tempDir+"/"+archiveName+idlPid, idlModel.ServiceName, tempDir)
	if err != nil {
		logger.Logger.Error("generate file failed", zap.Error(err))
		return &agent.GenerateCodeRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
		}, nil
	}

	fileContentMap := make(map[string][]byte)
	// parse the file and add it to the map
	if err := utils.ProcessFolders(fileContentMap, tempDir, "kitex_gen", "rpc"); err != nil {
		return &agent.GenerateCodeRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
		}, nil
	}

	// push files to the repository
	serviceRepositoryModel, err := s.svcCtx.DaoManager.Repository.GetRepository(s.ctx, idlModel.ServiceRepositoryId)
	if err != nil {
		return &agent.GenerateCodeRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
		}, nil
	}

	_, serviceRepoName, err := client.ParseRepoUrl(serviceRepositoryModel.RepositoryUrl)
	if err != nil {
		return &agent.GenerateCodeRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
		}, nil
	}

	err = client.PushFilesToRepository(fileContentMap, owner, serviceRepoName, consts.MainRef, "generated by cwgo")
	if err != nil {
		return &agent.GenerateCodeRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
		}, nil
	}

	return &agent.GenerateCodeRes{
		Code: 0,
		Msg:  successMsgGenerateCode,
	}, nil
}
