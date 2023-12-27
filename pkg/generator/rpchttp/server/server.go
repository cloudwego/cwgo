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

package server

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

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

	KitexIdlServiceName string

	GoFileImports imports.Map // handle .go files imports

	Extension
}

type Extension struct {
	Registry `yaml:"registry,omitempty"`
}

type Registry struct {
	RegistryName    string            `yaml:"registry_name,omitempty"`
	RegistryImports map[string]string `yaml:"registry_imports,omitempty"`
	RegistryBody    string            `yaml:"registry_body,omitempty"`
	RegistryAddress []string          `yaml:"registry_address,omitempty"`
	RegistryDocker  string            `yaml:"registry_docker"`
}

type fileInfo struct {
	mainGoContent        string
	confGoContent        string
	dockerComposeContent string
}

func NewGenerator(types string) (*Generator, error) {
	switch types {
	case consts.RPC:
		impts, err := imports.NewMap(consts.Server, consts.RPC)
		if err != nil {
			return nil, err
		}
		return &Generator{
			Render: Render{
				GoFileImports: impts,
				Extension: Extension{
					Registry{
						RegistryImports: map[string]string{},
					},
				},
			},
		}, nil

	case consts.HTTP:
		impts, err := imports.NewMap(consts.Server, consts.HTTP)
		if err != nil {
			return nil, err
		}
		return &Generator{
			Render: Render{
				GoFileImports: impts,
				Extension: Extension{
					Registry{
						RegistryImports: map[string]string{},
					},
				},
			},
		}, nil

	default:
		return nil, rhCommon.ErrTypeInput
	}
}

func ConvertGenerator(serverGen *Generator, args *config.ServerArgument) (err error) {
	// handle initial ServerGenerator arguments
	if err = serverGen.handleInitArguments(args); err != nil {
		return err
	}

	// handle initial go files imports
	if err = serverGen.handleInitImports(); err != nil {
		return err
	}

	// registry information
	if err = serverGen.handleRegistry(args.Registry); err != nil {
		return err
	}

	return nil
}

func (serverGen *Generator) handleInitArguments(args *config.ServerArgument) (err error) {
	serverGen.GoModule = args.GoMod
	serverGen.ServiceName = args.Service
	serverGen.CommunicationType = args.Type
	serverGen.CustomExtensionFile = args.CustomExtension
	serverGen.OutDir = args.OutDir
	serverGen.GoModPath = args.GoModPath
	serverGen.Codec, err = utils.GetIdlType(args.IdlPath)
	if err != nil {
		return err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	if serverGen.PackagePrefix, err = filepath.Rel(serverGen.GoModPath, cwd); err != nil {
		log.Warn("Get package prefix failed:", err.Error())
		os.Exit(1)
	}
	serverGen.PackagePrefix = filepath.Join(serverGen.GoModule, serverGen.PackagePrefix)
	if utils.IsWindows() {
		serverGen.PackagePrefix = strings.ReplaceAll(serverGen.PackagePrefix, consts.BackSlash, consts.Slash)
	}

	// handle custom extension
	if serverGen.CustomExtensionFile != "" {
		// parse custom extension file from yaml to go struct
		if err = serverGen.fromYAMLFile(serverGen.CustomExtensionFile); err != nil {
			return err
		}
		// check custom extension file
		if err = serverGen.checkCustomExtensionFile(); err != nil {
			return err
		}
	}

	if serverGen.IsNew && serverGen.CommunicationType == consts.RPC {
		handlerPath := filepath.Join(serverGen.OutDir, "handler.go")
		if isExist, _ := utils.PathExist(handlerPath); isExist {
			content, err := utils.ReadFileContent(handlerPath)
			if err != nil {
				return err
			}
			result, err := geneUtils.GetStructNames(string(content))
			if err != nil {
				return err
			}
			serverGen.KitexIdlServiceName = result[0][:len(result[0])-4]
		}
	}

	mainGoPath := filepath.Join(serverGen.OutDir, consts.Main)
	if isExist, _ := utils.PathExist(mainGoPath); isExist {
		content, err := utils.ReadFileContent(mainGoPath)
		if err != nil {
			return err
		}
		serverGen.mainGoContent = string(content)
	}

	confGoPath := filepath.Join(serverGen.OutDir, consts.ConfGo)
	if isExist, _ := utils.PathExist(confGoPath); isExist {
		content, err := utils.ReadFileContent(confGoPath)
		if err != nil {
			return err
		}
		serverGen.confGoContent = string(content)
	}

	dockerComposePath := filepath.Join(serverGen.OutDir, consts.DockerCompose)
	if isExist, _ := utils.PathExist(dockerComposePath); isExist {
		content, err := utils.ReadFileContent(dockerComposePath)
		if err != nil {
			return err
		}
		serverGen.dockerComposeContent = string(content)
	}

	return
}

func (serverGen *Generator) handleInitImports() (err error) {
	dalInitExtraImports := map[string]string{
		serverGen.PackagePrefix + "/biz/dal/mysql": "",
		serverGen.PackagePrefix + "/biz/dal/redis": "",
	}
	if err = serverGen.GoFileImports.AppendImports(consts.DalInitGo, dalInitExtraImports); err != nil {
		return err
	}

	mysqlInitExtraImports := map[string]string{
		serverGen.PackagePrefix + "/conf": "",
	}
	if err = serverGen.GoFileImports.AppendImports(consts.MysqlInit, mysqlInitExtraImports); err != nil {
		return err
	}

	redisInitExtraImports := map[string]string{
		serverGen.PackagePrefix + "/conf": "",
	}
	if err = serverGen.GoFileImports.AppendImports(consts.RedisInit, redisInitExtraImports); err != nil {
		return err
	}

	switch serverGen.CommunicationType {
	case consts.RPC:
		// set initial main.go imports
		dirs, err := utils.GetSubDirs(serverGen.OutDir, true)
		if err != nil {
			return err
		}

		var dir string
		for _, d := range dirs {
			if filepath.Base(d) == strings.ToLower(serverGen.KitexIdlServiceName) {
				dir, err = filepath.Rel(dirs[0], d)
				if err != nil {
					return err
				}
				if utils.IsWindows() {
					dir = strings.ReplaceAll(dir, consts.BackSlash, consts.Slash)
				}
				break
			}
		}

		mainExtraImports := map[string]string{
			serverGen.PackagePrefix + "/conf":   "",
			serverGen.PackagePrefix + "/" + dir: "",
		}
		if serverGen.Codec == "thrift" {
			mainExtraImports["github.com/cloudwego/kitex/pkg/transmeta"] = ""
		}
		if err = serverGen.GoFileImports.AppendImports(consts.Main, mainExtraImports); err != nil {
			return err
		}

	case consts.HTTP:
		// set initial main.go imports
		mainExtraImports := map[string]string{
			serverGen.PackagePrefix + "/biz/router": "",
			serverGen.PackagePrefix + "/conf":       "",
		}
		if err = serverGen.GoFileImports.AppendImports(consts.Main, mainExtraImports); err != nil {
			return
		}

	default:
		return rhCommon.ErrTypeInput
	}

	return
}

func (serverGen *Generator) compatibleOlderVersion() error {
	var (
		mvcTemplates           []template.Template
		appendInitRegistryFunc string
	)

	if serverGen.CommunicationType == consts.RPC {
		mvcTemplates = kitexServerMVCTemplates
		appendInitRegistryFunc = kitexNilAppendInitRegistryFunc
	} else {
		mvcTemplates = hzServerMVCTemplates
		appendInitRegistryFunc = hzNilAppendInitRegistryFunc
	}

	// compatible registry
	if serverGen.mainGoContent != "" {
		isExist, err := geneUtils.IsFuncExist(serverGen.mainGoContent, consts.FuncInitRegistry)
		if err != nil {
			return err
		}
		if !isExist {
			mvcTemplates[consts.FileServerMainIndex].Type = consts.Append
			mvcTemplates[consts.FileServerMainIndex].AppendContent += appendInitRegistryFunc + consts.LineBreak
		}
	}

	if serverGen.CommunicationType == consts.RPC {
		kitexServerMVCTemplates = mvcTemplates
	} else {
		hzServerMVCTemplates = mvcTemplates
	}

	return nil
}

func (s *Extension) fromYAMLFile(filename string) error {
	if s == nil {
		return nil
	}
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, s)
}

