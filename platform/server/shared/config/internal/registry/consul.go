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
	"github.com/cloudwego/kitex/pkg/discovery"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/go-cleanhttp"
)

type ConsulRegistryConfig struct {
	Address string
	Scheme  string
	Token   string
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

func (cm *ConsulRegistryConfigManager) init() {

}
