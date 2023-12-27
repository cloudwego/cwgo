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

package dynamic

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/cloudwego/cwgo/pkg/consts"
)

func commonQuestion() []*survey.Question {
	return []*survey.Question{
		{
			Name: consts.ServiceType,
			Prompt: &survey.Select{
				Message: "Select service type",
				Options: []string{consts.RPC, consts.HTTP},
			},
			Validate: survey.Required,
		},
		{
			Name: consts.Service,
			Prompt: &survey.Input{
				Message: "Please input service name",
			},
			Validate: survey.Required,
		},
		{
			Name: "GoMod",
			Prompt: &survey.Input{
				Message: "Please input module",
			},
		},
		{
			Name: "idlPath",
			Prompt: &survey.Input{
				Message: "Please input idlpath",
			},
			Validate: survey.Required,
		},
	}
}

type dfConfig struct {
	DefaultConfig bool
}

func defaultConfig() []*survey.Question {
	return []*survey.Question{
		{
			Name: "defaultConfig",
			Prompt: &survey.Confirm{
				Message: "Whether use default config to generate project",
			},
			Validate: survey.Required,
		},
	}
}

func protoSearch() []*survey.Question {
	return []*survey.Question{{
		Name: consts.ProtoSearchPath,
		Prompt: &survey.Input{
			Message: "Please input proto search path if exists, space as separator",
		},
	}}
}

func customConfig() []*survey.Question {
	return []*survey.Question{
		{
			Name: consts.Pass,
			Prompt: &survey.Input{
				Message: "Please input custom param",
			},
		},
	}
}

func registryConfig() []*survey.Question {
	return []*survey.Question{
		{
			Name: consts.Registry,
			Prompt: &survey.Select{
				Message: "Please select a registry",
				Options: []string{consts.Zk, consts.Polaris, consts.Etcd, consts.Nacos, consts.Consul, consts.Eureka, consts.ServiceComb},
			},
			Validate: survey.Required,
		},
	}
}

func resolverConfig() []*survey.Question {
	return []*survey.Question{
		{
			Name: consts.Resolver,
			Prompt: &survey.Select{
				Message: "Please select a resolver",
				Options: []string{consts.Zk, consts.Polaris, consts.Etcd, consts.Nacos, consts.Consul, consts.Eureka, consts.ServiceComb},
			},
			Validate: survey.Required,
		},
	}
}
