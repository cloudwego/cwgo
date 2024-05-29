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
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloudwego/cwgo/config"
	"github.com/cloudwego/cwgo/pkg/common/utils"
	"github.com/cloudwego/cwgo/pkg/consts"
)

func check(ca *config.ClientArgument) error {
	if ca.Type != consts.RPC && ca.Type != consts.HTTP {
		return errors.New("generate type not supported")
	}

	if ca.Registry != "" &&
		ca.Registry != consts.Zk &&
		ca.Registry != consts.Nacos &&
		ca.Registry != consts.Etcd &&
		ca.Registry != consts.Polaris {
		return errors.New("unsupported registry")
	}

	if ca.Service == "" && ca.ServerName == "" {
		return errors.New("must specify server name when use registry")
	}

	// handle cwd and output dir
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get current path failed: %s", err)
	}
	ca.Cwd = dir
	if ca.OutDir == "" {
		if strings.ToUpper(ca.Type) == consts.HTTP {
			ca.OutDir = consts.DefaultHZClientDir
		} else {
			ca.OutDir = dir
		}
	}
	if !filepath.IsAbs(ca.OutDir) {
		ap := filepath.Join(ca.Cwd, ca.OutDir)
		ca.OutDir = ap
	}

	gopath, err := utils.GetGOPATH()
	if err != nil {
		return fmt.Errorf("get gopath failed: %s", err)
	}
	if gopath == "" {
		return fmt.Errorf("GOPATH is not set")
	}

	ca.GoPath = gopath
	ca.GoSrc = filepath.Join(gopath, consts.Src)

	// Generate the project under gopath, use the relative path as the package name
	if strings.HasPrefix(ca.Cwd, ca.GoSrc) {
		if goPkg, err := filepath.Rel(ca.GoSrc, ca.Cwd); err != nil {
			return fmt.Errorf("get relative path to GOPATH/src failed: %s", err)
		} else {
			ca.GoPkg = goPkg
		}

		if ca.GoMod == "" {
			if utils.IsWindows() {
				ca.GoMod = strings.ReplaceAll(ca.GoPkg, consts.BackSlash, consts.Slash)
			} else {
				ca.GoMod = ca.GoPkg
			}
		}

		if ca.GoMod != "" {
			if utils.IsWindows() {
				goPkgSlash := strings.ReplaceAll(ca.GoPkg, consts.BackSlash, consts.Slash)
				if goPkgSlash != ca.GoMod {
					return fmt.Errorf("module name: %s is not the same with GoPkg under GoPath: %s", ca.GoMod, goPkgSlash)
				}
			} else {
				if ca.GoMod != ca.GoPkg {
					return fmt.Errorf("module name: %s is not the same with GoPkg under GoPath: %s", ca.GoMod, ca.GoPkg)
				}
			}
		}
	}
	return nil
}
