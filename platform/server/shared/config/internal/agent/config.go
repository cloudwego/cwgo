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

package agent

type Config struct {
	Addr           string `mapstructure:"addr"`
	MaxConnections int64  `mapstructure:"maxConnections"`
	MaxQPS         int64  `mapstructure:"maxQPS"`
	WorkerNum      int    `mapstructure:"workerNum"`
}

type Metadata struct {
	ServiceId string `yaml:"service_id"`
}

func (conf *Config) SetUp() {
	conf.setDefaults()
}

func (conf *Config) setDefaults() {
	if conf.Addr == "" {
		conf.Addr = "0.0.0.0:11010"
	}

	if conf.MaxConnections == 0 {
		conf.MaxConnections = 2000
	}

	if conf.MaxQPS == 0 {
		conf.MaxQPS = 500
	}
}
