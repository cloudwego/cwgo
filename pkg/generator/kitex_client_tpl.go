/*
 * Copyright 2022 CloudWeGo Authors
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

package generator

import (
	"github.com/cloudwego/cwgo/pkg/consts"
)

// related to service resolver
var (
	kitexNilResolverFuncBody = "{\n\treturn\n}"

	kitexAppendResolverFunc = `func GetResolverAddress() []string {
		e := os.Getenv("GO_KITEX_RESOLVER_[[ToUpper .ServiceName]]")
		if len(e) == 0 {
		  return []string{[[$lenSlice := len .ResolverAddress]][[range $key, $value := .ResolverAddress]]"[[$value]]"[[if eq $key (Sub $lenSlice 1)]][[else]], [[end]][[end]]}
	    }
	    return strings.Fields(e)
      }`

	kitexCommonResolverBody = `*options = append(*options, client.WithResolver(r))
	return nil`

	kitexEtcdClientImports = []string{"github.com/kitex-contrib/registry-etcd"}

	kitexEtcdClient = `r, err := etcd.NewEtcdResolver(rpc.GetResolverAddress())
	if err != nil {
		return err
	}` + consts.LineBreak + kitexCommonResolverBody

	kitexZKClientImports = []string{
		"github.com/kitex-contrib/registry-zookeeper/resolver",
		"time",
	}

	kitexZKClient = `r, err := resolver.NewZookeeperResolver(rpc.GetResolverAddress(), 40*time.Second)
    if err != nil {
		return err
    }` + consts.LineBreak + kitexCommonResolverBody

	kitexNacosClientImports = []string{"github.com/kitex-contrib/registry-nacos/resolver"}

	kitexNacosClient = `r, err := resolver.NewDefaultNacosResolver()
	if err != nil {
		return err
	}` + consts.LineBreak + kitexCommonResolverBody

	kitexPolarisClientImports = []string{
		"github.com/kitex-contrib/polaris",
	}

	kitexPolarisClient = `*options = append(*options, client.WithSuite(polaris.NewDefaultClientSuite()))
	return nil`

	kitexEurekaClientImports = []string{"github.com/kitex-contrib/registry-eureka/resolver"}

	kitexEurekaClient = `r := resolver.NewEurekaResolver(rpc.GetResolverAddress())` +
		consts.LineBreak + kitexCommonResolverBody

	kitexConsulClientImports = []string{"github.com/kitex-contrib/registry-consul"}

	kitexConsulClient = `r, err := consul.NewConsulResolver(rpc.GetResolverAddress()[0])
	if err != nil {
		return err
	}` + consts.LineBreak + kitexCommonResolverBody

	kitexServiceCombClientImports = []string{"github.com/kitex-contrib/registry-servicecomb/resolver"}

	kitexServiceCombClient = `r, err := resolver.NewDefaultSCResolver()
    if err != nil {
        return err
    }` + consts.LineBreak + kitexCommonResolverBody
)

var kitexClientMVCTemplates = []Template{
	{
		Path:   consts.InitGo,
		Delims: [2]string{consts.LeftDelimiter, consts.RightDelimiter},
		UpdateBehavior: UpdateBehavior{
			AppendRender: map[string]interface{}{},
			ReplaceFunc: ReplaceFunc{
				ReplaceFuncName:   make([]string, 0, 5),
				ReplaceFuncImport: make([][]string, 0, 15),
				ReplaceFuncBody:   make([]string, 0, 5),
			},
		},
		Body: `package {{.InitOptsPackage}}
  import (
     {{range $key, $value := .GoFileImports}}
	 {{if eq $key "init.go"}}
	 {{range $k, $v := $value}}
     {{if ne $k ""}}"{{$k}}"{{end}}{{end}}{{end}}{{end}}
  )

  var (
  	defaultClient     RPCClient
  	defaultDstService = "{{.ServiceName}}"
  	once       sync.Once
  )

  func init() {
  	DefaultClient()
  }

  func DefaultClient() RPCClient {
  	once.Do(func() {
  		defaultClient = newClient(defaultDstService)
  	})
  	return defaultClient
  }

  func newClient(dstService string, opts ...client.Option) RPCClient {
    options, err := initClientOpts()
      if err != nil {
      panic("failed to init client options: " + err.Error())
    }
  
    options = append(options, opts...)
  	c, err := NewRPCClient(dstService, options...)
  	if err != nil {
  		panic("failed to init client: " + err.Error())
  	}
  	return c
  }

  func InitClient(dstService string, opts ...client.Option) {
  	defaultClient = newClient(dstService, opts...)
  }
  
  func initClientOpts() (ops []client.Option, err error) {
    // todo edit custom config
  	ops = append(ops, client.WithHostPorts("127.0.0.1:8888"),
  		{{- if eq .Codec "thrift"}}
        client.WithMetaHandler(transmeta.ClientTTHeaderHandler),
        client.WithTransportProtocol(transport.TTHeader),
        {{- end}}
  	)

  	if err = initResolver(&ops); err != nil {
  		panic(err)
  	}

  	return
  }

  // If you do not use the service resolver function, do not edit this function.
  // Otherwise, you can customize and modify it.
  func initResolver(options *[]client.Option) (err error) {
  	{{if ne .ResolverName ""}}
    {{.ResolverBody}}
    {{else}}
    return
    {{end}}
  }`,
	},

	{
		Path:   consts.DefaultKitexClientDir + consts.Slash + consts.EnvGo,
		Delims: [2]string{"[[", "]]"},
		UpdateBehavior: UpdateBehavior{
			AppendRender: map[string]interface{}{},
		},
		CustomFunc: TemplateCustomFuncMap,
		Body: `// Code generated by cwgo generator. DO NOT EDIT.

	  package rpc
	  import (
		[[range $key, $value := .GoFileImports]]
	    [[if eq $key "env.go"]]
	    [[range $k, $v := $value]]
        [[if ne $k ""]]"[[$k]]"[[end]][[end]][[end]][[end]]
	  )

      [[if ne .ResolverName ""]]
      func GetResolverAddress() []string {
		e := os.Getenv("GO_KITEX_RESOLVER_[[ToUpper .ServiceName]]")
	    if len(e) == 0 {
		  return []string{[[$lenSlice := len .ResolverAddress]][[range $key, $value := .ResolverAddress]]"[[$value]]"[[if eq $key (Sub $lenSlice 1)]][[else]], [[end]][[end]]}
	    }
		return strings.Fields(e)
      }
	  [[end]]`,
	},
}
