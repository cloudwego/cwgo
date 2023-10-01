/*
* Copyright (c) 2023 lanshan team. All rights reserved.
 */

package service

import (
	"context"
	"github.com/cloudwego/cwgo/platform/server/cmd/agent/internal/svc"
	"github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/agent"
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
	// Finish your business logic.

	return
}
