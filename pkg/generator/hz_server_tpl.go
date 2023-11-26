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

// related to service registration
var (
	hzCommonRegistryBody = `opts = append(opts, server.WithRegistry(r, &registry.Info{
		ServiceName: conf.GetConf().Hertz.ServiceName,
		Addr:        utils.NewNetAddr("tcp", conf.GetConf().Hertz.Address),
		Weight:      10,
		Tags:        nil,
	}))`

	hzCommonRegistyImport = "github.com/cloudwego/hertz/pkg/app/server/registry"

	hzEtcdServerImports = []string{
		hzCommonRegistyImport,
		"github.com/hertz-contrib/registry/etcd",
	}

	hzEtcdServer = `r, err := etcd.NewEtcdRegistry(conf.GetConf().Registry.Address)
	if err != nil {
		panic(err)
	}` + consts.LineBreak + hzCommonRegistryBody

	hzNacosServerImports = []string{
		hzCommonRegistyImport,
		"github.com/hertz-contrib/registry/nacos",
	}

	hzNacosServer = `r, err := nacos.NewDefaultNacosRegistry()
    if err != nil {
        panic(err)
    }` + consts.LineBreak + hzCommonRegistryBody

	hzConsulServerImports = []string{
		hzCommonRegistyImport,
		"github.com/hashicorp/consul/api",
		"github.com/hertz-contrib/registry/consul",
	}

	hzConsulServer = `consulConfig := api.DefaultConfig()
    consulConfig.Address = conf.GetConf().Registry.Address[0]
    consulClient, err := api.NewClient(consulConfig)
    if err != nil {
        panic(err)
    }
    
    r := consul.NewConsulRegister(consulClient)` + consts.LineBreak + hzCommonRegistryBody

	hzEurekaServerImports = []string{
		hzCommonRegistyImport,
		"github.com/hertz-contrib/registry/eureka",
		"time",
	}

	hzEurekaServer = `r := eureka.NewEurekaRegistry(conf.GetConf().Registry.Address, 40*time.Second)` +
		consts.LineBreak + hzCommonRegistryBody

	hzPolarisServerImports = []string{
		hzCommonRegistyImport,
		"github.com/hertz-contrib/registry/polaris",
	}

	hzPolarisServer = `r, err := polaris.NewPolarisRegistry()
    if err != nil {
        panic(err)
    }
	opts = append(opts, server.WithRegistry(r, &registry.Info{
		ServiceName: conf.GetConf().Hertz.ServiceName,
		Addr:        utils.NewNetAddr("tcp", conf.GetConf().Hertz.Address),
		Tags: map[string]string{
            "namespace": "Polaris",
        },
	}))`

	hzServiceCombServerImports = []string{
		hzCommonRegistyImport,
		"github.com/hertz-contrib/registry/servicecomb",
	}

	hzServiceCombServer = `r, err := servicecomb.NewDefaultSCRegistry(conf.GetConf().Registry.Address)
    if err != nil {
        panic(err)
    }` + consts.LineBreak + hzCommonRegistryBody

	hzZKServerImports = []string{
		hzCommonRegistyImport,
		"github.com/hertz-contrib/registry/zookeeper",
		"time",
	}

	hzZKServer = `r, err := zookeeper.NewZookeeperRegistry(conf.GetConf().Registry.Address, 40*time.Second)
    if err != nil {
        panic(err)
    }` + consts.LineBreak + hzCommonRegistryBody
)

