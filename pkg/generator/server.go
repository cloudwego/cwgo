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
	"github.com/cloudwego/cwgo/config"
	"github.com/cloudwego/cwgo/meta"
	"github.com/cloudwego/cwgo/pkg/common/utils"
	"github.com/cloudwego/cwgo/pkg/consts"
	"github.com/cloudwego/kitex/tool/internal_pkg/generator"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
)

var errTypeInput = errors.New("input wrong type")

type ServerGenerator struct {
	CommonGenerator // common generator params

	ServerRender // for template render

	serverOptionFileInfo
}

type ServerRender struct {
	GoModule    string
	ServiceName string
	Codec       string

	KitexIdlServiceName string

	GoFileImports ImportsMap // handle .go files imports

	ServerExtension
}

type ServerExtension struct {
	Registry `yaml:"registry,omitempty"`
}

type Registry struct {
	RegistryName    string   `yaml:"registry_name,omitempty"`
	RegistryImports []string `yaml:"registry_imports,omitempty"`
	RegistryBody    string   `yaml:"registry_body,omitempty"`
	RegistryAddress []string `yaml:"registry_address,omitempty"`
	RegistryDocker  string   `yaml:"registry_docker"`
}

type serverOptionFileInfo struct {
	mainGoContent        string
	confGoContent        string
	dockerComposeContent string
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
		return nil, errTypeInput
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
	if err = serverGen.handleRegistry(args.Registry); err != nil {
		return err
	}

	// if serverGen.isNew == false, update manifest
	if !serverGen.isNew {
		serverGen.updateManifest()
	} else {
		serverGen.initManifest(consts.Server)
	}

	return nil
}

