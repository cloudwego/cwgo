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

package dynamic

import (
	"strconv"

	"github.com/AlecAivazis/survey/v2"
	"github.com/cloudwego/cwgo/config"
	"github.com/cloudwego/cwgo/pkg/client"
	"github.com/cloudwego/cwgo/pkg/common/utils"
	"github.com/cloudwego/cwgo/pkg/consts"
	"github.com/cloudwego/cwgo/pkg/model"
	"github.com/cloudwego/cwgo/pkg/server"
	"github.com/cloudwego/hertz/cmd/hz/meta"
	"github.com/urfave/cli/v2"
)

var generateType = []*survey.Question{
	{
		Name: consts.ServiceType,
		Prompt: &survey.MultiSelect{
			Message: "Select generate type",
			Options: []string{consts.Server, consts.Client, consts.DB},
		},
		Validate: survey.Required,
	},
}

func Terminal(*cli.Context) error {
	// Select which types of services to generate
	var idlProject []string
	err := survey.Ask(generateType, &idlProject)
	if err != nil {
		return err
	}
	typeMap := make(map[string]struct{})
	for _, i := range idlProject {
		typeMap[i] = struct{}{}
	}

	cfg := &dfConfig{}
	// ask whether generate server project
	if _, ok := typeMap[consts.Server]; ok {
		sa := config.NewServerArgument()
		err = survey.Ask(commonQuestion(), sa.CommonParam)
		if err != nil {
			return err
		}

		if t, _ := utils.GetIdlType(sa.IdlPath); t == meta.IdlProto {
			err = survey.Ask(protoSearch(), sa.SliceParam)
			if err != nil {
				return err
			}
		}
		err = survey.Ask(defaultConfig(), cfg)
		if err != nil {
			return err
		}
		if !cfg.DefaultConfig {
			if err = survey.Ask(registryConfig(), sa.CommonParam); err != nil {
				return err
			}
			err = survey.Ask(customConfig(), sa.SliceParam)
			if err != nil {
				return err
			}
		}
		err = server.Server(sa)
		if err != nil {
			return err
		}
	}

	num := &cNum{}
	// ask whether generate client project
	if _, ok := typeMap[consts.Client]; ok {
		err = survey.Ask(clientNum, num)
		if err != nil {
			return err
		}
		n, err := strconv.Atoi(num.ClientNum)
		if err != nil {
			return err
		}
		for i := 0; i < n; i++ {
			ca := config.NewClientArgument()
			err = survey.Ask(commonQuestion(), ca.CommonParam)
			if err != nil {
				return err
			}
			if t, _ := utils.GetIdlType(ca.IdlPath); t == meta.IdlProto {
				err = survey.Ask(protoSearch(), ca.SliceParam)
				if err != nil {
					return err
				}
			}
			err = survey.Ask(defaultConfig(), cfg)
			if err != nil {
				return err
			}
			if !cfg.DefaultConfig {
				if err = survey.Ask(resolverConfig(), ca.CommonParam); err != nil {
					return err
				}
				err = survey.Ask(customConfig(), ca.SliceParam)
				if err != nil {
					return err
				}
			}
			err = client.Client(ca)
			if err != nil {
				return err
			}
		}
	}

	// ask whether generate db project
	if _, ok := typeMap[consts.DB]; ok {
		da := config.NewModelArgument()
		err = survey.Ask(dbConfig, da)
		if err != nil {
			return err
		}
		err = survey.Ask(defaultConfig(), cfg)
		if err != nil {
			return err
		}
		if !cfg.DefaultConfig {
			p := &ps{}
			err = survey.Ask(pass, p)
			if err != nil {
				return err
			}
			err = parsePass(da, p.Pass)
			if err != nil {
				return err
			}
		}

		err = model.Model(da)
		if err != nil {
			return err
		}
	}

	return nil
}