var hzServerMVCTemplates = []Template{
	{
		Path:   consts.DevConf,
		Delims: [2]string{consts.LeftDelimiter, consts.RightDelimiter},
		Body: `hertz:
  service_name: "{{.ServiceName}}"
  address: ":8080"

middleware:
  enable_pprof: true
  enable_gzip: true
  enable_access_log: true

log:
  log_level: "info"
  log_file_name: "log/hertz.log"
  log_max_size: 10
  log_max_age: 3
  log_max_backups: 50

mysql:
  dsn: "gorm:gorm@tcp(127.0.0.1:3306)/gorm?charset=utf8&parseTime=True&loc=Local"

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
		Path:   consts.OnlineConf,
		Delims: [2]string{consts.LeftDelimiter, consts.RightDelimiter},
		Body: `hertz:
  service_name: "{{.ServiceName}}"
  address: ":8080"

middleware:
  enable_pprof: false
  enable_gzip: true
  enable_access_log: true

log:
  log_level: "info"
  log_file_name: "log/hertz.log"
  log_max_size: 10
  log_max_age: 3
  log_max_backups: 50

mysql:
  dsn: "gorm:gorm@tcp(127.0.0.1:3306)/gorm?charset=utf8&parseTime=True&loc=Local"

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
		Body: `hertz:
  service_name: "{{.ServiceName}}"
  address: ":8080"

middleware:
  enable_pprof: true
  enable_gzip: true
  enable_access_log: true

log:
  log_level: "info"
  log_file_name: "log/hertz.log"
  log_max_size: 10
  log_max_age: 3
  log_max_backups: 50

mysql:
  dsn: "gorm:gorm@tcp(127.0.0.1:3306)/gorm?charset=utf8&parseTime=True&loc=Local"

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
      	Env string

      	Hertz Hertz ` + "`yaml:\"hertz\"`" + `
		Middleware Middleware ` + "`yaml:\"middleware\"`" + `
		Log Log ` + "`yaml:\"log\"`" + `
        MySQL MySQL ` + "`yaml:\"mysql\"`" + `
        Redis Redis ` + "`yaml:\"redis\"`" + `
		{{if ne .RegistryName ""}}
		Registry Registry ` + "`yaml:\"registry\"`" + `
		{{end}}
      }

      type MySQL struct {
      	DSN string ` + "`yaml:\"dsn\"`" + `
      }

      type Redis struct {
      	Address  string ` + "`yaml:\"address\"`" + `
      	Password string ` + "`yaml:\"password\"`" + `
        Username string ` + "`yaml:\"username\"`" + `
        DB       int    ` + "`yaml:\"db\"`" + `
      }

      type Hertz struct {
      	Address       string ` + "`yaml:\"address\"`" + `
		ServiceName string ` + "`yaml:\"service_name\"`" + `
      }

      type Middleware struct {
		EnablePprof   bool   ` + "`yaml:\"enable_pprof\"`" + `
      	EnableGzip    bool   ` + "`yaml:\"enable_gzip\"`" + `
        EnableAccessLog bool ` + "`yaml:\"enable_access_log\"`" + `
	  }

	  type Log struct {
		LogLevel      string ` + "`yaml:\"log_level\"`" + `
      	LogFileName   string ` + "`yaml:\"log_file_name\"`" + `
      	LogMaxSize    int    ` + "`yaml:\"log_max_size\"`" + `
      	LogMaxBackups int    ` + "`yaml:\"log_max_backups\"`" + `
      	LogMaxAge     int    ` + "`yaml:\"log_max_age\"`" + `
	  }
	  {{if ne .RegistryName ""}}
      type Registry struct {
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
      }

      func LogLevel() hlog.Level {
      	level := GetConf().Log.LogLevel
      	switch level {
      	case "trace":
      		return hlog.LevelTrace
      	case "debug":
      		return hlog.LevelDebug
      	case "info":
      		return hlog.LevelInfo
      	case "notice":
      		return hlog.LevelNotice
      	case "warn":
      		return hlog.LevelWarn
      	case "error":
      		return hlog.LevelError
      	case "fatal":
      		return hlog.LevelFatal
      	default:
      		return hlog.LevelInfo
      	}
      }`,
	},

	{
		Path:   consts.Main,
		Delims: [2]string{consts.LeftDelimiter, consts.RightDelimiter},
		Body: `// Code generated by hertz generator.

      package main

      import (
        {{range $key, $value := .GoFileImports}}
	    {{if eq $key "main.go"}}
	    {{range $k, $v := $value}}
        {{if ne $k ""}}"{{$k}}"{{end}}{{end}}{{end}}{{end}}
      )

      func main() {
        // init dal
        // dal.Init()

      	h := server.Default(initServerOpts()...)

        registerMiddleware(h)

        // add a ping route to test
        h.GET("/ping", func(c context.Context, ctx *app.RequestContext) {
        	ctx.JSON(consts.StatusOK, utils.H{"ping": "pong"})
        })

      	router.GeneratedRegister(h)

      	h.Spin()
      }

      func initServerOpts() (opts []config.Option) {
		address := conf.GetConf().Hertz.Address
		opts = append(opts, server.WithHostPorts(address))
		{{if ne .RegistryName ""}}
        {{.RegistryBody}}
		{{end}}

		return opts
      }

      func registerMiddleware(h *server.Hertz) {
      	// log
      	logger := logrus.NewLogger()
      	hlog.SetLogger(logger)
      	hlog.SetLevel(conf.LogLevel())
      	hlog.SetOutput(&lumberjack.Logger{
      		Filename:   conf.GetConf().Log.LogFileName,
      		MaxSize:    conf.GetConf().Log.LogMaxSize,
      		MaxBackups: conf.GetConf().Log.LogMaxBackups,
      		MaxAge:     conf.GetConf().Log.LogMaxAge,
      	})

      	// pprof
      	if conf.GetConf().Middleware.EnablePprof {
      		pprof.Register(h)
      	}
      
      	// gzip
      	if conf.GetConf().Middleware.EnableGzip {
      		h.Use(gzip.Gzip(gzip.DefaultCompression))
      	}

        // access log
        if conf.GetConf().Middleware.EnableAccessLog {
            h.Use(accesslog.New())
        }

        // cors
        h.Use(cors.Default())
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
