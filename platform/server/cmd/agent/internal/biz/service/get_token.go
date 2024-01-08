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

	"github.com/cloudwego/cwgo/platform/server/cmd/agent/internal/svc"
	"github.com/cloudwego/cwgo/platform/server/shared/consts"
	agent "github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/agent"
	"github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/model"
)

type GetTokenService struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
} // NewGetTokenService new GetTokenService
func NewGetTokenService(ctx context.Context, svcCtx *svc.ServiceContext) *GetTokenService {
	return &GetTokenService{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Run create note info
func (s *GetTokenService) Run(req *agent.GetTokenReq) (resp *agent.GetTokenResp, err error) {
	tokens, total, err := s.svcCtx.DaoManager.Token.GetTokenList(s.ctx,
		model.Token{
			RepositoryType:   req.RepositoryType,
			RepositoryDomain: req.RepositoryDomain,
			Owner:            req.Owner,
		},
		req.Page, req.Limit, req.Order, req.OrderBy,
	)
	if err != nil {
		return &agent.GetTokenResp{
			Code: consts.ErrNumDatabase,
			Msg:  consts.ErrMsgDatabase,
			Data: nil,
		}, nil
	}

	return &agent.GetTokenResp{
		Code: 0,
		Msg:  "get tokens successfully",
		Data: &agent.GetTokenRespData{
			Tokens: tokens,
			Total:  int32(total),
		},
	}, nil
}
