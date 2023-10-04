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

package agent

import (
	"context"
	"errors"
	"github.com/cloudwego/cwgo/platform/server/cmd/agent/internal/svc"
	"github.com/cloudwego/cwgo/platform/server/shared/config"
	"github.com/cloudwego/cwgo/platform/server/shared/consts"
	"github.com/cloudwego/cwgo/platform/server/shared/dao"
	"github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/agent/agentservice"
	"github.com/cloudwego/cwgo/platform/server/shared/logger"
	"github.com/cloudwego/cwgo/platform/server/shared/repository"
	"go.uber.org/zap"
	"os"
)

func run(opts *setupOptions) error {
	var (
		serverMode consts.ServerMode
		configType consts.ConfigType
	)
	var ok bool
	// priority: command line > env > default
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
	err := config.InitManager(consts.ServerTypeNumAgent, serverMode, configType, metadata)
	if err != nil {
		return err
	}

	// init logger
	logger.InitLogger()

	// init dao
	daoManager, err := dao.NewDaoManager()
	if err != nil {
		return err
	}

	repoManager, err := repository.NewRepoManager(daoManager)
	if err != nil {
		logger.Logger.Fatal("service initialize repository manager failed", zap.Error(err))
	}

	ctx := context.Background()

	// get server options
	kitexServerOptions := config.GetManager().AgentConfigManager.GetKitexServerOptions()

	// start service
	svr := agentservice.NewServer(
		&AgentServiceImpl{
			ctx: ctx,
			svcCtx: &svc.ServiceContext{
				DaoManager:  daoManager,
				RepoManager: repoManager,
			},
		},
		kitexServerOptions...,
	)

	err = svr.Run()
	if err != nil {
		logger.Logger.Error("kitex server run failed", zap.Error(err))
	}

	return nil
}