func (serverGen *ServerGenerator) handleInitArguments(args *config.ServerArgument) (err error) {
	serverGen.GoModule = args.GoMod
	serverGen.ServiceName = args.Service
	serverGen.communicationType = args.Type
	serverGen.CustomExtensionFile = args.CustomExtension
	serverGen.OutDir = args.OutDir
	serverGen.Codec, err = utils.GetIdlType(args.IdlPath)
	if err != nil {
		return err
	}

	// handle manifest and determine if .cwgo exists
	isNew := utils.IsCwgoNew(args.OutDir)
	if isNew {
		serverGen.isNew = true
	} else {
		if err = serverGen.manifest.InitAndValidate(args.OutDir); err != nil {
			return err
		}

		if !(serverGen.manifest.CommandType == consts.Server && serverGen.manifest.CommunicationType == serverGen.communicationType) {
			serverGen.isNew = true
		}
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

	if serverGen.isNew && serverGen.communicationType == consts.RPC {
		handlerPath := filepath.Join(serverGen.OutDir, "handler.go")
		if isExist, _ := utils.PathExist(handlerPath); isExist {
			content, err := utils.ReadFileContent(handlerPath)
			if err != nil {
				return err
			}
			result, err := getStructNames(string(content))
			if err != nil {
				return err
			}
			serverGen.KitexIdlServiceName = result[0][:len(result[0])-4]
		}
	}

	if !serverGen.isNew {
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
	}

	return
}

func (serverGen *ServerGenerator) handleInitImports() (err error) {
	switch serverGen.communicationType {
	case consts.RPC:
		// set initial main.go imports
		dirs, err := utils.GetSubDirs(filepath.Join(serverGen.OutDir), true)
		if err != nil {
			return err
		}

		var dir string
		for _, d := range dirs {
			if filepath.Base(d) == serverGen.KitexIdlServiceName {
				dir, err = filepath.Rel(d, dirs[0])
				if err != nil {
					return err
				}
				break
			}
		}

		mainExtraImports := []string{
			serverGen.GoModule + "/conf",
			serverGen.GoModule + "/" + dir,
		}
		if serverGen.Codec == "thrift" {
			mainExtraImports = append(mainExtraImports, "github.com/cloudwego/kitex/pkg/transmeta")
		}
		if err = serverGen.GoFileImports.appendImports(consts.Main, mainExtraImports); err != nil {
			return err
		}

		// set initial conf.go imports
		confExtraImports := []string{""}
		if err = serverGen.GoFileImports.appendImports(consts.ConfGo, confExtraImports); err != nil {
			return err
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
		return errTypeInput
	}

	return
}

func (serverGen *ServerGenerator) handleRegistry(registryName string) (err error) {
	if serverGen.isNew {
		if err = serverGen.handleNewRegistry(registryName); err != nil {
			return err
		}
	} else {
		if err = serverGen.handleUpdateRegistry(registryName); err != nil {
			return err
		}
	}

	return nil
}

func (serverGen *ServerGenerator) handleNewRegistry(registryName string) (err error) {
	// custom server registry
	if serverGen.CustomExtensionFile != "" && serverGen.RegistryName != "" {
		if err = serverGen.GoFileImports.appendImports(consts.Main, serverGen.RegistryImports); err != nil {
			return
		}
		if err = serverGen.GoFileImports.appendImports(consts.ConfGo, envGoImports); err != nil {
			return err
		}
		return
	}

	serverGen.RegistryName = registryName

	switch serverGen.communicationType {
	case consts.RPC:
		switch serverGen.RegistryName {
		case consts.Nacos:
			if err = serverGen.handleNewRegistryTemplate(kitexNacosServer, nacosDocker, nacosServerAddr, kitexNacosServerImports); err != nil {
				return
			}
		case consts.Consul:
			if err = serverGen.handleNewRegistryTemplate(kitexConsulServer, consulDocker, consulServerAddr, kitexConsulServerImports); err != nil {
				return
			}
		case consts.Etcd:
			if err = serverGen.handleNewRegistryTemplate(kitexEtcdServer, etcdDocker, etcdServerAddr, kitexEtcdServerImports); err != nil {
				return
			}
		case consts.Eureka:
			if err = serverGen.handleNewRegistryTemplate(kitexEurekaServer, eurekaDocker, eurekaServerAddr, kitexEurekaServerImports); err != nil {
				return
			}
		case consts.Polaris:
			if err = serverGen.handleNewRegistryTemplate(kitexPolarisServer, polarisDocker, polarisServerAddr, kitexPolarisServerImports); err != nil {
				return
			}
		case consts.ServiceComb:
			if err = serverGen.handleNewRegistryTemplate(kitexServiceCombServer, serviceCombDocker, serviceCombServerAddr, kitexServiceCombServerImports); err != nil {
				return
			}
		case consts.Zk:
			if err = serverGen.handleNewRegistryTemplate(kitexZKServer, zkDocker, zkServerAddr, kitexZKServerImports); err != nil {
				return
			}
		default:
			return
		}

	case consts.HTTP:
		switch serverGen.RegistryName {
		case consts.Nacos:
			if err = serverGen.handleNewRegistryTemplate(hzNacosServer, nacosDocker, nacosServerAddr, hzNacosServerImports); err != nil {
				return
			}
		case consts.Consul:
			if err = serverGen.handleNewRegistryTemplate(hzConsulServer, consulDocker, consulServerAddr, hzConsulServerImports); err != nil {
				return
			}
		case consts.Etcd:
			if err = serverGen.handleNewRegistryTemplate(hzEtcdServer, etcdDocker, etcdServerAddr, hzEtcdServerImports); err != nil {
				return
			}
		case consts.Eureka:
			if err = serverGen.handleNewRegistryTemplate(hzEurekaServer, eurekaDocker, eurekaServerAddr, hzEurekaServerImports); err != nil {
				return
			}
		case consts.Polaris:
			if err = serverGen.handleNewRegistryTemplate(hzPolarisServer, polarisDocker, polarisServerAddr, hzPolarisServerImports); err != nil {
				return
			}
		case consts.ServiceComb:
			if err = serverGen.handleNewRegistryTemplate(hzServiceCombServer, serviceCombDocker, serviceCombServerAddr, hzServiceCombServerImports); err != nil {
				return
			}
		case consts.Zk:
			if err = serverGen.handleNewRegistryTemplate(hzZKServer, zkDocker, zkServerAddr, hzZKServerImports); err != nil {
				return
			}
		default:
		}

	default:
		return errTypeInput
	}

	return
}

func (serverGen *ServerGenerator) handleNewRegistryTemplate(body, docker string, addr, imports []string) (err error) {
	serverGen.RegistryBody = body
	serverGen.RegistryAddress = addr
	serverGen.RegistryDocker = docker

	if err = serverGen.GoFileImports.appendImports(consts.Main, imports); err != nil {
		return err
	}

	if err = serverGen.GoFileImports.appendImports(consts.ConfGo, envGoImports); err != nil {
		return err
	}

	return
}

func (serverGen *ServerGenerator) handleUpdateRegistry(registryName string) (err error) {
	if serverGen.CustomExtensionFile != "" && serverGen.RegistryName != "" {
		if err = serverGen.handleUpdateRegistryTemplate(serverGen.RegistryBody, serverGen.RegistryDocker, serverGen.RegistryAddress, serverGen.RegistryImports); err != nil {
			return
		}
		return
	}

	serverGen.RegistryName = registryName

	switch serverGen.communicationType {
	case consts.RPC:
		switch serverGen.RegistryName {
		case consts.Nacos:
			if err = serverGen.handleUpdateRegistryTemplate(kitexNacosServer, nacosDocker, nacosServerAddr, kitexNacosServerImports); err != nil {
				return
			}
		case consts.Consul:
			if err = serverGen.handleUpdateRegistryTemplate(kitexConsulServer, consulDocker, consulServerAddr, kitexConsulServerImports); err != nil {
				return
			}
		case consts.Etcd:
			if err = serverGen.handleUpdateRegistryTemplate(kitexEtcdServer, etcdDocker, etcdServerAddr, kitexEtcdServerImports); err != nil {
				return
			}
		case consts.Eureka:
			if err = serverGen.handleUpdateRegistryTemplate(kitexEurekaServer, eurekaDocker, eurekaServerAddr, kitexEurekaServerImports); err != nil {
				return
			}
		case consts.Polaris:
			if err = serverGen.handleUpdateRegistryTemplate(kitexPolarisServer, polarisDocker, polarisServerAddr, kitexPolarisServerImports); err != nil {
				return
			}
		case consts.ServiceComb:
			if err = serverGen.handleUpdateRegistryTemplate(kitexServiceCombServer, serviceCombDocker, serviceCombServerAddr, kitexServiceCombServerImports); err != nil {
				return
			}
		case consts.Zk:
			if err = serverGen.handleUpdateRegistryTemplate(kitexZKServer, zkDocker, zkServerAddr, kitexZKServerImports); err != nil {
				return
			}
		default:
			return
		}

	case consts.HTTP:
		switch serverGen.RegistryName {
		case consts.Nacos:
			if err = serverGen.handleUpdateRegistryTemplate(hzNacosServer, nacosDocker, nacosServerAddr, hzNacosServerImports); err != nil {
				return
			}
		case consts.Consul:
			if err = serverGen.handleUpdateRegistryTemplate(hzConsulServer, consulDocker, consulServerAddr, hzConsulServerImports); err != nil {
				return
			}
		case consts.Etcd:
			if err = serverGen.handleUpdateRegistryTemplate(hzEtcdServer, etcdDocker, etcdServerAddr, hzEtcdServerImports); err != nil {
				return
			}
		case consts.Eureka:
			if err = serverGen.handleUpdateRegistryTemplate(hzEurekaServer, eurekaDocker, eurekaServerAddr, hzEurekaServerImports); err != nil {
				return
			}
		case consts.Polaris:
			if err = serverGen.handleUpdateRegistryTemplate(hzPolarisServer, polarisDocker, polarisServerAddr, hzPolarisServerImports); err != nil {
				return
			}
		case consts.ServiceComb:
			if err = serverGen.handleUpdateRegistryTemplate(hzServiceCombServer, serviceCombDocker, serviceCombServerAddr, hzServiceCombServerImports); err != nil {
				return
			}
		case consts.Zk:
			if err = serverGen.handleUpdateRegistryTemplate(hzZKServer, zkDocker, zkServerAddr, hzZKServerImports); err != nil {
				return
			}
		default:
		}

	default:
		return errTypeInput
	}

	return
}

func (serverGen *ServerGenerator) handleUpdateRegistryTemplate(body, docker string, addr, imports []string) error {
	serverGen.RegistryAddress = addr

	var (
		mvcTemplates        []Template
		nilRegistryFuncBody string
		appendRegistryFunc  string
	)
	if serverGen.communicationType == consts.RPC {
		mvcTemplates = kitexServerMVCTemplates
		nilRegistryFuncBody = kitexNilRegistryFuncBody
		appendRegistryFunc = kitexAppendRegistryFunc
	} else {
		mvcTemplates = hzServerMVCTemplates
		nilRegistryFuncBody = hzNilRegistryFuncBody
		appendRegistryFunc = hzAppendRegistryFunc
	}

	equal, err := isFuncBodyEqual(serverGen.mainGoContent, consts.FuncInitRegistry, nilRegistryFuncBody)
	if err != nil {
		return err
	}

	if equal {
		mvcTemplates[consts.FileServerMainIndex].Type = consts.ReplaceFuncBody
		mvcTemplates[consts.FileServerMainIndex].ReplaceFuncName = append(mvcTemplates[consts.FileServerMainIndex].ReplaceFuncName, consts.FuncInitRegistry)
		mvcTemplates[consts.FileServerMainIndex].ReplaceFuncImport = append(mvcTemplates[consts.FileServerMainIndex].ReplaceFuncImport, imports)
		mvcTemplates[consts.FileServerMainIndex].ReplaceFuncBody = append(mvcTemplates[consts.FileServerMainIndex].ReplaceFuncBody, body)

		isExist, err := isFuncExist(serverGen.confGoContent, consts.FuncGetRegistryAddress)
		if err != nil {
			return err
		}

		if !isExist {
			mvcTemplates[consts.FileServerConfGoIndex].Type = consts.Append
			mvcTemplates[consts.FileServerConfGoIndex].AppendContent = appendRegistryFunc
			mvcTemplates[consts.FileServerConfGoIndex].AppendImport = envGoImports

			mvcTemplates[consts.FileServerDockerComposeIndex].Type = consts.Append
			mvcTemplates[consts.FileServerDockerComposeIndex].AppendContent = docker
		}
	}

	if serverGen.communicationType == consts.RPC {
		kitexServerMVCTemplates = mvcTemplates
	} else {
		hzServerMVCTemplates = mvcTemplates
	}

	return nil
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

func (s *ServerExtension) checkCustomExtensionFile() (err error) {
	// check resolver
	if s.RegistryName != "" {
		if s.RegistryImports == nil {
			return errors.New("please input RegistryImports")
		}
		if s.RegistryBody == "" {
			return errors.New("please input RegistryImports")
		}
	}

	return nil
}

func (serverGen *ServerGenerator) initManifest(commandType string) {
	serverGen.manifest.Version = meta.Version
	serverGen.manifest.CommandType = commandType
	serverGen.manifest.CommunicationType = serverGen.communicationType
	serverGen.manifest.Registry = serverGen.RegistryName
	serverGen.manifest.CustomExtensionFile = serverGen.CustomExtensionFile
}

func (serverGen *ServerGenerator) updateManifest() {
	serverGen.manifest.Registry = serverGen.RegistryName
	serverGen.manifest.CustomExtensionFile = serverGen.CustomExtensionFile
}
