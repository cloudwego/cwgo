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

package store

import (
	"github.com/cloudwego/cwgo/platform/server/shared/consts"
	"gorm.io/gorm"
)

type StoreConfigManager struct {
	config    Config
	storeType consts.StoreType
}

func NewStoreConfigManager(config Config) *StoreConfigManager {
	return &StoreConfigManager{
		config: config,
	}
}

func (cm *StoreConfigManager) GetStoreType() consts.StoreType {
	return cm.storeType
}

func (cm *StoreConfigManager) NewMysqlDb() (*gorm.DB, error) {
	return initMysqlDB(cm.config.Mysql.GetDsn())
}
