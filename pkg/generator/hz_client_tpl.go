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

import "github.com/cloudwego/cwgo/pkg/consts"

// related to service resolver
var (
	hzCommonResolverImport = "github.com/cloudwego/hertz/pkg/app/middlewares/client/sd"

	hzCommonResolverBody = `ops = append(ops, WithHertzClientMiddleware(sd.Discovery(r)))`

	hzEtcdClientImports = []string{
		hzCommonResolverImport,
		"github.com/hertz-contrib/registry/etcd",
	}

	hzEtcdClient = `r, err := etcd.NewEtcdResolver(conf.GetConf().Resolver.Address)
	if err != nil {
		return nil, err
	}` + consts.LineBreak + hzCommonResolverBody

	hzNacosClientImports = []string{
		hzCommonResolverImport,
		"github.com/hertz-contrib/registry/nacos",
	}

	hzNacosClient = `r, err := nacos.NewDefaultNacosResolver()
	if err != nil {
		return nil, err
	}` + consts.LineBreak + hzCommonResolverBody

	hzConsulClientImports = []string{
		hzCommonResolverImport,
		"github.com/hashicorp/consul/api",
		"github.com/hertz-contrib/registry/consul",
	}

	hzConsulClient = `consulConfig := api.DefaultConfig()
    consulConfig.Address = conf.GetConf().Resolver.Address[0]
    consulClient, err := api.NewClient(consulConfig)
    if err != nil {
        return nil, err
    }
    
    r := consul.NewConsulResolver(consulClient)` + consts.LineBreak + hzCommonResolverBody

	hzEurekaClientImports = []string{
		hzCommonResolverImport,
		"github.com/hertz-contrib/registry/eureka",
	}

	hzEurekaClient = `r := eureka.NewEurekaResolver(conf.GetConf().Resolver.Address)` +
		consts.LineBreak + hzCommonResolverBody

	hzPolarisClientImports = []string{
		hzCommonResolverImport,
		"github.com/hertz-contrib/registry/polaris",
	}

	hzPolarisClient = `r, err := polaris.NewPolarisResolver()
    if err != nil {
        return nil, err
    }` + consts.LineBreak + hzCommonResolverBody

	hzServiceCombClientImports = []string{
		hzCommonResolverImport,
		"github.com/hertz-contrib/registry/servicecomb",
	}

	hzServiceCombClient = `r, err := servicecomb.NewDefaultSCResolver(conf.GetConf().Resolver.Address)
    if err != nil {
        return nil, err
    }` + consts.LineBreak + hzCommonResolverBody

	hzZKClientImports = []string{
		hzCommonResolverImport,
		"github.com/hertz-contrib/registry/zookeeper",
		"time",
	}

	hzZKClient = `r, err := zookeeper.NewZookeeperResolver(conf.GetConf().Resolver.Address, 40*time.Second)
    if err != nil {
        return nil, err
    }` + consts.LineBreak + hzCommonResolverBody
)

