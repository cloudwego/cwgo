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
 */

package config

import (
	"fmt"
	"log"
	"net"

	"github.com/cloudwego/cwgo/platform/server/shared/consts"

	"github.com/cloudwego/cwgo/platform/server/shared/utils"
	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/registry"
	"github.com/cloudwego/kitex/pkg/remote/codec/thrift"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	"go.uber.org/zap"
)

type AgentConfig struct {
	Addr           string `mapstructure:"addr"`
	MaxConnections int64  `mapstructure:"maxConnections"`
	MaxQPS         int64  `mapstructure:"maxQPS"`
	WorkerNum      int    `mapstructure:"workerNum"`
}

type Metadata struct {
	ServiceId string `yaml:"service_id"`
}

func (conf *AgentConfig) Init() {
	if conf.Addr == "" {
		conf.Addr = "0.0.0.0:11010"
	}

	if conf.MaxConnections == 0 {
		conf.MaxConnections = 2000
	}

	if conf.MaxQPS == 0 {
		conf.MaxQPS = 500
	}
}

type AgentManager struct {
	config                AgentConfig
	RegistryConfigManager IRegistryConfigManager
	ServiceId             string
	ServiceName           string
}

func NewAgentManager(config AgentConfig, registryConfig RegistryConfig, storeConfig StoreConfig, serviceId string) *AgentManager {
	var registryConfigManager IRegistryConfigManager
	var err error

	switch registryConfig.Type {
	case consts.RegistryTypeBuiltin:
		registryConfigManager, err = NewBuiltinRegistryConfigManager(registryConfig.Builtin, storeConfig)
		if err != nil {
			panic(fmt.Sprintf("initialize registry failed, err: %v", err))
		}
	default:
		panic("not support registryConfigType")
	}

	return &AgentManager{
		config:                config,
		RegistryConfigManager: registryConfigManager,
		ServiceId:             serviceId,
		ServiceName:           fmt.Sprintf("%s-%s-%s", "cwgo", consts.ServerTypeAgent, serviceId),
	}
}

func (cm *AgentManager) GetKitexServerOptions() []server.Option {
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
		log.Fatal("resolve tcp addr failed", zap.Error(err), zap.String("addr", addr))
	}
	return tcpAddr
}

// GetRegistryAndInfo extracts registry-related logic
func (cm *AgentManager) getRegistryAndInfo() (registry.Registry, *registry.Info) {
	pubListenOn := utils.FigureOutListenOn(cm.config.Addr)
	return cm.RegistryConfigManager.GetKitexRegistry(cm.ServiceName, cm.ServiceId, pubListenOn)
}
