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

package client

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/cloudwego/cwgo/pkg/consts"

	"github.com/cloudwego/cwgo/pkg/common/utils"
	kargs "github.com/cloudwego/kitex/tool/cmd/kitex/args"
	"github.com/cloudwego/kitex/tool/internal_pkg/log"
	"github.com/cloudwego/kitex/tool/internal_pkg/pluginmode/thriftgo"

	"github.com/cloudwego/hertz/cmd/hz/app"

	"github.com/cloudwego/cwgo/config"

	"github.com/cloudwego/cwgo/pkg/generator"
	hzConfig "github.com/cloudwego/hertz/cmd/hz/config"
	"github.com/cloudwego/hertz/cmd/hz/meta"
	"github.com/cloudwego/hertz/cmd/hz/util/logs"
	"github.com/urfave/cli/v2"
)

func Client(c *config.ClientArgument) error {
	var err error
	err = check(c)
	if err != nil {
		return err
	}

	// check and install tools
	if _, err = utils.LookupTool(consts.Gofumpt); err != nil {
		return err
	}

	switch c.Type {
	case consts.RPC:
		var args kargs.Arguments
		log.Verbose = c.Verbose
		err = convertKitexArgs(c, &args)
		if err != nil {
			return err
		}

		out := new(bytes.Buffer)
		cmd := args.BuildCmd(out)
		err = cmd.Run()
		if err != nil {
			if args.Use != "" {
				out := strings.TrimSpace(out.String())
				if strings.HasSuffix(out, thriftgo.TheUseOptionMessage) {
					utils.ReplaceThriftVersion()
				}
			}
			os.Exit(1)
		}

		if c.Template == "" {
			// initialize cwgo side generator parameters
			clientGen, err := generator.NewClientGenerator(consts.RPC)
			if err != nil {
				return err
			}
			if err = generator.ConvertClientGenerator(clientGen, c); err != nil {
				return err
			}

			// generate cwgo side files
			if err = generator.GenerateClient(clientGen); err != nil {
				return cli.Exit(err, consts.GenerateCwgoError)
			}
		}

		utils.ReplaceThriftVersion()
	case consts.HTTP:
		args := hzConfig.NewArgument()
		utils.SetHzVerboseLog(c.Verbose)
		err = convertHzArgument(c, args)
		if err != nil {
			return err
		}

		module, path, ok := utils.SearchGoMod(consts.CurrentDir, false)
		if ok {
			// go.mod exists
			if module != c.GoMod {
				return fmt.Errorf("module name given by the '-module' option ('%s') is not consist with the name defined in go.mod ('%s' from %s)", c.GoMod, module, path)
			}
			c.GoMod = module
		} else {
			// generate go.mod file
			if err = utils.InitGoMod(c.GoMod); err != nil {
				return fmt.Errorf("init go mod failed: %s", err.Error())
			}
		}

		utils.ReplaceThriftVersion()

		args.CmdType = meta.CmdClient
		logs.Debugf("Args: %#v\n", args)
		err = app.TriggerPlugin(args)
		if err != nil {
			return cli.Exit(err, meta.PluginError)
		}

		if c.Template == "" {
			// initialize cwgo side generator parameters
			clientGen, err := generator.NewClientGenerator(consts.HTTP)
			if err != nil {
				return err
			}
			if err = generator.ConvertClientGenerator(clientGen, c); err != nil {
				return err
			}

			// generate cwgo side files
			if err = generator.GenerateClient(clientGen); err != nil {
				return cli.Exit(err, consts.GenerateCwgoError)
			}
		}
	}
	return nil
}
