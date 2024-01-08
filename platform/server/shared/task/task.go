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

package task

import (
	"github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/model"
	"github.com/cloudwego/cwgo/platform/server/shared/utils"
)

func NewTask(tp model.TaskType, scheduleTime string, idlID int64) *model.Task {
	taskId, _ := utils.NewTaskId()
	return &model.Task{
		ID:           taskId,
		ScheduleTime: scheduleTime,
		Type:         tp,
		IdlID:        idlID,
	}
}

type Message struct {
	Command string     `json:"command"`
	Task    model.Task `json:"task"`
}

const (
	AddTask    = "add"
	DeleteTask = "del"
)
