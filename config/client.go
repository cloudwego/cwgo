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

package config

import (
	"strings"

	"github.com/urfave/cli/v2"
)

type ClientArgument struct {
	// Common Param
	*CommonParam

	SliceParam *SliceParam

	Verbose  bool
	Template string
	Cwd      string
	GoSrc    string
	GoPkg    string
	GoPath   string
}

func NewClientArgument() *ClientArgument {
	return &ClientArgument{
		SliceParam:  &SliceParam{},
		CommonParam: &CommonParam{},
	}
}

func (c *ClientArgument) ParseCli(ctx *cli.Context) error {
	c.Type = strings.ToUpper(ctx.String(ServiceType))
	c.Registry = strings.ToUpper(ctx.String(Registry))
	c.Verbose = ctx.Bool(Verbose)
	c.SliceParam.ProtoSearchPath = ctx.StringSlice(ProtoSearchPath)
	c.SliceParam.Pass = ctx.StringSlice(Pass)
	return nil
}
