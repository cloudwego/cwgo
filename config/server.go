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

	"github.com/cloudwego/cwgo/pkg/consts"
	"github.com/urfave/cli/v2"
)

type ServerArgument struct {
	// Common Param
	*CommonParam

	Template   string
	Branch     string
	SliceParam *SliceParam
	Verbose    bool
	Hex        bool // add http listen for kitex

	Cwd    string
	GoSrc  string
	GoPkg  string
	GoPath string
}

type CommonParam struct {
	Service  string // service name
	Type     string // GenerateType: RPC or HTTP
	GoMod    string // Go Mod name
	IdlPath  string
	OutDir   string // output path
	Registry string
}

func NewServerArgument() *ServerArgument {
	return &ServerArgument{
		SliceParam:  &SliceParam{},
		CommonParam: &CommonParam{},
	}
}

func (s *ServerArgument) ParseCli(ctx *cli.Context) error {
	s.Type = strings.ToUpper(ctx.String(consts.ServiceType))
	s.Registry = strings.ToUpper(ctx.String(consts.Registry))
	s.Verbose = ctx.Bool(consts.Verbose)
	s.SliceParam.ProtoSearchPath = ctx.StringSlice(consts.ProtoSearchPath)
	s.SliceParam.Pass = ctx.StringSlice(consts.Pass)
	return nil
}

func (s *SliceParam) WriteAnswer(name string, value interface{}) error {
	if name == consts.Pass {
		s.Pass = strings.Split(value.(string), consts.BlackSpace)
	}
	if name == consts.ProtoSearchPath {
		s.ProtoSearchPath = strings.Split(value.(string), consts.BlackSpace)
	}
	return nil
}

type SliceParam struct {
	Pass            []string
	ProtoSearchPath []string
}
