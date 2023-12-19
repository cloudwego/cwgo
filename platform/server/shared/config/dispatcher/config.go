/*
 *
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
 *
 */

package dispatcher

import (
	"github.com/cloudwego/cwgo/platform/server/cmd/api/pkg/dispatcher"
	"github.com/cloudwego/cwgo/platform/server/shared/consts"
)

type Config struct {
	Type string `mapstructure:"type"`
}

func (conf *Config) SetUp() {
	conf.setDefaults()
}

func (conf *Config) setDefaults() {
	conf.Type = consts.DispatcherTypeHash
}

func (conf *Config) NewDispatcher() dispatcher.IDispatcher {
	dispatcherType, ok := consts.DispatcherMapToNum[conf.Type]
	if !ok {
		panic("invalid dispatcher type")
	}

	switch dispatcherType {
	case consts.DispatcherTypeNumHash:
		return dispatcher.NewConsistentHashDispatcher()
	default:
		return nil
	}
}
