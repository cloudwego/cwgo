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

package entity

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

const TableNameMysqlRepository = "repository"

// MysqlRepository mapped from table <repository>
type MysqlRepository struct {
	ID             int64                 `gorm:"column:id;primaryKey;autoIncrement:true;comment:repository id" json:"id"`  // repository id
	RepositoryType int32                 `gorm:"column:repository_type;not null;comment:repo type" json:"repository_type"` // repo type
	Domain         string                `gorm:"column:domain;not null;comment:repo domain" json:"domain"`                 // repo domain
	Owner          string                `gorm:"column:owner;not null;comment:repo owner" json:"owner"`                    // repo owner
	RepositoryName string                `gorm:"column:repository_name;not null;comment:repo name" json:"repository_name"` // repo name
	Branch         string                `gorm:"column:branch;not null;comment:repo branch" json:"branch"`                 // repo branch
	StoreType      int32                 `gorm:"column:store_type;not null;comment:store type" json:"store_type"`          // store type
	LastUpdateTime time.Time             `gorm:"column:last_update_time;comment:last update time" json:"last_update_time"` // last update time
	LastSyncTime   time.Time             `gorm:"column:last_sync_time;comment:last sync time" json:"last_sync_time"`       // last sync time
	TokenId        int64                 `gorm:"column:token_id;comment:repository token id" json:"token_id"`              // repository token id
	Status         int32                 `gorm:"column:status;default:1;comment:status" json:"status"`
	IsDeleted      soft_delete.DeletedAt `gorm:"column:is_deleted;softDelete:flag;not null;comment:is deleted" json:"is_deleted"`                    // is deleted
	CreateTime     time.Time             `gorm:"column:create_time;autoCreateTime;default:CURRENT_TIMESTAMP;comment:create time" json:"create_time"` // create time
	UpdateTime     time.Time             `gorm:"column:update_time;autoUpdateTime;default:CURRENT_TIMESTAMP;comment:update time" json:"update_time"` // update time
}

// TableName MysqlRepository's table name
func (*MysqlRepository) TableName() string {
	return TableNameMysqlRepository
}
