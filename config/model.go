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

type ModelArgument struct {
	DSN               string
	Type              string
	Tables            []string
	OnlyModel         bool
	OutPath           string
	OutFile           string
	WithUnitTest      bool
	ModelPkgName      string
	FieldNullable     bool
	FieldSignable     bool
	FieldWithIndexTag bool
	FieldWithTypeTag  bool
}

func NewModelArgument() *ModelArgument {
	return &ModelArgument{
		OutPath: "biz/dal/query",
		OutFile: "gen.go",
	}
}

func (c *ModelArgument) ParseCli(ctx *cli.Context) error {
	c.DSN = ctx.String(DSN)
	c.Type = strings.ToLower(ctx.String(DBType))
	c.Tables = ctx.StringSlice(Tables)
	c.OnlyModel = ctx.Bool(OnlyModel)
	c.OutPath = ctx.String(OutDir)
	c.OutFile = ctx.String(OutFile)
	c.WithUnitTest = ctx.Bool(UnitTest)
	c.ModelPkgName = ctx.String(ModelPkgName)
	c.FieldNullable = ctx.Bool(Nullable)
	c.FieldSignable = ctx.Bool(Signable)
	c.FieldWithIndexTag = ctx.Bool(IndexTag)
	c.FieldWithTypeTag = ctx.Bool(TypeTag)
	return nil
}
