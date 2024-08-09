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
	"github.com/cloudwego/cwgo/config"
	"github.com/cloudwego/cwgo/meta"
	"github.com/cloudwego/cwgo/pkg/api_list"
	"github.com/cloudwego/cwgo/pkg/client"
	"github.com/cloudwego/cwgo/pkg/consts"
	"github.com/cloudwego/cwgo/pkg/cronjob"
	"github.com/cloudwego/cwgo/pkg/curd/doc"
	"github.com/cloudwego/cwgo/pkg/fallback"
	"github.com/cloudwego/cwgo/pkg/job"
	"github.com/cloudwego/cwgo/pkg/model"
	"github.com/cloudwego/cwgo/pkg/server"
	"github.com/urfave/cli/v2"
)

func Init() *cli.App {
	globalArgs := config.GetGlobalArgs()
	verboseFlag := cli.BoolFlag{Name: "verbose,vv", Usage: "turn on verbose mode"}

	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.Name = meta.Name
	app.Usage = AppUsage
	app.Version = meta.Version
	// The default separator for multiple parameters is modified to ";"
	app.SliceFlagSeparator = consts.Comma

	// global flags
	app.Flags = []cli.Flag{
		&verboseFlag,
	}

	// Commands
	app.Commands = []*cli.Command{
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
			Name:  DocName,
			Usage: DocUsage,
			Flags: docFlags(),
			Action: func(c *cli.Context) error {
				if err := globalArgs.DocArgument.ParseCli(c); err != nil {
					return err
				}
				return doc.Doc(globalArgs.DocArgument)
			},
		},
		{
			Name:  JobName,
			Usage: JobUsage,
			Flags: jobFlags(),
			Action: func(c *cli.Context) error {
				if err := globalArgs.JobArgument.ParseCli(c); err != nil {
					return err
				}
				return job.Job(globalArgs.JobArgument)
			},
		},
		{
			Name:  CronJobName,
			Usage: CronJobUsage,
			Flags: cronjobFlags(),
			Action: func(c *cli.Context) error {
				if err := globalArgs.CronJobArgument.ParseCli(c); err != nil {
					return err
				}
				return cronjob.Cronjob(globalArgs.CronJobArgument)
			},
		},
		{
			Name:  ApiListName,
			Usage: ApiUsage,
			Flags: apiFlags(),
			Action: func(c *cli.Context) error {
				if err := globalArgs.ApiArgument.ParseCli(c); err != nil {
					return err
				}
				return api_list.Api(globalArgs.ApiArgument)
			},
		},
		{
			Name:  FallbackName,
			Usage: FallbackUsage,
			Action: func(c *cli.Context) error {
				if err := globalArgs.FallbackArgument.ParseCli(c); err != nil {
					return err
				}
				return fallback.Fallback(globalArgs.FallbackArgument)
			},
		},
		{
			Name:  CompletionName,
			Usage: CompletionUsage,
			Subcommands: []*cli.Command{
				{
					Name:  CompletionZshName,
					Usage: CompletionZshUsage,
					Action: func(context *cli.Context) error {
						context.App.Writer.Write([]byte(consts.ZshAutocomplete))
						return nil
					},
				},
				{
					Name:  CompletionBashName,
					Usage: CompletionBashUsage,
					Action: func(context *cli.Context) error {
						context.App.Writer.Write([]byte(consts.BashAutocomplete))
						return nil
					},
				},
				{
					Name:  CompletionPowershellName,
					Usage: CompletionPowershellUsage,
					Action: func(context *cli.Context) error {
						context.App.Writer.Write([]byte(consts.PowershellAutoComplete))
						return nil
					},
				},
			},
		},
	}
	return app
}

const (
	AppUsage = "All in one tools for CloudWeGo"

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
  cwgo  model --db_type mysql --dsn "gorm:gorm@tcp(localhost:9910)/gorm?charset=utf8mb4&parseTime=True&loc=Local"
`

	DocName  = "doc"
	DocUsage = `generate doc model

Examples:
  # Generate doc model code
  cwgo doc --name mongodb --idl {{path/to/IDL_file.thrift}}
`

	ApiListName = "api-list"
	ApiUsage    = `analyze router codes by golang ast

Examples:
  cwgo api --project_path ./
`
	JobName  = "job"
	JobUsage = `generate job code

Examples:
	cwgo job --job_name jobOne --job_name jobTwo --module my_job
`
	CronJobName  = "cronjob"
	CronJobUsage = `generate cronjob code

Examples:
	cwgo cronjob --job_name jobOne --job_name jobTwo --module my_cronjob
`
	FallbackName  = "fallback"
	FallbackUsage = "fallback to hz or kitex"

	CompletionName  = "completion"
	CompletionUsage = "Generate the autocompletion script for hugo for the specified shell"

	CompletionZshName  = "zsh"
	CompletionZshUsage = "Generate the autocompletion script for zsh"

	CompletionBashName  = "bash"
	CompletionBashUsage = "Generate the autocompletion script for bash"

	CompletionPowershellName  = "powershell"
	CompletionPowershellUsage = "Generate the autocompletion script for powershell"
)
