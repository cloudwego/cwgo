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
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloudwego/cwgo/config"
	"github.com/cloudwego/hertz/cmd/hz/util"
)

func check(sa *config.ServerArgument) error {
	if sa.Type != config.RPC && sa.Type != config.HTTP {
		return errors.New("generate type not supported")
	}

	if sa.Registry != "" &&
		sa.Registry != config.Zk &&
		sa.Registry != config.Nacos &&
		sa.Registry != config.Etcd &&
		sa.Registry != config.Polaris {
		return errors.New("unsupported registry")
	}

	if sa.Service == "" {
		return errors.New("must specify service name")
	}

	// handle cwd and output dir
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get current path failed: %s", err)
	}
	sa.Cwd = dir
	if sa.OutDir == "" {
		sa.OutDir = dir
	}
	if !filepath.IsAbs(sa.OutDir) {
		ap := filepath.Join(sa.Cwd, sa.OutDir)
		sa.OutDir = ap
	}

	gopath, err := util.GetGOPATH()
	if err != nil {
		return fmt.Errorf("get gopath failed: %s", err)
	}
	if gopath == "" {
		return fmt.Errorf("GOPATH is not set")
	}

	sa.GoPath = gopath
	sa.GoSrc = filepath.Join(gopath, "src")

	// Generate the project under gopath, use the relative path as the package name
	if strings.HasPrefix(sa.Cwd, sa.GoSrc) {
		if gopkg, err := filepath.Rel(sa.GoSrc, sa.Cwd); err != nil {
			return fmt.Errorf("get relative path to GOPATH/src failed: %s", err)
		} else {
			sa.GoPkg = gopkg
		}
		if sa.GoMod == "" {
			sa.GoMod = sa.GoPkg
		}
		if sa.GoMod != "" && sa.GoMod != sa.GoPkg {
			return fmt.Errorf("module name: %s is not the same with GoPkg under GoPath: %s", sa.GoMod, sa.GoPkg)
		}
		if sa.GoMod == "" {
			sa.GoMod = sa.GoPkg
		}
	}
	return nil
}
