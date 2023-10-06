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
	"github.com/bytedance/sonic"
	"github.com/cloudwego/cwgo/platform/server/cmd/agent/internal/svc"
	"github.com/cloudwego/cwgo/platform/server/cmd/agent/pkg/cron"
	agent "github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/agent"
	"github.com/cloudwego/cwgo/platform/server/shared/logger"
	"github.com/cloudwego/cwgo/platform/server/shared/task"
	"go.uber.org/zap"
	"net/http"
	"time"
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
	tasks := make([]*task.Task, 0, len(req.Tasks))

	for _, t := range req.Tasks {
		tp := task.Type(t.Type)
		switch tp {
		case task.SyncIdl:
			var data task.SyncIdlData
			err := sonic.Unmarshal([]byte(t.Data), &data)
			if err != nil {
				logger.Logger.Error("json unmarshal failed", zap.Error(err), zap.String("data", t.Data))
				return &agent.UpdateTasksRes{
					Code: http.StatusInternalServerError,
					Msg:  "internal err",
				}, nil
			}

			scheduleTime, _ := time.ParseDuration(t.ScheduleTime)

			tasks = append(tasks, &task.Task{
				Id:           t.Id,
				Type:         task.Type(t.Type),
				ScheduleTime: scheduleTime,
				Data:         data,
			})
		case task.SyncRepo:
			var data task.SyncRepoData
			err := sonic.Unmarshal([]byte(t.Data), &data)
			if err != nil {
				logger.Logger.Error("json unmarshal failed", zap.Error(err), zap.String("data", t.Data))
				return &agent.UpdateTasksRes{
					Code: http.StatusInternalServerError,
					Msg:  "internal err",
				}, nil
			}

			scheduleTime, _ := time.ParseDuration(t.ScheduleTime)

			tasks = append(tasks, &task.Task{
				Id:           t.Id,
				Type:         task.Type(t.Type),
				ScheduleTime: scheduleTime,
				Data:         data,
			})
		}
	}

	cron.CronInstance.Stop()

	cron.CronInstance.EmptyTask()

	for _, t := range tasks {
		cron.CronInstance.AddTask(t)
	}

	cron.CronInstance.Start()

	return &agent.UpdateTasksRes{
		Code: 0,
		Msg:  "update tasks successfully",
	}, nil
}
