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
)

type KubeArgument struct {
	Name            string
	Namespace       string
	Image           string
	Secret          string
	RequestCpu      int
	RequestMem      string
	LimitCpu        string
	LimitMem        string
	MinReplicas     int
	MaxReplicas     int
	ImagePullPolicy string
	ServiceAccount  string
}

func NewKubeArgument() *KubeArgument {
	return &KubeArgument{}
}

func (c *KubeArgument) ParseCli(ctx *cli.Context) error {
	c.Name = ctx.String(consts.Name)
	c.Namespace = ctx.String(consts.Namespace)
	c.Image = ctx.String(consts.Image)
	c.Secret = ctx.String(consts.Secret)
	c.RequestCpu = ctx.Int(consts.RequestCpu)
	c.RequestMem = ctx.String(consts.RequestMem)
	c.LimitCpu = ctx.String(consts.LimitCpu)
	c.LimitMem = ctx.String(consts.LimitMem)
	c.MinReplicas = ctx.Int(consts.MinReplicas)
	c.MaxReplicas = ctx.Int(consts.MaxReplicas)
	c.ImagePullPolicy = ctx.String(consts.ImagePullPolicy)
	c.ServiceAccount = ctx.String(consts.ServiceAccount)
	return nil
}
