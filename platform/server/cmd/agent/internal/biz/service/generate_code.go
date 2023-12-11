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
	"strconv"

	"github.com/cloudwego/cwgo/platform/server/cmd/agent/internal/svc"
	"github.com/cloudwego/cwgo/platform/server/shared/consts"
	"github.com/cloudwego/cwgo/platform/server/shared/errx"
	"github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/agent"
	"github.com/cloudwego/cwgo/platform/server/shared/logger"
	"github.com/cloudwego/cwgo/platform/server/shared/utils"
	"go.uber.org/zap"
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
			Code: consts.ErrNumDatabase,
			Msg:  "get idl info failed",
		}, nil
	}

	// get repository info by repository id
	repoModel, err := s.svcCtx.DaoManager.Repository.GetRepository(s.ctx, idlModel.IdlRepositoryId)
	if err != nil {
		logger.Logger.Error("get repo info failed", zap.Error(err))
		return &agent.GenerateCodeRes{
			Code: consts.ErrNumDatabase,
			Msg:  "get repo info failed",
		}, nil
	}

	// get repo client
	client, err := s.svcCtx.RepoManager.GetClient(repoModel.Id)
	if err != nil {
		if errx.GetCode(err) == consts.ErrNumTokenInvalid {
			// repo token is invalid or expired
			return &agent.GenerateCodeRes{
				Code: consts.ErrNumTokenInvalid,
				Msg:  err.Error(),
			}, nil
		}

		logger.Logger.Error(consts.ErrMsgRepoGetClient, zap.Error(err), zap.Int64("repo_id", repoModel.Id))
		return &agent.GenerateCodeRes{
			Code: consts.ErrNumRepoGetClient,
			Msg:  consts.ErrMsgRepoGetClient,
		}, nil
	}

	// parsing URLs to obtain information
	idlPid, owner, repoName, err := client.ParseFileUrl(
		utils.GetRepoFullUrl(
			repoModel.RepositoryType,
			fmt.Sprintf("https://%s/%s/%s",
				repoModel.RepositoryDomain,
				repoModel.RepositoryOwner,
				repoModel.RepositoryName,
			),
			repoModel.RepositoryBranch,
			idlModel.MainIdlPath,
		),
	)
	if err != nil {
		logger.Logger.Error(consts.ErrMsgParamRepositoryUrl, zap.Error(err))
		return &agent.GenerateCodeRes{
			Code: consts.ErrNumParamRepositoryUrl,
			Msg:  consts.ErrMsgParamRepositoryUrl,
		}, nil
	}

	// create temp dir
	tempDir, err := os.MkdirTemp(consts.TempDir, strconv.FormatInt(repoModel.Id, 10))
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(consts.TempDir, 0o700)
			if err != nil {
				logger.Logger.Error(consts.ErrMsgCommonCreateTempDir, zap.Error(err))
				return &agent.GenerateCodeRes{
					Code: consts.ErrNumCommonCreateTempDir,
					Msg:  consts.ErrMsgCommonCreateTempDir,
				}, nil
			}

			tempDir, err = os.MkdirTemp(consts.TempDir, strconv.FormatInt(repoModel.Id, 10))
			if err != nil {
				logger.Logger.Error(consts.ErrMsgCommonCreateTempDir, zap.Error(err))
				return &agent.GenerateCodeRes{
					Code: consts.ErrNumCommonCreateTempDir,
					Msg:  consts.ErrMsgCommonCreateTempDir,
				}, nil
			}
		} else {
			logger.Logger.Error(consts.ErrMsgCommonCreateTempDir, zap.Error(err))
			return &agent.GenerateCodeRes{
				Code: consts.ErrNumCommonCreateTempDir,
				Msg:  consts.ErrMsgCommonCreateTempDir,
			}, nil
		}
	}
	defer os.RemoveAll(tempDir)

	// get the entire repository archive
	archiveData, err := client.GetRepositoryArchive(owner, repoName, repoModel.RepositoryBranch)
	if err != nil {
		logger.Logger.Error(consts.ErrMsgRepoGetArchive, zap.Error(err))
		return &agent.GenerateCodeRes{
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
		return &agent.GenerateCodeRes{
			Code: consts.ErrNumRepoParseArchive,
			Msg:  consts.ErrMsgRepoParseArchive,
		}, nil
	}

	// generate code using cwgo
	err = s.svcCtx.Generator.Generate(tempDir+"/"+archiveName+idlPid, idlModel.ServiceName, tempDir)
	if err != nil {
		logger.Logger.Error(consts.ErrMsgCommonGenerateCode, zap.Error(err))
		return &agent.GenerateCodeRes{
			Code: consts.ErrNumCommonGenerateCode,
			Msg:  consts.ErrMsgCommonGenerateCode,
		}, nil
	}

	fileContentMap := make(map[string][]byte)
	// parse the file and add it to the map
	if err := utils.ProcessFolders(fileContentMap, tempDir, "kitex_gen", "rpc"); err != nil {
		return &agent.GenerateCodeRes{
			Code: consts.ErrNumCommonProcessFolders,
			Msg:  consts.ErrMsgCommonProcessFolders,
		}, nil
	}

	// push files to the repository
	serviceRepositoryModel, err := s.svcCtx.DaoManager.Repository.GetRepository(s.ctx, idlModel.ServiceRepositoryId)
	if err != nil {
		return &agent.GenerateCodeRes{
			Code: consts.ErrNumDatabase,
			Msg:  consts.ErrMsgDatabase,
		}, nil
	}

	err = client.PushFilesToRepository(fileContentMap, owner, serviceRepositoryModel.RepositoryName, repoModel.RepositoryBranch, "generated by cwgo")
	if err != nil {
		return &agent.GenerateCodeRes{
			Code: consts.ErrNumRepoPush,
			Msg:  consts.ErrMsgRepoPush,
		}, nil
	}

	return &agent.GenerateCodeRes{
		Code: 0,
		Msg:  "generate code successfully",
	}, nil
}
