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

func dockerFlags() []cli.Flag {
	globalArgs := config.GetGlobalArgs()
	return []cli.Flag{
		// Required Flags with Default Values
		&cli.StringFlag{
			Name:        consts.Template,
			Usage:       "The home path of the template",
			Value:       consts.Docker, // 使用默认值 consts.Docker
			Destination: &globalArgs.DockerArgument.Template,
		},
		&cli.StringFlag{
			Name:        consts.Exe,
			Usage:       "The executable name in the built image",
			Destination: &globalArgs.DockerArgument.ExeName,
		},
		&cli.StringFlag{
			Name:        consts.Branch,
			Value:       consts.MainBranch, // 默认值使用 consts.MainBranch
			Usage:       "The branch of the remote repo",
			Destination: &globalArgs.DockerArgument.Branch,
		},
		&cli.StringFlag{
			Name:        consts.Base,
			Value:       consts.Scratch, // 默认值 consts.Scratch
			Usage:       "The base image to build the docker image, default scratch",
			Destination: &globalArgs.DockerArgument.BaseImage,
		},
		&cli.StringFlag{
			Name:        consts.Go,
			Usage:       "The file that contains main function",
			Value:       consts.Main, // 默认值 consts.Main
			Destination: &globalArgs.DockerArgument.Main,
		},
		&cli.UintFlag{
			Name:        consts.Port,
			Usage:       "The port to expose, default 0 will not expose any port",
			Destination: &globalArgs.DockerArgument.Port,
		},
		&cli.StringFlag{
			Name:        consts.TZ,
			Value:       consts.AsizShangHai, // 默认值 consts.AsizShangHai
			Usage:       "The timezone of the container",
			Destination: &globalArgs.DockerArgument.TZ,
		},
		&cli.StringFlag{
			Name:        consts.Version,
			Value:       consts.Alpine, // 默认值 consts.Alpine
			Usage:       "The builder golang image version",
			Destination: &globalArgs.DockerArgument.Version,
		},
		&cli.StringSliceFlag{
			Name:  consts.Mirror,
			Usage: "The mirror site to use in go update",
		},
		&cli.StringSliceFlag{
			Name:  consts.Arguments,
			Usage: "Arguments will used in go run",
		},
		&cli.StringSliceFlag{
			Name:  consts.Etc,
			Usage: "The etc file dirs path of the project",
			Value: cli.NewStringSlice(consts.Etc), // 默认值 consts.Etc
		},
	}
}
