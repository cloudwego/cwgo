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
	"fmt"
	"github.com/cloudwego/cwgo/platform/server/shared/consts"
	"github.com/cloudwego/cwgo/platform/server/shared/registry"
	"github.com/cloudwego/cwgo/platform/server/shared/utils"
	"github.com/cloudwego/kitex/pkg/discovery"
	kitexregistry "github.com/cloudwego/kitex/pkg/registry"
)

type BuiltinRegistryConfig struct {
	Address string `mapstructure:"address"`
}

type BuiltinRegistryConfigManager struct {
	Config       BuiltinRegistryConfig
	RegistryType consts.RegistryType
	Registry     *registry.BuiltinRegistry
}

func NewBuiltinRegistryConfigManager(config BuiltinRegistryConfig, r *registry.BuiltinRegistry) (*BuiltinRegistryConfigManager, error) {
	return &BuiltinRegistryConfigManager{
		Config:       config,
		RegistryType: consts.RegistryTypeNumBuiltin,
		Registry:     r,
	}, nil
}

func (cm *BuiltinRegistryConfigManager) GetRegistryType() consts.RegistryType {
	return cm.RegistryType
}

func (cm *BuiltinRegistryConfigManager) GetKitexRegistry(serviceName, serviceId, addr string) (kitexregistry.Registry, *kitexregistry.Info) {
	registryClient, err := registry.NewBuiltinKitexRegistryClient(cm.Config.Address)
	if err != nil {
		panic(fmt.Sprintf("initialize builtin BuiltinRegistry client failed, err: %v", err))
	}

	registryInfo := &kitexregistry.Info{
		ServiceName: serviceName,
		Addr:        utils.NewNetAddr("tcp", addr),
		Tags: map[string]string{
			"service_id": serviceId,
		},
	}

	return registryClient, registryInfo
}

func (cm *BuiltinRegistryConfigManager) GetDiscoveryResolver() discovery.Resolver {
	resolver, err := registry.NewBuiltinRegistryResolver(cm.Registry)
	if err != nil {
		panic(fmt.Sprintf("initialize builtin BuiltinRegistry resolver failed, err: %v", err))
	}

	return resolver
}
