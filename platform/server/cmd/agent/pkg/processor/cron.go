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

package processor

import (
	"context"
	"github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/agent"
	"github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/model"
	"time"
)

type Worker struct {
	service agent.AgentService // sync methods

	// worker pool (write only)
	// push current worker's task chan into worker poll when worker is available
	workerPool chan<- chan model.Task

	taskChan chan model.Task // task queue

	stopChan chan struct{} // exit signal
}

func NewWorker(service agent.AgentService, workerPool chan<- chan model.Task) Worker {
	return Worker{
		service:    service,
		workerPool: workerPool,
		taskChan:   make(chan model.Task),
		stopChan:   make(chan struct{}),
	}
}

func (w Worker) Start() {
	go func() {
		w.workerPool <- w.taskChan // push worker into pool

		for {
			select {
			case t := <-w.taskChan:
				// process task with certain task type by calling service method
				switch t.Type {
				case model.Type_sync_idl_data:
					_, _ = w.service.SyncIDLsById(context.Background(), &agent.SyncIDLsByIdReq{
						Ids: []int64{t.Data.SyncIdlData.IdlId},
					})
				case model.Type_sync_repo_data:
					_, _ = w.service.SyncRepositoryById(context.Background(), &agent.SyncRepositoryByIdReq{
						Ids: []int64{t.Data.SyncRepoData.RepositoryId},
					})
				default:

				}

				w.workerPool <- w.taskChan // push worker into pool after finishing task
			case <-w.stopChan:
				// return when get exit signal
				return
			}
		}
	}()
}

func (w Worker) Stop() {
	w.stopChan <- struct{}{} // send exit signal
}

type Processor struct {
	service agent.AgentService // sync methods

	taskList []model.Task // sync tasks

	// worker pool (read only)
	// get available worker's task chan
	// and push task into this chan
	workerPool chan chan model.Task
	workerList []Worker

	stopChan chan struct{} // stop signal
}

const (
	defaultWorkerNum = 3  // worker num at initialization
	maxWorkerNum     = 16 // max worker num

	// sync time that controls current worker num
	maxSyncTime        = 60 * time.Minute // all task max sync time
	minSyncTime        = 10 * time.Second // all task min sync time
	adjustTimeDuration = 1 * time.Minute
)

var ProcessorInstance *Processor

func InitProcessor(service agent.AgentService) {
	// create worker pool
	workerPool := make(chan chan model.Task, defaultWorkerNum)

	// create workers
	workerList := make([]Worker, defaultWorkerNum)
	for i := 0; i < defaultWorkerNum; i++ {
		worker := NewWorker(service, workerPool)
		workerList[i] = worker
		worker.Start()
	}

	ProcessorInstance = &Processor{
		service:    service,
		taskList:   nil,
		workerPool: workerPool,
		workerList: workerList,
		stopChan:   make(chan struct{}),
	}
}

// Start dispatch tasks from current task list
func (c *Processor) Start() {
	var startTime time.Time
	var taskProcessedNum int64
	go func() {
		// worker adjust
		for {
			time.Sleep(adjustTimeDuration)
			if c.taskList == nil || len(c.taskList) == 0 {
				continue
			}

			if time.Now().Sub(startTime).Nanoseconds()/taskProcessedNum*int64(len(c.taskList)) > maxSyncTime.Nanoseconds() {
				if len(c.workerList) <= maxWorkerNum {
					// add worker when sync time exceed the max sync time
					worker := NewWorker(c.service, c.workerPool)
					c.workerList = append(c.workerList, worker)
					worker.Start()
				}
			}
		}
	}()
	go func() {
		// send task
		startTime = time.Now()
		for {
			for _, t := range c.taskList {
				select {
				case <-c.stopChan:
					// exit when get signal
					goto exit
				default:
					taskChan := <-c.workerPool // get available worker's task chan
					taskChan <- t              // push task into task chan
					taskProcessedNum++
				}
			}
			if time.Now().Sub(startTime) < minSyncTime && len(c.workerList) > defaultWorkerNum {
				// reduce worker
				c.workerList[0].Stop()
				c.workerList = c.workerList[1:]
			}
		}
	exit:
	}()
}

func (c *Processor) Stop() {
	// send exit signal
	c.stopChan <- struct{}{}
}

func (c *Processor) UpdateTasks(tasks []model.Task) {
	c.Stop()
	c.taskList = tasks // replace task list
	c.Start()
}
