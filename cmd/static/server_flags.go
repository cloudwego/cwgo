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
	"github.com/urfave/cli/v2"
)

func serverFlags() []cli.Flag {
	globalArgs := config.GetGlobalArgs()
	return []cli.Flag{
		&cli.StringFlag{Name: config.Service, Usage: "Specify the service name.", Destination: &globalArgs.ServerArgument.Service},
		&cli.StringFlag{Name: config.ServiceType, Usage: "Specify the generate type. (RPC or HTTP)", Value: config.RPC},
		&cli.StringFlag{Name: config.Module, Aliases: []string{"mod"}, Usage: "Specify the Go module name to generate go.mod.", Destination: &globalArgs.ServerArgument.GoMod},
		&cli.StringFlag{Name: config.IDLPath, Usage: "Specify the IDL file path. (.thrift or .proto)", Destination: &globalArgs.ServerArgument.IdlPath},
		&cli.StringFlag{Name: config.OutDir, Value: ".", Aliases: []string{"o"}, Usage: "Specify the output path. Currently cwgo supports git templates, such as `--template https://github.com/***/cwgo_template.git`", Destination: &globalArgs.ServerArgument.OutDir},
		&cli.StringFlag{Name: config.Template, Usage: "Specify the layout template.", Destination: &globalArgs.ServerArgument.Template},
		&cli.StringFlag{Name: config.Registry, Usage: "Specify the registry, default is None"},
		&cli.StringSliceFlag{Name: config.ProtoSearchPath, Aliases: []string{"I"}, Usage: "Add an IDL search path for includes. (Valid only if idl is protobuf)"},
		&cli.StringSliceFlag{Name: config.Pass, Usage: "pass param to hz or kitex"},
	}
}
