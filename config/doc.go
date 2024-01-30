/*
 * Copyright 2024 CloudWeGo Authors
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
	"fmt"
	"strings"

	"github.com/cloudwego/cwgo/pkg/consts"
	"github.com/cloudwego/hertz/cmd/hz/util"
	"github.com/urfave/cli/v2"
)

type DocArgument struct {
	GoMod           string
	PackagePrefix   string
	IdlPaths        []string
	IdlType         string
	OutDir          string
	Name            string
	ModelDir        string
	DaoDir          string
	Verbose         bool
	ProtoSearchPath []string
	ProtocOptions   []string // options to pass through to protoc
	ThriftOptions   []string // options to pass through to thriftgo for go flag
}

func NewDocArgument() *DocArgument {
	return &DocArgument{}
}

func (d *DocArgument) ParseCli(ctx *cli.Context) error {
	d.IdlPaths = ctx.StringSlice(consts.IDLPath)
	d.GoMod = ctx.String(consts.Module)
	d.OutDir = ctx.String(consts.OutDir)
	d.ModelDir = ctx.String(consts.ModelDir)
	d.DaoDir = ctx.String(consts.DaoDir)
	d.Name = ctx.String(consts.Name)
	d.Verbose = ctx.Bool(consts.Verbose)
	d.ProtoSearchPath = ctx.StringSlice(consts.ProtoSearchPath)
	d.ProtocOptions = ctx.StringSlice(consts.Protoc)
	d.ThriftOptions = ctx.StringSlice(consts.ThriftGo)
	return nil
}

func (d *DocArgument) Unpack(data []string) error {
	err := util.UnpackArgs(data, d)
	if err != nil {
		return fmt.Errorf("unpack argument failed: %s", err)
	}
	return nil
}

func (d *DocArgument) Pack() ([]string, error) {
	data, err := util.PackArgs(d)
	if err != nil {
		return nil, fmt.Errorf("pack argument failed: %s", err)
	}
	return data, nil
}

func (d *DocArgument) GetThriftgoOptions(prefix string) (string, error) {
	d.ThriftOptions = append(d.ThriftOptions, "package_prefix="+prefix)
	gas := "go:" + strings.Join(d.ThriftOptions, ",")
	return gas, nil
}
