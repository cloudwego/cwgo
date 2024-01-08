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

package svc

import (
	"github.com/cloudwego/cwgo/platform/server/cmd/api/pkg/manager"
	"github.com/cloudwego/cwgo/platform/server/shared/config"
	"github.com/cloudwego/cwgo/platform/server/shared/dao"
	"github.com/cloudwego/cwgo/platform/server/shared/log"
	"go.uber.org/zap"
)

type ServiceContext struct {
	Manager *manager.Manager
}

var Svc *ServiceContext

func InitServiceContext() {
	// init dao manager
	log.Info("init dao manager")
	daoManager, err := dao.NewDaoManager(config.GetManager().Config.Store)
	if err != nil {
		log.Fatal("init dao manager failed", zap.Error(err))
	}

	rdb, err := config.GetManager().Config.Store.NewRedisClient()
	if err != nil {
		log.Fatal("init redis failed", zap.Error(err))
	}

	log.Info("init dao manager successfully")
	Svc = &ServiceContext{
		Manager: manager.NewApiManager(
			config.GetManager().Config.App,
			config.GetManager().ServiceID,
			rdb,
			daoManager,
			config.GetManager().Config.Api.Dispatcher.NewDispatcher(),
			config.GetManager().ApiConfigManager.RegistryManager.GetRegistry(),
			config.GetManager().ApiConfigManager.RegistryManager.GetResolver(),
		),
	}
}
