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
	kitexNilRegistryFuncBody = "{\n\treturn\n}"

	kitexNewAppendRegistryAddrFunc = `func GetRegistryAddress() []string {
		e := os.Getenv("GO_KITEX_REGISTRY_[[ToUpper .ServiceName]]")
		if len(e) == 0 {
		  if GetConf().Registry.RegistryAddress != nil {
			return GetConf().Registry.RegistryAddress
		  } else {
			return []string{[[$lenSlice := len .RegistryAddress]][[range $key, $value := .RegistryAddress]]"[[$value]]"[[if eq $key (Sub $lenSlice 1)]][[else]], [[end]][[end]]}
		  }
	    }
	    return strings.Fields(e)
      }`

	kitexAppendRegistryAddrFunc = `type Registry struct {
        Address []string  ` + "`yaml:\"address\"`" + `
      }  
	
      func GetRegistryAddress() []string {
		e := os.Getenv("GO_KITEX_REGISTRY_[[ToUpper .ServiceName]]")
		if len(e) == 0 {
		  if GetConf().Registry.Address != nil {
			return GetConf().Registry.Address
		  } else {
			return []string{[[$lenSlice := len .RegistryAddress]][[range $key, $value := .RegistryAddress]]"[[$value]]"[[if eq $key (Sub $lenSlice 1)]][[else]], [[end]][[end]]}
		  }
	    }
	    return strings.Fields(e)
      }`

	kitexNilAppendInitRegistryFunc = `// If you do not use the service registry function, do not edit this function.
  // Otherwise, you can customize and modify it.
  func initRegistry(ops *[]server.Option) (err error) {
	return
  }`

	kitexAppendInitRegistryFunc = `// If you do not use the service registry function, do not edit this function.
  // Otherwise, you can customize and modify it.
  func initRegistry(ops *[]server.Option) (err error) {
	{{if ne .RegistryName ""}}
		{{.RegistryBody}}
		{{else}}
		return
        {{end}}
  }`

	kitexCommonRegisterBody = `*ops = append(*ops, server.WithRegistry(r))
	return nil`

	kitexEtcdServerImports = map[string]string{"github.com/kitex-contrib/registry-etcd": ""}

	kitexEtcdServer = `r, err := etcd.NewEtcdRegistry(conf.GetRegistryAddress())
	if err != nil {
		return err
	}` + consts.LineBreak + kitexCommonRegisterBody

	kitexZKServerImports = map[string]string{
		"github.com/kitex-contrib/registry-zookeeper/registry": "",
		"time": "",
	}

	kitexZKServer = `r, err := registry.NewZookeeperRegistry(conf.GetRegistryAddress(), 40*time.Second)
    if err != nil{
        return err
    }` + consts.LineBreak + kitexCommonRegisterBody

	kitexNacosServerImports = map[string]string{"github.com/kitex-contrib/registry-nacos/registry": ""}

	kitexNacosServer = `r, err := registry.NewDefaultNacosRegistry()
	if err != nil {
		return err
	}` + consts.LineBreak + kitexCommonRegisterBody

	kitexPolarisServerImports = map[string]string{
		"github.com/kitex-contrib/polaris":        "",
		"github.com/cloudwego/kitex/pkg/registry": "",
	}

	kitexPolarisServer = `so := polaris.ServerOptions{}
	r, err := polaris.NewPolarisRegistry(so)
	if err != nil {
		return err
	}
	info := &registry.Info{
		ServiceName: conf.GetConf().Kitex.ServiceName,
		Tags: map[string]string{
			"namespace": "Polaris",
		},
	}
	*ops = append(*ops, server.WithRegistry(r), server.WithRegistryInfo(info))
	return nil`

	kitexEurekaServerImports = map[string]string{
		"github.com/kitex-contrib/registry-eureka/registry": "",
		"time": "",
	}

	kitexEurekaServer = `r := registry.NewEurekaRegistry(conf.GetRegistryAddress(), 15*time.Second)` +
		consts.LineBreak + kitexCommonRegisterBody

	kitexConsulServerImports = map[string]string{
		"github.com/kitex-contrib/registry-consul": "",
		"github.com/cloudwego/kitex/pkg/registry":  "",
	}

	kitexConsulServer = `r, err := consul.NewConsulRegister(conf.GetRegistryAddress()[0])
	if err != nil {
		return err
	}
	info := &registry.Info{
		ServiceName: conf.GetConf().Kitex.ServiceName,
		Weight:      1, // weights must be greater than 0 in consul,else received error and exit.
	}
	*ops = append(*ops, server.WithRegistry(r), server.WithRegistryInfo(info))
	return nil`

	kitexServiceCombServerImports = map[string]string{"github.com/kitex-contrib/registry-servicecomb/registry": ""}

	kitexServiceCombServer = `r, err := registry.NewDefaultSCRegistry()
    if err != nil {
        return err
    }` + consts.LineBreak + kitexCommonRegisterBody
)
