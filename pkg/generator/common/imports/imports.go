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

package imports

import (
	"github.com/cloudwego/cwgo/pkg/consts"
	"github.com/cloudwego/cwgo/pkg/generator/rpchttp/common"
)

type Map map[string]map[string]string

func NewMap(cmdType, protocolType string) (Map, error) {
	switch cmdType {
	case consts.Server:
		switch protocolType {
		case consts.RPC:
			return kitexServerInitImports, nil
		case consts.HTTP:
			return hzServerInitImports, nil
		default:
			return nil, common.ErrTypeInput
		}
	case consts.Client:
		switch protocolType {
		case consts.RPC:
			return kitexClientInitImports, nil
		case consts.HTTP:
			return hzClientInitImports, nil
		default:
			return nil, common.ErrTypeInput
		}
	default:
		return nil, common.ErrTypeInput
	}
}

func (m Map) AppendImports(key string, imports map[string]string) (err error) {
	// check whether the key exists
	if _, ok := m[key]; !ok {
		return common.ErrKeyInput
	}

	for k, imp := range imports {
		if _, ok := m[key][imp]; ok {
			continue
		}
		m[key][k] = imp
	}

	return
}

var (
	kitexServerInitImports = map[string]map[string]string{
		consts.ConfGo: {
			"io/ioutil":                           "",
			"os":                                  "",
			"path/filepath":                       "",
			"sync":                                "",
			"github.com/cloudwego/kitex/pkg/klog": "",
			"github.com/kr/pretty":                "",
			"gopkg.in/validator.v2":               "",
			"gopkg.in/yaml.v2":                    "",
		},

		consts.Main: {
			"net":                                    "",
			"github.com/cloudwego/kitex/pkg/klog":    "",
			"github.com/cloudwego/kitex/pkg/rpcinfo": "",
			"github.com/cloudwego/kitex/server":      "",
			"github.com/kitex-contrib/obs-opentelemetry/logging/logrus": "kitexLogrus",
			"gopkg.in/natefinch/lumberjack.v2":                          "",
		},

		consts.DalInitGo: {},

		consts.MysqlInit: {
			"gorm.io/driver/mysql": "",
			"gorm.io/gorm":         "",
		},

		consts.RedisInit: {
			"context":                      "",
			"github.com/redis/go-redis/v9": "",
		},
	}

	kitexClientInitImports = map[string]map[string]string{
		consts.InitGo: {
			"sync":                              "",
			"github.com/cloudwego/kitex/client": "",
		},

		consts.EnvGo: {},
	}

	hzServerInitImports = map[string]map[string]string{
		consts.ConfGo: {
			"io/ioutil":     "",
			"os":            "",
			"path/filepath": "",
			"sync":          "",
			"github.com/cloudwego/hertz/pkg/common/hlog": "",
			"github.com/kr/pretty":                       "",
			"gopkg.in/validator.v2":                      "",
			"gopkg.in/yaml.v2":                           "",
		},

		consts.Main: {
			"context":                                        "",
			"github.com/cloudwego/hertz/pkg/app":             "",
			"github.com/cloudwego/hertz/pkg/app/server":      "",
			"github.com/cloudwego/hertz/pkg/common/config":   "",
			"github.com/cloudwego/hertz/pkg/common/hlog":     "",
			"github.com/cloudwego/hertz/pkg/common/utils":    "",
			"github.com/cloudwego/hertz/pkg/protocol/consts": "",
			"github.com/hertz-contrib/cors":                  "",
			"github.com/hertz-contrib/gzip":                  "",
			"github.com/hertz-contrib/logger/accesslog":      "",
			"github.com/hertz-contrib/logger/logrus":         "",
			"github.com/hertz-contrib/pprof":                 "",
			"gopkg.in/natefinch/lumberjack.v2":               "",
		},

		consts.DalInitGo: {},

		consts.MysqlInit: {
			"gorm.io/driver/mysql": "",
			"gorm.io/gorm":         "",
		},

		consts.RedisInit: {
			"context":                      "",
			"github.com/redis/go-redis/v9": "",
		},

		consts.RegisterGo: {
			"github.com/cloudwego/hertz/pkg/app/server": "",
		},

		consts.RespGo: {
			"context":                            "",
			"github.com/cloudwego/hertz/pkg/app": "",
		},
	}

	hzClientInitImports = map[string]map[string]string{
		consts.InitGo: {},

		consts.EnvGo: {},
	}
)
