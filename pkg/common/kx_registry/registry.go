// Copyright 2022 CloudWeGo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package kx_registry

import (
	"fmt"
	"os"
	"path"

	"github.com/cloudwego/cwgo/config"
	"github.com/cloudwego/cwgo/pkg/consts"
	"github.com/cloudwego/cwgo/tpl"
	"github.com/cloudwego/kitex/tool/internal_pkg/generator"
)

func HandleRegistry(ca *config.CommonParam, dir string) {
	te := &generator.TemplateExtension{
		Dependencies: map[string]string{
			ca.GoMod + "/conf":                       "conf",
			"github.com/cloudwego/kitex/pkg/klog":    "klog",
			"github.com/cloudwego/kitex/pkg/rpcinfo": "rpcinfo",
		},
	}

	importPath := []string{ca.GoMod + "/conf", "github.com/cloudwego/kitex/pkg/klog"}

	switch ca.Registry {
	case consts.Etcd:
		te.Dependencies["github.com/kitex-contrib/registry-etcd"] = "etcd"
		te.ExtendServer = &generator.APIExtension{
			ImportPaths:  append(importPath, "github.com/kitex-contrib/registry-etcd", "github.com/cloudwego/kitex/pkg/rpcinfo"),
			ExtendOption: fmt.Sprintf(etcdServer, ca.Service),
		}
		te.ExtendClient = &generator.APIExtension{
			ImportPaths:  append(importPath, "github.com/kitex-contrib/registry-etcd"),
			ExtendOption: etcdClient,
		}
	case consts.Zk:
		te.Dependencies["github.com/kitex-contrib/registry-zookeeper/registry"] = "zkregistry"
		te.Dependencies["github.com/kitex-contrib/registry-zookeeper/resolver"] = "zkresolver"
		te.Dependencies["time"] = "time"
		te.ExtendServer = &generator.APIExtension{
			ImportPaths:  append(importPath, "github.com/kitex-contrib/registry-zookeeper/registry", "time"),
			ExtendOption: zkServer,
		}
		te.ExtendClient = &generator.APIExtension{
			ImportPaths:  append(importPath, "github.com/kitex-contrib/registry-zookeeper/resolver", "time"),
			ExtendOption: zkClient,
		}
	case consts.Polaris:
		te.Dependencies["github.com/kitex-contrib/registry-polaris"] = "polaris"
		te.Dependencies["github.com/cloudwego/kitex/pkg/registry"] = "registry"
		te.ExtendServer = &generator.APIExtension{
			ImportPaths:  []string{"github.com/cloudwego/kitex/pkg/registry", "github.com/kitex-contrib/registry-polaris", "github.com/cloudwego/kitex/pkg/klog"},
			ExtendOption: fmt.Sprintf(polarisServer, ca.Service),
		}
		te.ExtendClient = &generator.APIExtension{
			ImportPaths:  []string{"github.com/cloudwego/kitex/pkg/registry", "github.com/kitex-contrib/registry-polaris", "github.com/cloudwego/kitex/pkg/klog"},
			ExtendOption: fmt.Sprintf(polarisClient, ca.Service),
		}
	case consts.Nacos:
		te.Dependencies["github.com/kitex-contrib/registry-nacos/registry"] = "registry"
		te.Dependencies["github.com/kitex-contrib/registry-nacos/resolver"] = "resolver"
		te.ExtendServer = &generator.APIExtension{
			ImportPaths:  []string{"github.com/cloudwego/kitex/pkg/klog", "github.com/kitex-contrib/registry-nacos/registry", "github.com/cloudwego/kitex/pkg/rpcinfo"},
			ExtendOption: fmt.Sprintf(nacosServer, ca.Service),
		}
		te.ExtendClient = &generator.APIExtension{
			ImportPaths:  []string{"github.com/cloudwego/kitex/pkg/klog", "github.com/kitex-contrib/registry-nacos/resolver"},
			ExtendOption: nacosClient,
		}
	default:
		RemoveExtension()
		return
	}

	path := path.Join(dir, consts.KitexExtensionYaml)
	te.ToYAMLFile(path)
}

func RemoveExtension() {
	path := tpl.KitexDir + consts.KitexExtensionYaml
	os.RemoveAll(path)
}

const etcdServer = `
	r, err := etcd.NewEtcdRegistryWithAuth(conf.GetConf().Registry.RegistryAddress, conf.GetConf().Registry.Username, conf.GetConf().Registry.Password)
	if err != nil {
		klog.Fatal(err)
	}
	 options = append(options, server.WithRegistry(r), server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
		ServiceName: "%s",
	}))
`

const etcdClient = `
	r, err := etcd.NewEtcdResolverWithAuth(conf.GetConf().Registry.RegistryAddress, conf.GetConf().Registry.Username, conf.GetConf().Registry.Password)
	if err != nil {
		klog.Fatal(err)
	}
	options = append(options, client.WithResolver(r))
`

const zkServer = `
	r, err := zkregistry.NewZookeeperRegistryWithAuth(conf.GetConf().Registry.RegistryAddress, 30*time.Second, conf.GetConf().Registry.Username, conf.GetConf().Registry.Password)
    if err != nil{
        klog.Fatal(err)
    }
	options = append(options, server.WithRegistry(r))
`

const zkClient = `
	r, err := zkresolver.NewZookeeperResolverWithAuth(conf.GetConf().Registry.RegistryAddress, 30*time.Second, conf.GetConf().Registry.Username, conf.GetConf().Registry.Password)
    if err != nil {
		klog.Fatal(err)
    }
	options = append(options, client.WithResolver(r))
`

const nacosServer = `
	r, err := registry.NewDefaultNacosRegistry()
	if err != nil {
		klog.Fatal(err)
	}
	options = append(options, server.WithRegistry(r), server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
		ServiceName: "%s",
	}))
`

const nacosClient = `
	r, err := resolver.NewDefaultNacosResolver()
	if err != nil {
		klog.Fatal(err)
	}
	options = append(options, client.WithResolver(r))
`

const polarisServer = `
	r, err := polaris.NewPolarisRegistry()
	if err != nil {
		klog.Fatal(err)
	}
	Info := &registry.Info{
		ServiceName: %s",
		Tags: map[string]string{
			"namespace": "Polaris",
		},
	}
	options = append(options, server.WithRegistry(r))
`

const polarisClient = `
	r, err := polaris.NewPolarisRegistry()
	if err != nil {
		klog.Fatal(err)
	}
	Info := &registry.Info{
		ServiceName: %s",
		Tags: map[string]string{
			"namespace": "Polaris",
		},
	}
	options = append(options, client.WithResolver(r))
`
