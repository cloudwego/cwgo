/*
 * Copyright 2024 CloudWeGo Authors
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

package doc

import (
	"errors"
	"path/filepath"

	"github.com/cloudwego/cwgo/config"
	"github.com/cloudwego/cwgo/pkg/consts"
	"github.com/cloudwego/cwgo/pkg/doc/mongo/plugin"
	"github.com/cloudwego/hertz/cmd/hz/util/logs"
)

func Doc(c *config.DocArgument) error {
	if err := check(c); err != nil {
		return err
	}

	switch c.Name {
	case consts.MongoDb:
		setLogVerbose(c.Verbose)
		if err := plugin.MongoTriggerPlugin(c); err != nil {
			return err
		}
	default:
	}

	return nil
}

func check(c *config.DocArgument) error {
	if c.Name == "" {
		c.Name = consts.MongoDb
	}
	if c.Name != consts.MongoDb {
		return errors.New("doc name not supported")
	}
	if c.IdlPath == "" {
		return errors.New("must specify idl path")
	}

	if c.ModelDir == "" {
		c.ModelDir = consts.DefaultDocModelOutDir
	}
	c.ModelDir = filepath.Join(c.OutDir, c.ModelDir)

	if c.DaoDir == "" {
		c.DaoDir = consts.DefaultDocDaoOutDir
	}
	c.DaoDir = filepath.Join(c.OutDir, c.DaoDir)

	return nil
}

func setLogVerbose(verbose bool) {
	if verbose {
		logs.SetLevel(logs.LevelDebug)
	} else {
		logs.SetLevel(logs.LevelWarn)
	}
}
