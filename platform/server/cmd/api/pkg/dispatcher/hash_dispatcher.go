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

package dispatcher

import (
	"errors"
	"hash/fnv"
	"sync"

	"github.com/cloudwego/cwgo/platform/server/shared/log"
	"go.uber.org/zap"

	"github.com/cloudwego/cwgo/platform/server/shared/consts"

	"github.com/buraksezer/consistent"
	"github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/model"
)

type hasher struct{}

func (h hasher) Sum64(data []byte) uint64 {
	hash := fnv.New64a()
	_, _ = hash.Write(data)
	return hash.Sum64()
}

type member string

func (m member) String() string {
	return string(m)
}

// TaskMap is a map of taskID to task
type TaskMap map[string]*model.Task

type ConsistentHashDispatcher struct {
	sync.Mutex

	hasher *consistent.Consistent
	// taskMap is a map of taskID to task
	taskMap          TaskMap
	service2TasksMap map[string]TaskMap
	// convertIDMap is used to convert idlID to taskID
	convertIDMap map[int64]string
}

func NewConsistentHashDispatcher() *ConsistentHashDispatcher {
	return &ConsistentHashDispatcher{
		Mutex: sync.Mutex{},
		hasher: consistent.New(nil,
			consistent.Config{
				Hasher:            hasher{},
				PartitionCount:    5000,
				ReplicationFactor: 5,
				Load:              1.25,
			},
		),
		taskMap:          make(map[string]*model.Task),
		service2TasksMap: make(map[string]TaskMap),
		convertIDMap:     make(map[int64]string),
	}
}

func (c *ConsistentHashDispatcher) AddService(serviceID string) error {
	c.Lock()
	defer c.Unlock()

	c.hasher.Add(member(serviceID))

	members := c.hasher.GetMembers()
	serviceTasksMap := make(map[string]TaskMap, len(members))
	for _, m := range c.hasher.GetMembers() {
		serviceTasksMap[m.String()] = make(map[string]*model.Task)
	}

	for taskId, t := range c.taskMap {
		srvID := c.hasher.LocateKey([]byte(taskId)).String()
		serviceTasksMap[srvID][taskId] = t
	}
	c.service2TasksMap = serviceTasksMap

	log.Debug("", zap.Reflect("service2TasksMap", c.service2TasksMap), zap.Reflect("taskMap", c.taskMap))
	return nil
}

func (c *ConsistentHashDispatcher) DelService(serviceID string) error {
	c.Lock()
	defer c.Unlock()

	if _, ok := c.service2TasksMap[serviceID]; !ok {
		return errors.New("service not found")
	}

	c.hasher.Remove(serviceID)

	for taskID, task := range c.service2TasksMap[serviceID] {
		memberKey := c.hasher.LocateKey([]byte(taskID)).String()
		if memberKey != "" {
			if taskMap, ok := c.service2TasksMap[memberKey]; ok {
				taskMap[taskID] = task
			}
		}
	}

	delete(c.service2TasksMap, serviceID)

	return nil
}

func (c *ConsistentHashDispatcher) AddTask(task *model.Task) error {
	c.Lock()
	defer c.Unlock()

	// taskID --- task information
	c.taskMap[task.ID] = task

	// idlID --- [taskID]
	switch task.Type {
	case consts.Sync:
		c.convertIDMap[task.IdlID] = task.ID
	}
	svrID := c.hasher.LocateKey([]byte(task.ID))
	if svrID != nil {
		c.service2TasksMap[svrID.String()][task.ID] = task
	}
	log.Debug("", zap.Reflect("service2TasksMap", c.service2TasksMap), zap.Reflect("taskMap", c.taskMap))

	return nil
}

func (c *ConsistentHashDispatcher) removeTaskByID(taskId string) error {
	c.Lock()
	defer c.Unlock()

	delete(c.taskMap, taskId)
	memberKey := c.hasher.LocateKey([]byte(taskId))
	if memberKey != nil {
		delete(c.service2TasksMap[memberKey.String()], taskId)
	}

	return nil
}

func (c *ConsistentHashDispatcher) RemoveTaskByIdlID(IdlID int64) error {
	c.Lock()
	defer c.Unlock()

	// convert idlID to taskID
	taskID, ok := c.convertIDMap[IdlID]
	if ok {
		delete(c.convertIDMap, IdlID)
		return c.removeTaskByID(taskID)
	}

	return nil
}

func (c *ConsistentHashDispatcher) DelTask(taskId string) error {
	c.Lock()
	defer c.Unlock()

	delete(c.taskMap, taskId)
	m := c.hasher.LocateKey([]byte(taskId))
	if m != nil {
		delete(c.service2TasksMap[m.String()], taskId)
	}

	return nil
}

func (c *ConsistentHashDispatcher) GetAllTasks() []*model.Task {
	c.Lock()
	defer c.Unlock()

	tasks := make([]*model.Task, 0, len(c.taskMap))
	for _, t := range c.taskMap {
		tasks = append(tasks, t)
	}

	return tasks
}

func (c *ConsistentHashDispatcher) GetTasksByServiceID(serviceID string) []*model.Task {
	tasks := make([]*model.Task, 0)
	for _, v := range c.service2TasksMap[serviceID] {
		tasks = append(tasks, v)
	}

	return tasks
}