func (s *Extension) checkCustomExtensionFile() (err error) {
	// check reso
	if s.RegistryName != "" {
		if s.RegistryImports == nil {
			return errors.New("please input RegistryImports")
		}
		if s.RegistryBody == "" {
			return errors.New("please input RegistryImports")
		}
		if len(s.RegistryAddress) == 0 {
			return errors.New("please input RegistryAddress")
		}
	}

	return nil
}

// GenerateServer generate cwgo side server files
func GenerateServer(serverGen *Generator) error {
	tg := &template.Generator{
		Files: make([]template.File, 0, 10),
	}

	// render template
	if err := renderServer(tg, serverGen); err != nil {
		return err
	}

	// generate files
	if err := tg.Persist(); err != nil {
		return err
	}

	return nil
}

func renderServer(tg *template.Generator, serverGen *Generator) (err error) {
	var mvcTemplates []template.Template

	switch serverGen.CommunicationType {
	case consts.RPC:
		mvcTemplates = kitexServerMVCTemplates
	case consts.HTTP:
		mvcTemplates = hzServerMVCTemplates
	default:
		return rhCommon.ErrTypeInput
	}

	for _, tpl := range mvcTemplates {
		// handle append render
		if tpl.Type == consts.Append || tpl.Type == consts.ReplaceFuncBody {
			handleAppendRender(&tpl, serverGen)
		}

		// render cwgo server update
		if err = tg.RenderCwgoTemplateFile(&tpl, serverGen.Render); err != nil {
			return err
		}
	}

	return nil
}

func handleAppendRender(tpl *template.Template, serverGen *Generator) {
	tpl.AppendRender["GoModule"] = serverGen.GoModule
	tpl.AppendRender["PackagePrefix"] = serverGen.PackagePrefix
	tpl.AppendRender["ServiceName"] = serverGen.ServiceName
	if serverGen.RegistryName != "" {
		tpl.AppendRender["RegistryName"] = serverGen.RegistryName
	}
	if serverGen.RegistryBody != "" {
		tpl.AppendRender["RegistryBody"] = serverGen.RegistryBody
	}
	if serverGen.RegistryAddress != nil {
		tpl.AppendRender["RegistryAddress"] = serverGen.RegistryAddress
	}
	if serverGen.RegistryDocker != "" {
		tpl.AppendRender["RegistryDocker"] = serverGen.RegistryDocker
	}
}
