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
	"github.com/cloudwego/cwgo/platform/server/shared/task"
	"github.com/go-co-op/gocron"
	"time"
)

type ICron interface {
	AddTask()
	DeleteTask()
	GetTasks()
}

type Cron struct {
	scheduler *gocron.Scheduler
}

func NewCron() *Cron {
	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.TagsUnique()

	return &Cron{
		scheduler: scheduler,
	}
}

func (c *Cron) AddTask(t *task.Task) {
	switch t.Type {
	case task.SyncRepo:
		c.scheduler.Every(t.ScheduleTime).Tag(t.Id).Do(func() {

		})
	default:

	}
}
