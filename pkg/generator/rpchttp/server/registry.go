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
	"fmt"

	"github.com/cloudwego/cwgo/pkg/common/utils"
	"github.com/cloudwego/cwgo/pkg/consts"
	"github.com/cloudwego/cwgo/pkg/generator/common/template"
	geneUtils "github.com/cloudwego/cwgo/pkg/generator/common/utils"
	rhCommon "github.com/cloudwego/cwgo/pkg/generator/rpchttp/common"
)

func (serverGen *Generator) handleRegistry(registryName string) (err error) {
	if serverGen.IsNew {
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

func (serverGen *Generator) handleNewRegistry(registryName string) (err error) {
	// custom server registry
	if serverGen.CustomExtensionFile != "" && serverGen.RegistryName != "" {
		if err = serverGen.handleNewRegistryTemplate(serverGen.RegistryBody, serverGen.RegistryDocker, serverGen.RegistryAddress, serverGen.RegistryImports); err != nil {
			return
		}
		return
	}

	serverGen.RegistryName = registryName

	switch serverGen.CommunicationType {
	case consts.RPC:
		switch serverGen.RegistryName {
		case consts.Nacos:
			if err = serverGen.handleNewRegistryTemplate(kitexNacosServer, rhCommon.NacosDocker, rhCommon.NacosServerAddr, kitexNacosServerImports); err != nil {
				return
			}
		case consts.Consul:
			if err = serverGen.handleNewRegistryTemplate(kitexConsulServer, rhCommon.ConsulDocker, rhCommon.ConsulServerAddr, kitexConsulServerImports); err != nil {
				return
			}
		case consts.Etcd:
			if err = serverGen.handleNewRegistryTemplate(kitexEtcdServer, rhCommon.EtcdDocker, rhCommon.EtcdServerAddr, kitexEtcdServerImports); err != nil {
				return
			}
		case consts.Eureka:
			if err = serverGen.handleNewRegistryTemplate(kitexEurekaServer, rhCommon.EurekaDocker, rhCommon.EurekaServerAddr, kitexEurekaServerImports); err != nil {
				return
			}
		case consts.Polaris:
			if err = serverGen.handleNewRegistryTemplate(kitexPolarisServer, rhCommon.PolarisDocker, rhCommon.PolarisServerAddr, kitexPolarisServerImports); err != nil {
				return
			}
		case consts.ServiceComb:
			if err = serverGen.handleNewRegistryTemplate(kitexServiceCombServer, rhCommon.ServiceCombDocker, rhCommon.ServiceCombServerAddr, kitexServiceCombServerImports); err != nil {
				return
			}
		case consts.Zk:
			if err = serverGen.handleNewRegistryTemplate(kitexZKServer, rhCommon.ZkDocker, rhCommon.ZkServerAddr, kitexZKServerImports); err != nil {
				return
			}
		default:
			if err = serverGen.compatibleOlderVersion(); err != nil {
				return
			}
		}

	case consts.HTTP:
		switch serverGen.RegistryName {
		case consts.Nacos:
			if err = serverGen.handleNewRegistryTemplate(hzNacosServer, rhCommon.NacosDocker, rhCommon.NacosServerAddr, hzNacosServerImports); err != nil {
				return
			}
		case consts.Consul:
			if err = serverGen.handleNewRegistryTemplate(hzConsulServer, rhCommon.ConsulDocker, rhCommon.ConsulServerAddr, hzConsulServerImports); err != nil {
				return
			}
		case consts.Etcd:
			if err = serverGen.handleNewRegistryTemplate(hzEtcdServer, rhCommon.EtcdDocker, rhCommon.EtcdServerAddr, hzEtcdServerImports); err != nil {
				return
			}
		case consts.Eureka:
			if err = serverGen.handleNewRegistryTemplate(hzEurekaServer, rhCommon.EurekaDocker, rhCommon.EurekaServerAddr, hzEurekaServerImports); err != nil {
				return
			}
		case consts.Polaris:
			if err = serverGen.handleNewRegistryTemplate(hzPolarisServer, rhCommon.PolarisDocker, rhCommon.PolarisServerAddr, hzPolarisServerImports); err != nil {
				return
			}
		case consts.ServiceComb:
			if err = serverGen.handleNewRegistryTemplate(hzServiceCombServer, rhCommon.ServiceCombDocker, rhCommon.ServiceCombServerAddr, hzServiceCombServerImports); err != nil {
				return
			}
		case consts.Zk:
			if err = serverGen.handleNewRegistryTemplate(hzZKServer, rhCommon.ZkDocker, rhCommon.ZkServerAddr, hzZKServerImports); err != nil {
				return
			}
		default:
			if err = serverGen.compatibleOlderVersion(); err != nil {
				return
			}
		}

	default:
		return rhCommon.ErrTypeInput
	}

	return
}

func (serverGen *Generator) handleNewRegistryTemplate(body, docker string, addr []string, imports map[string]string) (err error) {
	var (
		mvcTemplates           []template.Template
		appendInitRegistryFunc string
		appendRegistryAddrFunc string
		disableAddConf         bool
	)

	if serverGen.CommunicationType == consts.RPC {
		mvcTemplates = kitexServerMVCTemplates
		appendInitRegistryFunc = kitexAppendInitRegistryFunc
		if serverGen.confGoContent != "" {
			isExist, err := geneUtils.IsStructExist(serverGen.confGoContent, "Registry")
			if err != nil {
				return err
			}
			if isExist {
				disableAddConf = true
				appendRegistryAddrFunc = kitexNewAppendRegistryAddrFunc
			} else {
				appendRegistryAddrFunc = kitexAppendRegistryAddrFunc
			}
		} else {
			appendRegistryAddrFunc = kitexAppendRegistryAddrFunc
		}
	} else {
		mvcTemplates = hzServerMVCTemplates
		appendInitRegistryFunc = hzAppendInitRegistryFunc
		appendRegistryAddrFunc = hzAppendRegistryAddrFunc
	}

	if serverGen.CommunicationType == consts.HTTP {
		serverGen.RegistryBody = fmt.Sprintf(body, serverGen.ServiceName)
	} else {
		serverGen.RegistryBody = body
	}
	serverGen.RegistryAddress = addr
	serverGen.RegistryDocker = docker

	if serverGen.mainGoContent != "" {
		isExist, err := geneUtils.IsFuncExist(serverGen.mainGoContent, consts.FuncInitRegistry)
		if err != nil {
			return err
		}
		if !isExist {
			mvcTemplates[consts.FileServerMainIndex].Type = consts.Append
			mvcTemplates[consts.FileServerMainIndex].AppendContent += appendInitRegistryFunc + consts.LineBreak
			geneUtils.Add2MapStrStr(imports, mvcTemplates[consts.FileServerMainIndex].AppendImport)
		}
	} else {
		if err = serverGen.GoFileImports.AppendImports(consts.Main, imports); err != nil {
			return err
		}
	}

	if serverGen.confGoContent != "" {
		isExist, err := geneUtils.IsFuncExist(serverGen.confGoContent, consts.FuncGetRegistryAddress)
		if err != nil {
			return err
		}

		if !isExist {
			content, err := geneUtils.InsertField2Struct(serverGen.confGoContent, "Config", rhCommon.RegistryStructField, "Registry")
			if err != nil {
				return err
			}
			if err = utils.CreateFile(consts.ConfGo, content); err != nil {
				return err
			}

			mvcTemplates[consts.FileServerConfGoIndex].Type = consts.Append
			mvcTemplates[consts.FileServerConfGoIndex].AppendContent += appendRegistryAddrFunc + consts.LineBreak
			geneUtils.Add2MapStrStr(rhCommon.EnvGoImports, mvcTemplates[consts.FileServerConfGoIndex].AppendImport)

			if serverGen.dockerComposeContent != "" {
				mvcTemplates[consts.FileServerDockerComposeIndex].Type = consts.Append
				mvcTemplates[consts.FileServerDockerComposeIndex].AppendContent += docker + consts.LineBreak
			}

			if !disableAddConf {
				mvcTemplates[consts.FileServerDevConf].Type = consts.Append
				mvcTemplates[consts.FileServerDevConf].AppendContent += rhCommon.RegistryConfYaml + consts.LineBreak
				mvcTemplates[consts.FileServerOnlineConf].Type = consts.Append
				mvcTemplates[consts.FileServerOnlineConf].AppendContent += rhCommon.RegistryConfYaml + consts.LineBreak
				mvcTemplates[consts.FileServerTestConf].Type = consts.Append
				mvcTemplates[consts.FileServerTestConf].AppendContent += rhCommon.RegistryConfYaml + consts.LineBreak
			}
		}
	} else {
		if err = serverGen.GoFileImports.AppendImports(consts.ConfGo, rhCommon.EnvGoImports); err != nil {
			return err
		}
	}

	if serverGen.CommunicationType == consts.RPC {
		kitexServerMVCTemplates = mvcTemplates
	} else {
		hzServerMVCTemplates = mvcTemplates
	}

	return
}

func (serverGen *Generator) handleUpdateRegistry(registryName string) (err error) {
	if serverGen.CustomExtensionFile != "" && serverGen.RegistryName != "" {
		if err = serverGen.handleUpdateRegistryTemplate(serverGen.RegistryBody, serverGen.RegistryDocker, serverGen.RegistryAddress, serverGen.RegistryImports); err != nil {
			return
		}
		return
	}

	serverGen.RegistryName = registryName

	switch serverGen.CommunicationType {
	case consts.RPC:
		switch serverGen.RegistryName {
		case consts.Nacos:
			if err = serverGen.handleUpdateRegistryTemplate(kitexNacosServer, rhCommon.NacosDocker, rhCommon.NacosServerAddr, kitexNacosServerImports); err != nil {
				return
			}
		case consts.Consul:
			if err = serverGen.handleUpdateRegistryTemplate(kitexConsulServer, rhCommon.ConsulDocker, rhCommon.ConsulServerAddr, kitexConsulServerImports); err != nil {
				return
			}
		case consts.Etcd:
			if err = serverGen.handleUpdateRegistryTemplate(kitexEtcdServer, rhCommon.EtcdDocker, rhCommon.EtcdServerAddr, kitexEtcdServerImports); err != nil {
				return
			}
		case consts.Eureka:
			if err = serverGen.handleUpdateRegistryTemplate(kitexEurekaServer, rhCommon.EurekaDocker, rhCommon.EurekaServerAddr, kitexEurekaServerImports); err != nil {
				return
			}
		case consts.Polaris:
			if err = serverGen.handleUpdateRegistryTemplate(kitexPolarisServer, rhCommon.PolarisDocker, rhCommon.PolarisServerAddr, kitexPolarisServerImports); err != nil {
				return
			}
		case consts.ServiceComb:
			if err = serverGen.handleUpdateRegistryTemplate(kitexServiceCombServer, rhCommon.ServiceCombDocker, rhCommon.ServiceCombServerAddr, kitexServiceCombServerImports); err != nil {
				return
			}
		case consts.Zk:
			if err = serverGen.handleUpdateRegistryTemplate(kitexZKServer, rhCommon.ZkDocker, rhCommon.ZkServerAddr, kitexZKServerImports); err != nil {
				return
			}
		default:
			return
		}

	case consts.HTTP:
		switch serverGen.RegistryName {
		case consts.Nacos:
			if err = serverGen.handleUpdateRegistryTemplate(hzNacosServer, rhCommon.NacosDocker, rhCommon.NacosServerAddr, hzNacosServerImports); err != nil {
				return
			}
		case consts.Consul:
			if err = serverGen.handleUpdateRegistryTemplate(hzConsulServer, rhCommon.ConsulDocker, rhCommon.ConsulServerAddr, hzConsulServerImports); err != nil {
				return
			}
		case consts.Etcd:
			if err = serverGen.handleUpdateRegistryTemplate(hzEtcdServer, rhCommon.EtcdDocker, rhCommon.EtcdServerAddr, hzEtcdServerImports); err != nil {
				return
			}
		case consts.Eureka:
			if err = serverGen.handleUpdateRegistryTemplate(hzEurekaServer, rhCommon.EurekaDocker, rhCommon.EurekaServerAddr, hzEurekaServerImports); err != nil {
				return
			}
		case consts.Polaris:
			if err = serverGen.handleUpdateRegistryTemplate(hzPolarisServer, rhCommon.PolarisDocker, rhCommon.PolarisServerAddr, hzPolarisServerImports); err != nil {
				return
			}
		case consts.ServiceComb:
			if err = serverGen.handleUpdateRegistryTemplate(hzServiceCombServer, rhCommon.ServiceCombDocker, rhCommon.ServiceCombServerAddr, hzServiceCombServerImports); err != nil {
				return
			}
		case consts.Zk:
			if err = serverGen.handleUpdateRegistryTemplate(hzZKServer, rhCommon.ZkDocker, rhCommon.ZkServerAddr, hzZKServerImports); err != nil {
				return
			}
		default:
		}

	default:
		return rhCommon.ErrTypeInput
	}

	return
}

func (serverGen *Generator) handleUpdateRegistryTemplate(body, docker string, addr []string, imports map[string]string) error {
	serverGen.RegistryAddress = addr
	if serverGen.CommunicationType == consts.HTTP {
		body = fmt.Sprintf(body, serverGen.ServiceName)
	}

	var (
		mvcTemplates        []template.Template
		nilRegistryFuncBody string
		appendRegistryFunc  string
		disableAddConf      bool
	)

	if serverGen.CommunicationType == consts.RPC {
		mvcTemplates = kitexServerMVCTemplates
		nilRegistryFuncBody = kitexNilRegistryFuncBody
		isExist, err := geneUtils.IsStructExist(serverGen.confGoContent, "Registry")
		if err != nil {
			return err
		}
		if isExist {
			disableAddConf = true
			appendRegistryFunc = kitexNewAppendRegistryAddrFunc
		} else {
			appendRegistryFunc = kitexAppendRegistryAddrFunc
		}
	} else {
		mvcTemplates = hzServerMVCTemplates
		nilRegistryFuncBody = hzNilRegistryFuncBody
		appendRegistryFunc = hzAppendRegistryAddrFunc
	}

	equal, err := geneUtils.IsFuncBodyEqual(serverGen.mainGoContent, consts.FuncInitRegistry, nilRegistryFuncBody)
	if err != nil {
		return err
	}

	if equal {
		mvcTemplates[consts.FileServerMainIndex].Type = consts.ReplaceFuncBody
		mvcTemplates[consts.FileServerMainIndex].ReplaceFuncName = append(mvcTemplates[consts.FileServerMainIndex].ReplaceFuncName, consts.FuncInitRegistry)
		mvcTemplates[consts.FileServerMainIndex].ReplaceFuncAppendImport = append(mvcTemplates[consts.FileServerMainIndex].ReplaceFuncAppendImport, imports)
		mvcTemplates[consts.FileServerMainIndex].ReplaceFuncBody = append(mvcTemplates[consts.FileServerMainIndex].ReplaceFuncBody, body)

		isExist, err := geneUtils.IsFuncExist(serverGen.confGoContent, consts.FuncGetRegistryAddress)
		if err != nil {
			return err
		}

		if !isExist {
			content, err := geneUtils.InsertField2Struct(serverGen.confGoContent, "Config", rhCommon.RegistryStructField, "Registry")
			if err != nil {
				return err
			}
			if err = utils.CreateFile(consts.ConfGo, content); err != nil {
				return err
			}

			mvcTemplates[consts.FileServerConfGoIndex].Type = consts.Append
			mvcTemplates[consts.FileServerConfGoIndex].AppendContent += appendRegistryFunc + consts.LineBreak
			geneUtils.Add2MapStrStr(rhCommon.EnvGoImports, mvcTemplates[consts.FileServerConfGoIndex].AppendImport)

			mvcTemplates[consts.FileServerDockerComposeIndex].Type = consts.Append
			mvcTemplates[consts.FileServerDockerComposeIndex].AppendContent += docker + consts.LineBreak

			if !disableAddConf {
				mvcTemplates[consts.FileServerDevConf].Type = consts.Append
				mvcTemplates[consts.FileServerDevConf].AppendContent += rhCommon.RegistryConfYaml + consts.LineBreak
				mvcTemplates[consts.FileServerOnlineConf].Type = consts.Append
				mvcTemplates[consts.FileServerOnlineConf].AppendContent += rhCommon.RegistryConfYaml + consts.LineBreak
				mvcTemplates[consts.FileServerTestConf].Type = consts.Append
				mvcTemplates[consts.FileServerTestConf].AppendContent += rhCommon.RegistryConfYaml + consts.LineBreak
			}
		}
	}

	if serverGen.CommunicationType == consts.RPC {
		kitexServerMVCTemplates = mvcTemplates
	} else {
		hzServerMVCTemplates = mvcTemplates
	}

	return nil
}
