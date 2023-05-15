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

	"github.com/cloudwego/cwgo/config"
	"github.com/cloudwego/cwgo/pkg/common/kx_registry"
	"github.com/cloudwego/cwgo/pkg/common/utils"
	"github.com/cloudwego/hertz/cmd/hz/app"
	hzConfig "github.com/cloudwego/hertz/cmd/hz/config"
	"github.com/cloudwego/hertz/cmd/hz/meta"
	"github.com/cloudwego/hertz/cmd/hz/util"
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

	switch c.Type {
	case config.RPC:
		var args kargs.Arguments
		log.Verbose = c.Verbose
		err = convertKitexArgs(c, &args)
		if err != nil {
			return err
		}
		kx_registry.HandleRegistry(c.CommonParam, args.TemplateDir)
		defer kx_registry.RemoveExtension()

		out := new(bytes.Buffer)
		cmd := args.BuildCmd(out)
		err = cmd.Run()
		if err != nil {
			if args.Use != "" {
				out := strings.TrimSpace(out.String())
				if strings.HasSuffix(out, thriftgo.TheUseOptionMessage) {
					replaceThriftVersion(&args)
				}
			}
			os.Exit(1)
		}
		replaceThriftVersion(&args)
	case config.HTTP:
		args := hzConfig.NewArgument()
		utils.SetHzVerboseLog(c.Verbose)
		err = convertHzArgument(c, args)
		if err != nil {
			return err
		}

		if utils.IsHzNew(c.OutDir) {
			args.CmdType = meta.CmdNew
			if c.GoMod == "" {
				return fmt.Errorf("output directory %s is not under GOPATH/src. Please specify a module name with the '-module' flag", c.Cwd)
			}
			module, path, ok := util.SearchGoMod(".", false)
			if ok {
				// go.mod exists
				if module != c.GoMod {
					return fmt.Errorf("module name given by the '-module' option ('%s') is not consist with the name defined in go.mod ('%s' from %s)", c.GoMod, module, path)
				}
				c.GoMod = module
			} else {
				args.NeedGoMod = true
			}
			err = app.GenerateLayout(args)
			if err != nil {
				return cli.Exit(err, meta.GenerateLayoutError)
			}
			defer func() {
				// ".hz" file converges to the hz tool
				manifest := new(meta.Manifest)
				args.InitManifest(manifest)
				err = manifest.Persist(args.OutDir)
				if err != nil {
					err = cli.Exit(fmt.Errorf("persist manifest failed: %v", err), meta.PersistError)
				}
			}()
		} else {
			args.CmdType = meta.CmdUpdate
			manifest := new(meta.Manifest)
			err = manifest.InitAndValidate(args.OutDir)
			if err != nil {
				return cli.Exit(err, meta.LoadError)
			}

			module, path, ok := util.SearchGoMod(".", false)
			if ok {
				// go.mod exists
				if c.GoMod != "" && module != c.GoMod {
					return fmt.Errorf("module name given by the '-module' option ('%s') is not consist with the name defined in go.mod ('%s' from %s)", c.GoMod, module, path)
				}
				args.Gomod = module
			} else {
				workPath, err := filepath.Abs(".")
				if err != nil {
					return fmt.Errorf(err.Error())
				}
				return fmt.Errorf("go.mod not found in %s", workPath)
			}

			// update argument by ".hz", can automatically get "handler_dir"/"model_dir"/"router_dir"
			args.UpdateByManifest(manifest)

			defer func() {
				// If the "handler_dir"/"model_dir" is updated, write it back to ".hz"
				args.UpdateManifest(manifest)
				err = manifest.Persist(args.OutDir)
				if err != nil {
					err = cli.Exit(fmt.Errorf("persist manifest failed: %v", err), meta.PersistError)
				}
			}()
		}

		err = app.TriggerPlugin(args)
		if err != nil {
			return cli.Exit(err, meta.PluginError)
		}
	}

	return nil
}
