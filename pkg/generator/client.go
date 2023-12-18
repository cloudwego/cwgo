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
	"fmt"
	"github.com/cloudwego/cwgo/config"
	"github.com/cloudwego/cwgo/meta"
	"github.com/cloudwego/cwgo/pkg/common/utils"
	"github.com/cloudwego/cwgo/pkg/consts"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

type ClientGenerator struct {
	CommonGenerator // common generator params

	ClientRender // for template render

	clientOptionFileInfo
}

type ClientRender struct {
	GoModule    string
	ServiceName string
	Codec       string

	InitOptsPackage string

	GoFileImports ImportsMap // handle .go files imports

	ClientExtension
}

type ClientExtension struct {
	Resolver `yaml:"resolver,omitempty"`
}

type Resolver struct {
	ResolverName    string   `yaml:"resolver_name,omitempty"`
	ResolverImports []string `yaml:"resolver_imports,omitempty"`
	ResolverBody    string   `yaml:"resolver_body,omitempty"`
	ResolverAddress []string `yaml:"resolver_address,omitempty"`
}

type clientOptionFileInfo struct {
	initGoContents []string
	initGoPaths    []string
	envGoContent   string
	envGoPath      string
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
				manifest: new(meta.Manifest),
			},
			ClientRender: ClientRender{
				GoFileImports: imports,
			},
			clientOptionFileInfo: clientOptionFileInfo{
				initGoPaths:    make([]string, 0, 5),
				initGoContents: make([]string, 0, 5),
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
				GoFileImports: imports,
			},
			clientOptionFileInfo: clientOptionFileInfo{
				initGoPaths:    make([]string, 0, 5),
				initGoContents: make([]string, 0, 5),
			},
		}, nil

	default:
		return nil, errTypeInput
	}
}

func ConvertClientGenerator(clientGen *ClientGenerator, args *config.ClientArgument) (err error) {
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

	// if clientGen.isNew == false, update manifest
	if !clientGen.isNew {
		clientGen.updateManifest()
	} else {
		clientGen.initManifest(consts.Client)
	}

	return nil
}

func (clientGen *ClientGenerator) handleInitArguments(args *config.ClientArgument) (err error) {
	clientGen.GoModule = args.GoMod
	clientGen.ServiceName = args.Service
	clientGen.communicationType = args.Type
	clientGen.CustomExtensionFile = args.CustomExtension
	clientGen.OutDir = args.OutDir
	clientGen.Codec, err = utils.GetIdlType(args.IdlPath)
	if err != nil {
		return err
	}

	// handle manifest and determine if .cwgo exists
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get current path failed: %s", err)
	}

	isNew := utils.IsCwgoNew(dir)
	if isNew {
		clientGen.isNew = true
	} else {
		if err = clientGen.manifest.InitAndValidate(dir); err != nil {
			return err
		}

		if !(clientGen.manifest.CommandType == consts.Client && clientGen.manifest.CommunicationType == clientGen.communicationType) {
			clientGen.isNew = true
		}
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

	if !clientGen.isNew {
		bizDir := ""
		if clientGen.communicationType == consts.HTTP {
			bizDir = clientGen.OutDir
		} else {
			bizDir = filepath.Join(clientGen.OutDir, consts.DefaultKitexClientDir)
		}

		subDirs, err := utils.GetSubDirs(bizDir, false)
		if err != nil {
			return err
		}
		for _, subDir := range subDirs {
			filePath := filepath.Join(subDir, consts.InitGo)
			if isExist, _ := utils.PathExist(filePath); isExist {
				content, err := utils.ReadFileContent(filePath)
				if err != nil {
					return err
				}
				clientGen.initGoContents = append(clientGen.initGoContents, string(content))
				clientGen.initGoPaths = append(clientGen.initGoPaths, filePath)
			}
		}

		clientGen.envGoPath = filepath.Join(bizDir, consts.EnvGo)
		if isExist, _ := utils.PathExist(clientGen.envGoPath); isExist {
			content, err := utils.ReadFileContent(clientGen.envGoPath)
			if err != nil {
				return err
			}
			clientGen.envGoContent = string(content)
		}
	}

	return
}

