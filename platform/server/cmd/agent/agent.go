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

import (
	"context"
	"errors"
	"os"

	"github.com/cloudwego/cwgo/platform/server/cmd/agent/pkg/processor"

	"github.com/cloudwego/cwgo/platform/server/shared/args"

	"github.com/cloudwego/cwgo/platform/server/cmd/agent/handler"
	"github.com/cloudwego/cwgo/platform/server/cmd/agent/internal/svc"
	"github.com/cloudwego/cwgo/platform/server/cmd/agent/pkg/generator"
	"github.com/cloudwego/cwgo/platform/server/shared/config"
	"github.com/cloudwego/cwgo/platform/server/shared/consts"
	"github.com/cloudwego/cwgo/platform/server/shared/dao"
	"github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/agent/agentservice"
	"github.com/cloudwego/cwgo/platform/server/shared/log"
	"github.com/cloudwego/cwgo/platform/server/shared/repository"
	"go.uber.org/zap"
)

func run(args *args.AgentArgs) error {
	var (
		serverMode consts.ServerMode
		configType consts.ConfigType
		metadata   any
		err        error
	)

	metadata, serverMode, configType, err = validateArgs(args)
	if err != nil {
		return err
	}
	// init config
	err = config.InitManager(consts.ServerTypeNumAgent, serverMode, configType, metadata)
	if err != nil {
		return err
	}

	// init logger
	loggerConfig := config.GetManager().Config.Logger
	log.InitLogger(
		log.Config{
			SavePath:     loggerConfig.SavePath,
			EncoderType:  loggerConfig.EncoderType,
			EncodeLevel:  loggerConfig.EncodeLevel,
			EncodeCaller: loggerConfig.EncodeCaller,
		},
		config.GetManager().ServerType,
		config.GetManager().ServiceID,
		config.GetManager().ServerMode,
	)

	// init dao manager
	log.Info("init dao manager")
	daoManager, err := dao.NewDaoManager(config.GetManager().Config.Store)
	if err != nil {
		log.Error("init dao manager failed", zap.Error(err))
		return err
	}
	log.Info("init dao manager successfully")

	log.Info("init dao manager")
	repoManager, err := repository.NewRepoManager(daoManager)
	if err != nil {
		log.Fatal("init repository manager failed", zap.Error(err))
	}
	log.Info("init dao manager successfully")

	ctx := context.Background()

	// get server options
	log.Info("getting kitex server options")
	kitexServerOptions := config.GetManager().AgentConfigManager.GetKitexServerOptions()
	log.Info("getting kitex server options successfully")

	// init agent service
	log.Info("init agent service impl")
	agentService := handler.NewAgentServiceImpl(
		ctx,
		&svc.ServiceContext{
			DaoManager:  daoManager,
			RepoManager: repoManager,
			Generator:   generator.NewCwgoGenerator(),
		},
	)
	log.Info("init agent service impl successfully")

	// init processor
	log.Info("init processor")

	processor.InitProcessor(agentService)

	log.Info("init processor successfully")

	// start service
	log.Info("register agent service")
	svr := agentservice.NewServer(
		agentService,
		kitexServerOptions...,
	)

	log.Info("start running agent service...")
	err = svr.Run()
	if err != nil {
		log.Error("kitex server run failed", zap.Error(err))
	}

	// graceful shutdown
	defer func() {
		err = svr.Stop()
		if err != nil {
			log.Error("kitex server gracefully shutdown failed", zap.Error(err))
		}
	}()

	return nil
}

func validateArgs(opts *args.AgentArgs) (metaData interface{}, serverMode consts.ServerMode, configType consts.ConfigType, err error) {
	var ok bool

	// priority: command line > env > config  > default
	if opts.ServerMode != "" {
		serverMode, ok = consts.ServerModeMapToNum[opts.ServerMode]
		if !ok {
			err = errors.New("invalid server_mode")
			return
		}
	}
	if serverMode == 0 {
		if serverModeStr := os.Getenv(consts.ServerTypeEnvName); serverModeStr != "" {
			serverMode, ok = consts.ServerModeMapToNum[serverModeStr]
			if !ok {
				err = errors.New("invalid server_mode")
				return
			}
		}
		if serverMode == 0 {
			serverMode = consts.ServerModeNumProd
		}
	}

	if opts.ConfigType != "" {
		configType, ok = consts.ConfigTypeMapToNum[opts.ConfigType]
		if !ok {
			err = errors.New("invalid config_type")
			return
		}
	}
	if configType == 0 {
		if configTypeStr := os.Getenv(consts.ConfigTypeEnvName); configTypeStr != "" {
			configType, ok = consts.ConfigTypeMapToNum[configTypeStr]
			if ok {
				err = errors.New("invalid config_type")
				return
			}
		}
		if configType == 0 {
			configType = consts.ConfigTypeNumFile
		}
	}

	switch configType {
	case consts.ConfigTypeNumFile:
		var configPath string

		if opts.ConfigPath != "" {
			configPath = opts.ConfigPath
		} else if configPath = os.Getenv(consts.ConfigPathEnvName); configPath == "" {
			configPath = consts.ConfigDefaultPath
		}

		metaData = config.FileConfig{
			Path: configPath,
		}
	default:

	}
	return
}
