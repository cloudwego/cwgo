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
	"github.com/cloudwego/cwgo/pkg/consts"
	"github.com/cloudwego/cwgo/pkg/generator/common/template"
	geneUtils "github.com/cloudwego/cwgo/pkg/generator/common/utils"
	rhCommon "github.com/cloudwego/cwgo/pkg/generator/rpchttp/common"
)

func (clientGen *Generator) handleResolver(resolverName string) (err error) {
	if clientGen.IsNew {
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

func (clientGen *Generator) handleNewResolver(resolverName string) (err error) {
	// custom server resolver
	if clientGen.CustomExtensionFile != "" && clientGen.ResolverName != "" {
		if err = clientGen.handleNewResolverTemplate(clientGen.ResolverBody, clientGen.ResolverAddress, clientGen.ResolverImports, false, false); err != nil {
			return
		}
		return
	}

	clientGen.ResolverName = resolverName

	switch clientGen.CommunicationType {
	case consts.RPC:
		switch clientGen.ResolverName {
		case consts.Nacos:
			if err = clientGen.handleNewResolverTemplate(kitexNacosClient, rhCommon.NacosServerAddr, kitexNacosClientImports, false, false); err != nil {
				return
			}
		case consts.Consul:
			if err = clientGen.handleNewResolverTemplate(kitexConsulClient, rhCommon.ConsulServerAddr, kitexConsulClientImports, true, false); err != nil {
				return
			}
		case consts.Etcd:
			if err = clientGen.handleNewResolverTemplate(kitexEtcdClient, rhCommon.EtcdServerAddr, kitexEtcdClientImports, true, false); err != nil {
				return
			}
		case consts.Eureka:
			if err = clientGen.handleNewResolverTemplate(kitexEurekaClient, rhCommon.EurekaServerAddr, kitexEurekaClientImports, true, false); err != nil {
				return
			}
		case consts.Polaris:
			if err = clientGen.handleNewResolverTemplate(kitexPolarisClient, rhCommon.PolarisServerAddr, kitexPolarisClientImports, false, false); err != nil {
				return
			}
		case consts.ServiceComb:
			if err = clientGen.handleNewResolverTemplate(kitexServiceCombClient, rhCommon.ServiceCombServerAddr, kitexServiceCombClientImports, false, false); err != nil {
				return
			}
		case consts.Zk:
			if err = clientGen.handleNewResolverTemplate(kitexZKClient, rhCommon.ZkServerAddr, kitexZKClientImports, true, false); err != nil {
				return
			}
		default:
			if err = clientGen.compatibleOlderVersion(); err != nil {
				return
			}
		}

	case consts.HTTP:
		switch clientGen.ResolverName {
		case consts.Nacos:
			if err = clientGen.handleNewResolverTemplate(hzNacosClient, rhCommon.NacosServerAddr, hzNacosClientImports, false, false); err != nil {
				return
			}
		case consts.Consul:
			if err = clientGen.handleNewResolverTemplate(hzConsulClient, rhCommon.ConsulServerAddr, hzConsulClientImports, false, true); err != nil {
				return
			}
		case consts.Etcd:
			if err = clientGen.handleNewResolverTemplate(hzEtcdClient, rhCommon.EtcdServerAddr, hzEtcdClientImports, false, true); err != nil {
				return
			}
		case consts.Eureka:
			if err = clientGen.handleNewResolverTemplate(hzEurekaClient, rhCommon.EurekaServerAddr, hzEurekaClientImports, false, true); err != nil {
				return
			}
		case consts.Polaris:
			if err = clientGen.handleNewResolverTemplate(hzPolarisClient, rhCommon.PolarisServerAddr, hzPolarisClientImports, false, false); err != nil {
				return
			}
		case consts.ServiceComb:
			if err = clientGen.handleNewResolverTemplate(hzServiceCombClient, rhCommon.ServiceCombServerAddr, hzServiceCombClientImports, false, true); err != nil {
				return
			}
		case consts.Zk:
			if err = clientGen.handleNewResolverTemplate(hzZKClient, rhCommon.ZkServerAddr, hzZKClientImports, false, true); err != nil {
				return
			}
		default:
			if err = clientGen.compatibleOlderVersion(); err != nil {
				return
			}
		}

	default:
		return rhCommon.ErrTypeInput
	}

	return
}

func (clientGen *Generator) handleNewResolverTemplate(body string, addr []string, imports map[string]string, needRpc, needHttp bool) (err error) {
	var (
		mvcTemplates           []template.Template
		appendInitResolverFunc string
		appendResolverAddrFunc string
	)

	if clientGen.CommunicationType == consts.RPC {
		mvcTemplates = kitexClientMVCTemplates
		appendInitResolverFunc = kitexAppendInitResolverFunc
		appendResolverAddrFunc = kitexAppendResolverAddrFunc
	} else {
		mvcTemplates = hzClientMVCTemplates
		appendInitResolverFunc = hzAppendInitResolverFunc
		appendResolverAddrFunc = hzAppendResolverAddrFunc
	}

	clientGen.ResolverBody = body
	clientGen.ResolverAddress = addr

	if clientGen.CustomExtensionFile == "" {
		if clientGen.CommunicationType == consts.HTTP {
			if needHttp {
				imports[clientGen.PackagePrefix+consts.Slash+consts.DefaultHZClientDir] = "hzHttp"
			}
		} else {
			if needRpc {
				imports[clientGen.PackagePrefix+consts.Slash+consts.DefaultKitexClientDir] = "kitexRpc"
			}
		}
	}

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
				geneUtils.Add2MapStrStr(imports, mvcTemplates[consts.FileClientInitIndex].AppendImport)
			}
			if !isExist && index != 0 {
				t := template.Template{
					Path: clientGen.initGoPaths[index],
					UpdateBehavior: template.UpdateBehavior{
						Type: consts.Append,
						Append: template.Append{
							AppendContent: appendInitResolverFunc,
							AppendImport:  imports,
						},
						AppendRender: map[string]interface{}{},
					},
				}

				mvcTemplates = append(mvcTemplates, t)
			}
		}
	} else {
		if err = clientGen.GoFileImports.AppendImports(consts.InitGo, imports); err != nil {
			return err
		}
	}

	if clientGen.envGoContent != "" {
		isExist, err := geneUtils.IsFuncExist(clientGen.envGoContent, consts.FuncGetResolverAddress)
		if err != nil {
			return err
		}
		if !isExist {
			mvcTemplates[consts.FileClientEnvIndex].Type = consts.Append
			mvcTemplates[consts.FileClientEnvIndex].AppendContent += appendResolverAddrFunc + consts.LineBreak
			geneUtils.Add2MapStrStr(rhCommon.EnvGoImports, mvcTemplates[consts.FileClientEnvIndex].AppendImport)
		}
	} else {
		if err = clientGen.GoFileImports.AppendImports(consts.EnvGo, rhCommon.EnvGoImports); err != nil {
			return err
		}
	}

	if clientGen.CommunicationType == consts.RPC {
		kitexClientMVCTemplates = mvcTemplates
	} else {
		hzClientMVCTemplates = mvcTemplates
	}

	return
}

