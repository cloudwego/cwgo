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

package dispatcher

import "github.com/cloudwego/cwgo/platform/server/shared/task"

type IDispatcher interface {
	AddService(serviceId string) error
	DelService(serviceId string) error

	AddTask(task *task.Task) error
	RemoveTask(taskId string) error

	GetTaskByServiceId(serviceId string) []*task.Task
	GetTotalTaskNum() int
}