var hzClientMVCTemplates = []Template{
	{
		Path:   `{{.OutDir}}/{{.CurrentIDLServiceName}}/init.go`,
		Delims: [2]string{consts.LeftDelimiter, consts.RightDelimiter},
		Body: `package {{.CurrentIDLServiceName}}
      import (
		{{range $key, $value := .GoFileImports}}
	    {{if eq $key "init.go"}}
	    {{range $k, $v := $value}}
        {{if ne $k ""}}"{{$k}}"{{end}}{{end}}{{end}}{{end}}
	  )

	  func initClientOpts(hostUrl string) (ops []Option, err error) {
		ops = append(ops, withHostUrl(hostUrl))
        {{if ne .ResolverName ""}}
		{{.ResolverBody}}
        {{end}}

		return
	  }`,
	},

	{
		Path:   consts.DevConf,
		Delims: [2]string{consts.LeftDelimiter, consts.RightDelimiter},
		Body: `hertz_client:
  host_url: "127.0.0.1:8080"
{{if ne .ResolverName ""}}
resolver:
  address: {{range .DefaultResolverAddress}}
	- {{.}}{{end}}
{{end}}`,
	},

	{
		Path:   consts.OnlineConf,
		Delims: [2]string{consts.LeftDelimiter, consts.RightDelimiter},
		Body: `hertz_client:
  host_url: "127.0.0.1:8080"
{{if ne .ResolverName ""}}
resolver:
  address: {{range .DefaultResolverAddress}}
	- {{.}}{{end}}
{{end}}`,
	},

	{
		Path:   consts.TestConf,
		Delims: [2]string{consts.LeftDelimiter, consts.RightDelimiter},
		Body: `hertz_client:
  host_url: "127.0.0.1:8080"
{{if ne .ResolverName ""}}
resolver:
  address: {{range .DefaultResolverAddress}}
	- {{.}}{{end}}
{{end}}`,
	},

	{
		Path:   consts.ConfGo,
		Delims: [2]string{consts.LeftDelimiter, consts.RightDelimiter},
		Body: `package conf
      import (
        {{range $key, $value := .GoFileImports}}
	    {{if eq $key "conf/conf.go"}}
	    {{range $k, $v := $value}}
        {{if ne $k ""}}"{{$k}}"{{end}}{{end}}{{end}}{{end}}
      )

      var (
      	conf *Config
      	once sync.Once
      )

      type Config struct {
      	Env string

      	HertzClient HertzClient ` + "`yaml:\"hertz_client\"`" + `
		{{if ne .ResolverName ""}}
		Resolver Resolver ` + "`yaml:\"resolver\"`" + `
		{{end}}
      }

      type HertzClient struct {
      	HostUrl       string ` + "`yaml:\"host_url\"`" + `
      }
	  {{if ne .ResolverName ""}}
      type Resolver struct {
		Address []string  ` + "`yaml:\"address\"`" + `
      }   
	  {{end}}

      // GetConf gets configuration instance
      func GetConf() *Config {
      	once.Do(initConf)
      	return conf
      }

      func initConf() {
      	prefix := "conf"
        confFileRelPath := filepath.Join(prefix, filepath.Join(GetEnv(), "conf.yaml"))
      	content, err := ioutil.ReadFile(confFileRelPath)
      	if err != nil {
      		panic(err)
      	}

      	conf = new(Config)
      	err = yaml.Unmarshal(content, conf)
      	if err != nil {
      		hlog.Error("parse yaml error - %v", err)
      		panic(err)
      	}
      	if err := validator.Validate(conf); err != nil {
      		hlog.Error("validate config error - %v", err)
      		panic(err)
      	}

      	conf.Env = GetEnv()

      	pretty.Printf("%+v\n", conf)
      }

      func GetEnv() string {
      	e := os.Getenv("GO_ENV")
      	if len(e) == 0 {
      		return "test"
      	}
      	return e
      }`,
	},

	{
		Path:   consts.Main,
		Delims: [2]string{consts.LeftDelimiter, consts.RightDelimiter},
		Body: `// Code generated by cwgo generator.

      package main
	  import (
		{{range $key, $value := .GoFileImports}}
	    {{if eq $key "main.go"}}
	    {{range $k, $v := $value}}
        {{if ne $k ""}}"{{$k}}"{{end}}{{end}}{{end}}{{end}}
	  )

	  func main() {
		{{$snakeServiceNames := .SnakeServiceNames}}
	    {{range $index, $value := .CamelServiceNames}}
		{{$value}}Client, err := {{index $snakeServiceNames $index}}.New{{$value}}Client(conf.GetConf().HertzClient.HostUrl)
		if err != nil {
		  panic(err)
    	}

		fmt.Printf("%v\n", {{$value}}Client)
		
		// todo your custom code
	    {{end}}
	  }`,
	},
}
