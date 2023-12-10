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

import (
	"errors"
	"os"

	"github.com/cloudwego/cwgo/platform/server/cmd/api/internal/svc"
	"github.com/cloudwego/cwgo/platform/server/shared/config"
	"github.com/cloudwego/cwgo/platform/server/shared/consts"
	"github.com/cloudwego/cwgo/platform/server/shared/logger"
)

func run(opts *setupOptions) error {
	var (
		serverMode     consts.ServerMode
		configType     consts.ConfigType
		staticFilePath string
	)
	var ok bool
	// priority: command line > env > config > default
	if opts.ServerMode != "" {
		serverMode, ok = consts.ServerModeMapToNum[opts.ServerMode]
		if !ok {
			return errors.New("invalid server_mode")
		}
	}
	if serverMode == 0 {
		if serverModeStr := os.Getenv(consts.ServerTypeEnvName); serverModeStr != "" {
			serverMode, ok = consts.ServerModeMapToNum[serverModeStr]
			if !ok {
				return errors.New("invalid server_mode")
			}
		}
		if serverMode == 0 {
			serverMode = consts.ServerModeNumPro
		}
	}

	if opts.ConfigType != "" {
		configType, ok = consts.ConfigTypeMapToNum[opts.ConfigType]
		if !ok {
			return errors.New("invalid config_type")
		}
	}
	if configType == 0 {
		if configTypeStr := os.Getenv(consts.ConfigTypeEnvName); configTypeStr != "" {
			configType, ok = consts.ConfigTypeMapToNum[configTypeStr]
			if ok {
				return errors.New("invalid config_type")
			}
		}
		if configType == 0 {
			configType = consts.ConfigTypeNumFile
		}
	}
	if opts.StaticFilePath != "" {
		staticFilePath = opts.StaticFilePath
	}
	if staticFilePath == "" {
		if staticFilePathStr := os.Getenv(consts.StaticFilePathEnvName); staticFilePathStr != "" {
			staticFilePath = staticFilePathStr
		}
		if staticFilePath == "" {
			staticFilePath = consts.StaticFileDefaultPath
		}
	}

	var metadata interface{}
	switch configType {
	case consts.ConfigTypeNumFile:
		var configPath string

		if opts.ConfigPath != "" {
			configPath = opts.ConfigPath
		} else if configPath = os.Getenv(consts.ConfigPathEnvName); configPath == "" {
			configPath = consts.ConfigDefaultPath
		}

		metadata = config.FileConfig{
			Path: configPath,
		}
	}

	// init config
	err := config.InitManager(consts.ServerTypeNumApi, serverMode, configType, metadata)
	if err != nil {
		return err
	}

	// init logger
	loggerConfig := config.GetManager().Config.Logger
	logger.InitLogger(
		logger.Config{
			SavePath:     loggerConfig.SavePath,
			EncoderType:  loggerConfig.EncoderType,
			EncodeLevel:  loggerConfig.EncodeLevel,
			EncodeCaller: loggerConfig.EncodeCaller,
		},
		config.GetManager().ServerType,
		config.GetManager().ServiceId,
		config.GetManager().ServerMode,
	)

	// init service context
	svc.InitServiceContext()

	// start api service
	logger.Logger.Info("register api service")
	register(config.GetManager().ApiConfigManager.Server, staticFilePath)

	logger.Logger.Info("start running api service...")
	config.GetManager().ApiConfigManager.Server.Spin()

	return nil
}
