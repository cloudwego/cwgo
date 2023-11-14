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
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/cloudwego/cwgo/config"
	"github.com/cloudwego/cwgo/meta"
	"github.com/cloudwego/cwgo/pkg/common/utils"
	"github.com/cloudwego/cwgo/pkg/consts"
	"github.com/cloudwego/kitex/tool/internal_pkg/generator"
)

type ClientGenerator struct {
	CommonGenerator

	ClientRender
}

type ClientRender struct {
	GoModule    string
	ServiceName string

	OutDir                string
	CurrentIDLServiceName string
	SnakeServiceNames     []string
	CamelServiceNames     []string

	GoFileImports ImportsMap

	ClientExtension
}

type ClientExtension struct {
	Resolver
}

type Resolver struct {
	ResolverName        string
	ResolverBody        string
	DefaultResolverPort string
}

func NewClientGenerator(types string) (*ClientGenerator, error) {
	switch types {
	case consts.RPC:
		imports, err := newImportsMap(consts.Client, consts.RPC)
		if err != nil {
			return nil, err
		}
		return &ClientGenerator{
			CommonGenerator: CommonGenerator{
				kitexExtension: &generator.TemplateExtension{
					Dependencies: map[string]string{},
					ExtendClient: &generator.APIExtension{},
				},
				manifest: new(meta.Manifest),
			},
			ClientRender: ClientRender{
				GoFileImports: imports,
			},
		}, nil

	case consts.HTTP:
		imports, err := newImportsMap(consts.Client, consts.HTTP)
		if err != nil {
			return nil, err
		}
		return &ClientGenerator{
			CommonGenerator: CommonGenerator{
				manifest: new(meta.Manifest),
			},
			ClientRender: ClientRender{
				SnakeServiceNames: make([]string, 0, 5),
				CamelServiceNames: make([]string, 0, 5),
				GoFileImports:     imports,
			},
		}, nil

	default:
		return nil, typeInputErr
	}
}

func ConvertClientGenerator(clientGen *ClientGenerator, args *config.ClientArgument) (err error) {
	// handle initial ClientGenerator arguments
	if err = clientGen.handleInitArguments(args); err != nil {
		return err
	}

	// handle initial go files imports
	if err = clientGen.handleInitImports(); err != nil {
		return err
	}

	// handle resolve information
	if err = clientGen.handleResolver(); err != nil {
		return err
	}

	// if clientGen.isNew == false, update manifest
	if !clientGen.isNew {
		clientGen.updateManifest()
	}

	return nil
}

func (clientGen *ClientGenerator) setKitexExtension(key, extendOption string) (err error) {
	if _, ok := clientGen.GoFileImports[key]; !ok {
		return keyInputErr
	}

	for impt := range clientGen.GoFileImports[key] {
		value := strings.Split(impt, consts.Slash)
		// To avoid reporting errors in special circumstances, for example: registry-etcd.
		valueFinal := strings.Split(value[len(value)-1], consts.TheCrossed)
		if _, ok := clientGen.kitexExtension.Dependencies[impt]; ok {
			continue
		}
		clientGen.kitexExtension.Dependencies[impt] = valueFinal[len(valueFinal)-1]
		clientGen.kitexExtension.ExtendClient.ImportPaths = append(clientGen.kitexExtension.ExtendClient.ImportPaths, impt)
	}

	if clientGen.kitexExtension.ExtendClient.ExtendOption == "" {
		clientGen.kitexExtension.ExtendClient.ExtendOption = extendOption
	} else {
		clientGen.kitexExtension.ExtendClient.ExtendOption += consts.LineBreak + extendOption
	}

	return nil
}

