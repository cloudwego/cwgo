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

package processor

import (
	"context"

	"github.com/cloudwego/cwgo/platform/server/shared/consts"
	"github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/agent"
	"github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/model"
	"github.com/cloudwego/cwgo/platform/server/shared/log"
	"go.uber.org/zap"
)

type taskChan chan model.Task

type worker struct {
	srv        agent.AgentService
	taskChan   taskChan
	workerChan chan<- taskChan
	stopChan   chan any
}

func newWorker(srv agent.AgentService, workerChan chan taskChan) worker {
	return worker{
		srv:        srv,
		workerChan: workerChan,
		taskChan:   make(taskChan),
		stopChan:   make(chan any),
	}
}

func (w worker) start() {
	handleTask := func(_task model.Task) {
		switch _task.Type {
		case consts.Sync:
			log.Debug("process task",
				zap.Reflect("task", _task),
			)
			// sync idl generate code
			resp, err := w.srv.SyncIDLsById(context.Background(), &agent.SyncIDLsByIdReq{
				Ids: []int64{_task.IdlID},
			})
			if err != nil {
				log.Error("sync IDL generate code fail", zap.Error(err))
			}
			log.Debug("sync IDL generate code success", zap.Reflect("resp", resp))
		}
	}

	go func() {
		for {
			w.workerChan <- w.taskChan
			select {
			case _task := <-w.taskChan:
				handleTask(_task)
				w.workerChan <- w.taskChan
			case <-w.stopChan:
				return
			}
		}
	}()
}
