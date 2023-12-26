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

var hzServerMVCTemplates = []template.Template{
	{
		Path:   consts.DevConf,
		Delims: [2]string{consts.LeftDelimiter, consts.RightDelimiter},
		UpdateBehavior: template.UpdateBehavior{
			AppendRender: map[string]interface{}{},
			Append: template.Append{
				AppendImport: map[string]string{},
			},
		},
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
		Path:   consts.OnlineConf,
		Delims: [2]string{consts.LeftDelimiter, consts.RightDelimiter},
		UpdateBehavior: template.UpdateBehavior{
			AppendRender: map[string]interface{}{},
			Append: template.Append{
				AppendImport: map[string]string{},
			},
		},
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
      	Env string

      	Hertz Hertz ` + "`yaml:\"hertz\"`" + `
		Middleware Middleware ` + "`yaml:\"middleware\"`" + `
		Log Log ` + "`yaml:\"log\"`" + `
        MySQL MySQL ` + "`yaml:\"mysql\"`" + `
        Redis Redis ` + "`yaml:\"redis\"`" + `
		[[if ne .RegistryName ""]]
		Registry Registry ` + "`yaml:\"registry\"`" + `
		[[end]]
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
      }

	  [[if ne .RegistryName ""]]
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
      }
	  [[end]]

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
		Body: `// Code generated by hertz generator.

      package main

      import (
        {{range $key, $value := .GoFileImports}}
	    {{if eq $key "main.go"}}
	    {{range $k, $v := $value}}
        {{if ne $k ""}}{{if ne $v ""}}{{$v}} "{{$k}}"{{else}}"{{$k}}"{{end}}{{end}}{{end}}{{end}}{{end}}
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
		opts = append(opts, server.WithHostPorts(conf.GetConf().Hertz.Address))
		
		if err := initRegistry(&opts); err != nil {
		  panic(err)
        }

		return opts
      }

	  // If you do not use the service registry function, do not edit this function.
	  // Otherwise, you can customize and modify it.
	  func initRegistry(ops *[]config.Option) (err error) {
		{{if ne .RegistryName ""}}
		{{.RegistryBody}}
		{{else}}
		return
        {{end}}
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
		UpdateBehavior: template.UpdateBehavior{
			AppendRender: map[string]interface{}{},
			Append: template.Append{
				AppendImport: map[string]string{},
			},
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
		Path:   consts.RegisterGo,
		Delims: [2]string{consts.LeftDelimiter, consts.RightDelimiter},
		Body: `// Code generated by hertz generator. DO NOT EDIT.

      package router

      import (
      	{{range $key, $value := .GoFileImports}}
	    {{if eq $key "biz/router/register.go"}}
	    {{range $k, $v := $value}}
        {{if ne $k ""}}{{if ne $v ""}}{{$v}} "{{$k}}"{{else}}"{{$k}}"{{end}}{{end}}{{end}}{{end}}{{end}}
      )

      // GeneratedRegister registers routers generated by IDL.
      func GeneratedRegister(r *server.Hertz){
      	//INSERT_POINT: DO NOT DELETE THIS LINE!
      }`,
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

      var RedisClient *redis.Client

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

- Use the [Hertz](https://github.com/cloudwego/hertz/) framework
- Integration of pprof, cors, recovery, access_log, gzip and other extensions of Hertz.
- Generating the base code for unit tests.
- Provides basic profile functions.
- Provides the most basic MVC code hierarchy.

## Directory structure

|  catalog   | introduce  |
|  ----  | ----  |
| conf  | Configuration files |
| main.go  | Startup file |
| hertz_gen  | Hertz generated model |
| biz/handler  | Used for request processing, validation and return of response. |
| biz/service  | The actual business logic. |
| biz/dal  | Logic for operating the storage layer |
| biz/route  | Routing and middleware registration |
| biz/utils  | Wrapped some common methods |

## How to run` + "\n\n```shell\nsh build.sh\nsh output/bootstrap.sh\n```",
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
		Path:   consts.RespGo,
		Delims: [2]string{consts.LeftDelimiter, consts.RightDelimiter},
		Body: `package utils

      import (
      	{{range $key, $value := .GoFileImports}}
	    {{if eq $key "biz/utils/resp.go"}}
	    {{range $k, $v := $value}}
        {{if ne $k ""}}{{if ne $v ""}}{{$v}} "{{$k}}"{{else}}"{{$k}}"{{end}}{{end}}{{end}}{{end}}{{end}}
      )

      // SendErrResponse  pack error response
      func SendErrResponse(ctx context.Context, c *app.RequestContext, code int, err error) {
      	// todo edit custom code
      	c.String(code, err.Error())
      }

      // SendSuccessResponse  pack success response
      func SendSuccessResponse(ctx context.Context, c *app.RequestContext, code int, data interface{}) {
      	// todo edit custom code
      	c.JSON(code, data)
      }`,
	},

	{
		Path:   consts.BuildSh,
		Delims: [2]string{consts.LeftDelimiter, consts.RightDelimiter},
		Body: `#!/bin/bash
RUN_NAME={{.ServiceName}}
mkdir -p output/bin output/conf
cp script/bootstrap.sh output 2>/dev/null
chmod +x output/bootstrap.sh
cp -r conf/* output/conf
go build -o output/bin/${RUN_NAME}`,
	},

	{
		Path:   consts.BootstrapSh,
		Delims: [2]string{consts.LeftDelimiter, consts.RightDelimiter},
		Body: `#!/bin/bash
CURDIR=$(cd $(dirname $0); pwd)
BinaryName={{.ServiceName}}
echo "$CURDIR/bin/${BinaryName}"
exec $CURDIR/bin/${BinaryName}`,
	},
}
