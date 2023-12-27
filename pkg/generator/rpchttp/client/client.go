/*
 * Copyright 2023 CloudWeGo Authors
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
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/cloudwego/cwgo/config"
	"github.com/cloudwego/cwgo/pkg/common/utils"
	"github.com/cloudwego/cwgo/pkg/consts"
	"github.com/cloudwego/cwgo/pkg/generator/common/imports"
	"github.com/cloudwego/cwgo/pkg/generator/common/template"
	geneUtils "github.com/cloudwego/cwgo/pkg/generator/common/utils"
	rhCommon "github.com/cloudwego/cwgo/pkg/generator/rpchttp/common"
	"github.com/cloudwego/kitex/tool/internal_pkg/log"
	"gopkg.in/yaml.v2"
)

type Generator struct {
	rhCommon.ServerClientComGen // common generator params

	Render // for template render

	fileInfo
}

type Render struct {
	GoModule      string
	PackagePrefix string
	ServiceName   string
	Codec         string

	InitOptsPackage string

	GoFileImports imports.Map // handle .go files imports

	Extension
}

type Extension struct {
	Resolver `yaml:"resolver,omitempty"`
}

type Resolver struct {
	ResolverName    string            `yaml:"resolver_name,omitempty"`
	ResolverImports map[string]string `yaml:"resolver_imports,omitempty"`
	ResolverBody    string            `yaml:"resolver_body,omitempty"`
	ResolverAddress []string          `yaml:"resolver_address,omitempty"`
}

type fileInfo struct {
	once           sync.Once
	bizDir         string
	subDirs        []string
	initGoContents []string
	initGoPaths    []string
	envGoContent   string
	envGoPath      string
}

func NewGenerator(types string) (*Generator, error) {
	switch types {
	case consts.RPC:
		impts, err := imports.NewMap(consts.Client, consts.RPC)
		if err != nil {
			return nil, err
		}
		return &Generator{
			Render: Render{
				GoFileImports: impts,
				Extension: Extension{
					Resolver{
						ResolverImports: map[string]string{},
					},
				},
			},
			fileInfo: fileInfo{
				initGoPaths:    make([]string, 0, 5),
				initGoContents: make([]string, 0, 5),
			},
		}, nil

	case consts.HTTP:
		impts, err := imports.NewMap(consts.Client, consts.HTTP)
		if err != nil {
			return nil, err
		}
		return &Generator{
			Render: Render{
				GoFileImports: impts,
				Extension: Extension{
					Resolver{
						ResolverImports: map[string]string{},
					},
				},
			},
			fileInfo: fileInfo{
				initGoPaths:    make([]string, 0, 5),
				initGoContents: make([]string, 0, 5),
			},
		}, nil

	default:
		return nil, rhCommon.ErrTypeInput
	}
}

func ConvertGenerator(clientGen *Generator, args *config.ClientArgument) (err error) {
	// handle initial ClientGenerator arguments
	if err = clientGen.handleInitArguments(args); err != nil {
		return err
	}

	// handle initial go files imports when new
	if err = clientGen.handleInitImports(); err != nil {
		return err
	}

	// handle resolve information
	if err = clientGen.handleResolver(args.Resolver); err != nil {
		return err
	}

	return nil
}

func (clientGen *Generator) handleInitArguments(args *config.ClientArgument) (err error) {
	clientGen.GoModule = args.GoMod
	clientGen.ServiceName = args.Service
	clientGen.CommunicationType = args.Type
	clientGen.CustomExtensionFile = args.CustomExtension
	clientGen.OutDir = args.OutDir
	clientGen.GoModPath = args.GoModPath
	clientGen.Codec, err = utils.GetIdlType(args.IdlPath)
	if err != nil {
		return err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	if clientGen.PackagePrefix, err = filepath.Rel(clientGen.GoModPath, cwd); err != nil {
		log.Warn("Get package prefix failed:", err.Error())
		os.Exit(1)
	}
	clientGen.PackagePrefix = filepath.Join(clientGen.GoModule, clientGen.PackagePrefix)
	if utils.IsWindows() {
		clientGen.PackagePrefix = strings.ReplaceAll(clientGen.PackagePrefix, consts.BackSlash, consts.Slash)
	}

	// handle custom extension
	if clientGen.CustomExtensionFile != "" {
		// parse custom extension file from yaml to go struct
		if err = clientGen.fromYAMLFile(clientGen.CustomExtensionFile); err != nil {
			return err
		}
		// check custom extension file
		if err = clientGen.checkCustomExtensionFile(); err != nil {
			return err
		}
	}

	if err := clientGen.getDirInfo(); err != nil {
		return err
	}

	for _, subDir := range clientGen.subDirs {
		name := filepath.Base(subDir)
		filePath := filepath.Join(subDir, name+"_"+consts.InitGo)
		if isExist, _ := utils.PathExist(filePath); isExist {
			content, err := utils.ReadFileContent(filePath)
			if err != nil {
				return err
			}
			clientGen.initGoContents = append(clientGen.initGoContents, string(content))
			clientGen.initGoPaths = append(clientGen.initGoPaths, filePath)
		}
	}

	if len(clientGen.initGoPaths) != 0 {
		if clientGen.CommunicationType == consts.RPC {
			kitexClientMVCTemplates[consts.FileClientInitIndex].Path = clientGen.initGoPaths[0]
		} else {
			hzClientMVCTemplates[consts.FileClientInitIndex].Path = clientGen.initGoPaths[0]
		}
	}

	clientGen.envGoPath = filepath.Join(clientGen.bizDir, consts.EnvGo)
	if isExist, _ := utils.PathExist(clientGen.envGoPath); isExist {
		content, err := utils.ReadFileContent(clientGen.envGoPath)
		if err != nil {
			return err
		}
		clientGen.envGoContent = string(content)
	}

	return
}

func (clientGen *Generator) handleInitImports() (err error) {
	switch clientGen.CommunicationType {
	case consts.RPC:
		// set initial init.go imports
		if clientGen.Codec == "thrift" {
			initExtraImports := map[string]string{
				"github.com/cloudwego/kitex/pkg/transmeta": "",
				"github.com/cloudwego/kitex/transport":     "",
			}
			if err = clientGen.GoFileImports.AppendImports(consts.InitGo, initExtraImports); err != nil {
				return err
			}
		}
	case consts.HTTP:
	default:
		return rhCommon.ErrTypeInput
	}

	return
}

func (clientGen *Generator) compatibleOlderVersion() error {
	var (
		mvcTemplates           []template.Template
		appendInitResolverFunc string
	)

	if clientGen.CommunicationType == consts.RPC {
		mvcTemplates = kitexClientMVCTemplates
		appendInitResolverFunc = kitexNilAppendInitResolverFunc
	} else {
		mvcTemplates = hzClientMVCTemplates
		appendInitResolverFunc = hzNilAppendInitResolverFunc
	}

	// compatible registry
	if len(clientGen.initGoContents) != 0 {
		for index, content := range clientGen.initGoContents {
			isExist, err := geneUtils.IsFuncExist(content, consts.FuncInitResolver)
			if err != nil {
				return err
			}
			if !isExist && index == 0 {
				mvcTemplates[consts.FileClientInitIndex].Path = clientGen.initGoPaths[index]
				mvcTemplates[consts.FileClientInitIndex].Type = consts.Append
				mvcTemplates[consts.FileClientInitIndex].AppendContent += appendInitResolverFunc + consts.LineBreak
			}
			if !isExist && index != 0 {
				t := template.Template{
					Path: clientGen.initGoPaths[index],
					UpdateBehavior: template.UpdateBehavior{
						Type: consts.Append,
						Append: template.Append{
							AppendContent: appendInitResolverFunc,
						},
						AppendRender: map[string]interface{}{},
					},
				}

				mvcTemplates = append(mvcTemplates, t)
			}
		}
	}

	if clientGen.CommunicationType == consts.RPC {
		kitexClientMVCTemplates = mvcTemplates
	} else {
		hzClientMVCTemplates = mvcTemplates
	}

	return nil
}

func (c *Extension) fromYAMLFile(filename string) error {
	if c == nil {
		return nil
	}
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, c)
}

func (c *Extension) checkCustomExtensionFile() (err error) {
	// check resolver
	if c.ResolverName != "" {
		if c.ResolverImports == nil {
			return errors.New("please input ResolverImports")
		}
		if c.ResolverBody == "" {
			return errors.New("please input ResolverImports")
		}
		if len(c.ResolverAddress) == 0 {
			return errors.New("please input ResolverAddress")
		}
	}

	return nil
}

// GenerateClient generate cwgo side client files
func GenerateClient(clientGen *Generator) error {
	tg := &template.Generator{
		Files: make([]template.File, 0, 10),
	}

	// render template
	if err := renderClient(tg, clientGen); err != nil {
		return err
	}

	// generate files
	if err := tg.Persist(); err != nil {
		return err
	}

	return nil
}

func renderClient(tg *template.Generator, clientGen *Generator) error {
	var mvcTemplates []template.Template

	switch clientGen.CommunicationType {
	case consts.RPC:
		mvcTemplates = kitexClientMVCTemplates
	case consts.HTTP:
		mvcTemplates = hzClientMVCTemplates
	default:
		return rhCommon.ErrTypeInput
	}

	for index, tpl := range mvcTemplates {
		// render init.go
		if index == consts.FileClientInitIndex && tpl.Path == consts.InitGo && clientGen.IsNew {
			if err := clientGen.getDirInfo(); err != nil {
				return err
			}

			for _, name := range clientGen.subDirs {
				clientGen.InitOptsPackage = filepath.Base(name)

				// render body
				data, err := template.Render(name, tpl.Body, clientGen.Render, &tpl)
				if err != nil {
					return err
				}

				file := template.File{Path: filepath.Join(name, clientGen.InitOptsPackage+"_"+consts.InitGo), Content: data.Bytes()}
				tg.Files = append(tg.Files, file)
			}
			continue
		}

		// handle append render
		if tpl.Type == consts.Append || tpl.Type == consts.ReplaceFuncBody {
			handleAppendRender(&tpl, clientGen)
		}

		// render cwgo client update
		if err := tg.RenderCwgoTemplateFile(&tpl, clientGen.Render); err != nil {
			return err
		}
	}

	return nil
}

func handleAppendRender(tpl *template.Template, clientGen *Generator) {
	tpl.AppendRender["GoModule"] = clientGen.GoModule
	tpl.AppendRender["PackagePrefix"] = clientGen.PackagePrefix
	tpl.AppendRender["ServiceName"] = clientGen.ServiceName
	if clientGen.ResolverName != "" {
		tpl.AppendRender["ResolverName"] = clientGen.ResolverName
	}
	if clientGen.ResolverBody != "" {
		tpl.AppendRender["ResolverBody"] = clientGen.ResolverBody
	}
	if clientGen.ResolverAddress != nil {
		tpl.AppendRender["ResolverAddress"] = clientGen.ResolverAddress
	}
}

func (clientGen *Generator) getDirInfo() (err error) {
	clientGen.once.Do(func() {
		if clientGen.CommunicationType == consts.HTTP {
			clientGen.bizDir = clientGen.OutDir
		} else {
			clientGen.bizDir = filepath.Join(clientGen.OutDir, consts.DefaultKitexClientDir)
		}

		clientGen.subDirs, err = utils.GetSubDirs(clientGen.bizDir, false)
		return
	})

	return
}
