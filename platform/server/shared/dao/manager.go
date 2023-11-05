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

package dao

import (
	"fmt"
	"github.com/cloudwego/cwgo/platform/server/shared/config/store"
	"github.com/cloudwego/cwgo/platform/server/shared/consts"
	"github.com/cloudwego/cwgo/platform/server/shared/dao/internal/idl"
	"github.com/cloudwego/cwgo/platform/server/shared/dao/internal/repository"
	"github.com/cloudwego/cwgo/platform/server/shared/dao/internal/template"
	"github.com/cloudwego/cwgo/platform/server/shared/logger"
	"go.uber.org/zap"
)

type Manager struct {
	Idl        idl.IIdlDaoManager
	Repository repository.IRepositoryDaoManager
	Template   template.ITemplateDaoManager
}

func NewDaoManager(conf store.Config) (*Manager, error) {
	switch conf.GetStoreType() {
	case consts.StoreTypeNumMysql:
		logger.Logger.Info("initializing mysql")
		mysqlDb, err := conf.NewMysqlDB()
		if err != nil {
			logger.Logger.Error("initializing mysql failed", zap.Error(err))
			return nil, err
		}
		logger.Logger.Info("initialize mysql successfully")

		idlDaoManager := idl.NewMysqlIDL(mysqlDb)
		repositoryDaoManager := repository.NewMysqlRepository(mysqlDb)
		templateDaoManager := template.NewMysqlTemplate(mysqlDb)

		return &Manager{
			Idl:        idlDaoManager,
			Repository: repositoryDaoManager,
			Template:   templateDaoManager,
		}, nil

	case consts.StoreTypeNumMongo:
		panic("to be implemented")

	case consts.StoreTypeNumRedis:
		panic("to be implemented")

	default:
		return nil, fmt.Errorf("invalid store type")
	}
}
