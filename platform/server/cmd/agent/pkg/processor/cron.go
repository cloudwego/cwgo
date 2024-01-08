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
	"time"

	"github.com/cloudwego/cwgo/platform/server/shared/config"
	"github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/agent"
	"github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/model"
)

const (
	maxWorkerNum = 16               // max worker num
	minSyncTime  = 30 * time.Second // all task min sync time
)

type Processor struct {
	srv        agent.AgentService
	tasks      []model.Task
	workers    []worker
	workerChan chan taskChan
	stopChan   chan any
}

var Srv *Processor

func InitProcessor(srv agent.AgentService) {
	workerNum := config.GetManager().Config.Agent.WorkerNum
	if workerNum > maxWorkerNum {
		workerNum = maxWorkerNum
	}
	workerChan := make(chan taskChan, workerNum)
	workers := make([]worker, workerNum)
	for i := 0; i < workerNum; i++ {
		w := newWorker(srv, workerChan)
		workers = append(workers, w)
		w.start()
	}

	Srv = &Processor{
		srv:        srv,
		workers:    workers,
		workerChan: workerChan,
		stopChan:   make(chan any),
	}
}

func (p *Processor) start() {
	var startTime time.Time

	go func() {
		// send task
		for {
			startTime = time.Now()

			for _, t := range p.tasks {
				select {
				case <-p.stopChan:
					// exit when get signal
					goto exit
				default:
					_taskChan := <-p.workerChan // get available worker's task chan
					_taskChan <- t              // push task into task chan
				}
			}

			if time.Since(startTime) < minSyncTime {
				time.Sleep(minSyncTime - time.Since(startTime))
			}
		}
	exit:
	}()
}

func (p *Processor) stop() {
	// send exit signal
	p.stopChan <- struct{}{}
}

func (p *Processor) UpdateTasks(tasks []model.Task) {
	if len(p.tasks) != 0 {
		p.stop()
	}

	p.tasks = tasks // replace task list

	if len(p.tasks) != 0 {
		p.start()
	}
}
