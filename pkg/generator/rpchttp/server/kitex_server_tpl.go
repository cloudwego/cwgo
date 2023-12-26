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

import (
	"github.com/cloudwego/cwgo/pkg/consts"
	"github.com/cloudwego/cwgo/pkg/generator/common/template"
)

var kitexServerMVCTemplates = []template.Template{
	{
		Path:   consts.DevConf,
		Delims: [2]string{consts.LeftDelimiter, consts.RightDelimiter},
		UpdateBehavior: template.UpdateBehavior{
			AppendRender: map[string]interface{}{},
			Append: template.Append{
				AppendImport: map[string]string{},
			},
		},
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
  address: {{range .RegistryAddress}}
	- {{.}}{{end}}
{{end}}
`,
	},

	{
		Path:   consts.OnlineConf,
		Delims: [2]string{consts.LeftDelimiter, consts.RightDelimiter},
		UpdateBehavior: template.UpdateBehavior{
			AppendRender: map[string]interface{}{},
			Append: template.Append{
				AppendImport: map[string]string{},
			},
		},
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
  address: {{range .RegistryAddress}}
	- {{.}}{{end}}
{{end}}`,
	},

	{
		Path:   consts.TestConf,
		Delims: [2]string{consts.LeftDelimiter, consts.RightDelimiter},
		UpdateBehavior: template.UpdateBehavior{
			AppendRender: map[string]interface{}{},
			Append: template.Append{
				AppendImport: map[string]string{},
			},
		},
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
  address: {{range .RegistryAddress}}
	- {{.}}{{end}}
{{end}}`,
	},

	{
		Path:   consts.ConfGo,
		Delims: [2]string{"[[", "]]"},
		UpdateBehavior: template.UpdateBehavior{
			AppendRender: map[string]interface{}{},
			Append: template.Append{
				AppendImport: map[string]string{},
			},
		},
		CustomFunc: template.CustomFuncMap,
		Body: `package conf

  import (
    [[range $key, $value := .GoFileImports]]
	[[if eq $key "conf/conf.go"]]
	[[range $k, $v := $value]]
    [[if ne $k ""]][[if ne $v ""]][[$v]] "[[$k]]"[[else]]"[[$k]]"[[end]][[end]][[end]][[end]][[end]]
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
  	[[if ne .RegistryName ""]]
	Registry Registry ` + "`yaml:\"registry\"`" + `
	[[end]]
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
  [[if ne .RegistryName ""]]
  type Registry struct {
    Address []string  ` + "`yaml:\"address\"`" + `
  }   
  [[end]]

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

  [[if ne .RegistryName ""]]
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
  }
  [[end]]

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
		Path:   consts.Main,
		Delims: [2]string{consts.LeftDelimiter, consts.RightDelimiter},
		UpdateBehavior: template.UpdateBehavior{
			AppendRender: map[string]interface{}{},
			Append: template.Append{
				AppendImport: map[string]string{},
			},
			ReplaceFunc: template.ReplaceFunc{
				ReplaceFuncName:         make([]string, 0, 5),
				ReplaceFuncAppendImport: make([]map[string]string, 0, 10),
				ReplaceFuncDeleteImport: make([]map[string]string, 0, 10),
				ReplaceFuncBody:         make([]string, 0, 5),
			},
		},
		CustomFunc: template.CustomFuncMap,
		Body: `package main

  import (
    {{range $key, $value := .GoFileImports}}
	{{if eq $key "main.go"}}
	{{range $k, $v := $value}}
    {{if ne $k ""}}{{if ne $v ""}}{{$v}} "{{$k}}"{{else}}"{{$k}}"{{end}}{{end}}{{end}}{{end}}{{end}}
  )

  func main() {
    opts := kitexInit()

    svr := {{ToLower .KitexIdlServiceName}}.NewServer(new({{.KitexIdlServiceName}}Impl), opts...)

    err := svr.Run()
    if err != nil {
      klog.Error(err.Error())
    }
  }

  func kitexInit() (opts []server.Option) {
    // address
    addr, err := net.ResolveTCPAddr("tcp", conf.GetConf().Kitex.Address)
    if err != nil {
      panic(err)
    }
    opts = append(opts, server.WithServiceAddr(addr))

    // service info
    opts = append(opts, server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
	  ServiceName: conf.GetConf().Kitex.ServiceName,
    }))
  
    {{- if eq .Codec "thrift"}}
    // thrift meta handler
    opts = append(opts, server.WithMetaHandler(transmeta.ServerTTHeaderHandler))
    {{- end}}

	if err = initRegistry(&opts); err != nil {
	  panic(err)
    }

    // klog
    logger := kitexLogrus.NewLogger()
    klog.SetLogger(logger)
    klog.SetLevel(conf.LogLevel())
    klog.SetOutput(&lumberjack.Logger{
          		Filename:   conf.GetConf().Log.LogFileName,
          		MaxSize:    conf.GetConf().Log.LogMaxSize,
          		MaxBackups: conf.GetConf().Log.LogMaxBackups,
          		MaxAge:     conf.GetConf().Log.LogMaxAge,
          	})
    return
  }
  
  // If you do not use the service registry function, do not edit this function.
  // Otherwise, you can customize and modify it.
  func initRegistry(ops *[]server.Option) (err error) {
	{{if ne .RegistryName ""}}
		{{.RegistryBody}}
		{{else}}
		return
        {{end}}
  }`,
	},

	{
		Path:   consts.DockerCompose,
		Delims: [2]string{consts.LeftDelimiter, consts.RightDelimiter},
		UpdateBehavior: template.UpdateBehavior{
			AppendRender: map[string]interface{}{},
		},
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

	{
		Path:   consts.DalInitGo,
		Delims: [2]string{consts.LeftDelimiter, consts.RightDelimiter},
		Body: `package dal
  
  import (
    {{range $key, $value := .GoFileImports}}
	{{if eq $key "biz/dal/init.go"}}
	{{range $k, $v := $value}}
    {{if ne $k ""}}{{if ne $v ""}}{{$v}} "{{$k}}"{{else}}"{{$k}}"{{end}}{{end}}{{end}}{{end}}{{end}}
  )

  func Init() {
    redis.Init()
    mysql.Init()
  }`,
	},

	{
		Path:   consts.Gitignore,
		Delims: [2]string{consts.LeftDelimiter, consts.RightDelimiter},
		Body: `*.o
*.a
*.so
_obj
_test
*.[568vq]
[568vq].out
*.cgo1.go
*.cgo2.c
_cgo_defun.c
_cgo_gotypes.go
_cgo_export.*
_testmain.go
*.exe
*.exe~
*.test
*.prof
*.rar
*.zip
*.gz
*.psd
*.bmd
*.cfg
*.pptx
*.log
*nohup.out
*settings.pyc
*.sublime-project
*.sublime-workspace
!.gitkeep
.DS_Store
/.idea
/.vscode
/output
*.local.yml`,
	},

	{
		Path:   consts.MysqlInit,
		Delims: [2]string{consts.LeftDelimiter, consts.RightDelimiter},
		Body: `package mysql
  
  import (
    {{range $key, $value := .GoFileImports}}
	{{if eq $key "biz/dal/mysql/init.go"}}
	{{range $k, $v := $value}}
    {{if ne $k ""}}{{if ne $v ""}}{{$v}} "{{$k}}"{{else}}"{{$k}}"{{end}}{{end}}{{end}}{{end}}{{end}}
  )

  var (
    DB  *gorm.DB
    err error
  )

  func Init() {
    DB, err = gorm.Open(mysql.Open(conf.GetConf().MySQL.DSN),
      &gorm.Config{
        PrepareStmt:            true,
        SkipDefaultTransaction: true,
      },
    )
    if err != nil {
      panic(err)
    }
  }`,
	},

	{
		Path:   consts.RedisInit,
		Delims: [2]string{consts.LeftDelimiter, consts.RightDelimiter},
		Body: `package redis
  
  import (
    {{range $key, $value := .GoFileImports}}
	{{if eq $key "biz/dal/redis/init.go"}}
	{{range $k, $v := $value}}
    {{if ne $k ""}}{{if ne $v ""}}{{$v}} "{{$k}}"{{else}}"{{$k}}"{{end}}{{end}}{{end}}{{end}}{{end}}
  )

  var (
    RedisClient *redis.Client
  )

  func Init() {
    RedisClient = redis.NewClient(&redis.Options{
      Addr:     conf.GetConf().Redis.Address,
      Username: conf.GetConf().Redis.Username,
      Password: conf.GetConf().Redis.Password,
      DB:       conf.GetConf().Redis.DB,
    })
    if err := RedisClient.Ping(context.Background()).Err(); err != nil {
      panic(err)
    }
  }`,
	},

	{
		Path:   consts.Readme,
		Delims: [2]string{consts.LeftDelimiter, consts.RightDelimiter},
		Body: `# *** Project

## introduce

- Use the [Kitex](https://github.com/cloudwego/kitex/) framework
- Generating the base code for unit tests.
- Provides basic config functions
- Provides the most basic MVC code hierarchy.

## Directory structure

|  catalog   | introduce  |
|  ----  | ----  |
| conf  | Configuration files |
| main.go  | Startup file |
| handler.go  | Used for request processing return of response. |
| kitex_gen  | kitex generated code |
| biz/service  | The actual business logic. |
| biz/dal  | Logic for operating the storage layer |

## How to run` + "\n\n```shell\nsh build.sh\nsh output/bootstrap.sh\n```",
	},

	{
		Path:   consts.BuildSh,
		Delims: [2]string{consts.LeftDelimiter, consts.RightDelimiter},
		Body: `#!/usr/bin/env bash
RUN_NAME="{{.ServiceName}}"
mkdir -p output/bin output/conf
cp script/* output/
cp -r conf/* output/conf
chmod +x output/bootstrap.sh
go build -o output/bin/${RUN_NAME}`,
	},

	{
		Path: consts.BootstrapSh,
		Body: `#! /usr/bin/env bash
CURDIR=$(cd $(dirname $0); pwd)
echo "$CURDIR/bin/{{.ServiceName}}"
exec "$CURDIR/bin/{{.ServiceName}}"`,
	},
}
