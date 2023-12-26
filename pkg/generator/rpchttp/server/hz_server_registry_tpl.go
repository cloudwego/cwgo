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

package server

import "github.com/cloudwego/cwgo/pkg/consts"

// related to service registration
var (
	hzNilRegistryFuncBody = "{\n\treturn\n}"

	hzAppendRegistryAddrFunc = `type Registry struct {
		Address []string  ` + "`yaml:\"address\"`" + `
      }

      func GetRegistryAddress() []string {
		e := os.Getenv("GO_HERTZ_REGISTRY_[[ToUpper .ServiceName]]")
		if len(e) == 0 {
		  if GetConf().Registry.Address != nil {
			return GetConf().Registry.Address
		  } else {
			return []string{[[$lenSlice := len .RegistryAddress]][[range $key, $value := .RegistryAddress]]"[[$value]]"[[if eq $key (Sub $lenSlice 1)]][[else]], [[end]][[end]]}
		  }
	    }
	    return strings.Fields(e)
      }`

	hzNilAppendInitRegistryFunc = `// If you do not use the service registry function, do not edit this function.
	  // Otherwise, you can customize and modify it.
	  func initRegistry(ops *[]config.Option) (err error) {
		return
	  }`

	hzAppendInitRegistryFunc = `// If you do not use the service registry function, do not edit this function.
	  // Otherwise, you can customize and modify it.
	  func initRegistry(ops *[]config.Option) (err error) {
		{{if ne .RegistryName ""}}
		{{.RegistryBody}}
		{{else}}
		return
        {{end}}
	  }`

	hzCommonRegistryBody = `*ops = append(*ops, server.WithRegistry(r, &registry.Info{
		ServiceName: "%s",
		Addr:        utils.NewNetAddr("tcp", conf.GetConf().Hertz.Address),
		Weight:      10,
		Tags:        nil,
	}))
    return nil`

	hzCommonRegistryImport = "github.com/cloudwego/hertz/pkg/app/server/registry"

	hzEtcdServerImports = map[string]string{
		hzCommonRegistryImport:                   "",
		"github.com/hertz-contrib/registry/etcd": "",
	}

	hzEtcdServer = `r, err := etcd.NewEtcdRegistry(conf.GetRegistryAddress())
	if err != nil {
		return err
	}` + consts.LineBreak + hzCommonRegistryBody

	hzNacosServerImports = map[string]string{
		hzCommonRegistryImport:                    "",
		"github.com/hertz-contrib/registry/nacos": "",
	}

	hzNacosServer = `r, err := nacos.NewDefaultNacosRegistry()
    if err != nil {
        return err
    }` + consts.LineBreak + hzCommonRegistryBody

	hzConsulServerImports = map[string]string{
		hzCommonRegistryImport:                     "",
		"github.com/hashicorp/consul/api":          "",
		"github.com/hertz-contrib/registry/consul": "",
	}

	hzConsulServer = `consulConfig := api.DefaultConfig()
    consulConfig.Address = conf.GetRegistryAddress()[0]
    consulClient, err := api.NewClient(consulConfig)
    if err != nil {
        return err
    }
    
    r := consul.NewConsulRegister(consulClient)` + consts.LineBreak + hzCommonRegistryBody

	hzEurekaServerImports = map[string]string{
		hzCommonRegistryImport:                     "",
		"github.com/hertz-contrib/registry/eureka": "",
		"time": "",
	}

	hzEurekaServer = `r := eureka.NewEurekaRegistry(conf.GetRegistryAddress(), 40*time.Second)` +
		consts.LineBreak + hzCommonRegistryBody

	hzPolarisServerImports = map[string]string{
		hzCommonRegistryImport:                      "",
		"github.com/hertz-contrib/registry/polaris": "",
	}

	hzPolarisServer = `r, err := polaris.NewPolarisRegistry()
    if err != nil {
        return err
    }
	*ops = append(*ops, server.WithRegistry(r, &registry.Info{
		ServiceName: conf.GetConf().Hertz.ServiceName,
		Addr:        utils.NewNetAddr("tcp", conf.GetConf().Hertz.Address),
		Tags: map[string]string{
            "namespace": "Polaris",
        },
	}))
	return nil`

	hzServiceCombServerImports = map[string]string{
		hzCommonRegistryImport:                          "",
		"github.com/hertz-contrib/registry/servicecomb": "",
	}

	hzServiceCombServer = `r, err := servicecomb.NewDefaultSCRegistry(conf.GetRegistryAddress())
    if err != nil {
        return err
    }` + consts.LineBreak + hzCommonRegistryBody

	hzZKServerImports = map[string]string{
		hzCommonRegistryImport:                        "",
		"github.com/hertz-contrib/registry/zookeeper": "",
		"time": "",
	}

	hzZKServer = `r, err := zookeeper.NewZookeeperRegistry(conf.GetRegistryAddress(), 40*time.Second)
    if err != nil {
        return err
    }` + consts.LineBreak + hzCommonRegistryBody
)
