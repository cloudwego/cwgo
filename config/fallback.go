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
	"fmt"

	"github.com/cloudwego/cwgo/pkg/consts"
	"github.com/urfave/cli/v2"
)

type FallbackArgument struct {
	ToolType consts.ToolType
	Args     []string
}

func NewFallbackArgument() *FallbackArgument {
	return &FallbackArgument{}
}

func (c *FallbackArgument) ParseCli(ctx *cli.Context) error {
	args := ctx.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("please input tool type")
	}

	c.ToolType = consts.ToolType(args[0])
	switch consts.ToolType(args[0]) {
	case consts.Hz:
		c.ToolType = consts.Hz
	case consts.KitexTool:
		c.ToolType = consts.KitexTool
	default:
		return fmt.Errorf("tool type is not supported")
	}

	c.Args = args
	return nil
}
