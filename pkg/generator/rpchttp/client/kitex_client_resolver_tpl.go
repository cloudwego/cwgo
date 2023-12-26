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
	kitexNilResolverFuncBody = "{\n\treturn\n}"

	kitexAppendResolverAddrFunc = `func GetResolverAddress() []string {
		e := os.Getenv("GO_KITEX_REGISTRY_[[ToUpper .ServiceName]]")
		if len(e) == 0 {
		  return []string{[[$lenSlice := len .ResolverAddress]][[range $key, $value := .ResolverAddress]]"[[$value]]"[[if eq $key (Sub $lenSlice 1)]][[else]], [[end]][[end]]}
	    }
	    return strings.Fields(e)
      }`

	kitexNilAppendInitResolverFunc = `// If you do not use the service resolver function, do not edit this function.
  // Otherwise, you can customize and modify it.
  func initResolver(options *[]client.Option) (err error) {
  	return
  }`

	kitexAppendInitResolverFunc = `// If you do not use the service resolver function, do not edit this function.
  // Otherwise, you can customize and modify it.
  func initResolver(options *[]client.Option) (err error) {
  	{{if ne .ResolverName ""}}
    {{.ResolverBody}}
    {{else}}
    return
    {{end}}
  }`

	kitexCommonResolverBody = `*options = append(*options, client.WithResolver(r))
	return nil`

	kitexEtcdClientImports = map[string]string{"github.com/kitex-contrib/registry-etcd": ""}

	kitexEtcdClient = `r, err := etcd.NewEtcdResolver(kitexRpc.GetResolverAddress())
	if err != nil {
		return err
	}` + consts.LineBreak + kitexCommonResolverBody

	kitexZKClientImports = map[string]string{
		"github.com/kitex-contrib/registry-zookeeper/resolver": "",
		"time": "",
	}

	kitexZKClient = `r, err := resolver.NewZookeeperResolver(kitexRpc.GetResolverAddress(), 40*time.Second)
    if err != nil {
		return err
    }` + consts.LineBreak + kitexCommonResolverBody

	kitexNacosClientImports = map[string]string{"github.com/kitex-contrib/registry-nacos/resolver": ""}

	kitexNacosClient = `r, err := resolver.NewDefaultNacosResolver()
	if err != nil {
		return err
	}` + consts.LineBreak + kitexCommonResolverBody

	kitexPolarisClientImports = map[string]string{
		"github.com/kitex-contrib/polaris": "",
	}

	kitexPolarisClient = `*options = append(*options, client.WithSuite(polaris.NewDefaultClientSuite()))
	return nil`

	kitexEurekaClientImports = map[string]string{"github.com/kitex-contrib/registry-eureka/resolver": ""}

	kitexEurekaClient = `r := resolver.NewEurekaResolver(kitexRpc.GetResolverAddress())` +
		consts.LineBreak + kitexCommonResolverBody

	kitexConsulClientImports = map[string]string{"github.com/kitex-contrib/registry-consul": ""}

	kitexConsulClient = `r, err := consul.NewConsulResolver(kitexRpc.GetResolverAddress()[0])
	if err != nil {
		return err
	}` + consts.LineBreak + kitexCommonResolverBody

	kitexServiceCombClientImports = map[string]string{"github.com/kitex-contrib/registry-servicecomb/resolver": ""}

	kitexServiceCombClient = `r, err := resolver.NewDefaultSCResolver()
    if err != nil {
        return err
    }` + consts.LineBreak + kitexCommonResolverBody
)
