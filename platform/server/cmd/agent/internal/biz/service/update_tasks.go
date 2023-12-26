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
	"github.com/cloudwego/cwgo/platform/server/cmd/agent/pkg/processor"
	agent "github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/agent"
	"github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/model"
)

type UpdateTasksService struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
} // NewUpdateTasksService new UpdateTasksService
func NewUpdateTasksService(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTasksService {
	return &UpdateTasksService{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Run create note info
func (s *UpdateTasksService) Run(req *agent.UpdateTasksReq) (resp *agent.UpdateTasksRes, err error) {
	tasks := make([]model.Task, 0, len(req.Tasks))

	for _, t := range req.Tasks {
		tasks = append(tasks, model.Task{
			Id:           t.Id,
			Type:         t.Type,
			ScheduleTime: t.ScheduleTime,
			Data:         t.Data,
		})
	}

	processor.ProcessorInstance.UpdateTasks(tasks)

	return &agent.UpdateTasksRes{
		Code: 0,
		Msg:  "update tasks successfully",
	}, nil
}
