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

package agent

import (
	"fmt"
	"net"

	registryconfig "github.com/cloudwego/cwgo/platform/server/shared/config/internal/registry"
	"github.com/cloudwego/cwgo/platform/server/shared/config/store"
	"github.com/cloudwego/cwgo/platform/server/shared/consts"
	"github.com/cloudwego/cwgo/platform/server/shared/logger"
	"github.com/cloudwego/cwgo/platform/server/shared/utils"
	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	"go.uber.org/zap"
)

type ConfigManager struct {
	config                Config
	RegistryConfigManager registryconfig.IRegistryConfigManager
	ServiceId             string
	ServiceName           string
}

func NewConfigManager(config Config, registryConfig registryconfig.Config, storeConfig store.Config, serviceId string) *ConfigManager {
	var registryConfigManager registryconfig.IRegistryConfigManager
	var err error

	switch registryConfig.Type {
	case consts.RegistryTypeBuiltin:
		registryConfigManager, err = registryconfig.NewBuiltinRegistryConfigManager(registryConfig.Builtin, storeConfig)
		if err != nil {
			panic(fmt.Sprintf("initialize registry failed, err: %v", err))
		}
	default:
		panic("not support registryConfigType")
	}

	return &ConfigManager{
		config:                config,
		RegistryConfigManager: registryConfigManager,
		ServiceId:             serviceId,
		ServiceName:           fmt.Sprintf("%s-%s-%s", "cwgo", consts.ServerTypeAgent, serviceId),
	}
}

func (cm *ConfigManager) GetKitexServerOptions() []server.Option {
	var KitexServerOptions []server.Option
	addr, err := net.ResolveTCPAddr("tcp", cm.config.Addr)
	if err != nil {
		logger.Logger.Fatal("resolve tcp addr failed", zap.Error(err), zap.String("addr", cm.config.Addr))
	} else {
		KitexServerOptions = append(KitexServerOptions, server.WithServiceAddr(addr))
	}

	pubListenOn := utils.FigureOutListenOn(addr.String())

	kitexRegistry, kitexRegistryInfo := cm.RegistryConfigManager.GetKitexRegistry(
		cm.ServiceName,
		cm.ServiceId,
		pubListenOn,
	)

	KitexServerOptions = append(KitexServerOptions, server.WithRegistry(kitexRegistry))
	KitexServerOptions = append(KitexServerOptions, server.WithRegistryInfo(kitexRegistryInfo))

	KitexServerOptions = append(KitexServerOptions, server.WithLimit(&limit.Option{MaxConnections: 2000, MaxQPS: 500}))

	KitexServerOptions = append(KitexServerOptions, server.WithSuite(tracing.NewServerSuite()))
	KitexServerOptions = append(KitexServerOptions, server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
		ServiceName: cm.ServiceName,
	}))

	// thrift meta handler
	KitexServerOptions = append(KitexServerOptions, server.WithMetaHandler(transmeta.ServerTTHeaderHandler))

	return KitexServerOptions
}
