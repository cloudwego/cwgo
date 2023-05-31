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

package config

import (
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

var globalArgs = NewArgument()

func GetGlobalArgs() *Argument {
	return globalArgs
}

type Argument struct {
	Verbose bool

	*ServerArgument
	*ClientArgument
	*ModelArgument
	*FallbackArgument
}

func NewArgument() *Argument {
	return &Argument{
		ServerArgument:   NewServerArgument(),
		ClientArgument:   NewClientArgument(),
		ModelArgument:    NewModelArgument(),
		FallbackArgument: NewFallbackArgument(),
	}
}

const (
	RPC  = "RPC"
	HTTP = "HTTP"
)

const (
	Standard = "standard"
)

const (
	Zk      = "ZK"
	Nacos   = "NACOS"
	Etcd    = "ETCD"
	Polaris = "POLARIS"
)

type DataBaseType string

const (
	MySQL     DataBaseType = "mysql"
	SQLServer DataBaseType = "sqlserver"
	Sqlite    DataBaseType = "sqlite"
	Postgres  DataBaseType = "postgres"
)

type DialectorFunc func(string) gorm.Dialector

var OpenTypeFuncMap = map[DataBaseType]DialectorFunc{
	MySQL:     mysql.Open,
	SQLServer: sqlserver.Open,
	Sqlite:    sqlite.Open,
	Postgres:  postgres.Open,
}

type ToolType string

const (
	Hz    ToolType = "hz"
	Kitex ToolType = "kitex"
)

const (
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

	DSN          = "dsn"
	DBType       = "db_type"
	Tables       = "tables"
	OnlyModel    = "only_model"
	OutFile      = "out_file"
	UnitTest     = "unittest"
	ModelPkgName = "model_pkg"
	Nullable     = "nullable"
	Signable     = "signable"
	IndexTag     = "index_tag"
	TypeTag      = "type_tag"
)