func (clientGen *ClientGenerator) handleInitImports() (err error) {
	switch clientGen.communicationType {
	case consts.RPC:
		// set initial init.go imports
		if clientGen.Codec == "thrift" {
			initExtraImports := []string{
				"github.com/cloudwego/kitex/pkg/transmeta",
				"github.com/cloudwego/kitex/transport",
			}
			if err = clientGen.GoFileImports.appendImports(consts.InitGo, initExtraImports); err != nil {
				return err
			}
		}
	case consts.HTTP:
	default:
		return errTypeInput
	}

	return
}

func (clientGen *ClientGenerator) handleResolver(resolverName string) (err error) {
	if clientGen.isNew {
		if err = clientGen.handleNewResolver(resolverName); err != nil {
			return err
		}
	} else {
		if err = clientGen.handleUpdateResolver(resolverName); err != nil {
			return err
		}
	}

	return nil
}

func (clientGen *ClientGenerator) handleNewResolver(resolverName string) (err error) {
	// custom server resolver
	if clientGen.CustomExtensionFile != "" && clientGen.ResolverName != "" {
		if err = clientGen.GoFileImports.appendImports(consts.InitGo, clientGen.ResolverImports); err != nil {
			return
		}
		if err = clientGen.GoFileImports.appendImports(consts.EnvGo, envGoImports); err != nil {
			return err
		}
		return
	}

	clientGen.ResolverName = resolverName

	switch clientGen.communicationType {
	case consts.RPC:
		switch clientGen.ResolverName {
		case consts.Nacos:
			if err = clientGen.handleNewResolverTemplate(kitexNacosClient, nacosServerAddr, kitexNacosClientImports); err != nil {
				return
			}
		case consts.Consul:
			if err = clientGen.handleNewResolverTemplate(kitexConsulClient, consulServerAddr, kitexConsulClientImports); err != nil {
				return
			}
		case consts.Etcd:
			if err = clientGen.handleNewResolverTemplate(kitexEtcdClient, etcdServerAddr, kitexEtcdClientImports); err != nil {
				return
			}
		case consts.Eureka:
			if err = clientGen.handleNewResolverTemplate(kitexEurekaClient, eurekaServerAddr, kitexEurekaClientImports); err != nil {
				return
			}
		case consts.Polaris:
			if err = clientGen.handleNewResolverTemplate(kitexPolarisClient, polarisServerAddr, kitexPolarisClientImports); err != nil {
				return
			}
		case consts.ServiceComb:
			if err = clientGen.handleNewResolverTemplate(kitexServiceCombClient, serviceCombServerAddr, kitexServiceCombClientImports); err != nil {
				return
			}
		case consts.Zk:
			if err = clientGen.handleNewResolverTemplate(kitexZKClient, zkServerAddr, kitexZKClientImports); err != nil {
				return
			}
		default:
		}

	case consts.HTTP:
		switch clientGen.ResolverName {
		case consts.Nacos:
			if err = clientGen.handleNewResolverTemplate(hzNacosClient, nacosServerAddr, hzNacosClientImports); err != nil {
				return
			}
		case consts.Consul:
			if err = clientGen.handleNewResolverTemplate(hzConsulClient, consulServerAddr, hzConsulClientImports); err != nil {
				return
			}
		case consts.Etcd:
			if err = clientGen.handleNewResolverTemplate(hzEtcdClient, etcdServerAddr, hzEtcdClientImports); err != nil {
				return
			}
		case consts.Eureka:
			if err = clientGen.handleNewResolverTemplate(hzEurekaClient, eurekaServerAddr, hzEurekaClientImports); err != nil {
				return
			}
		case consts.Polaris:
			if err = clientGen.handleNewResolverTemplate(hzPolarisClient, polarisServerAddr, hzPolarisClientImports); err != nil {
				return
			}
		case consts.ServiceComb:
			if err = clientGen.handleNewResolverTemplate(hzServiceCombClient, serviceCombServerAddr, hzServiceCombClientImports); err != nil {
				return
			}
		case consts.Zk:
			if err = clientGen.handleNewResolverTemplate(hzZKClient, zkServerAddr, hzZKClientImports); err != nil {
				return
			}
		default:
		}

	default:
		return errTypeInput
	}

	return
}

