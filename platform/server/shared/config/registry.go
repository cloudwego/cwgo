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

package config

import (
	"github.com/cloudwego/cwgo/platform/server/shared/log"
	"go.uber.org/zap"

	"github.com/cloudwego/cwgo/platform/server/shared/consts"
	"github.com/cloudwego/cwgo/platform/server/shared/registry"
	"github.com/cloudwego/cwgo/platform/server/shared/utils"
	"github.com/cloudwego/kitex/pkg/discovery"
	kitexregistry "github.com/cloudwego/kitex/pkg/registry"
)

type IRegistryConfigManager interface {
	GetRegistryType() consts.RegistryType
	GetRegistry() registry.IRegistry
	NewKitexRegistryInfo(serviceName, serviceId, addr string) *kitexregistry.Info
	GetResolver() discovery.Resolver
}

type RegistryConfig struct {
	Type                string `mapstructure:"type"`
	RedisRegistryConfig `mapstructure:"redis"`
}

type RedisRegistryConfig struct{}

func (c *RegistryConfig) Init() {
	if c.Type == "" {
		c.Type = consts.RegistryTypeRedis
	}
}

type RedisRegistryManager struct {
	Config       RedisRegistryConfig
	storeConfig  StoreConfig
	RegistryType consts.RegistryType
	Registry     registry.IRegistry
}

func NewRedisRegistryManager(config RedisRegistryConfig, storeConfig StoreConfig) (*RedisRegistryManager, error) {
	m := &RedisRegistryManager{
		Config:       config,
		storeConfig:  storeConfig,
		RegistryType: consts.RegistryTypeNumRedis,
		Registry:     nil,
	}
	rdb, err := m.storeConfig.NewRedisClient()
	if err != nil {
		return nil, err
	}
	m.Registry = registry.NewRedisRegistry(rdb)
	return m, nil
}

func (m *RedisRegistryManager) GetRegistryType() consts.RegistryType {
	return consts.RegistryTypeNumRedis
}

func (m *RedisRegistryManager) GetRegistry() registry.IRegistry {
	return m.Registry
}

func (m *RedisRegistryManager) NewKitexRegistryInfo(serviceName, serviceID, addr string) *kitexregistry.Info {
	netAddr := utils.NewNetAddr("tcp", addr)
	registryInfo := &kitexregistry.Info{
		ServiceName: consts.ServiceNameAgent,
		Addr:        netAddr,
		Tags: map[string]string{
			"service_id": serviceID,
		},
	}
	return registryInfo
}

func (m *RedisRegistryManager) GetResolver() discovery.Resolver {
	rdb, err := m.storeConfig.NewRedisClient()
	if err != nil {
		log.Fatal("init redis failed", zap.Error(err))
	}
	return registry.NewRedisResolver(rdb)
}
