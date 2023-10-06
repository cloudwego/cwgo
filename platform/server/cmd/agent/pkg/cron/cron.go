/*
 * Copyright 2022 CloudWeGo Authors
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
 */

package cron

import (
	"github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/agent"
	"github.com/cloudwego/cwgo/platform/server/shared/task"
	"github.com/go-co-op/gocron"
	"time"
)

type Cron struct {
	scheduler *gocron.Scheduler
	service   agent.AgentService
}

var CronInstance *Cron

func InitCron(service agent.AgentService) {
	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.TagsUnique()

	CronInstance = &Cron{
		scheduler: scheduler,
		service:   service,
	}
}

func (c *Cron) AddTask(t *task.Task) {
	switch t.Type {
	case task.SyncIdl:
		_, _ = c.scheduler.Every(t.ScheduleTime).Tag(t.Id).Do(func() {

		})
	case task.SyncRepo:
		_, _ = c.scheduler.Every(t.ScheduleTime).Tag(t.Id).Do(func() {

		})
	default:

	}
}

func (c *Cron) EmptyTask() {
	c.scheduler.Clear()
}
