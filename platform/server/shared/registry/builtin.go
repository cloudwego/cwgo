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

package registry

import (
	"errors"
	"fmt"
	"github.com/cloudwego/cwgo/platform/server/shared/service"
	"sync"
	"time"
)

type BuiltinRegistry struct {
	sync.Mutex
	agents        map[string]*service.BuiltinService
	cleanInterval time.Duration
	manager       *Manager
}

type Manager struct {
	agents      []*service.BuiltinService
	currentSize int
	expireTime  time.Duration
	mutex       sync.Mutex
}

func (sw *Manager) add(agentService *service.BuiltinService, serviceNum int) {
	sw.mutex.Lock()
	defer sw.mutex.Unlock()

	if sw.currentSize < serviceNum {
		if sw.currentSize == cap(sw.agents) {
			newAgents := make([]*service.BuiltinService, cap(sw.agents)<<1)
			copy(newAgents, sw.agents)
			sw.agents = newAgents
		}

		sw.agents[sw.currentSize] = agentService
		sw.currentSize++
	} else {
		copy(sw.agents, sw.agents[serviceNum-sw.currentSize:])
		sw.agents[serviceNum-1] = agentService
	}
}

func (sw *Manager) getExpiredServiceIds() []string {
	sw.mutex.Lock()
	defer sw.mutex.Unlock()

	expiredServiceIds := make([]string, 0)
	for _, agentService := range sw.agents {
		if agentService.LastUpdateTime.Add(sw.expireTime).Before(time.Now()) {
			expiredServiceIds = append(expiredServiceIds, agentService.Id)
		} else {
			break
		}
	}

	return expiredServiceIds
}

var _ IRegistry = (*BuiltinRegistry)(nil)

const (
	minCleanInterval = 100 * time.Millisecond
)

func NewBuiltinRegistry() *BuiltinRegistry {
	registry := &BuiltinRegistry{
		Mutex:         sync.Mutex{},
		agents:        make(map[string]*service.BuiltinService),
		cleanInterval: 3 * time.Second,
		manager: &Manager{
			agents:      make([]*service.BuiltinService, 0),
			currentSize: 0,
			mutex:       sync.Mutex{},
			expireTime:  time.Minute,
		},
	}

	go registry.CleanUp()

	return registry
}

func (r *BuiltinRegistry) Register(serviceId string, host string, port int) error {
	r.Lock()
	defer r.Unlock()

	agentService, err := service.NewBuiltinService(serviceId, fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return err
	}

	r.agents[serviceId] = agentService

	r.manager.add(agentService, r.Count())

	return nil
}

func (r *BuiltinRegistry) Unregister(id string) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.agents[id]; !ok {
		return errors.New("service not found")
	}

	delete(r.agents, id)

	return nil
}

func (r *BuiltinRegistry) Update(serviceId string) error {
	r.Lock()
	defer r.Unlock()

	if agentService, ok := r.agents[serviceId]; !ok {
		return errors.New("service not found")
	} else {
		agentService.LastUpdateTime = time.Now()
		r.manager.add(agentService, r.Count())
		return nil
	}
}

func (r *BuiltinRegistry) CleanUp() {
	for {
		time.Sleep(r.cleanInterval)

		expiredServiceIds := r.manager.getExpiredServiceIds()

		r.Mutex.Lock()
		for _, serviceId := range expiredServiceIds {
			if _, ok := r.agents[serviceId]; ok {
				delete(r.agents, serviceId)
			}
		}
		r.Mutex.Unlock()
	}
}

func (r *BuiltinRegistry) Count() int {
	return len(r.agents)
}

func (r *BuiltinRegistry) GetServiceById(serviceId string) (service.IService, error) {
	if agentService, ok := r.agents[serviceId]; !ok {
		return nil, errors.New("service not found")
	} else {
		return agentService, nil
	}
}

func (r *BuiltinRegistry) ServiceExists(serviceId string) bool {
	_, ok := r.agents[serviceId]

	return ok
}