func (clientGen *Generator) handleUpdateResolver(resolverName string) (err error) {
	if clientGen.CustomExtensionFile != "" && clientGen.ResolverName != "" {
		if err = clientGen.handleUpdateResolverTemplate(clientGen.ResolverBody, clientGen.ResolverAddress, clientGen.ResolverImports, false, false); err != nil {
			return err
		}
		return
	}

	clientGen.ResolverName = resolverName

	switch clientGen.CommunicationType {
	case consts.RPC:
		switch clientGen.ResolverName {
		case consts.Nacos:
			if err = clientGen.handleUpdateResolverTemplate(kitexNacosClient, rhCommon.NacosServerAddr, kitexNacosClientImports, false, false); err != nil {
				return err
			}
		case consts.Consul:
			if err = clientGen.handleUpdateResolverTemplate(kitexConsulClient, rhCommon.ConsulServerAddr, kitexConsulClientImports, true, false); err != nil {
				return err
			}
		case consts.Etcd:
			if err = clientGen.handleUpdateResolverTemplate(kitexEtcdClient, rhCommon.EtcdServerAddr, kitexEtcdClientImports, true, false); err != nil {
				return err
			}
		case consts.Eureka:
			if err = clientGen.handleUpdateResolverTemplate(kitexEurekaClient, rhCommon.EurekaServerAddr, kitexEurekaClientImports, true, false); err != nil {
				return err
			}
		case consts.Polaris:
			if err = clientGen.handleUpdateResolverTemplate(kitexPolarisClient, rhCommon.PolarisServerAddr, kitexPolarisClientImports, false, false); err != nil {
				return err
			}
		case consts.ServiceComb:
			if err = clientGen.handleUpdateResolverTemplate(kitexServiceCombClient, rhCommon.ServiceCombServerAddr, kitexServiceCombClientImports, false, false); err != nil {
				return err
			}
		case consts.Zk:
			if err = clientGen.handleUpdateResolverTemplate(kitexZKClient, rhCommon.ZkServerAddr, kitexZKClientImports, true, false); err != nil {
				return err
			}
		default:
		}

	case consts.HTTP:
		switch clientGen.ResolverName {
		case consts.Nacos:
			if err = clientGen.handleUpdateResolverTemplate(hzNacosClient, rhCommon.NacosServerAddr, hzNacosClientImports, false, false); err != nil {
				return err
			}
		case consts.Consul:
			if err = clientGen.handleUpdateResolverTemplate(hzConsulClient, rhCommon.ConsulServerAddr, hzConsulClientImports, false, true); err != nil {
				return err
			}
		case consts.Etcd:
			if err = clientGen.handleUpdateResolverTemplate(hzEtcdClient, rhCommon.EtcdServerAddr, hzEtcdClientImports, false, true); err != nil {
				return err
			}
		case consts.Eureka:
			if err = clientGen.handleUpdateResolverTemplate(hzEurekaClient, rhCommon.EurekaServerAddr, hzEurekaClientImports, false, true); err != nil {
				return err
			}
		case consts.Polaris:
			if err = clientGen.handleUpdateResolverTemplate(hzPolarisClient, rhCommon.PolarisServerAddr, hzPolarisClientImports, false, false); err != nil {
				return err
			}
		case consts.ServiceComb:
			if err = clientGen.handleUpdateResolverTemplate(hzServiceCombClient, rhCommon.ServiceCombServerAddr, hzServiceCombClientImports, false, true); err != nil {
				return err
			}
		case consts.Zk:
			if err = clientGen.handleUpdateResolverTemplate(hzZKClient, rhCommon.ZkServerAddr, hzZKClientImports, false, true); err != nil {
				return err
			}
		default:
		}

	default:
		return rhCommon.ErrTypeInput
	}

	return nil
}

