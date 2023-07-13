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

package static

import (
	"github.com/cloudwego/cwgo/cmd/dynamic"
	"github.com/cloudwego/cwgo/config"
	"github.com/cloudwego/cwgo/meta"
	"github.com/cloudwego/cwgo/pkg/client"
	"github.com/cloudwego/cwgo/pkg/fallback"
	"github.com/cloudwego/cwgo/pkg/model"
	"github.com/cloudwego/cwgo/pkg/server"
	"github.com/urfave/cli/v2"
)

func Init() *cli.App {
	globalArgs := config.GetGlobalArgs()
	verboseFlag := cli.BoolFlag{Name: "verbose,vv", Usage: "turn on verbose mode"}

	app := cli.NewApp()
	app.Name = "cwgo"
	app.Usage = "All in one tools for CloudWeGo"
	app.Version = meta.Version
	// The default separator for multiple parameters is modified to ";"
	app.SliceFlagSeparator = ";"

	// global flags
	app.Flags = []cli.Flag{
		&verboseFlag,
	}

	// Commands
	app.Commands = []*cli.Command{
		{
			Name:   InitName,
			Usage:  InitUsage,
			Action: dynamic.Terminal,
		},
		{
			Name:  ServerName,
			Usage: ServerUsage,
			Flags: serverFlags(),
			Action: func(c *cli.Context) error {
				err := globalArgs.ServerArgument.ParseCli(c)
				if err != nil {
					return err
				}

				return server.Server(globalArgs.ServerArgument)
			},
		},
		{
			Name:  ClientName,
			Usage: ClientUsage,
			Flags: clientFlags(),
			Action: func(c *cli.Context) error {
				err := globalArgs.ClientArgument.ParseCli(c)
				if err != nil {
					return err
				}
				return client.Client(globalArgs.ClientArgument)
			},
		},
		{
			Name:  ModelName,
			Usage: ModelUsage,
			Flags: modelFlags(),
			Action: func(c *cli.Context) error {
				if err := globalArgs.ModelArgument.ParseCli(c); err != nil {
					return err
				}
				return model.Model(globalArgs.ModelArgument)
			},
		},
		{
			Name:  "fallback",
			Usage: "fallback to hz or kitex",
			Action: func(c *cli.Context) error {
				if err := globalArgs.FallbackArgument.ParseCli(c); err != nil {
					return err
				}
				return fallback.Fallback(globalArgs.FallbackArgument)
			},
		},
	}
	return app
}

const (
	ServerName  = "server"
	ServerUsage = `generate RPC or HTTP server

Examples:
  # Generate RPC server code 
  cwgo server --type RPC --idl  {{path/to/IDL_file.thrift}} --service {{svc_name}}
  
  # Generate HTTP server code 
  cwgo server --type HTTP --idl  {{path/to/IDL_file.thrift}} --service {{svc_name}}
`

	ClientName  = "client"
	ClientUsage = `generate RPC or HTTP client

Examples:
  # Generate RPC client code 
  cwgo client --type RPC --idl  {{path/to/IDL_file.thrift}} --service {{svc_name}}
  
  # Generate HTTP client code 
  cwgo client --type HTTP --idl  {{path/to/IDL_file.thrift}} --service {{svc_name}}
`

	ModelName  = "model"
	ModelUsage = `generate DB model

Examples:
  # Generate DB model code 
  cwgo  model --db_type mysql --dsn "gorm:gorm@tcp(localhost:9910)/gorm?charset=utf8&parseTime=True&loc=Local"
`

	InitName  = "init"
	InitUsage = `interactive terminals provide a more user-friendly experience for generating code`
)
