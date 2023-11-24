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

// related to basic options
var (
	kitexClientBasicImports = []string{
		"github.com/cloudwego/kitex/pkg/klog",
	}

	kitexClientBasicOpts = ``
)

// related to service resolver
var (
	kitexCommonResolverBody = `options = append(options, client.WithResolver(r))`

	kitexEtcdClientImports = []string{"github.com/kitex-contrib/registry-etcd"}

	kitexEtcdClient = `r, err := etcd.NewEtcdResolver(conf.GetConf().Resolver.Address)
	if err != nil {
		klog.Fatal(err)
	}` + consts.LineBreak + kitexCommonResolverBody

	kitexZKClientImports = []string{
		"github.com/kitex-contrib/registry-zookeeper/resolver",
		"time",
	}

	kitexZKClient = `r, err := resolver.NewZookeeperResolver(conf.GetConf().Registry.Address, 40*time.Second)
    if err != nil {
		klog.Fatal(err)
    }` + consts.LineBreak + kitexCommonResolverBody

	kitexNacosClientImports = []string{"github.com/kitex-contrib/registry-nacos/resolver"}

	kitexNacosClient = `r, err := resolver.NewDefaultNacosResolver()
	if err != nil {
		klog.Fatal(err)
	}` + consts.LineBreak + kitexCommonResolverBody

	kitexPolarisClientImports = []string{
		"github.com/kitex-contrib/registry-polaris",
	}

	kitexPolarisClient = `options = append(options, client.WithSuite(polaris.NewDefaultClientSuite()))`

	kitexEurekaClientImports = []string{"github.com/kitex-contrib/registry-eureka/resolver"}

	kitexEurekaClient = `r := resolver.NewEurekaResolver(conf.GetConf().Registry.Address)` +
		consts.LineBreak + kitexCommonResolverBody

	kitexConsulClientImports = []string{"github.com/kitex-contrib/registry-consul"}

	kitexConsulClient = `r, err := consul.NewConsulResolver("127.0.0.1:8500")
	if err != nil {
		klog.Fatal(err)
	}` + consts.LineBreak + kitexCommonResolverBody

	kitexServiceCombClientImports = []string{"github.com/kitex-contrib/registry-servicecomb/resolver"}

	kitexServiceCombClient = `r, err := resolver.NewDefaultSCResolver()
    if err != nil {
        klog.Fatal(err)
    }` + consts.LineBreak + kitexCommonResolverBody
)

var kitexClientMVCTemplates = []Template{
	{
		Path:   consts.DevConf,
		Delims: [2]string{consts.LeftDelimiter, consts.RightDelimiter},
		Body: `kitex_client:
  service_name: "{{.ServiceName}}"
{{if ne .ResolverName ""}}
resolver:
  address: {{range .DefaultResolverAddress}}
	- {{.}}{{end}}
{{end}}`,
	},

	{
		Path:   consts.OnlineConf,
		Delims: [2]string{consts.LeftDelimiter, consts.RightDelimiter},
		Body: `kitex_client:
  service_name: "{{.ServiceName}}"
{{if ne .ResolverName ""}}
resolver:
  address: {{range .DefaultResolverAddress}}
	- {{.}}{{end}}
{{end}}`,
	},

	{
		Path:   consts.TestConf,
		Delims: [2]string{consts.LeftDelimiter, consts.RightDelimiter},
		Body: `kitex_client:
  service_name: "{{.ServiceName}}"
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

      	KitexClient KitexClient ` + "`yaml:\"kitex_client\"`" + `
		{{if ne .ResolverName ""}}
		Resolver Resolver ` + "`yaml:\"resolver\"`" + `
		{{end}}
      }

      type KitexClient struct {
      	ServiceName       string ` + "`yaml:\"service_name\"`" + `
      }
	  {{if ne .ResolverName ""}}
      type Resolver struct {
		Address []string  ` + "`yaml:\"address\"`" + `
      }   
	  {{end}}

      type BindMainDir struct {
        Dir string ` + "`json:\"Dir\"`" + `
      }

      // GetConf gets configuration instance
      func GetConf() *Config {
      	once.Do(initConf)
      	return conf
      }

      func initConf() {
      	confFileRelPath := getConfAbsPath()
      	content, err := ioutil.ReadFile(confFileRelPath)
      	if err != nil {
      		panic(err)
      	}

      	conf = new(Config)
      	err = yaml.Unmarshal(content, conf)
      	if err != nil {
      		klog.Error("parse yaml error - %v", err)
      		panic(err)
      	}
      	if err := validator.Validate(conf); err != nil {
      		klog.Error("validate config error - %v", err)
      		panic(err)
      	}

      	conf.Env = GetEnv()

      	pretty.Printf("%+v\n", conf)
      }

      func getConfAbsPath() string {
        cmd := exec.Command("go", "list", "-m", "-json")

        var out bytes.Buffer
        cmd.Stdout = &out
        cmd.Stderr = &out
        if err := cmd.Run(); err != nil {
          panic(err)
        }

        bindDir := &BindMainDir{}
        if err := sonic.Unmarshal(out.Bytes(), bindDir); err != nil {
          panic(err)
        }

        prefix := "conf"
        return filepath.Join(bindDir.Dir, prefix, filepath.Join(GetEnv(), "conf.yaml"))
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
		cli, err := {{.ServiceName}}.NewRPCClient(conf.GetConf().KitexClient.ServiceName)
		if err != nil {
		  panic(err)
    	}

		fmt.Printf("%v\n", cli)
		
		// todo your custom code
	  }`,
	},
}