func (clientGen *Generator) handleUpdateResolverTemplate(body string, addr []string, imports map[string]string, needRpc, needHttp bool) error {
	clientGen.ResolverAddress = addr

	var (
		mvcTemplates        []template.Template
		nilResolverFuncBody string
		appendResolverFunc  string
	)
	if clientGen.CommunicationType == consts.RPC {
		if needRpc {
			imports[clientGen.PackagePrefix+consts.Slash+consts.DefaultKitexClientDir] = "kitexRpc"
		}
		mvcTemplates = kitexClientMVCTemplates
		nilResolverFuncBody = kitexNilResolverFuncBody
		appendResolverFunc = kitexAppendResolverAddrFunc
	} else {
		if needHttp {
			imports[clientGen.PackagePrefix+consts.Slash+consts.DefaultHZClientDir] = "hzHttp"
		}
		mvcTemplates = hzClientMVCTemplates
		nilResolverFuncBody = hzNilResolverFuncBody
		appendResolverFunc = hzAppendResolverAddrFunc
	}

	flag := 0
	for index, content := range clientGen.initGoContents {
		equal, err := geneUtils.IsFuncBodyEqual(content, consts.FuncInitResolver, nilResolverFuncBody)
		if err != nil {
			return err
		}

		if equal {
			if index == 0 {
				mvcTemplates[consts.FileClientInitIndex].Path = clientGen.initGoPaths[index]
				mvcTemplates[consts.FileClientInitIndex].Type = consts.ReplaceFuncBody
				mvcTemplates[consts.FileClientInitIndex].ReplaceFuncName = append(mvcTemplates[consts.FileClientInitIndex].ReplaceFuncName, consts.FuncInitResolver)
				mvcTemplates[consts.FileClientInitIndex].ReplaceFuncAppendImport = append(mvcTemplates[consts.FileClientInitIndex].ReplaceFuncAppendImport, imports)
				mvcTemplates[consts.FileClientInitIndex].ReplaceFuncBody = append(mvcTemplates[consts.FileClientInitIndex].ReplaceFuncBody, body)
			} else {
				t := template.Template{
					Path: clientGen.initGoPaths[index],
					UpdateBehavior: template.UpdateBehavior{
						Type: consts.ReplaceFuncBody,
						ReplaceFunc: template.ReplaceFunc{
							ReplaceFuncName:         []string{consts.FuncInitResolver},
							ReplaceFuncAppendImport: []map[string]string{imports},
							ReplaceFuncBody:         []string{body},
						},
						AppendRender: map[string]interface{}{},
					},
				}

				mvcTemplates = append(mvcTemplates, t)
			}

			if flag == 0 {
				isExist, err := geneUtils.IsFuncExist(clientGen.envGoContent, consts.FuncGetResolverAddress)
				if err != nil {
					return err
				}

				if !isExist {
					mvcTemplates[consts.FileClientEnvIndex].Path = clientGen.envGoPath
					mvcTemplates[consts.FileClientEnvIndex].Type = consts.Append
					mvcTemplates[consts.FileClientEnvIndex].AppendContent += appendResolverFunc + consts.LineBreak
					geneUtils.Add2MapStrStr(rhCommon.EnvGoImports, mvcTemplates[consts.FileClientEnvIndex].AppendImport)
				}

				flag++
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
