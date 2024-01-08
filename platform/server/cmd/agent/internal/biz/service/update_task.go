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
	"github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/model"
	task "github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/task"
)

type UpdateTaskService struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
} // NewUpdateTaskService new UpdateTaskService
func NewUpdateTaskService(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTaskService {
	return &UpdateTaskService{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Run create note info
func (s *UpdateTaskService) Run(req *task.UpdateTaskReq) (resp *task.UpdateTaskResp, err error) {
	tasks := make([]model.Task, 0)
	for _, _task := range req.Tasks {
		tasks = append(tasks, *_task)
	}

	processor.Srv.UpdateTasks(tasks)
	return &task.UpdateTaskResp{
		Code: 0,
		Msg:  "update tasks successfully",
	}, nil
}
