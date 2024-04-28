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

func convertHzArgument(ca *config.ClientArgument, hzArgument *hzConfig.Argument) (err error) {
	// Common commands
	abPath, err := filepath.Abs(ca.IdlPath)
	if err != nil {
		return fmt.Errorf("idl path %s is not absolute", ca.IdlPath)
	}

	if strings.HasSuffix(ca.Template, consts.SuffixGit) {
		err = utils.GitClone(ca.Template, path.Join(tpl.HertzDir, consts.Client))
		if err != nil {
			return err
		}
		gitPath, err := utils.GitPath(ca.Template)
		if err != nil {
			return err
		}
		gitPath = path.Join(tpl.HertzDir, consts.Client, gitPath)
		if err = utils.GitCheckout(ca.Branch, gitPath); err != nil {
			return err
		}
		hzArgument.CustomizePackage = path.Join(gitPath, consts.PackageLayoutFile)
	} else {
		if len(ca.Template) != 0 {
			hzArgument.CustomizePackage = path.Join(ca.Template, consts.PackageLayoutFile)
		} else {
			hzArgument.CustomizePackage = path.Join(tpl.HertzDir, consts.Client, consts.Standard, consts.PackageLayoutFile)
		}
	}

	hzArgument.IdlPaths = []string{abPath}
	hzArgument.Gomod = ca.GoMod
	hzArgument.ServiceName = ca.Service
	hzArgument.Includes = ca.SliceParam.ProtoSearchPath
	hzArgument.Cwd = ca.Cwd
	hzArgument.Gosrc = ca.GoSrc
	hzArgument.Gopkg = ca.GoPkg
	hzArgument.Gopath = ca.GoPath
	hzArgument.Verbose = ca.Verbose
	hzArgument.OutDir = ca.Cwd
	// Automatic judgment param
	hzArgument.IdlType, err = utils.GetIdlType(abPath)
	if err != nil {
		return
	}

	// specific commands from -pass param
	f := flag.NewFlagSet("", flag.ContinueOnError)
	f.StringVar(&hzArgument.HandlerDir, "handler_dir", "", "")
	f.StringVar(&hzArgument.ModelDir, "model_dir", consts.DefaultHZModelDir, "")
	f.StringVar(&hzArgument.ClientDir, "client_dir", ca.OutDir, "")
	f.StringVar(&hzArgument.Use, "use", "", "")
	f.StringVar(&hzArgument.BaseDomain, "base_domain", "", "")
	var excludeFile, thriftgo, protoc, thriftPlugins, protocPlugins utils.FlagStringSlice
	f.Var(&excludeFile, "exclude_file", "")
	f.Var(&thriftgo, "thriftgo", "")
	f.Var(&protoc, "protoc", "")
	f.Var(&thriftPlugins, "thrift-plugins", "")
	f.Var(&protocPlugins, "protoc-plugins", "")
	f.BoolVar(&hzArgument.NoRecurse, "no_recurse", false, "")
	f.BoolVar(&hzArgument.JSONEnumStr, "json_enumstr", false, "")
	f.BoolVar(&hzArgument.UnsetOmitempty, "unset_omitempty", false, "")
	f.BoolVar(&hzArgument.ProtobufCamelJSONTag, "pb_camel_json_tag", false, "")
	f.BoolVar(&hzArgument.SnakeName, "snake_tag", false, "")
	f.BoolVar(&hzArgument.HandlerByMethod, "handler_by_method", false, "")

	err = f.Parse(utils.StringSliceSpilt(ca.SliceParam.Pass))
	if err != nil {
		return err
	}
	hzArgument.Excludes = excludeFile
	hzArgument.ThriftOptions = thriftgo
	hzArgument.ProtocOptions = protoc
	hzArgument.ThriftPlugins = thriftPlugins
	hzArgument.ProtobufPlugins = protocPlugins
	return nil
}
