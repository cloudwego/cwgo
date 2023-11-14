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

package generator

import (
	"errors"
	"io/ioutil"
	"path"
	"strings"

	"github.com/cloudwego/cwgo/config"
	"github.com/cloudwego/cwgo/meta"
	"github.com/cloudwego/cwgo/pkg/common/utils"
	"github.com/cloudwego/cwgo/pkg/consts"
	"github.com/cloudwego/kitex/tool/internal_pkg/generator"
	"gopkg.in/yaml.v2"
)

var (
	typeInputErr = errors.New("input wrong type")
)

type ServerGenerator struct {
	CommonGenerator

	ServerRender
}

type ServerRender struct {
	GoModule    string
	ServiceName string

	GoFileImports ImportsMap

	ServerExtension
}

type ServerExtension struct {
	Registry `yaml:"registry,omitempty"`
}

type Registry struct {
	// todo docker 内容和 imports 内容
	RegistryName        string `yaml:"registry_name,omitempty"`
	RegistryBody        string `yaml:"registry_body,omitempty"`
	DefaultRegistryPort string `yaml:"default_registry_port,omitempty"`
}

func NewServerGenerator(types string) (*ServerGenerator, error) {
	switch types {
	case consts.RPC:
		imports, err := newImportsMap(consts.Server, consts.RPC)
		if err != nil {
			return nil, err
		}
		return &ServerGenerator{
			CommonGenerator: CommonGenerator{
				kitexExtension: &generator.TemplateExtension{
					Dependencies: map[string]string{},
					ExtendServer: &generator.APIExtension{},
				},
				manifest: new(meta.Manifest),
			},
			ServerRender: ServerRender{
				GoFileImports: imports,
			},
		}, nil

	case consts.HTTP:
		imports, err := newImportsMap(consts.Server, consts.HTTP)
		if err != nil {
			return nil, err
		}
		return &ServerGenerator{
			CommonGenerator: CommonGenerator{
				manifest: new(meta.Manifest),
			},
			ServerRender: ServerRender{
				GoFileImports: imports,
			},
		}, nil

	default:
		return nil, typeInputErr
	}
}

func ConvertServerGenerator(serverGen *ServerGenerator, args *config.ServerArgument) (err error) {
	// handle initial ServerGenerator arguments
	if err = serverGen.handleInitArguments(args); err != nil {
		return err
	}

	// handle initial go files imports
	if err = serverGen.handleInitImports(); err != nil {
		return err
	}

	// registry information
	if err = serverGen.handleRegistry(); err != nil {
		return err
	}

	// if serverGen.isNew == false, update manifest
	if !serverGen.isNew {
		serverGen.updateManifest()
	}

	return nil
}

func (serverGen *ServerGenerator) setKitexExtension(key, extendOption string) (err error) {
	if _, ok := serverGen.GoFileImports[key]; !ok {
		return keyInputErr
	}

	for impt := range serverGen.GoFileImports[key] {
		value := strings.Split(impt, consts.Slash)
		// To avoid reporting errors in special circumstances, for example: registry-etcd.
		valueFinal := strings.Split(value[len(value)-1], consts.TheCrossed)
		if _, ok := serverGen.kitexExtension.Dependencies[impt]; ok {
			continue
		}
		serverGen.kitexExtension.Dependencies[impt] = valueFinal[len(valueFinal)-1]
		serverGen.kitexExtension.ExtendServer.ImportPaths = append(serverGen.kitexExtension.ExtendServer.ImportPaths, impt)
	}

	if serverGen.kitexExtension.ExtendServer.ExtendOption == "" {
		serverGen.kitexExtension.ExtendServer.ExtendOption = extendOption
	} else {
		serverGen.kitexExtension.ExtendServer.ExtendOption += consts.LineBreak + extendOption
	}

	return nil
}

