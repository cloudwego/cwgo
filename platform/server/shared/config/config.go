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

package config

type Cfg interface {
	Init()
}

type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Logger   LoggerConfig   `mapstructure:"logger"`
	Registry RegistryConfig `mapstructure:"registry"`
	Store    StoreConfig    `mapstructure:"store"`
	Api      ApiConfig      `mapstructure:"api"`
	Agent    AgentConfig    `mapstructure:"agent"`
}

func (conf *Config) Init() {
	conf.App.Init()
	conf.Logger.Init()
	conf.Registry.Init()
	conf.Store.Init()
	conf.Api.Init()
	conf.Agent.Init()
}
