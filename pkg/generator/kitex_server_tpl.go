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
	kitexEtcdServerImports = []string{"github.com/kitex-contrib/registry-etcd"}

	kitexEtcdServer = `r, err := etcd.NewEtcdRegistry([]string{conf.GetConf().Registry.Address})
	if err != nil {
		klog.Fatal(err)
	}
    options = append(options, server.WithRegistry(r))`
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
  dsn: "gorm:gorm@tcp(127.0.0.1:3306)/gorm?charset=utf8&parseTime=True&loc=Local"
  
redis:
  address: "127.0.0.1:6379"
  username: ""
  password: ""
  db: 0
{{if ne .RegistryName ""}}
registry:
  address: "127.0.0.1:{{.DefaultRegistryPort}}"
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
  dsn: "gorm:gorm@tcp(127.0.0.1:3306)/gorm?charset=utf8&parseTime=True&loc=Local"
  
redis:
  address: "127.0.0.1:6379"
  username: ""
  password: ""
  db: 0
{{if ne .RegistryName ""}}
registry:
  address: "127.0.0.1:{{.DefaultRegistryPort}}"
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
  dsn: "gorm:gorm@tcp(127.0.0.1:3306)/gorm?charset=utf8&parseTime=True&loc=Local"
  
redis:
  address: "127.0.0.1:6379"
  username: ""
  password: ""
  db: 0
{{if ne .RegistryName ""}}
registry:
  address: "127.0.0.1:{{.DefaultRegistryPort}}"
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
    Address string  ` + "`yaml:\"address\"`" + `
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
  {{if eq .RegistryName "ETCD"}}
  Etcd:
    image: 'bitnami/etcd:latest'
    environment:
      - ALLOW_NONE_AUTHENTICATION=yes
      - ETCD_ADVERTISE_CLIENT_URLS=http://etcd:2379
    ports:
      - "2379:2379"
      - "2380:2380"	
  {{end}}`,
	},
}
