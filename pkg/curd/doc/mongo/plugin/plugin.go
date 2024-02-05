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
	"github.com/cloudwego/cwgo/pkg/common/parser"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/cloudwego/cwgo/pkg/curd/doc/mongo/codegen"
	"github.com/cloudwego/cwgo/pkg/curd/extract"
	"github.com/cloudwego/cwgo/pkg/curd/parse"
	"github.com/cloudwego/cwgo/pkg/curd/template"

	"github.com/cloudwego/cwgo/config"
	"github.com/cloudwego/cwgo/pkg/common/utils"
	"github.com/cloudwego/cwgo/pkg/consts"
	"github.com/cloudwego/hertz/cmd/hz/meta"
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

	if c.IdlType == meta.IdlProto {
		info := &extract.PbUsedInfo{
			DocArgs: c,
		}
		rawStructs, err := info.ParsePbIdl()
		if err != nil {
			return err
		}
		operations, err := parse.HandleOperations(rawStructs)
		if err != nil {
			return err
		}
		methodRenders := codegen.HandleCodegen(operations)
		if err = info.GeneratePbFile(); err != nil {
			return err
		}
		if err = generatePbMongoFile(rawStructs, methodRenders, info); err != nil {
			return err
		}
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

	path, err := utils.LookupTool(args.IdlType)
	if err != nil {
		return nil, err
	}
	cmd := &exec.Cmd{
		Path: path,
	}

	if args.IdlType == meta.IdlThrift {
		os.Setenv(consts.CwgoDocPluginMode, consts.ThriftCwgoDocPluginName)

		thriftOpt, err := args.GetThriftgoOptions(args.PackagePrefix)
		if err != nil {
			return nil, err
		}

		cmd.Args = append(cmd.Args, meta.TpCompilerThrift)
		if args.Verbose {
			cmd.Args = append(cmd.Args, "-v")
		}
		cmd.Args = append(cmd.Args,
			"-o", args.ModelDir,
			"-p", "cwgo-doc="+exe+":"+kas,
			"-g", thriftOpt,
			"-r",
			args.IdlPath,
		)
	} else {
		cmd.Args = append(cmd.Args, meta.TpCompilerProto)
		for _, inc := range args.ProtoSearchPath {

			idlParser := parser.NewProtoParser()
			importBaseDirPath, importPaths, _ := idlParser.GetDependentFilePaths(inc, args.IdlPath)
			cmd.Args = append(cmd.Args, "-I", importBaseDirPath)
			cmd.Args = append(cmd.Args, importPaths...)
		}
		cmd.Args = append(cmd.Args, "--go_out="+args.ModelDir)
		for _, kv := range args.ProtocOptions {
			cmd.Args = append(cmd.Args, "--"+kv)
		}

		cmd.Args = append(cmd.Args, args.IdlPath)
	}

	return cmd, err
}

func MongoPluginMode() {
	mode := os.Getenv(consts.CwgoDocPluginMode)
	if len(os.Args) <= 1 && mode != "" {
		switch mode {
		case consts.ThriftCwgoDocPluginName:
			os.Exit(thriftPluginRun())
		}
	}
}

func generatePbMongoFile(structs []*extract.IdlExtractStruct, methodRenders [][]*template.MethodRender, info *extract.PbUsedInfo) error {
	for index, st := range structs {
		// get base render
		baseRender := getBaseRender(st)
		// get fileMongoName and fileIfName
		fileMongoName, fileIfName := extract.GetFileName(st.Name, info.DocArgs.DaoDir)
		if isExist, _ := utils.PathExist(filepath.Dir(fileMongoName)); !isExist {
			if err := os.MkdirAll(filepath.Dir(fileMongoName), 0o755); err != nil {
				return err
			}
		}
		if isExist, _ := utils.PathExist(filepath.Dir(fileIfName)); !isExist {
			if err := os.MkdirAll(filepath.Dir(fileIfName), 0o755); err != nil {
				return err
			}
		}

		if st.Update {
			// build update mongo file
			formattedCode, err := getUpdateMongoCode(methodRenders[index], string(st.UpdateCurdFileContent))
			if err != nil {
				return err
			}
			formattedCode, err = codegen.AddMongoImports(formattedCode)
			if err != nil {
				return err
			}
			formattedCode, err = extract.AddMongoModelImports(formattedCode, info.ImportPaths)
			if err != nil {
				return err
			}
			if err = utils.CreateFile(fileMongoName, formattedCode); err != nil {
				return err
			}

			// build update interface file
			formattedCode, err = getUpdateIfCode(st, baseRender)
			if err != nil {
				return err
			}
			formattedCode, err = codegen.AddMongoImports(formattedCode)
			if err != nil {
				return err
			}
			formattedCode, err = extract.AddMongoModelImports(formattedCode, info.ImportPaths)
			if err != nil {
				return err
			}
			if err = utils.CreateFile(fileIfName, formattedCode); err != nil {
				return err
			}
		} else {
			// build new mongo file
			formattedCode, err := getNewMongoCode(methodRenders[index], st, baseRender)
			if err != nil {
				return err
			}
			formattedCode, err = codegen.AddMongoImports(formattedCode)
			if err != nil {
				return err
			}
			formattedCode, err = extract.AddMongoModelImports(formattedCode, info.ImportPaths)
			if err != nil {
				return err
			}
			if err = utils.CreateFile(fileMongoName, formattedCode); err != nil {
				return err
			}

			// build new interface file
			formattedCode, err = getNewIfCode(st, baseRender)
			if err != nil {
				return err
			}
			formattedCode, err = codegen.AddMongoImports(formattedCode)
			if err != nil {
				return err
			}
			formattedCode, err = extract.AddMongoModelImports(formattedCode, info.ImportPaths)
			if err != nil {
				return err
			}
			if err = utils.CreateFile(fileIfName, formattedCode); err != nil {
				return err
			}
		}
	}

	return nil
}
