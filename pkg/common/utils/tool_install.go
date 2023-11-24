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

package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/cloudwego/cwgo/pkg/consts"
	"github.com/cloudwego/hertz/cmd/hz/util"
	"github.com/cloudwego/hertz/cmd/hz/util/logs"
)

// installAndCheckGofumpt will automatically install gofumpt and judge whether it is installed successfully.
func installAndCheckGofumpt() error {
	exe, err := exec.LookPath("go")
	if err != nil {
		return fmt.Errorf("can not find tool 'go': %v", err)
	}
	var buf strings.Builder
	cmd := &exec.Cmd{
		Path: exe,
		Args: []string{
			exe, "install", "mvdan.cc/gofumpt@latest",
		},
		Stdin:  os.Stdin,
		Stdout: &buf,
		Stderr: &buf,
	}

	done := make(chan error)
	logs.Infof("installing gofumpt automatically")
	go func() {
		done <- cmd.Run()
	}()
	select {
	case err = <-done:
		if err != nil {
			return fmt.Errorf("can not install gofumpt, err: %v", cmd.Stderr)
		}
	case <-time.After(time.Second * 30):
		return fmt.Errorf("install gofumpt time out.Please install it manual")
	}

	return nil
}

func LookupTool(tool string) (string, error) {
	path, err := exec.LookPath(tool)
	if err != nil {
		goPath, err := util.GetGOPATH()
		if err != nil {
			return "", fmt.Errorf("get 'GOPATH' failed for find %s : %v", tool, path)
		}
		path = filepath.Join(goPath, "bin", tool)
	}

	isExist, err := util.PathExist(path)
	if err != nil {
		return "", fmt.Errorf("check '%s' path error: %v", path, err)
	}

	if !isExist {
		if tool == consts.Gofumpt {
			err = installAndCheckGofumpt()
			if err != nil {
				return "", fmt.Errorf("can't install '%s' automatically, please install it manually, err : %v", tool, err)
			}
		}
	}

	return path, nil
}
