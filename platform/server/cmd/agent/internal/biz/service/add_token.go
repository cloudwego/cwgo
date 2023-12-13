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
	"time"

	"github.com/cloudwego/cwgo/platform/server/cmd/agent/internal/svc"
	"github.com/cloudwego/cwgo/platform/server/shared/consts"
	"github.com/cloudwego/cwgo/platform/server/shared/errx"
	agent "github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/agent"
	"github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/model"
	"github.com/cloudwego/cwgo/platform/server/shared/utils"
)

type AddTokenService struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
} // NewAddTokenService new AddTokenService
func NewAddTokenService(ctx context.Context, svcCtx *svc.ServiceContext) *AddTokenService {
	return &AddTokenService{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Run create note info
func (s *AddTokenService) Run(req *agent.AddTokenReq) (resp *agent.AddTokenRes, err error) {
	var owner string
	var expirationTime time.Time

	switch req.RepositoryType {
	case consts.RepositoryTypeNumGitLab:
		client, err := utils.NewGitlabClient(req.Token, "https://"+req.RepositoryDomain)
		if err != nil {
			return &agent.AddTokenRes{
				Code: errx.GetCode(err),
				Msg:  err.Error(),
				Data: nil,
			}, nil
		}

		owner, expirationTime, err = utils.GetGitLabTokenInfo(client)
		if err != nil {
			return &agent.AddTokenRes{
				Code: errx.GetCode(err),
				Msg:  err.Error(),
				Data: nil,
			}, nil
		}
	case consts.RepositoryTypeNumGithub:
		client, err := utils.NewGithubClient(req.Token)
		if err != nil {
			return &agent.AddTokenRes{
				Code: errx.GetCode(err),
				Msg:  err.Error(),
				Data: nil,
			}, nil
		}

		owner, expirationTime, err = utils.GetGitHubTokenInfo(client)
		if err != nil {
			return &agent.AddTokenRes{
				Code: errx.GetCode(err),
				Msg:  err.Error(),
				Data: nil,
			}, nil
		}
	}

	tokenModel := model.Token{
		RepositoryType:   req.RepositoryType,
		RepositoryDomain: req.RepositoryDomain,
		Owner:            owner,
		Token:            req.Token,
		Status:           consts.TokenStatusNumValid,
		ExpirationTime:   expirationTime.String(),
	}

	_, err = s.svcCtx.DaoManager.Token.AddToken(s.ctx, tokenModel)
	if err != nil {
		return &agent.AddTokenRes{
			Code: consts.ErrNumDatabase,
			Msg:  consts.ErrMsgDatabase,
		}, nil
	}

	return &agent.AddTokenRes{
		Code: 0,
		Msg:  "add token successfully",
		Data: &agent.AddTokenResData{
			Owner:          owner,
			ExpirationTime: expirationTime.Format(time.DateTime),
		},
	}, nil
}
