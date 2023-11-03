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

package consts

import "runtime"

const (
	Kitex = "kitex"
	Hertz = "hertz"
)

const (
	RPC  = "RPC"
	HTTP = "HTTP"
)

const (
	Server = "server"
	Client = "client"
	DB     = "db"
)

const (
	Thrift   = "thrift"
	Protobuf = "protobuf"
)

// SysType is the running program's operating system type
const SysType = runtime.GOOS

const WindowsOS = "windows"

const (
	Slash      = "/"
	BackSlash  = "\\"
	BlackSpace = " "
	Comma      = ";"
	Tilde      = "~"
	LineBreak  = "\n"
)

// Package Name
const (
	Src                  = "src"
	DefaultHZModelDir    = "hertz_gen"
	DefaultHZClientDir   = "biz/http"
	DefaultKitexModelDir = "kitex_gen"
	DefaultDbOutDir      = "biz/dal/query"
	Standard             = "standard"
	CurrentDir           = "."
)

// File Name
const (
	KitexExtensionYaml = "extensions.yaml"
	LayoutFile         = "layout.yaml"
	PackageLayoutFile  = "package.yaml"
	SuffixGit          = ".git"
	DefaultDbOutFile   = "gen.go"
	Main               = "main.go"
	GoMod              = "go.mod"
	HzFile             = ".hz"
)

// Registration Center
const (
	Zk      = "ZK"
	Nacos   = "NACOS"
	Etcd    = "ETCD"
	Polaris = "POLARIS"
)

type DataBaseType string

// DataBase Name
const (
	MySQL     DataBaseType = "mysql"
	SQLServer DataBaseType = "sqlserver"
	Sqlite    DataBaseType = "sqlite"
	Postgres  DataBaseType = "postgres"
)

type ToolType string

// Tool Name
const (
	Hz        ToolType = "hz"
	KitexTool ToolType = "kitex"
)

const (
	Go     = "go"
	GOPATH = "GOPATH"
	Env    = "env"
	Mod    = "mod"
	Init   = "init"

	OutDir   = "out_dir"
	Verbose  = "verbose"
	Template = "template"

	Service         = "service"
	ServiceType     = "type"
	Module          = "module"
	IDLPath         = "idl"
	Registry        = "registry"
	Pass            = "pass"
	ProtoSearchPath = "proto_search_path"

	DSN           = "dsn"
	DBType        = "db_type"
	Tables        = "tables"
	ExcludeTables = "exclude_tables"
	OnlyModel     = "only_model"
	OutFile       = "out_file"
	UnitTest      = "unittest"
	ModelPkgName  = "model_pkg"
	Nullable      = "nullable"
	Signable      = "signable"
	IndexTag      = "index_tag"
	TypeTag       = "type_tag"
	HexTag        = "hex"
)
