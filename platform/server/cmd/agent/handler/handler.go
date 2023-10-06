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

package handler

import (
	"context"
	"github.com/cloudwego/cwgo/platform/server/cmd/agent/internal/biz/service"
	"github.com/cloudwego/cwgo/platform/server/cmd/agent/internal/svc"
	"github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/agent"
)

// AgentServiceImpl implements the last service interface defined in the IDL.
type AgentServiceImpl struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAgentServiceImpl(ctx context.Context, svcCtx *svc.ServiceContext) *AgentServiceImpl {
	return &AgentServiceImpl{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// AddRepository implements the AgentServiceImpl interface.
func (s *AgentServiceImpl) AddRepository(ctx context.Context, req *agent.AddRepositoryReq) (resp *agent.AddRepositoryRes, err error) {
	resp, err = service.NewAddRepositoryService(ctx, s.svcCtx).Run(req)

	return resp, err
}

// UpdateRepositoryStatus implements the AgentServiceImpl interface.
func (s *AgentServiceImpl) UpdateRepositoryStatus(ctx context.Context, req *agent.UpdateRepositoryStatusReq) (resp *agent.UpdateRepositoryStatusRes, err error) {
	resp, err = service.NewUpdateRepositoryStatusService(ctx, s.svcCtx).Run(req)

	return resp, err
}

// SyncRepositoryById implements the AgentServiceImpl interface.
func (s *AgentServiceImpl) SyncRepositoryById(ctx context.Context, req *agent.SyncRepositoryByIdReq) (resp *agent.SyncRepositoryByIdRes, err error) {
	resp, err = service.NewSyncRepositoryByIdService(ctx, s.svcCtx).Run(req)

	return resp, err
}

// AddIDL implements the AgentServiceImpl interface.
func (s *AgentServiceImpl) AddIDL(ctx context.Context, req *agent.AddIDLReq) (resp *agent.AddIDLRes, err error) {
	resp, err = service.NewAddIDLService(ctx, s.svcCtx).Run(req)

	return resp, err
}

// SyncIDLsById implements the AgentServiceImpl interface.
func (s *AgentServiceImpl) SyncIDLsById(ctx context.Context, req *agent.SyncIDLsByIdReq) (resp *agent.SyncIDLsByIdRes, err error) {
	resp, err = service.NewSyncIDLsByIdService(ctx, s.svcCtx).Run(req)

	return resp, err
}

// UpdateTasks implements the AgentServiceImpl interface.
func (s *AgentServiceImpl) UpdateTasks(ctx context.Context, req *agent.UpdateTasksReq) (resp *agent.UpdateTasksRes, err error) {
	resp, err = service.NewUpdateTasksService(ctx, s.svcCtx).Run(req)

	return resp, err
}

// GenerateCode implements the AgentServiceImpl interface.
func (s *AgentServiceImpl) GenerateCode(ctx context.Context, req *agent.GenerateCodeReq) (resp *agent.GenerateCodeRes, err error) {
	resp, err = service.NewGenerateCodeService(ctx, s.svcCtx).Run(req)

	return resp, err
}
