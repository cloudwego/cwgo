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
	CwgoDocPluginMode       = "CWGO_DOC_PLUGIN_DOC"
	ThriftCwgoDocPluginName = "thrift-gen-cwgo-doc"
)

const (
	HertzRepoDefaultUrl = "github.com/cloudwego/hertz"
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
	Proto    = "proto"
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
	Src                   = "src"
	DefaultHZModelDir     = "hertz_gen"
	DefaultHZClientDir    = "biz/http"
	DefaultKitexModelDir  = "kitex_gen"
	DefaultDbOutDir       = "biz/dal/query"
	DefaultDocModelOutDir = "biz/doc/model"
	DefaultDocDaoOutDir   = "biz/doc/dao"
	Standard              = "standard"
	RestApi               = "rest_api"
	CurrentDir            = "."
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
	Branch   = "branch"
	Name     = "name"

	ModelDir = "model_dir"
	DaoDir   = "dao_dir"

	Service         = "service"
	ServerName      = "server_name"
	ServiceType     = "type"
	Module          = "module"
	IDLPath         = "idl"
	Registry        = "registry"
	Pass            = "pass"
	ProtoSearchPath = "proto_search_path"
	ThriftGo        = "thriftgo"
	Protoc          = "protoc"
	GenBase         = "gen_base"

	ProjectPath   = "project_path"
	HertzRepoUrl  = "hertz_repo_url"
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
	SQLDir        = "sql_dir"
)

const (
	MongoDb = "mongodb"
)

const (
	BashAutocomplete = `#! /bin/bash

# Macs have bash3 for which the bash-completion package doesn't include
# _init_completion. This is a minimal version of that function.
_cli_init_completion() {
  COMPREPLY=()
  _get_comp_words_by_ref "$@" cur prev words cword
}

_cli_bash_autocomplete() {
  if [[ "${COMP_WORDS[0]}" != "source" ]]; then
    local cur opts base words
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    if declare -F _init_completion >/dev/null 2>&1; then
      _init_completion -n "=:" || return
    else
      _cli_init_completion -n "=:" || return
    fi
    words=("${words[@]:0:$cword}")
    if [[ "$cur" == "-"* ]]; then
      requestComp="${words[*]} ${cur} --generate-bash-completion"
    else
      requestComp="${words[*]} --generate-bash-completion"
    fi
    opts=$(eval "${requestComp}" 2>/dev/null)
    COMPREPLY=($(compgen -W "${opts}" -- ${cur}))
    return 0
  fi
}

complete -o bashdefault -o default -o nospace -F _cli_bash_autocomplete cwgo`

	PowershellAutoComplete = `$fn = $($MyInvocation.MyCommand.Name)
$name = $fn -replace "(.*)\.ps1$", '$1'
Register-ArgumentCompleter -Native -CommandName $name -ScriptBlock {
     param($commandName, $wordToComplete, $cursorPosition)
     $other = "$wordToComplete --generate-bash-completion"
         Invoke-Expression $other | ForEach-Object {
            [System.Management.Automation.CompletionResult]::new($_, $_, 'ParameterValue', $_)
         }
 }`
	ZshAutocomplete = `#compdef cwgo

_cli_zsh_autocomplete() {
  local -a opts
  local cur
  cur=${words[-1]}
  if [[ "$cur" == "-"* ]]; then
    opts=("${(@f)$(${words[@]:0:#words[@]-1} ${cur} --generate-bash-completion)}")
  else
    opts=("${(@f)$(${words[@]:0:#words[@]-1} --generate-bash-completion)}")
  fi

  if [[ "${opts[1]}" != "" ]]; then
    _describe 'values' opts
  else
    _files
  fi
}

compdef _cli_zsh_autocomplete cwgo`
)
