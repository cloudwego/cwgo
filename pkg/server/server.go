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

package server

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloudwego/cwgo/pkg/generator/rpchttp/server"

	"github.com/cloudwego/cwgo/config"
	cwgoMeta "github.com/cloudwego/cwgo/meta"
	"github.com/cloudwego/cwgo/pkg/common/utils"
	"github.com/cloudwego/cwgo/pkg/consts"
	"github.com/cloudwego/hertz/cmd/hz/app"
	hzConfig "github.com/cloudwego/hertz/cmd/hz/config"
	"github.com/cloudwego/hertz/cmd/hz/meta"
	kargs "github.com/cloudwego/kitex/tool/cmd/kitex/args"
	"github.com/cloudwego/kitex/tool/internal_pkg/log"
	"github.com/cloudwego/kitex/tool/internal_pkg/pluginmode/thriftgo"
	"github.com/urfave/cli/v2"
)

func Server(c *config.ServerArgument) error {
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

	serverGen := new(server.Generator)
	cwgoManifest := new(cwgoMeta.Manifest)

	if c.Type == consts.RPC {
		serverGen, err = server.NewGenerator(consts.RPC)
		if err != nil {
			return err
		}
	} else {
		serverGen, err = server.NewGenerator(consts.HTTP)
		if err != nil {
			return err
		}
	}

	// handle manifest and determine if .cwgo exists
	isNew := utils.IsCwgoNew(c.OutDir)
	if isNew {
		serverGen.IsNew = true
	} else {
		if err = cwgoManifest.InitAndValidate(c.OutDir); err != nil {
			return err
		}
		if !(cwgoManifest.CommandType == consts.Server && cwgoManifest.CommunicationType == c.Type) {
			serverGen.IsNew = true
		}
	}

	cwgoManifest.Version = cwgoMeta.Version
	cwgoManifest.CommandType = consts.Server
	cwgoManifest.CommunicationType = c.Type

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
		if c.Hex { // add http listen for kitex
			hzArgs, err := hzArgsForHex(c)
			if err != nil {
				return err
			}
			err = app.TriggerPlugin(hzArgs)
			if err != nil {
				return err
			}
			err = generateHexFile(c)
			if err != nil {
				return err
			}
			err = addHexOptions()
			if err != nil {
				log.Warn("please add \"opts = append(opts,server.WithTransHandlerFactory(&mixTransHandlerFactory{nil}))\", to your kitex options")
			}
		}

		cwgoManifest.KitexInfo.Version = args.Version
		cwgoManifest.KitexInfo.ServiceName = args.ServiceName
		cwgoManifest.HzInfo = meta.Manifest{}

		if c.Template == "" {
			// initialize cwgo side generator parameters
			if err = server.ConvertGenerator(serverGen, c); err != nil {
				return err
			}

			// generate cwgo side files
			if err = server.GenerateServer(serverGen); err != nil {
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
		manifest := new(meta.Manifest)

		if serverGen.IsNew {
			args.CmdType = meta.CmdNew
			if c.GoMod == "" {
				return fmt.Errorf("output directory %s is not under GOPATH/src. Please specify a module name with the '-module' flag", c.Cwd)
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
			err = app.GenerateLayout(args)
			if err != nil {
				return cli.Exit(err, meta.GenerateLayoutError)
			}

			args.InitManifest(manifest)
		} else {
			args.CmdType = meta.CmdUpdate

			module, path, ok := utils.SearchGoMod(workPath, true)
			if ok {
				// go.mod exists
				if c.GoMod != "" && module != c.GoMod {
					return fmt.Errorf("module name given by the '-module' option ('%s') is not consist with the name defined in go.mod ('%s' from %s)", c.GoMod, module, path)
				}
				args.Gomod = module
				c.GoModPath = path
			} else {
				return fmt.Errorf("go.mod not found in %s", workPath)
			}

			args.UpdateByManifest(manifest)
		}

		err = app.TriggerPlugin(args)
		if err != nil {
			return cli.Exit(err, meta.PluginError)
		}

		cwgoManifest.HzInfo = *manifest
		cwgoManifest.KitexInfo = cwgoMeta.KitexInfo{}

		if c.Template == "" {
			// initialize cwgo side generator parameters
			if err = server.ConvertGenerator(serverGen, c); err != nil {
				return err
			}

			// generate cwgo side files
			if err = server.GenerateServer(serverGen); err != nil {
				return cli.Exit(err, consts.GenerateCwgoError)
			}
		}
	}

	// generate .cwgo file
	if err = cwgoManifest.Persist(c.OutDir); err != nil {
		return err
	}

	utils.ReplaceThriftVersion(c.GoModPath)
	return nil
}
