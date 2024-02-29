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

package api_list

import (
	"fmt"
	"path/filepath"

	"github.com/cloudwego/cwgo/config"
	"github.com/cloudwego/cwgo/pkg/consts"
)

func Api(c *config.ApiArgument) error {
	if c.ProjectPath == "" {
		curPath, err := filepath.Abs(".")
		if err != nil {
			return fmt.Errorf("get current path failed, err: %v", err)
		}
		c.ProjectPath = curPath
	}

	if c.HertzRepoUrl == "" {
		c.HertzRepoUrl = consts.HertzRepoDefaultUrl
	}

	parser, err := NewParser(c.ProjectPath, c.HertzRepoUrl)
	if err != nil {
		return err
	}

	moduleName, err := getModuleName(c.ProjectPath)
	if err != nil {
		return err
	}

	fmt.Printf("found module name: %s\n", parser.moduleName)

	err = parser.searchFunc(moduleName, "main", make(map[string]*Var), nil)
	if err != nil {
		return err
	}

	parser.PrintRouters()

	return nil
}
