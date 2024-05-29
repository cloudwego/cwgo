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
	"flag"
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/cloudwego/cwgo/config"
	"github.com/cloudwego/cwgo/pkg/common/utils"
	"github.com/cloudwego/cwgo/pkg/consts"
	"github.com/cloudwego/cwgo/tpl"
	hzConfig "github.com/cloudwego/hertz/cmd/hz/config"
)

func convertHzArgument(sa *config.ServerArgument, hzArgument *hzConfig.Argument) (err error) {
	// Common commands
	abPath, err := filepath.Abs(sa.IdlPath)
	if err != nil {
		return fmt.Errorf("idl path %s is not absolute", sa.IdlPath)
	}

	if strings.HasSuffix(sa.Template, consts.SuffixGit) {
		err = utils.GitClone(sa.Template, path.Join(tpl.HertzDir, consts.Server))
		if err != nil {
			return err
		}
		gitPath, err := utils.GitPath(sa.Template)
		if err != nil {
			return err
		}
		gitPath = path.Join(tpl.HertzDir, consts.Server, gitPath)
		if err = utils.GitCheckout(sa.Branch, gitPath); err != nil {
			return err
		}
		hzArgument.CustomizeLayout = path.Join(gitPath, consts.LayoutFile)
		hzArgument.CustomizePackage = path.Join(gitPath, consts.PackageLayoutFile)
		layoutDataPath := path.Join(gitPath, "render.json")
		isExist, _ := utils.PathExist(layoutDataPath)
		if isExist {
			hzArgument.CustomizeLayoutData = layoutDataPath
		}
	} else {
		if len(sa.Template) != 0 {
			hzArgument.CustomizeLayout = path.Join(sa.Template, consts.LayoutFile)
			hzArgument.CustomizePackage = path.Join(sa.Template, consts.PackageLayoutFile)
			layoutDataPath := path.Join(sa.Template, "render.json")
			isExist, _ := utils.PathExist(layoutDataPath)
			if isExist {
				hzArgument.CustomizeLayoutData = layoutDataPath
			}
		} else {
			hzArgument.CustomizeLayout = path.Join(tpl.HertzDir, consts.Server, consts.Standard, consts.LayoutFile)
			hzArgument.CustomizePackage = path.Join(tpl.HertzDir, consts.Server, consts.Standard, consts.PackageLayoutFile)
		}
	}

	hzArgument.IdlPaths = []string{abPath}
	hzArgument.Gomod = sa.GoMod
	hzArgument.ServiceName = sa.ServerName
	hzArgument.OutDir = sa.OutDir
	hzArgument.Includes = sa.SliceParam.ProtoSearchPath
	hzArgument.Cwd = sa.Cwd
	hzArgument.Gosrc = sa.GoSrc
	hzArgument.Gopkg = sa.GoPkg
	hzArgument.Gopath = sa.GoPath
	hzArgument.Verbose = sa.Verbose
	// Automatic judgment param
	hzArgument.IdlType, err = utils.GetIdlType(abPath)
	if err != nil {
		return
	}

	// specific commands from -pass param
	f := flag.NewFlagSet("", flag.ContinueOnError)
	handlerDir := f.String("handler_dir", "", "")
	modelDir := f.String("model_dir", consts.DefaultHZModelDir, "")
	routerDir := f.String("router_dir", "", "")
	use := f.String("use", "", "")
	var excludeFile, thriftgo, protoc, thriftPlugins, protocPlugins utils.FlagStringSlice
	f.Var(&excludeFile, "exclude_file", "")
	f.Var(&thriftgo, "thriftgo", "")
	f.Var(&protoc, "protoc", "")
	f.Var(&thriftPlugins, "thrift-plugins", "")
	f.Var(&protocPlugins, "protoc-plugins", "")
	noRecurse := f.Bool("no_recurse", false, "")
	JSONEnumStr := f.Bool("json_enumstr", false, "")
	UnsetOmitempty := f.Bool("unset_omitempty", false, "")
	pbCamelJSONTag := f.Bool("pb_camel_json_tag", false, "")
	snakeTag := f.Bool("snake_tag", false, "")
	handlerByMethod := f.Bool("handler_by_method", false, "")

	err = f.Parse(utils.StringSliceSpilt(sa.SliceParam.Pass))
	if err != nil {
		return err
	}
	hzArgument.HandlerDir = *handlerDir
	hzArgument.ModelDir = *modelDir
	hzArgument.RouterDir = *routerDir
	hzArgument.Use = *use
	hzArgument.Excludes = excludeFile
	hzArgument.ThriftOptions = thriftgo
	hzArgument.ProtocOptions = protoc
	hzArgument.ThriftPlugins = thriftPlugins
	hzArgument.ProtobufPlugins = protocPlugins
	hzArgument.NoRecurse = *noRecurse
	hzArgument.JSONEnumStr = *JSONEnumStr
	hzArgument.UnsetOmitempty = *UnsetOmitempty
	hzArgument.ProtobufCamelJSONTag = *pbCamelJSONTag
	hzArgument.SnakeName = *snakeTag
	hzArgument.HandlerByMethod = *handlerByMethod
	return nil
}