func (clientGen *ClientGenerator) handleNewResolverTemplate(body string, addr, imports []string) (err error) {
	clientGen.ResolverBody = body
	clientGen.ResolverAddress = addr

	if clientGen.communicationType == consts.HTTP {
		imports = append(imports, clientGen.GoModule+consts.Slash+consts.DefaultHZClientDir)
	} else {
		imports = append(imports, clientGen.GoModule+consts.Slash+consts.DefaultKitexClientDir)
	}
	if err = clientGen.GoFileImports.appendImports(consts.InitGo, imports); err != nil {
		return err
	}

	if err = clientGen.GoFileImports.appendImports(consts.EnvGo, envGoImports); err != nil {
		return err
	}

	return
}

func (clientGen *ClientGenerator) handleUpdateResolver(resolverName string) (err error) {
	if clientGen.CustomExtensionFile != "" && clientGen.ResolverName != "" {
		if err = clientGen.handleUpdateResolverTemplate(clientGen.ResolverBody, clientGen.ResolverAddress, clientGen.ResolverImports); err != nil {
			return err
		}
		return
	}

	clientGen.ResolverName = resolverName

	switch clientGen.communicationType {
	case consts.RPC:
		switch clientGen.ResolverName {
		case consts.Nacos:
			if err = clientGen.handleUpdateResolverTemplate(kitexNacosClient, nacosServerAddr, kitexNacosClientImports); err != nil {
				return err
			}
		case consts.Consul:
			if err = clientGen.handleUpdateResolverTemplate(kitexConsulClient, consulServerAddr, kitexConsulClientImports); err != nil {
				return err
			}
		case consts.Etcd:
			if err = clientGen.handleUpdateResolverTemplate(kitexEtcdClient, etcdServerAddr, kitexEtcdClientImports); err != nil {
				return err
			}
		case consts.Eureka:
			if err = clientGen.handleUpdateResolverTemplate(kitexEurekaClient, eurekaServerAddr, kitexEurekaClientImports); err != nil {
				return err
			}
		case consts.Polaris:
			if err = clientGen.handleUpdateResolverTemplate(kitexPolarisClient, polarisServerAddr, kitexPolarisClientImports); err != nil {
				return err
			}
		case consts.ServiceComb:
			if err = clientGen.handleUpdateResolverTemplate(kitexServiceCombClient, serviceCombServerAddr, kitexServiceCombClientImports); err != nil {
				return err
			}
		case consts.Zk:
			if err = clientGen.handleUpdateResolverTemplate(kitexZKClient, zkServerAddr, kitexZKClientImports); err != nil {
				return err
			}
		default:
		}

	case consts.HTTP:
		switch clientGen.ResolverName {
		case consts.Nacos:
			if err = clientGen.handleUpdateResolverTemplate(hzNacosClient, nacosServerAddr, hzNacosClientImports); err != nil {
				return err
			}
		case consts.Consul:
			if err = clientGen.handleUpdateResolverTemplate(hzConsulClient, consulServerAddr, hzConsulClientImports); err != nil {
				return err
			}
		case consts.Etcd:
			if err = clientGen.handleUpdateResolverTemplate(hzEtcdClient, etcdServerAddr, hzEtcdClientImports); err != nil {
				return err
			}
		case consts.Eureka:
			if err = clientGen.handleUpdateResolverTemplate(hzEurekaClient, eurekaServerAddr, hzEurekaClientImports); err != nil {
				return err
			}
		case consts.Polaris:
			if err = clientGen.handleUpdateResolverTemplate(hzPolarisClient, polarisServerAddr, hzPolarisClientImports); err != nil {
				return err
			}
		case consts.ServiceComb:
			if err = clientGen.handleUpdateResolverTemplate(hzServiceCombClient, serviceCombServerAddr, hzServiceCombClientImports); err != nil {
				return err
			}
		case consts.Zk:
			if err = clientGen.handleUpdateResolverTemplate(hzZKClient, zkServerAddr, hzZKClientImports); err != nil {
				return err
			}
		default:
		}

	default:
		return errTypeInput
	}

	return nil
}

