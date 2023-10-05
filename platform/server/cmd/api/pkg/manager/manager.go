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
	"github.com/cloudwego/cwgo/platform/server/cmd/api/pkg/dispatcher"
	"github.com/cloudwego/cwgo/platform/server/shared/logger"
	"github.com/cloudwego/cwgo/platform/server/shared/registry"
	"github.com/cloudwego/cwgo/platform/server/shared/service"
	"go.uber.org/zap"
	"time"
)

type Manager struct {
	agents []*service.Service

	updateInterval time.Duration
	dispatcher     dispatcher.IDispatcher
	registry       registry.IRegistry
}

const (
	DefaultUpdateInterval = time.Second * 3
)

func NewManager(dispatcher dispatcher.IDispatcher, registry registry.IRegistry, updateInterval time.Duration) *Manager {
	return &Manager{
		agents:         make([]*service.Service, 0),
		updateInterval: updateInterval,
		dispatcher:     dispatcher,
		registry:       registry,
	}
}

func (m *Manager) StartUpdate() {
	for {
		time.Sleep(m.updateInterval)

		services, err := m.registry.GetAllService()
		if err != nil {
			logger.Logger.Error("get registry service failed", zap.Error(err))
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
				logger.Logger.Error("add service to dispatcher failed", zap.Error(err))
				continue
			}
		}

		for _, serviceId := range delServicesIds {
			err = m.dispatcher.DelService(serviceId)
			if err != nil {
				logger.Logger.Error("del service to dispatcher failed", zap.Error(err))
				continue
			}
		}

		m.agents = services
	}
}

func (m *Manager) GetDispatcher() dispatcher.IDispatcher {
	return m.dispatcher
}

func (m *Manager) GetRegistry() registry.IRegistry {
	return m.registry
}
