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

package repository

import (
	"context"
	"net/url"

	"github.com/cloudwego/cwgo/platform/server/cmd/api/internal/svc"
	"github.com/cloudwego/cwgo/platform/server/shared/consts"
	"github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/agent"
	"github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/repository"
	"github.com/cloudwego/cwgo/platform/server/shared/log"
	"github.com/cloudwego/cwgo/platform/server/shared/utils"
	"go.uber.org/zap"
)

const (
	successMsgAddRepository = "add repository successfully"
)

type AddRepositoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddRepositoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddRepositoryLogic {
	return &AddRepositoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddRepositoryLogic) AddRepository(req *repository.AddRepositoryReq) (res *repository.AddRepositoryRes) {
	// validate repository url
	_, err := url.Parse(req.RepositoryUrl)
	if err != nil {
		return &repository.AddRepositoryRes{
			Code: consts.ErrNumParamRepositoryUrl,
			Msg:  consts.ErrMsgParamRepositoryUrl,
		}
	}

	if _, ok := consts.RepositoryTypeNumMap[int(req.RepositoryType)]; !ok {
		return &repository.AddRepositoryRes{
			Code: consts.ErrNumParamRepositoryType,
			Msg:  consts.ErrMsgParamRepositoryType,
		}
	}

	// parse repository url
	domain, owner, repoName, err := utils.ParseRepoUrl(req.RepositoryUrl)
	if err != nil {
		return &repository.AddRepositoryRes{
			Code: consts.ErrNumParamRepositoryUrl,
			Msg:  consts.ErrMsgParamRepositoryUrl,
		}
	}

	if req.RepositoryType == consts.RepositoryTypeNumGithub {
		if domain != consts.GitHubDomain {
			return &repository.AddRepositoryRes{
				Code: consts.ErrNumParamRepositoryUrl,
				Msg:  "invalid github repository url",
			}
		}
	}

	client, err := l.svcCtx.Manager.GetAgentClient()
	if err != nil {
		log.Error(consts.ErrMsgRpcGetClient, zap.Error(err))
		return &repository.AddRepositoryRes{
			Code: consts.ErrNumRpcGetClient,
			Msg:  consts.ErrMsgRpcGetClient,
		}
	}

	rpcRes, err := client.AddRepository(l.ctx, &agent.AddRepositoryReq{
		RepositoryType:   req.RepositoryType,
		RepositoryDomain: domain,
		RepositoryOwner:  owner,
		RepositoryName:   repoName,
		Branch:           req.Branch,
		StoreType:        req.StoreType,
	})
	if err != nil {
		log.Error(consts.ErrMsgRpcConnectClient, zap.Error(err))
		return &repository.AddRepositoryRes{
			Code: consts.ErrNumRpcConnectClient,
			Msg:  consts.ErrMsgRpcConnectClient,
		}
	}
	if rpcRes.Code != 0 {
		return &repository.AddRepositoryRes{
			Code: rpcRes.Code,
			Msg:  rpcRes.Msg,
		}
	}

	return &repository.AddRepositoryRes{
		Code: 0,
		Msg:  successMsgAddRepository,
	}
}
