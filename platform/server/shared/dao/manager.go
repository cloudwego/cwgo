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

package dao

import (
	"fmt"

	"github.com/cloudwego/cwgo/platform/server/shared/consts"
	"github.com/cloudwego/cwgo/platform/server/shared/log"
	"github.com/cloudwego/cwgo/platform/server/shared/config"

	"go.uber.org/zap"
)

type Manager struct {
	Idl        IIdlDaoManager
	Repository IRepositoryDaoManager
	Template   ITemplateDaoManager
	Token      ITokenDaoManager
}

func NewDaoManager(conf config.StoreConfig) (*Manager, error) {
	switch conf.GetStoreType() {
	case consts.StoreTypeNumMysql:
		log.Info("initializing mysql")
		mysqlDb, err := conf.NewMysqlDB()
		if err != nil {
			log.Error("initializing mysql failed", zap.Error(err))
			return nil, err
		}
		log.Info("initialize mysql successfully")

		idlDaoManager := NewMysqlIDL(mysqlDb)
		repositoryDaoManager := NewMysqlRepository(mysqlDb)
		templateDaoManager := NewMysqlTemplate(mysqlDb)
		tokenDaoManager := NewMysqlToken(mysqlDb)

		return &Manager{
			Idl:        idlDaoManager,
			Repository: repositoryDaoManager,
			Template:   templateDaoManager,
			Token:      tokenDaoManager,
		}, nil

	default:
		return nil, fmt.Errorf("invalid store type")
	}
}
