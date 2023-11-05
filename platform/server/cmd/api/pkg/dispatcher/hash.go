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

package dispatcher

import (
	"errors"
	"github.com/buraksezer/consistent"
	"github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/model"
	"hash/fnv"
	"sync"
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

type ConsistentHashDispatcher struct {
	mutex sync.Mutex

	hasher           *consistent.Consistent
	Tasks            map[string]*model.Task
	ServiceWithTasks map[string]map[string]*model.Task
}

func NewConsistentHashDispatcher() *ConsistentHashDispatcher {
	consistentHasher := consistent.New(
		nil,
		consistent.Config{
			Hasher:            hasher{},
			PartitionCount:    5000,
			ReplicationFactor: 5,
			Load:              1.25,
		},
	)

	return &ConsistentHashDispatcher{
		mutex:            sync.Mutex{},
		hasher:           consistentHasher,
		Tasks:            make(map[string]*model.Task),
		ServiceWithTasks: make(map[string]map[string]*model.Task),
	}
}

func (c *ConsistentHashDispatcher) AddService(serviceId string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.hasher.Add(member(serviceId))

	members := c.hasher.GetMembers()

	serviceWithTasks := make(map[string]map[string]*model.Task, len(members))

	for taskId, t := range c.Tasks {
		m := c.hasher.LocateKey([]byte(taskId)).String()
		_, ok := serviceWithTasks[m]
		if !ok {
			serviceWithTasks[m] = make(map[string]*model.Task)
		}

		serviceWithTasks[m][taskId] = t
	}
	c.ServiceWithTasks = serviceWithTasks

	return nil
}

func (c *ConsistentHashDispatcher) DelService(serviceId string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, ok := c.ServiceWithTasks[serviceId]; !ok {
		return errors.New("service not found")
	}

	c.hasher.Remove(serviceId)

	for taskId, t := range c.ServiceWithTasks[serviceId] {
		m := c.hasher.LocateKey([]byte(taskId)).String()
		c.ServiceWithTasks[m][taskId] = t
	}

	delete(c.ServiceWithTasks, serviceId)

	return nil
}

func (c *ConsistentHashDispatcher) AddTask(task *model.Task) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.Tasks[task.Id] = task
	m := c.hasher.LocateKey([]byte(task.Id))
	if m != nil {
		c.ServiceWithTasks[m.String()][task.Id] = task
	}

	return nil
}

func (c *ConsistentHashDispatcher) RemoveTask(taskId string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.Tasks, taskId)
	m := c.hasher.LocateKey([]byte(taskId))
	if m != nil {
		delete(c.ServiceWithTasks[m.String()], taskId)
	}

	return nil
}

func (c *ConsistentHashDispatcher) DelTask(taskId string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.Tasks, taskId)
	m := c.hasher.LocateKey([]byte(taskId))
	if m != nil {
		delete(c.ServiceWithTasks[m.String()], taskId)
	}

	return nil
}

func (c *ConsistentHashDispatcher) GetTaskByServiceId(serviceId string) []*model.Task {
	tasks := make([]*model.Task, 0, len(c.ServiceWithTasks[serviceId]))

	for _, t := range c.ServiceWithTasks[serviceId] {
		tasks = append(tasks, t)
	}

	return tasks
}

func (c *ConsistentHashDispatcher) GetServiceIdByTaskId(taskId string) string {
	serviceId := c.hasher.LocateKey([]byte(taskId))
	if serviceId != nil {
		return serviceId.String()
	}
	return ""
}

func (c *ConsistentHashDispatcher) GetTotalTaskNum() int {
	return len(c.Tasks)
}
