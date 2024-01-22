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

package plugin

import (
	"fmt"
	"github.com/cloudwego/cwgo/config"
	"github.com/cloudwego/cwgo/pkg/common/utils"
	"github.com/cloudwego/cwgo/pkg/consts"
	"github.com/cloudwego/hertz/cmd/hz/meta"
	"os"
	"os/exec"
	"strings"
)

func MongoTriggerPlugin(c *config.DocArgument) error {
	cmd, err := buildPluginCmd(c)
	if err != nil {
		return fmt.Errorf("build plugin command failed: %v", err)
	}

	buf, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("plugin cwgo-doc returns error: %v, cause:\n%v", err, string(buf))
	}

	// If len(buf) != 0, the plugin returned the log.
	if len(buf) != 0 {
		fmt.Println(string(buf))
	}
	return nil
}

func buildPluginCmd(args *config.DocArgument) (*exec.Cmd, error) {
	exe, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("failed to detect current executable, err: %v", err)
	}

	argPacks, err := args.Pack()
	if err != nil {
		return nil, err
	}
	kas := strings.Join(argPacks, ",")

	idlType, err := utils.GetIdlType(args.IdlPath)
	if err != nil {
		return nil, err
	}

	path, err := utils.LookupTool(idlType)
	if err != nil {
		return nil, err
	}
	cmd := &exec.Cmd{
		Path: path,
	}

	if idlType == meta.IdlThrift {
		os.Setenv(consts.CwgoDocPluginMode, consts.ThriftCwgoDocPluginName)

		cmd.Args = append(cmd.Args, meta.TpCompilerThrift)
		if args.Verbose {
			cmd.Args = append(cmd.Args, "-v")
		}
		cmd.Args = append(cmd.Args,
			"-o", args.ModelDir,
			"-p", "cwgo-doc="+exe+":"+kas,
			"-g", "go",
			"-r",
			args.IdlPath,
		)
	}

	return cmd, err
}

func MongoPluginMode() {
	mode := os.Getenv(consts.CwgoDocPluginMode)
	if len(os.Args) <= 1 && mode != "" {
		switch mode {
		case consts.ThriftCwgoDocPluginName:
			os.Exit(pluginRun())
		}
	}
}