func (clientGen *ClientGenerator) handleInitArguments(args *config.ClientArgument) (err error) {
	clientGen.GoModule = args.GoMod
	clientGen.ServiceName = args.Service
	clientGen.communicationType = args.Type
	clientGen.ResolverName = args.Resolver

	// handle manifest
	isNew := utils.IsCwgoNew(args.OutDir)
	if isNew {
		clientGen.isNew = true
		clientGen.initManifest(consts.Client)
	} else {
		if err = clientGen.manifest.InitAndValidate(args.OutDir); err != nil {
			return err
		}

		if !(clientGen.manifest.CommandType == consts.Client && clientGen.manifest.CommunicationType == clientGen.communicationType) {
			clientGen.isNew = true
			clientGen.initManifest(consts.Client)
		}
	}

	switch clientGen.communicationType {
	case consts.RPC:
		clientGen.templateDir = args.TemplateDir
	case consts.HTTP:
		// get current dir
		currentDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("get current path failed: %s", err)
		}

		// get relative out dir
		dir, err := filepath.Rel(currentDir, args.OutDir)
		if err != nil {
			return fmt.Errorf("get relative out path to current path failed: %s", err)
		}
		if utils.IsWindows() {
			dir = strings.ReplaceAll(dir, consts.BackSlash, consts.Slash)
		}
		clientGen.OutDir = dir

		clientGen.SnakeServiceNames = args.SnakeServiceNames
		for _, s := range clientGen.SnakeServiceNames {
			clientGen.CamelServiceNames = append(clientGen.CamelServiceNames, utils.SnakeToCamel(s))
		}
	default:
		return typeInputErr
	}

	return
}

func (clientGen *ClientGenerator) handleInitImports() (err error) {
	switch clientGen.communicationType {
	case consts.RPC:
		// set initial main.go imports
		mainExtraImports := []string{
			clientGen.GoModule + "/conf",
			clientGen.GoModule + "/biz/rpc/" + clientGen.ServiceName,
		}
		if err = clientGen.GoFileImports.appendImports(consts.Main, mainExtraImports); err != nil {
			return err
		}

		// set initial conf.go imports
		confExtraImports := []string{""}
		if err = clientGen.GoFileImports.appendImports(consts.ConfGo, confExtraImports); err != nil {
			return err
		}

		// set kitex client basic options for client.go
		kitexClientBasicImports = append(kitexClientBasicImports, clientGen.GoModule+"/conf")
		if err = clientGen.GoFileImports.appendImports(consts.KitexExtensionClient, kitexClientBasicImports); err != nil {
			return err
		}
		if err = clientGen.setKitexExtension(consts.KitexExtensionClient, kitexClientBasicOpts); err != nil {
			return err
		}
	case consts.HTTP:
		// set initial main.go imports
		mainExtraImports := []string{
			clientGen.GoModule + "/conf",
		}
		for _, name := range clientGen.SnakeServiceNames {
			mainExtraImports = append(mainExtraImports, clientGen.GoModule+consts.Slash+clientGen.OutDir+consts.Slash+name)
		}
		if err = clientGen.GoFileImports.appendImports(consts.Main, mainExtraImports); err != nil {
			return err
		}

		// set initial conf.go imports
		confExtraImports := []string{""}
		if err = clientGen.GoFileImports.appendImports(consts.ConfGo, confExtraImports); err != nil {
			return err
		}
	default:
		return typeInputErr
	}

	return
}

func (clientGen *ClientGenerator) handleResolver() (err error) {
	if !clientGen.isNew && clientGen.ResolverName == "" {
		clientGen.ResolverName = clientGen.manifest.Resolver
	}

	switch clientGen.communicationType {
	case consts.RPC:
		switch clientGen.ResolverName {
		case consts.Etcd:
			clientGen.DefaultResolverPort = consts.EtcdPort

			if err = clientGen.GoFileImports.appendImports(consts.KitexExtensionClient, kitexEtcdClientImports); err != nil {
				return
			}
			if err = clientGen.setKitexExtension(consts.KitexExtensionClient, kitexEtcdClient); err != nil {
				return
			}
		// todo 其他注册中心
		default:
			utils.RemoveKitexExtension()
			return
		}

		p := path.Join(clientGen.templateDir, consts.KitexExtensionYaml)
		if err = clientGen.kitexExtension.ToYAMLFile(p); err != nil {
			return
		}

	case consts.HTTP:
		switch clientGen.ResolverName {
		case consts.Etcd:
			clientGen.ResolverBody = hzEtcdClient
			clientGen.DefaultResolverPort = consts.EtcdPort

			hzEtcdClientImports = append(hzEtcdClientImports, clientGen.GoModule+"/conf")
			if err = clientGen.GoFileImports.appendImports(consts.InitGo, hzEtcdClientImports); err != nil {
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

func (clientGen *ClientGenerator) initManifest(commandType string) {
	clientGen.manifest.Version = meta.Version
	clientGen.manifest.CommandType = commandType
	clientGen.manifest.CommunicationType = clientGen.communicationType
	clientGen.manifest.Resolver = clientGen.ResolverName
}

func (clientGen *ClientGenerator) updateManifest() {
	clientGen.manifest.Resolver = clientGen.ResolverName
}
