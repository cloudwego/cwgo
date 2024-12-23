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
	"github.com/cloudwego/cwgo/tpl"

	"github.com/urfave/cli/v2"
)

func dockerFlags() []cli.Flag {
	globalArgs := config.GetGlobalArgs()
	return []cli.Flag{
		&cli.StringFlag{Name: consts.Template, Usage: "The home path of the template", Value: tpl.DockerDir, Destination: &globalArgs.DockerArgument.Template},
		&cli.StringFlag{Name: consts.Exe, Usage: "The executable name in the built image", Destination: &globalArgs.DockerArgument.ExeName},
		&cli.StringFlag{Name: consts.Branch, Value: consts.MainBranch, Usage: "The branch of the remote repo", Destination: &globalArgs.DockerArgument.Branch},
		&cli.BoolFlag{Name: consts.Verbose, Usage: "The Dockerfile path, default is Dockerfile", Destination: &globalArgs.DockerArgument.Verbose},
		&cli.StringFlag{Name: consts.Base, Value: consts.Scratch, Usage: "The base image to build the docker image, default scratch (default \"scratch\")", Destination: &globalArgs.DockerArgument.BaseImage},
		&cli.StringFlag{Name: consts.Go, Usage: "The file that contains main function", Value: consts.Main, Destination: &globalArgs.DockerArgument.Main},
		&cli.UintFlag{Name: consts.Port, Usage: "The port to expose, default 0 will not expose any port", Destination: &globalArgs.DockerArgument.Port},
		&cli.StringFlag{Name: consts.TZ, Value: consts.AsizShangHai, Usage: "The timezone of the container (default \"Asia/Shanghai\")"},
		&cli.StringFlag{Name: consts.Version, Usage: "The builder golang image version"},
		&cli.StringSliceFlag{Name: consts.Mirror, Usage: "The mirror site to use in go update"},
		&cli.StringSliceFlag{Name: consts.Etc, Usage: "The ets file dirs path of the project", Value: cli.NewStringSlice(consts.CurrentDir)},
		&cli.StringSliceFlag{Name: consts.Arguments, Usage: "cmd arguments also checkout if /etc has yaml file with -f)"},
	}
}
