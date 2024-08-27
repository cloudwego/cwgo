// Copyright 2024 CloudWeGo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package static

import (
	"github.com/cloudwego/cwgo/pkg/consts"
	"github.com/cloudwego/cwgo/tpl"
	"github.com/urfave/cli/v2"
)

func dockerFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     consts.GoVersion,
			Usage:    "Specify the go version.",
			Aliases:  []string{"go"},
			Required: true,
		},
		&cli.BoolFlag{
			Name:        consts.EnableGoProxy,
			Usage:       "Enable go proxy.",
			DefaultText: "false",
		},
		&cli.StringFlag{
			Name:        consts.GoProxy,
			Usage:       "Specify the go proxy.",
			DefaultText: "https://goproxy.cn,direct",
		},
		&cli.StringFlag{
			Name:        consts.Timezone,
			Usage:       "Specify the timezone.",
			DefaultText: "Asia/Shanghai",
		},
		&cli.StringFlag{
			Name:        consts.BaseImage,
			Usage:       "Specify the base image.",
			DefaultText: "scratch",
		},
		&cli.IntFlag{
			Name:        consts.Port,
			Usage:       "Specify the port.",
			DefaultText: "0",
		},
		&cli.StringFlag{
			Name:     consts.GoFileName,
			Usage:    "Specify the go file name.",
			Aliases:  []string{"f"},
			Required: true,
		},
		&cli.StringFlag{
			Name:  consts.ExeFileName,
			Usage: "Specify the exe file name.",
		},
		&cli.StringFlag{
			Name:        consts.Template,
			Usage:       "Specify the template path. Currently cwgo supports git templates, such as `--template https://github.com/***/cwgo_template.git`",
			Aliases:     []string{"t"},
			DefaultText: tpl.DockerDir,
		},
		&cli.StringFlag{
			Name:  consts.Branch,
			Usage: "Specify the git template's branch, default is main branch.",
		},
		&cli.StringSliceFlag{
			Name:  consts.RunArgs,
			Usage: "Specify the run args.",
		},
	}
}
