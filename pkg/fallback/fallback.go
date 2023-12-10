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

package fallback

import (
	"bytes"
	"os"
	"strings"

	"github.com/cloudwego/cwgo/config"
	"github.com/cloudwego/cwgo/pkg/consts"
	"github.com/cloudwego/hertz/cmd/hz/app"
	"github.com/cloudwego/hertz/cmd/hz/util/logs"
	"github.com/cloudwego/kitex"
	kargs "github.com/cloudwego/kitex/tool/cmd/kitex/args"
	"github.com/cloudwego/kitex/tool/internal_pkg/pluginmode/thriftgo"
)

func Fallback(c *config.FallbackArgument) error {
	switch c.ToolType {
	case consts.KitexTool:
		os.Args = c.Args
		var args kargs.Arguments
		args.ParseArgs(kitex.Version)

		out := new(bytes.Buffer)
		cmd := args.BuildCmd(out)
		err := cmd.Run()
		if err != nil {
			if args.Use != "" {
				out := strings.TrimSpace(out.String())
				if strings.HasSuffix(out, thriftgo.TheUseOptionMessage) {
					os.Exit(0)
				}
			}
			os.Exit(1)
		}
	case consts.Hz:
		os.Args = c.Args
		defer func() {
			logs.Flush()
		}()

		cli := app.Init()
		err := cli.Run(os.Args)
		if err != nil {
			logs.Errorf("%v\n", err)
		}
	}
	return nil
}