func (clientGen *ClientGenerator) handleUpdateResolverTemplate(body string, addr, imports []string) error {
	clientGen.ResolverAddress = addr

	var (
		mvcTemplates        []Template
		nilResolverFuncBody string
		appendResolverFunc  string
	)
	if clientGen.communicationType == consts.RPC {
		imports = append(imports, clientGen.GoModule+consts.Slash+consts.DefaultKitexClientDir)
		mvcTemplates = kitexClientMVCTemplates
		nilResolverFuncBody = kitexNilResolverFuncBody
		appendResolverFunc = kitexAppendResolverFunc
	} else {
		imports = append(imports, clientGen.GoModule+consts.Slash+consts.DefaultHZClientDir)
		mvcTemplates = hzClientMVCTemplates
		nilResolverFuncBody = hzNilResolverFuncBody
		appendResolverFunc = hzAppendResolverFunc
	}

	flag := 0
	for index, content := range clientGen.initGoContents {
		equal, err := isFuncBodyEqual(content, consts.FuncInitResolver, nilResolverFuncBody)
		if err != nil {
			return err
		}

		if equal {
			if index == 0 {
				mvcTemplates[consts.FileClientInitIndex].Path = clientGen.initGoPaths[index]
				mvcTemplates[consts.FileClientInitIndex].Type = consts.ReplaceFuncBody
				mvcTemplates[consts.FileClientInitIndex].ReplaceFuncName = append(mvcTemplates[consts.FileClientInitIndex].ReplaceFuncName, consts.FuncInitResolver)
				mvcTemplates[consts.FileClientInitIndex].ReplaceFuncImport = append(mvcTemplates[consts.FileClientInitIndex].ReplaceFuncImport, imports)
				mvcTemplates[consts.FileClientInitIndex].ReplaceFuncBody = append(mvcTemplates[consts.FileClientInitIndex].ReplaceFuncBody, body)
			} else {
				template := Template{
					Path: clientGen.initGoPaths[index],
					UpdateBehavior: UpdateBehavior{
						Type: consts.ReplaceFuncBody,
						ReplaceFunc: ReplaceFunc{
							ReplaceFuncName:   []string{consts.FuncInitResolver},
							ReplaceFuncImport: [][]string{imports},
							ReplaceFuncBody:   []string{body},
						},
						AppendRender: map[string]interface{}{},
					},
				}
				mvcTemplates = append(mvcTemplates, template)
			}

			if flag == 0 {
				isExist, err := isFuncExist(clientGen.envGoContent, consts.FuncGetResolverAddress)
				if err != nil {
					return err
				}

				if !isExist {
					mvcTemplates[consts.FileClientEnvIndex].Path = clientGen.envGoPath
					mvcTemplates[consts.FileClientEnvIndex].Type = consts.Append
					mvcTemplates[consts.FileClientEnvIndex].AppendContent = appendResolverFunc
					mvcTemplates[consts.FileClientEnvIndex].AppendImport = envGoImports
				}

				flag++
			}
		}
	}

	if clientGen.communicationType == consts.RPC {
		kitexClientMVCTemplates = mvcTemplates
	} else {
		hzClientMVCTemplates = mvcTemplates
	}

	return nil
}

func (c *ClientExtension) fromYAMLFile(filename string) error {
	if c == nil {
		return nil
	}
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, c)
}

func (c *ClientExtension) checkCustomExtensionFile() (err error) {
	// check resolver
	if c.ResolverName != "" {
		if c.ResolverImports == nil {
			return errors.New("please input ResolverImports")
		}
		if c.ResolverBody == "" {
			return errors.New("please input ResolverImports")
		}
	}

	return nil
}

func (clientGen *ClientGenerator) initManifest(commandType string) {
	clientGen.manifest.Version = meta.Version
	clientGen.manifest.CommandType = commandType
	clientGen.manifest.CommunicationType = clientGen.communicationType
	clientGen.manifest.Resolver = clientGen.ResolverName
	clientGen.manifest.CustomExtensionFile = clientGen.CustomExtensionFile
}

func (clientGen *ClientGenerator) updateManifest() {
	clientGen.manifest.Resolver = clientGen.ResolverName
	clientGen.manifest.CustomExtensionFile = clientGen.CustomExtensionFile
}
