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

	"github.com/cloudwego/cwgo/pkg/consts"
)

var errKeyInput = errors.New("input wrong key")

type ImportsMap map[string]map[string]struct{}

func newImportsMap(cmdType, protocolType string) (ImportsMap, error) {
	switch cmdType {
	case consts.Server:
		switch protocolType {
		case consts.RPC:
			return kitexServerInitImports, nil
		case consts.HTTP:
			return hzServerInitImports, nil
		default:
			return nil, errTypeInput
		}
	case consts.Client:
		switch protocolType {
		case consts.RPC:
			return kitexClientInitImports, nil
		case consts.HTTP:
			return hzClientInitImports, nil
		default:
			return nil, errTypeInput
		}
	default:
		return nil, errTypeInput
	}
}

func (m ImportsMap) appendImports(key string, imports []string) (err error) {
	// check whether the key exists
	if _, ok := m[key]; !ok {
		return errKeyInput
	}

	for _, imp := range imports {
		if _, ok := m[key][imp]; ok {
			continue
		}
		m[key][imp] = struct{}{}
	}

	return
}

var (
	kitexServerInitImports = map[string]map[string]struct{}{
		consts.ConfGo: {
			"io/ioutil":                           struct{}{},
			"os":                                  struct{}{},
			"path/filepath":                       struct{}{},
			"sync":                                struct{}{},
			"github.com/cloudwego/kitex/pkg/klog": struct{}{},
			"github.com/kr/pretty":                struct{}{},
			"gopkg.in/validator.v2":               struct{}{},
			"gopkg.in/yaml.v2":                    struct{}{},
		},

		consts.Main: {
			"net":                                    struct{}{},
			"github.com/cloudwego/kitex/pkg/klog":    struct{}{},
			"github.com/cloudwego/kitex/pkg/rpcinfo": struct{}{},
			"github.com/cloudwego/kitex/server":      struct{}{},
			"github.com/kitex-contrib/obs-opentelemetry/logging/logrus": struct{}{},
			"gopkg.in/natefinch/lumberjack.v2":                          struct{}{},
		},
	}

	kitexClientInitImports = map[string]map[string]struct{}{
		consts.InitGo: {
			"sync":                              struct{}{},
			"github.com/cloudwego/kitex/client": struct{}{},
		},

		consts.EnvGo: {},
	}

	hzServerInitImports = map[string]map[string]struct{}{
		consts.ConfGo: {
			"io/ioutil":     struct{}{},
			"os":            struct{}{},
			"path/filepath": struct{}{},
			"sync":          struct{}{},
			"github.com/cloudwego/hertz/pkg/common/hlog": struct{}{},
			"github.com/kr/pretty":                       struct{}{},
			"gopkg.in/validator.v2":                      struct{}{},
			"gopkg.in/yaml.v2":                           struct{}{},
		},

		consts.Main: {
			"context":                                        struct{}{},
			"github.com/cloudwego/hertz/pkg/app":             struct{}{},
			"github.com/cloudwego/hertz/pkg/app/server":      struct{}{},
			"github.com/cloudwego/hertz/pkg/common/config":   struct{}{},
			"github.com/cloudwego/hertz/pkg/common/hlog":     struct{}{},
			"github.com/cloudwego/hertz/pkg/common/utils":    struct{}{},
			"github.com/cloudwego/hertz/pkg/protocol/consts": struct{}{},
			"github.com/hertz-contrib/cors":                  struct{}{},
			"github.com/hertz-contrib/gzip":                  struct{}{},
			"github.com/hertz-contrib/logger/accesslog":      struct{}{},
			"github.com/hertz-contrib/logger/logrus":         struct{}{},
			"github.com/hertz-contrib/pprof":                 struct{}{},
			"gopkg.in/natefinch/lumberjack.v2":               struct{}{},
		},
	}

	hzClientInitImports = map[string]map[string]struct{}{
		consts.InitGo: {},

		consts.EnvGo: {},
	}
)
