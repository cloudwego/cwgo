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

// related to basic options
var (
	kitexServiceBasicImports = []string{
		"net",
		"github.com/cloudwego/kitex/pkg/klog",
		"github.com/cloudwego/kitex/pkg/rpcinfo",
	}

	kitexServiceBasicOpts = `// address
    addr, err := net.ResolveTCPAddr("tcp", conf.GetConf().Kitex.Address)
    if err != nil {
      klog.Fatal(err)
    }
    options = append(options, server.WithServiceAddr(addr))

	// service info
    options = append(options, server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
    	ServiceName: conf.GetConf().Kitex.ServiceName,
    }))`
)

// related to service registration
var (
	kitexCommonRegisterBody = `options = append(options, server.WithRegistry(r))`

	kitexEtcdServerImports = []string{"github.com/kitex-contrib/registry-etcd"}

	kitexEtcdServer = `r, err := etcd.NewEtcdRegistry(conf.GetConf().Registry.Address)
	if err != nil {
		klog.Fatal(err)
	}` + consts.LineBreak + kitexCommonRegisterBody

	kitexZKServerImports = []string{
		"github.com/kitex-contrib/registry-zookeeper/registry",
		"time",
	}

	kitexZKServer = `r, err := registry.NewZookeeperRegistry(conf.GetConf().Registry.Address, 40*time.Second)
    if err != nil{
        klog.Fatal(err)
    }` + consts.LineBreak + kitexCommonRegisterBody

	kitexNacosServerImports = []string{"github.com/kitex-contrib/registry-nacos/registry"}

	kitexNacosServer = `r, err := registry.NewDefaultNacosRegistry()
	if err != nil {
		klog.Fatal(err)
	}` + consts.LineBreak + kitexCommonRegisterBody

	kitexPolarisServerImports = []string{
		"github.com/kitex-contrib/polaris",
		"github.com/cloudwego/kitex/pkg/registry",
	}

	kitexPolarisServer = `so := polaris.ServerOptions{}
	r, err := polaris.NewPolarisRegistry(so)
	if err != nil {
		klog.Fatal(err)
	}
	info := &registry.Info{
		ServiceName: conf.GetConf().Kitex.ServiceName,
		Tags: map[string]string{
			"namespace": "Polaris",
		},
	}
	options = append(options, server.WithRegistry(r), server.WithRegistryInfo(info))`

	kitexEurekaServerImports = []string{
		"github.com/kitex-contrib/registry-eureka/registry",
		"time",
	}

	kitexEurekaServer = `r := registry.NewEurekaRegistry(conf.GetConf().Registry.Address, 15*time.Second)` +
		consts.LineBreak + kitexCommonRegisterBody

	kitexConsulServerImports = []string{
		"github.com/kitex-contrib/registry-consul",
		"github.com/cloudwego/kitex/pkg/registry",
	}

	kitexConsulServer = `r, err := consul.NewConsulRegister(conf.GetConf().Registry.Address[0])
	if err != nil {
		klog.Fatal(err)
	}
	info := &registry.Info{
		ServiceName: conf.GetConf().Kitex.ServiceName,
		Weight:      1, // weights must be greater than 0 in consul,else received error and exit.
	}
	options = append(options, server.WithRegistry(r), server.WithRegistryInfo(info))`

	kitexServiceCombServerImports = []string{"github.com/kitex-contrib/registry-servicecomb/registry"}

	kitexServiceCombServer = `r, err := registry.NewDefaultSCRegistry()
    if err != nil {
        klog.Fatal(err)
    }` + consts.LineBreak + kitexCommonRegisterBody
)

var (
	etcdServerAddr        = []string{"127.0.0.1:2379"}
	nacosServerAddr       = []string{"127.0.0.1:8848"}
	consulServerAddr      = []string{"127.0.0.1:8500"}
	eurekaServerAddr      = []string{"http://127.0.0.1:8761/eureka"}
	polarisServerAddr     = []string{"127.0.0.1:8090"}
	serviceCombServerAddr = []string{"127.0.0.1:30100"}
	zkServerAddr          = []string{"127.0.0.1:2181"}

	etcdDocker = `Etcd:
    image: 'bitnami/etcd:latest'
    ports:
      - "2379:2379"
      - "2380:2380"	`

	zkDocker = `zookeeper:
    image: zookeeper
    ports:
      - "2181:2181"`

	nacosDocker = `nacos:
    image: nacos/nacos-server:latest
    ports:
      - "8848:8848"`

	polarisDocker = `polaris:
    image: polarismesh/polaris-server:latest
    ports:
      - "8090:8090"`

	eurekaDocker = `eureka:
    image: 'xdockerh/eureka-server:latest'
    ports:
      - 8761:8761`

	consulDocker = `consul:
    image: consul:latest
    ports:
      - "8500:8500"`

	serviceCombDocker = `service-center:
    image: 'servicecomb/service-center:latest'
    ports:
      - "30100:30100"`
)

