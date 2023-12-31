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

	"github.com/cloudwego/kitex/pkg/remote/codec/thrift"

	"github.com/cloudwego/kitex/pkg/registry"

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
	tcpAddr := getTCPAddr(cm.config.Addr)
	kxRegistry, registryInfo := cm.getRegistryAndInfo()

	return []server.Option{
		server.WithServiceAddr(tcpAddr),
		server.WithRegistry(kxRegistry),
		// open frugal
		server.WithPayloadCodec(thrift.NewThriftCodecWithConfig(thrift.FrugalRead | thrift.FrugalWrite)),
		server.WithRegistryInfo(registryInfo),
		server.WithLimit(&limit.Option{MaxConnections: 2000, MaxQPS: 500}),
		server.WithSuite(tracing.NewServerSuite()),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: cm.ServiceName}),
		server.WithMetaHandler(transmeta.ServerTTHeaderHandler),
	}
}

// getTCPAddr function remains the same
func getTCPAddr(addr string) *net.TCPAddr {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		logger.Logger.Fatal("resolve tcp addr failed", zap.Error(err), zap.String("addr", addr))
	}
	return tcpAddr
}

// GetRegistryAndInfo extracts registry-related logic
func (cm *ConfigManager) getRegistryAndInfo() (registry.Registry, *registry.Info) {
	pubListenOn := utils.FigureOutListenOn(cm.config.Addr)
	return cm.RegistryConfigManager.GetKitexRegistry(cm.ServiceName, cm.ServiceId, pubListenOn)
}
