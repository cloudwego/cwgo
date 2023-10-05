/*
 *
 *  * Copyright 2022 CloudWeGo Authors
 *  *
 *  * Licensed under the Apache License, Version 2.0 (the "License");
 *  * you may not use this file except in compliance with the License.
 *  * You may obtain a copy of the License at
 *  *
 *  *     http://www.apache.org/licenses/LICENSE-2.0
 *  *
 *  * Unless required by applicable law or agreed to in writing, software
 *  * distributed under the License is distributed on an "AS IS" BASIS,
 *  * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  * See the License for the specific language governing permissions and
 *  * limitations under the License.
 *
 */

package api

import "github.com/cloudwego/cwgo/platform/server/shared/config/internal/dispatcher"

type Config struct {
	Host       string                `mapstructure:"host"`
	Port       int                   `mapstructure:"port"`
	Tracing    TracerConf            `mapstructure:"tracing"`
	MetricsUrl string                `mapstructure:"metricsUrl"`
	RpcClients map[int]RpcClientConf `mapstructure:"rpcClients"`
	Dispatcher dispatcher.Config     `mapstructure:"dispatcher"`
}

type TracerConf struct {
	Endpoint string  `mapstructure:"endpoint"`
	Sampler  float64 `mapstructure:"sampler"`
}

type RpcClientConf struct {
	Name          string `mapstructure:"name" json:"name"`
	MuxConnection int    `mapstructure:"muxConnection" json:"mux_connection,default=1"`
}
