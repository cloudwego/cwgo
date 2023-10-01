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
	"github.com/cloudwego/cwgo/platform/server/shared/logger"
	"github.com/cloudwego/cwgo/platform/server/shared/utils"
	hertzregistry "github.com/cloudwego/hertz/pkg/app/server/registry"
	"github.com/cloudwego/kitex/pkg/discovery"
	kitexregistry "github.com/cloudwego/kitex/pkg/registry"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/go-cleanhttp"
	hertzconsul "github.com/hertz-contrib/registry/consul"
	kitexconsul "github.com/kitex-contrib/registry-consul"
	"go.uber.org/zap"
	"net"
	"strconv"
)

type ConsulRegistryConfig struct {
	Address string `mapstructure:"address"`
	Token   string `mapstructure:"token"`
}

func (c ConsulRegistryConfig) GetConsulApiConfig() *consulapi.Config {
	return &consulapi.Config{
		Address:    c.Address,
		Scheme:     c.Address,
		PathPrefix: "",
		Datacenter: "",
		Transport:  cleanhttp.DefaultTransport(),
		HttpClient: nil,
		HttpAuth:   nil,
		WaitTime:   0,
		Token:      c.Token,
		TokenFile:  "",
		Namespace:  "",
		Partition:  "",
		TLSConfig:  consulapi.TLSConfig{},
	}
}

type ConsulRegistryConfigManager struct {
	consulApiConfig *consulapi.Config
	consulClient    *consulapi.Client
	consulResolver  discovery.Resolver
}

func NewConsulRegistryConfigManager(config ConsulRegistryConfig) (*ConsulRegistryConfigManager, error) {
	consulApiConfig := &consulapi.Config{
		Address:    config.Address,
		Scheme:     config.Address,
		PathPrefix: "",
		Datacenter: "",
		Transport:  cleanhttp.DefaultTransport(),
		HttpClient: nil,
		HttpAuth:   nil,
		WaitTime:   0,
		Token:      config.Token,
		TokenFile:  "",
		Namespace:  "",
		Partition:  "",
		TLSConfig:  consulapi.TLSConfig{},
	}

	consulClient, err := consulapi.NewClient(consulApiConfig)
	if err != nil {
		return nil, fmt.Errorf("initialize consul client failed, err: %v", err)
	}

	consulResolver, err := kitexconsul.NewConsulResolverWithConfig(consulApiConfig)
	if err != nil {
		return nil, fmt.Errorf("initialize consul resolver failed, err: %v", err)
	}

	return &ConsulRegistryConfigManager{
		consulApiConfig: consulApiConfig,
		consulClient:    consulClient,
		consulResolver:  consulResolver,
	}, nil
}

func (cm *ConsulRegistryConfigManager) GetHertzRegistry(serviceId int, serviceName string, Host string, port int) (hertzregistry.Registry, *hertzregistry.Info) {
	registry := hertzconsul.NewConsulRegister(
		cm.consulClient,
		hertzconsul.WithCheck(&consulapi.AgentServiceCheck{
			CheckID:                        "",
			Name:                           "",
			Args:                           nil,
			DockerContainerID:              "",
			Shell:                          "",
			Interval:                       "7s",
			Timeout:                        "5s",
			TTL:                            "",
			HTTP:                           "",
			Header:                         nil,
			Method:                         "",
			Body:                           "",
			TCP:                            "",
			UDP:                            "",
			Status:                         "",
			Notes:                          "",
			TLSServerName:                  "",
			TLSSkipVerify:                  false,
			GRPC:                           "",
			GRPCUseTLS:                     false,
			H2PING:                         "",
			H2PingUseTLS:                   false,
			AliasNode:                      "",
			AliasService:                   "",
			SuccessBeforePassing:           0,
			FailuresBeforeWarning:          0,
			FailuresBeforeCritical:         0,
			DeregisterCriticalServiceAfter: "15s",
		}),
	)

	registryInfo := &hertzregistry.Info{
		ServiceName: serviceName,
		Addr:        utils.NewNetAddr("tcp", net.JoinHostPort(Host, strconv.Itoa(port))),
		Weight:      hertzregistry.DefaultWeight,
		Tags: map[string]string{
			"service_id": strconv.Itoa(serviceId),
		},
	}

	return registry, registryInfo
}

func (cm *ConsulRegistryConfigManager) GetKitexRegistry(serviceId int, serviceName string, addr string) (kitexregistry.Registry, *kitexregistry.Info) {
	registry, err := kitexconsul.NewConsulRegisterWithConfig(
		cm.consulApiConfig,
		kitexconsul.WithCheck(&consulapi.AgentServiceCheck{
			CheckID:                        "",
			Name:                           "",
			Args:                           nil,
			DockerContainerID:              "",
			Shell:                          "",
			Interval:                       "7s",
			Timeout:                        "5s",
			TTL:                            "",
			HTTP:                           "",
			Header:                         nil,
			Method:                         "",
			Body:                           "",
			TCP:                            "",
			UDP:                            "",
			Status:                         "",
			Notes:                          "",
			TLSServerName:                  "",
			TLSSkipVerify:                  false,
			GRPC:                           "",
			GRPCUseTLS:                     false,
			H2PING:                         "",
			H2PingUseTLS:                   false,
			AliasNode:                      "",
			AliasService:                   "",
			SuccessBeforePassing:           0,
			FailuresBeforeWarning:          0,
			FailuresBeforeCritical:         0,
			DeregisterCriticalServiceAfter: "15s",
		}),
	)
	if err != nil {
		logger.Logger.Fatal("initialize consul registry failed.", zap.Error(err))
	}

	registryInfo := &kitexregistry.Info{
		ServiceName: serviceName,
		Addr:        utils.NewNetAddr("tcp", addr),
		Tags: map[string]string{
			"service_id": strconv.Itoa(serviceId),
		},
	}

	return registry, registryInfo
}
