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

package handler

import (
	"context"

	"github.com/cloudwego/cwgo/platform/server/cmd/agent/internal/biz/service"
	"github.com/cloudwego/cwgo/platform/server/cmd/agent/internal/svc"
	"github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/agent"
	"github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/task"
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
func (s *AgentServiceImpl) AddRepository(ctx context.Context, req *agent.AddRepositoryReq) (resp *agent.AddRepositoryResp, err error) {
	resp, err = service.NewAddRepositoryService(ctx, s.svcCtx).Run(req)

	return resp, err
}

// AddIDL implements the AgentServiceImpl interface.
func (s *AgentServiceImpl) AddIDL(ctx context.Context, req *agent.AddIDLReq) (resp *agent.AddIDLResp, err error) {
	resp, err = service.NewAddIDLService(ctx, s.svcCtx, s).Run(req)

	return resp, err
}

// SyncIDLsById implements the AgentServiceImpl interface.
func (s *AgentServiceImpl) SyncIDLsById(ctx context.Context, req *agent.SyncIDLsByIdReq) (resp *agent.SyncIDLsByIdResp, err error) {
	resp, err = service.NewSyncIDLsByIdService(ctx, s.svcCtx, s).Run(req)

	return resp, err
}

// GetRepositories implements the AgentServiceImpl interface.
func (s *AgentServiceImpl) GetRepositories(ctx context.Context, req *agent.GetRepositoriesReq) (resp *agent.GetRepositoriesResp, err error) {
	resp, err = service.NewGetRepositoriesService(ctx, s.svcCtx).Run(req)

	return resp, err
}

// DeleteIDL implements the AgentServiceImpl interface.
func (s *AgentServiceImpl) DeleteIDL(ctx context.Context, req *agent.DeleteIDLsReq) (resp *agent.DeleteIDLsResp, err error) {
	resp, err = service.NewDeleteIDLService(ctx, s.svcCtx).Run(req)

	return resp, err
}

// UpdateIDL implements the AgentServiceImpl interface.
func (s *AgentServiceImpl) UpdateIDL(ctx context.Context, req *agent.UpdateIDLReq) (resp *agent.UpdateIDLResp, err error) {
	resp, err = service.NewUpdateIDLService(ctx, s.svcCtx).Run(req)

	return resp, err
}

// GetIDLs implements the AgentServiceImpl interface.
func (s *AgentServiceImpl) GetIDLs(ctx context.Context, req *agent.GetIDLsReq) (resp *agent.GetIDLsResp, err error) {
	resp, err = service.NewGetIDLsService(ctx, s.svcCtx).Run(req)

	return resp, err
}

// DeleteRepositories implements the AgentServiceImpl interface.
func (s *AgentServiceImpl) DeleteRepositories(ctx context.Context, req *agent.DeleteRepositoriesReq) (resp *agent.DeleteRepositoriesResp, err error) {
	resp, err = service.NewDeleteRepositoriesService(ctx, s.svcCtx).Run(req)

	return resp, err
}

// UpdateRepository implements the AgentServiceImpl interface.
func (s *AgentServiceImpl) UpdateRepository(ctx context.Context, req *agent.UpdateRepositoryReq) (resp *agent.UpdateRepositoryResp, err error) {
	resp, err = service.NewUpdateRepositoryService(ctx, s.svcCtx).Run(req)

	return resp, err
}

// AddTemplate implements the AgentServiceImpl interface.
func (s *AgentServiceImpl) AddTemplate(ctx context.Context, req *agent.AddTemplateReq) (resp *agent.AddTemplateResp, err error) {
	resp, err = service.NewAddTemplateService(ctx, s.svcCtx).Run(req)

	return resp, err
}

// DeleteTemplate implements the AgentServiceImpl interface.
func (s *AgentServiceImpl) DeleteTemplate(ctx context.Context, req *agent.DeleteTemplateReq) (resp *agent.DeleteTemplateResp, err error) {
	resp, err = service.NewDeleteTemplateService(ctx, s.svcCtx).Run(req)

	return resp, err
}

// UpdateTemplate implements the AgentServiceImpl interface.
func (s *AgentServiceImpl) UpdateTemplate(ctx context.Context, req *agent.UpdateTemplateReq) (resp *agent.UpdateTemplateResp, err error) {
	resp, err = service.NewUpdateTemplateService(ctx, s.svcCtx).Run(req)

	return resp, err
}

// GetTemplates implements the AgentServiceImpl interface.
func (s *AgentServiceImpl) GetTemplates(ctx context.Context, req *agent.GetTemplatesReq) (resp *agent.GetTemplatesResp, err error) {
	resp, err = service.NewGetTemplatesService(ctx, s.svcCtx).Run(req)

	return resp, err
}

// AddTemplateItem implements the AgentServiceImpl interface.
func (s *AgentServiceImpl) AddTemplateItem(ctx context.Context, req *agent.AddTemplateItemReq) (resp *agent.AddTemplateItemResp, err error) {
	resp, err = service.NewAddTemplateItemService(ctx, s.svcCtx).Run(req)

	return resp, err
}

// DeleteTemplateItem implements the AgentServiceImpl interface.
func (s *AgentServiceImpl) DeleteTemplateItem(ctx context.Context, req *agent.DeleteTemplateItemReq) (resp *agent.DeleteTemplateItemResp, err error) {
	resp, err = service.NewDeleteTemplateItemService(ctx, s.svcCtx).Run(req)

	return resp, err
}

// UpdateTemplateItem implements the AgentServiceImpl interface.
func (s *AgentServiceImpl) UpdateTemplateItem(ctx context.Context, req *agent.UpdateTemplateItemReq) (resp *agent.UpdateTemplateItemResp, err error) {
	resp, err = service.NewUpdateTemplateItemService(ctx, s.svcCtx).Run(req)

	return resp, err
}

// GetTemplateItems implements the AgentServiceImpl interface.
func (s *AgentServiceImpl) GetTemplateItems(ctx context.Context, req *agent.GetTemplateItemsReq) (resp *agent.GetTemplateItemsResp, err error) {
	resp, err = service.NewGetTemplateItemsService(ctx, s.svcCtx).Run(req)

	return resp, err
}

// AddToken implements the AgentServiceImpl interface.
func (s *AgentServiceImpl) AddToken(ctx context.Context, req *agent.AddTokenReq) (resp *agent.AddTokenResp, err error) {
	resp, err = service.NewAddTokenService(ctx, s.svcCtx).Run(req)

	return resp, err
}

// DeleteToken implements the AgentServiceImpl interface.
func (s *AgentServiceImpl) DeleteToken(ctx context.Context, req *agent.DeleteTokenReq) (resp *agent.DeleteTokenResp, err error) {
	resp, err = service.NewDeleteTokenService(ctx, s.svcCtx).Run(req)

	return resp, err
}

// GetToken implements the AgentServiceImpl interface.
func (s *AgentServiceImpl) GetToken(ctx context.Context, req *agent.GetTokenReq) (resp *agent.GetTokenResp, err error) {
	resp, err = service.NewGetTokenService(ctx, s.svcCtx).Run(req)

	return resp, err
}

// Ping implements the AgentServiceImpl interface.
func (s *AgentServiceImpl) Ping(ctx context.Context, req *agent.PingReq) (resp *agent.PingResp, err error) {
	resp, err = service.NewPingService(ctx, s.svcCtx).Run(req)

	return resp, err
}

// UpdateTask implements the AgentServiceImpl interface.
func (s *AgentServiceImpl) UpdateTask(ctx context.Context, req *task.UpdateTaskReq) (resp *task.UpdateTaskResp, err error) {
	resp, err = service.NewUpdateTaskService(ctx, s.svcCtx).Run(req)

	return resp, err
}
