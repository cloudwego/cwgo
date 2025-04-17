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

package config

import (
	"github.com/cloudwego/cwgo/pkg/consts"
	"github.com/urfave/cli/v2"
	"strings"
)

type DockerArgument struct {
	GoVersion     string // compile image go version
	EnableGoProxy bool   // enable go proxy
	GoProxy       string // go proxy url
	Timezone      string // image timezone
	BaseImage     string // service run image
	Port          int    // docker image expose port
	GoFileName    string // go project main file name
	ExeFileName   string // build exe file name
	Template      string // specify local or remote template path
	Branch        string // remote template branch
	RunArgs       string // ext run args
}

func NewDockerArgument() *DockerArgument {
	return &DockerArgument{}
}

func (c *DockerArgument) ParseCli(ctx *cli.Context) error {
	c.GoVersion = ctx.String(consts.GoVersion)
	c.EnableGoProxy = ctx.Bool(consts.EnableGoProxy)
	c.GoProxy = ctx.String(consts.GoProxy)
	c.Timezone = ctx.String(consts.Timezone)
	c.BaseImage = ctx.String(consts.BaseImage)
	c.Port = ctx.Int(consts.Port)
	c.GoFileName = ctx.String(consts.GoFileName)
	c.ExeFileName = ctx.String(consts.ExeFileName)
	c.Template = ctx.String(consts.Template)
	var builder strings.Builder
	for _, arg := range ctx.StringSlice(consts.RunArgs) {
		builder.WriteString(`, "` + arg + `"`)
	}
	c.RunArgs = builder.String()

	return nil
}
