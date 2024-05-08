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
	"github.com/cloudwego/cwgo/pkg/consts"
	"github.com/urfave/cli/v2"
)

func clientFlags() []cli.Flag {
	globalArgs := config.GetGlobalArgs()
	return []cli.Flag{
		&cli.StringFlag{Name: consts.Service, Usage: "Specify the service name.", Destination: &globalArgs.ClientArgument.Service},
		&cli.StringFlag{Name: consts.ServiceType, Usage: "Specify the generate type. (RPC or HTTP)", Value: consts.RPC},
		&cli.StringFlag{Name: consts.Module, Aliases: []string{"mod"}, Usage: "Specify the Go module name to generate go.mod.", Destination: &globalArgs.ClientArgument.GoMod},
		&cli.StringFlag{Name: consts.IDLPath, Usage: "Specify the IDL file path. (.thrift or .proto)", Destination: &globalArgs.ClientArgument.IdlPath},
		&cli.StringFlag{Name: consts.Template, Usage: "Specify the template path. Currently cwgo supports git templates, such as `--template https://github.com/***/cwgo_template.git`", Destination: &globalArgs.ClientArgument.Template},
		&cli.StringFlag{Name: consts.Branch, Usage: "Specify the git template's branch, default is main branch.", Destination: &globalArgs.ClientArgument.Branch},
		&cli.StringFlag{Name: consts.Registry, Usage: "Specify the registry, default is None"},
		&cli.StringSliceFlag{Name: consts.ProtoSearchPath, Aliases: []string{"I"}, Usage: "Add an IDL search path for includes. (Valid only if idl is protobuf)"},
		&cli.StringSliceFlag{Name: consts.Pass, Usage: "pass param to hz or kitex"},
		&cli.BoolFlag{Name: consts.Verbose, Usage: "Turn on verbose mode."},
	}
}
