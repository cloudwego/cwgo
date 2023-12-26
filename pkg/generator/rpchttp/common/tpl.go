/*
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

package common

var (
	EnvGoImports = map[string]string{
		"os":      "",
		"strings": "",
	}

	RegistryStructField = "Registry Registry `yaml:\"registry\"`"
	RegistryConfYaml    = `registry:
  address: {{range .RegistryAddress}}
	- {{.}}{{end}}`

	EtcdServerAddr        = []string{"127.0.0.1:2379"}
	NacosServerAddr       = []string{"127.0.0.1:8848"}
	ConsulServerAddr      = []string{"127.0.0.1:8500"}
	EurekaServerAddr      = []string{"http://127.0.0.1:8761/eureka"}
	PolarisServerAddr     = []string{"127.0.0.1:8090"}
	ServiceCombServerAddr = []string{"127.0.0.1:30100"}
	ZkServerAddr          = []string{"127.0.0.1:2181"}

	EtcdDocker = `  Etcd:
    image: 'bitnami/etcd:latest'
    ports:
      - "2379:2379"
      - "2380:2380"	`

	ZkDocker = `  zookeeper:
    image: zookeeper
    ports:
      - "2181:2181"`

	NacosDocker = `  nacos:
    image: nacos/nacos-server:latest
    ports:
      - "8848:8848"`

	PolarisDocker = `  polaris:
    image: polarismesh/polaris-server:latest
    ports:
      - "8090:8090"`

	EurekaDocker = `  eureka:
    image: 'xdockerh/eureka-server:latest'
    ports:
      - 8761:8761`

	ConsulDocker = `  consul:
    image: consul:latest
    ports:
      - "8500:8500"`

	ServiceCombDocker = `  service-center:
    image: 'servicecomb/service-center:latest'
    ports:
      - "30100:30100"`
)