func (serverGen *ServerGenerator) handleInitArguments(args *config.ServerArgument) (err error) {
	serverGen.GoModule = args.GoMod
	serverGen.ServiceName = args.Service
	serverGen.RegistryName = args.Registry
	serverGen.communicationType = args.Type

	// handle manifest
	isNew := utils.IsCwgoNew(args.OutDir)
	if isNew {
		serverGen.isNew = true
		serverGen.initManifest(consts.Server)
	} else {
		if err = serverGen.manifest.InitAndValidate(args.OutDir); err != nil {
			return err
		}

		if !(serverGen.manifest.CommandType == consts.Server && serverGen.manifest.CommunicationType == serverGen.communicationType) {
			serverGen.isNew = true
			serverGen.initManifest(consts.Server)
		}
	}

	// handle custom extension
	if args.CustomExtension != "" {
		if err = serverGen.ServerExtension.fromYAMLFile(args.CustomExtension); err != nil {
			return err
		}
	}

	switch serverGen.communicationType {
	case consts.RPC:
		serverGen.templateDir = args.TemplateDir

	case consts.HTTP:

	default:
		return typeInputErr
	}

	return
}

func (serverGen *ServerGenerator) handleInitImports() (err error) {
	switch serverGen.communicationType {
	case consts.RPC:
		// set initial conf.go imports
		confExtraImports := []string{""}
		if err = serverGen.GoFileImports.appendImports(consts.ConfGo, confExtraImports); err != nil {
			return err
		}

		// set kitex server basic options for server.go
		kitexServiceBasicImports = append(kitexServiceBasicImports, serverGen.GoModule+"/conf")
		if err = serverGen.GoFileImports.appendImports(consts.KitexExtensionServer, kitexServiceBasicImports); err != nil {
			return
		}
		if err = serverGen.setKitexExtension(consts.KitexExtensionServer, kitexServiceBasicOpts); err != nil {
			return
		}

	case consts.HTTP:
		// set initial main.go imports
		mainExtraImports := []string{
			serverGen.GoModule + "/biz/router",
			serverGen.GoModule + "/conf",
		}
		if err = serverGen.GoFileImports.appendImports(consts.Main, mainExtraImports); err != nil {
			return
		}

		// set initial conf.go imports
		confExtraImports := []string{""}
		if err = serverGen.GoFileImports.appendImports(consts.ConfGo, confExtraImports); err != nil {
			return err
		}

	default:
		return typeInputErr
	}

	return
}

func (serverGen *ServerGenerator) handleRegistry() (err error) {
	if !serverGen.isNew && serverGen.RegistryName == "" {
		serverGen.RegistryName = serverGen.manifest.Registry
	}

	switch serverGen.communicationType {
	case consts.RPC:
		switch serverGen.RegistryName {
		case consts.Etcd:
			serverGen.DefaultRegistryPort = consts.EtcdPort

			if err = serverGen.GoFileImports.appendImports(consts.KitexExtensionServer, kitexEtcdServerImports); err != nil {
				return
			}
			if err = serverGen.setKitexExtension(consts.KitexExtensionServer, kitexEtcdServer); err != nil {
				return
			}
		// todo 其他注册中心
		default:
			utils.RemoveKitexExtension()
			return
		}

		p := path.Join(serverGen.templateDir, consts.KitexExtensionYaml)
		if err = serverGen.kitexExtension.ToYAMLFile(p); err != nil {
			return
		}

	case consts.HTTP:
		switch serverGen.RegistryName {
		case consts.Etcd:
			serverGen.RegistryBody = hzEtcdServer
			serverGen.DefaultRegistryPort = consts.EtcdPort

			if err = serverGen.GoFileImports.appendImports(consts.Main, hzEtcdServerImports); err != nil {
				return
			}
		// todo 其他注册中心
		default:
		}

	default:
		return typeInputErr
	}

	return
}

func (s *ServerExtension) fromYAMLFile(filename string) error {
	if s == nil {
		return nil
	}
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, s)
}

func (serverGen *ServerGenerator) initManifest(commandType string) {
	serverGen.manifest.Version = meta.Version
	serverGen.manifest.CommandType = commandType
	serverGen.manifest.CommunicationType = serverGen.communicationType
	serverGen.manifest.Registry = serverGen.RegistryName
}

func (serverGen *ServerGenerator) updateManifest() {
	serverGen.manifest.Registry = serverGen.RegistryName
}
