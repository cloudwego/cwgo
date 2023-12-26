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

package client

import "github.com/cloudwego/cwgo/pkg/consts"

// related to service resolver
var (
	hzNilResolverFuncBody = "{\n\treturn\n}"

	hzAppendResolverAddrFunc = `func GetResolverAddress() []string {
		e := os.Getenv("GO_HERTZ_REGISTRY_[[ToUpper .ServiceName]]")
		if len(e) == 0 {
		  return []string{[[$lenSlice := len .ResolverAddress]][[range $key, $value := .ResolverAddress]]"[[$value]]"[[if eq $key (Sub $lenSlice 1)]][[else]], [[end]][[end]]}
	    }
	    return strings.Fields(e)
      }`

	hzNilAppendInitResolverFunc = `// If you do not use the service resolver function, do not edit this function.
	  // Otherwise, you can customize and modify it.
	  func initResolver(ops *[]Option) (err error) {
		return
	  }`

	hzAppendInitResolverFunc = `// If you do not use the service resolver function, do not edit this function.
	  // Otherwise, you can customize and modify it.
	  func initResolver(ops *[]Option) (err error) {
		{{if ne .ResolverName ""}}
		{{.ResolverBody}}
		{{else}}
		return
        {{end}}
	  }`

	hzCommonResolverImport = "github.com/cloudwego/hertz/pkg/app/middlewares/client/sd"

	hzCommonResolverBody = `*ops = append(*ops, WithHertzClientMiddleware(sd.Discovery(r)))
	return nil`

	hzEtcdClientImports = map[string]string{
		hzCommonResolverImport:                   "",
		"github.com/hertz-contrib/registry/etcd": "",
	}

	hzEtcdClient = `r, err := etcd.NewEtcdResolver(hzHttp.GetResolverAddress())
	if err != nil {
		return err
	}` + consts.LineBreak + hzCommonResolverBody

	hzNacosClientImports = map[string]string{
		hzCommonResolverImport:                    "",
		"github.com/hertz-contrib/registry/nacos": "",
	}

	hzNacosClient = `r, err := nacos.NewDefaultNacosResolver()
	if err != nil {
		return err
	}` + consts.LineBreak + hzCommonResolverBody

	hzConsulClientImports = map[string]string{
		hzCommonResolverImport:                     "",
		"github.com/hashicorp/consul/api":          "",
		"github.com/hertz-contrib/registry/consul": "",
	}

	hzConsulClient = `consulConfig := api.DefaultConfig()
    consulConfig.Address = hzHttp.GetResolverAddress()[0]
    consulClient, err := api.NewClient(consulConfig)
    if err != nil {
        return err
    }
    
    r := consul.NewConsulResolver(consulClient)` + consts.LineBreak + hzCommonResolverBody

	hzEurekaClientImports = map[string]string{
		hzCommonResolverImport:                     "",
		"github.com/hertz-contrib/registry/eureka": "",
	}

	hzEurekaClient = `r := eureka.NewEurekaResolver(hzHttp.GetResolverAddress())` +
		consts.LineBreak + hzCommonResolverBody

	hzPolarisClientImports = map[string]string{
		hzCommonResolverImport:                      "",
		"github.com/hertz-contrib/registry/polaris": "",
	}

	hzPolarisClient = `r, err := polaris.NewPolarisResolver()
    if err != nil {
        return err
    }` + consts.LineBreak + hzCommonResolverBody

	hzServiceCombClientImports = map[string]string{
		hzCommonResolverImport:                          "",
		"github.com/hertz-contrib/registry/servicecomb": "",
	}

	hzServiceCombClient = `r, err := servicecomb.NewDefaultSCResolver(hzHttp.GetResolverAddress())
    if err != nil {
        return err
    }` + consts.LineBreak + hzCommonResolverBody

	hzZKClientImports = map[string]string{
		hzCommonResolverImport:                        "",
		"github.com/hertz-contrib/registry/zookeeper": "",
		"time": "",
	}

	hzZKClient = `r, err := zookeeper.NewZookeeperResolver(hzHttp.GetResolverAddress(), 40*time.Second)
    if err != nil {
        return err
    }` + consts.LineBreak + hzCommonResolverBody
)
