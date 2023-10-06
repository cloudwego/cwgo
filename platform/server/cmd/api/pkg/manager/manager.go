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

package manager

import (
	"context"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/cwgo/platform/server/cmd/api/pkg/dispatcher"
	"github.com/cloudwego/cwgo/platform/server/shared/consts"
	"github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/agent"
	"github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/agent/agentservice"
	"github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/base"
	"github.com/cloudwego/cwgo/platform/server/shared/registry"
	"github.com/cloudwego/cwgo/platform/server/shared/service"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/discovery"
	"sync"
	"time"
)

type Manager struct {
	agents []*service.Service

	updateInterval time.Duration
	dispatcher     dispatcher.IDispatcher
	registry       registry.IRegistry
	resolver       discovery.Resolver
}

const (
	DefaultUpdateInterval = time.Second * 3
)

func NewManager(dispatcher dispatcher.IDispatcher, registry registry.IRegistry, resolver discovery.Resolver, updateInterval time.Duration) *Manager {
	// TODO: init dispatcher tasks

	return &Manager{
		agents:         make([]*service.Service, 0),
		updateInterval: updateInterval,
		dispatcher:     dispatcher,
		registry:       registry,
		resolver:       resolver,
	}
}

func (m *Manager) StartUpdate() {
	for {
		time.Sleep(m.updateInterval)

		services, err := m.registry.GetAllService()
		if err != nil {
			continue
		}

		seta := make(map[string]struct{})
		setb := make(map[string]struct{})
		var addServiceIds, delServicesIds []string

		for _, svr := range m.agents {
			seta[svr.Id] = struct{}{}
		}
		for _, svr := range services {
			setb[svr.Id] = struct{}{}
		}

		for serviceId := range seta {
			if _, ok := setb[serviceId]; !ok {
				delServicesIds = append(delServicesIds, serviceId)
			}
		}

		for serviceId := range setb {
			if _, ok := seta[serviceId]; !ok {
				addServiceIds = append(addServiceIds, serviceId)
			}
		}

		for _, serviceId := range addServiceIds {
			err = m.dispatcher.AddService(serviceId)
			if err != nil {
				continue
			}
		}

		for _, serviceId := range delServicesIds {
			err = m.dispatcher.DelService(serviceId)
			if err != nil {
				continue
			}
		}

		m.agents = services

		if len(addServiceIds) != 0 || len(delServicesIds) != 0 {
			// service changed, update cron
			var wg sync.WaitGroup
			for _, svr := range m.agents {
				wg.Add(1)
				go func(serviceId string) {
					defer wg.Done()

					c, _ := agentservice.NewClient(
						consts.ProjectName+"-"+consts.ServerTypeAgent,
						client.WithResolver(m.resolver),
						client.WithTag("service_id", serviceId),
					)

					tasks := m.dispatcher.GetTaskByServiceId(serviceId)

					tasksModels := make([]*base.Task, len(tasks))
					for i, t := range tasks {
						data, _ := sonic.MarshalString(t.Data)
						tasksModels[i] = &base.Task{
							Id:           t.Id,
							Type:         int32(t.Type),
							ScheduleTime: t.ScheduleTime.String(),
							Data:         data,
						}
					}

					_, _ = c.UpdateTasks(context.Background(), &agent.UpdateTasksReq{Tasks: tasksModels})
				}(svr.Id)
			}
			wg.Wait()
		}
	}
}

func (m *Manager) GetDispatcher() dispatcher.IDispatcher {
	return m.dispatcher
}

func (m *Manager) GetRegistry() registry.IRegistry {
	return m.registry
}
