/*
 * Copyright 2022 CloudWeGo Authors
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
 */

package config

import "gorm.io/gorm"

type IRepository interface {
	GetTokenByID(id int64) string
}

type MysqlRepository struct {
	db *gorm.DB
}

type Repository struct {
	id             int64
	repositoryUrl  string
	lastUpdateTime string
	lastSyncTime   string
	token          string
	status         string
}

func (sql *MysqlRepository) GetTokenByID(id int64) string {
	var repo Repository
	sql.db.Model(&repo).Where("id = ?", id).First(&repo)
	return repo.token
}
