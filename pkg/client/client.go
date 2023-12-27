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
	"path/filepath"
	"strings"

	"github.com/cloudwego/cwgo/pkg/generator/rpchttp/client"

	"github.com/cloudwego/cwgo/pkg/consts"

	cwgoMeta "github.com/cloudwego/cwgo/meta"
	"github.com/cloudwego/cwgo/pkg/common/utils"
	kargs "github.com/cloudwego/kitex/tool/cmd/kitex/args"
	"github.com/cloudwego/kitex/tool/internal_pkg/log"
	"github.com/cloudwego/kitex/tool/internal_pkg/pluginmode/thriftgo"

	"github.com/cloudwego/hertz/cmd/hz/app"

	"github.com/cloudwego/cwgo/config"

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

	workPath, err := filepath.Abs(consts.CurrentDir)
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	clientGen := new(client.Generator)
	cwgoManifest := new(cwgoMeta.Manifest)

	if c.Type == consts.RPC {
		clientGen, err = client.NewGenerator(consts.RPC)
		if err != nil {
			return err
		}
	} else {
		clientGen, err = client.NewGenerator(consts.HTTP)
		if err != nil {
			return err
		}
	}

	// handle manifest and determine if .cwgo exists
	var dir string
	dir, err = os.Getwd()
	if err != nil {
		return fmt.Errorf("get current path failed: %s", err)
	}

	isNew := utils.IsCwgoNew(dir)
	if isNew {
		clientGen.IsNew = true
	} else {
		if err = cwgoManifest.InitAndValidate(dir); err != nil {
			return err
		}
		if !(cwgoManifest.CommandType == consts.Client && cwgoManifest.CommunicationType == c.Type) {
			clientGen.IsNew = true
		}
	}

	cwgoManifest.Version = cwgoMeta.Version
	cwgoManifest.CommandType = consts.Client
	cwgoManifest.CommunicationType = c.Type
	cwgoManifest.HzInfo = meta.Manifest{}
	cwgoManifest.KitexInfo = cwgoMeta.KitexInfo{}

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
					utils.ReplaceThriftVersion(c.GoModPath)
				}
			}
			os.Exit(1)
		}

		if c.Template == "" {
			// initialize cwgo side generator parameters
			if err = client.ConvertGenerator(clientGen, c); err != nil {
				return err
			}

			// generate cwgo side files
			if err = client.GenerateClient(clientGen); err != nil {
				return cli.Exit(err, consts.GenerateCwgoError)
			}
		}

	case consts.HTTP:
		args := hzConfig.NewArgument()
		utils.SetHzVerboseLog(c.Verbose)
		err = convertHzArgument(c, args)
		if err != nil {
			return err
		}

		module, path, ok := utils.SearchGoMod(workPath, true)
		if ok {
			// go.mod exists
			if module != c.GoMod {
				return fmt.Errorf("module name given by the '-module' option ('%s') is not consist with the name defined in go.mod ('%s' from %s)", c.GoMod, module, path)
			}
			c.GoMod = module
			c.GoModPath = path
		} else {
			if err = utils.InitGoMod(c.GoMod); err != nil {
				log.Warn("Init go mod failed:", err.Error())
				os.Exit(1)
			}
			c.GoModPath = workPath
		}

		args.CmdType = meta.CmdClient
		logs.Debugf("Args: %#v\n", args)
		err = app.TriggerPlugin(args)
		if err != nil {
			return cli.Exit(err, meta.PluginError)
		}

		if c.Template == "" {
			// initialize cwgo side generator parameters
			if err = client.ConvertGenerator(clientGen, c); err != nil {
				return err
			}

			// generate cwgo side files
			if err = client.GenerateClient(clientGen); err != nil {
				return cli.Exit(err, consts.GenerateCwgoError)
			}
		}
	}

	// generate .cwgo file
	if err = cwgoManifest.Persist(dir); err != nil {
		return err
	}

	utils.ReplaceThriftVersion(c.GoModPath)
	return nil
}