var kitexServerMVCTemplates = []Template{
	{
		Path:   consts.DevConf,
		Delims: [2]string{consts.LeftDelimiter, consts.RightDelimiter},
		Body: `kitex:
  service_name: "{{.ServiceName}}"
  address: ":8888"

log:
  log_level: info
  log_file_name: "log/kitex.log"
  log_max_size: 10
  log_max_age: 3
  log_max_backups: 50

mysql:
  dsn: "gorm:gorm@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"
  
redis:
  address: "127.0.0.1:6379"
  username: ""
  password: ""
  db: 0
{{if ne .RegistryName ""}}
registry:
  address: {{range .DefaultRegistryAddress}}
	- {{.}}{{end}}
{{end}}
`,
	},

	{
		Path:   consts.OnlineConf,
		Delims: [2]string{consts.LeftDelimiter, consts.RightDelimiter},
		Body: `kitex:
  service_name: "{{.ServiceName}}"
  address: ":8888"

log:
  log_level: info
  log_file_name: "log/kitex.log"
  log_max_size: 10
  log_max_age: 3
  log_max_backups: 50

mysql:
  dsn: "gorm:gorm@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"
  
redis:
  address: "127.0.0.1:6379"
  username: ""
  password: ""
  db: 0
{{if ne .RegistryName ""}}
registry:
  address: {{range .DefaultRegistryAddress}}
	- {{.}}{{end}}
{{end}}`,
	},

	{
		Path:   consts.TestConf,
		Delims: [2]string{consts.LeftDelimiter, consts.RightDelimiter},
		Body: `kitex:
  service_name: "{{.ServiceName}}"
  address: ":8888"

log:
  log_level: info
  log_file_name: "log/kitex.log"
  log_max_size: 10
  log_max_age: 3
  log_max_backups: 50

mysql:
  dsn: "gorm:gorm@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"
  
redis:
  address: "127.0.0.1:6379"
  username: ""
  password: ""
  db: 0
{{if ne .RegistryName ""}}
registry:
  address: {{range .DefaultRegistryAddress}}
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
  	Env      string

  	Kitex    Kitex    ` + "`yaml:\"kitex\"`" + `
  	MySQL    MySQL    ` + "`yaml:\"mysql\"`" + `
  	Redis    Redis    ` + "`yaml:\"redis\"`" + `
	Log      Log      ` + "`yaml:\"log\"`" + `
  	{{if ne .RegistryName ""}}
	Registry Registry ` + "`yaml:\"registry\"`" + `
	{{end}}
  }

  type MySQL struct {
    DSN string ` + "`yaml:\"dsn\"`" + `
  }

  type Redis struct {
    Address  string ` + "`yaml:\"address\"`" + `
    Username string ` + "`yaml:\"username\"`" + `
    Password string ` + "`yaml:\"password\"`" + `
    DB       int    ` + "`yaml:\"db\"`" + `
  }

  type Kitex struct {
    ServiceName     string   ` + "`yaml:\"service_name\"`" + `
    Address         string   ` + "`yaml:\"address\"`" + `
  }

  type Log struct {
    LogLevel        string   ` + "`yaml:\"log_level\"`" + `
    LogFileName     string   ` + "`yaml:\"log_file_name\"`" + `
    LogMaxSize      int      ` + "`yaml:\"log_max_size\"`" + `
    LogMaxBackups   int      ` + "`yaml:\"log_max_backups\"`" + `
    LogMaxAge       int      ` + "`yaml:\"log_max_age\"`" + `
  }
  {{if ne .RegistryName ""}}
  type Registry struct {
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

  func GetEnv() string {
    e := os.Getenv("GO_ENV")
    if len(e) == 0 {
      return "test"
    }
    return e
  }

  func LogLevel() klog.Level {
    level := GetConf().Log.LogLevel
    switch level {
    case "trace":
      return klog.LevelTrace
    case "debug":
      return klog.LevelDebug
    case "info":
      return klog.LevelInfo
    case "notice":
      return klog.LevelNotice
    case "warn":
      return klog.LevelWarn
    case "error":
      return klog.LevelError
    case "fatal":
      return klog.LevelFatal
    default:
      return klog.LevelInfo
    }
  }`,
	},

	{
		Path:   consts.DockerCompose,
		Delims: [2]string{consts.LeftDelimiter, consts.RightDelimiter},
		Body: `version: '3'

services:
  mysql:
    image: 'mysql:latest'
    ports:
      - 3306:3306
    environment:
      - MYSQL_DATABASE=gorm
      - MYSQL_USER=gorm
      - MYSQL_PASSWORD=gorm
      - MYSQL_RANDOM_ROOT_PASSWORD="yes"

  redis:
    image: 'redis:latest'
    ports:
      - 6379:6379
  
  {{.RegistryDocker}}
`,
	},
}
