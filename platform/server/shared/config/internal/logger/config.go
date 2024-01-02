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

package logger

import "github.com/cloudwego/cwgo/platform/server/shared/log"

type Config struct {
	SavePath     string `mapstructure:"savePath"`
	EncoderType  string `mapstructure:"encoderType"`
	EncodeLevel  string `mapstructure:"encodeLevel"`
	EncodeCaller string `mapstructure:"encodeCaller"`
}

func (conf *Config) SetUp() {
	conf.setDefaults()
}

func (conf *Config) setDefaults() {
	if conf.SavePath == "" {
		conf.SavePath = "log"
	}

	if conf.EncoderType == "" {
		conf.EncoderType = log.ConsoleEncoder
	}

	if conf.EncodeLevel == "" {
		conf.EncodeLevel = log.CapitalLevelEncoder
	}

	if conf.EncodeCaller == "" {
		conf.EncodeCaller = log.FullCallerEncoder
	}
}
