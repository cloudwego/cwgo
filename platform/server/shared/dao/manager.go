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
	"github.com/cloudwego/cwgo/platform/server/shared/config"
	"github.com/cloudwego/cwgo/platform/server/shared/consts"
	"github.com/cloudwego/cwgo/platform/server/shared/dao/idl"
	"github.com/cloudwego/cwgo/platform/server/shared/dao/repository"
	"github.com/cloudwego/cwgo/platform/server/shared/dao/template"
)

type Manager struct {
	Idl        idl.IIdlDaoManager
	Repository repository.IRepositoryDaoManager
	Template   template.ITemplateDaoManager
}

func NewDaoManager() (*Manager, error) {
	switch config.GetManager().StoreConfigManager.GetStoreType() {
	case consts.StoreTypeMysql:
		mysqlDb, err := config.GetManager().StoreConfigManager.NewMysqlDb()
		if err != nil {
			return nil, err
		}

		idlDaoManager := idl.NewMysqlIDL(mysqlDb)
		repositoryDaoManager := repository.NewMysqlRepository(mysqlDb)
		templateDaoManager := template.NewMysqlTemplate(mysqlDb)

		return &Manager{
			Idl:        idlDaoManager,
			Repository: repositoryDaoManager,
			Template:   templateDaoManager,
		}, nil

	case consts.StoreTypeMongo:
		panic("to be implemented")

	case consts.StoreTypeRedis:
		panic("to be implemented")

	default:
		return nil, fmt.Errorf("invalid store type")
	}
}
