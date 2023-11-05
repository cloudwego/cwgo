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

package entity

import (
	"gorm.io/plugin/soft_delete"
	"time"
)

const TableNameMysqlRepository = "repository"

// MysqlRepository mapped from table <repository>
type MysqlRepository struct {
	ID             int64                 `gorm:"column:id;primaryKey;comment:repository id" json:"id"`                                               // repository id
	RepositoryType int32                 `gorm:"column:repository_type;not null;comment:repo type" json:"repository_type"`                           // repo type
	StoreType      int32                 `gorm:"column:store_type;not null;comment:store type" json:"store_type"`                                    // store type
	RepositoryURL  string                `gorm:"column:repository_url;not null;comment:repository URL" json:"repository_url"`                        // repository URL
	LastUpdateTime time.Time             `gorm:"column:last_update_time;comment:last update time" json:"last_update_time"`                           // last update time
	LastSyncTime   time.Time             `gorm:"column:last_sync_time;comment:last sync time" json:"last_sync_time"`                                 // last sync time
	Token          string                `gorm:"column:token;comment:repository token" json:"token"`                                                 // repository token
	Status         int32                 `gorm:"column:status;default:1;comment:status" json:"status"`                                               // status
	IsDeleted      soft_delete.DeletedAt `gorm:"column:is_deleted;softDelete:flag;not null;comment:is deleted" json:"is_deleted"`                    // is deleted
	CreateTime     time.Time             `gorm:"column:create_time;autoCreateTime;default:CURRENT_TIMESTAMP;comment:create time" json:"create_time"` // create time
	UpdateTime     time.Time             `gorm:"column:update_time;autoUpdateTime;default:CURRENT_TIMESTAMP;comment:update time" json:"update_time"` // update time
}

// TableName MysqlRepository's table name
func (*MysqlRepository) TableName() string {
	return TableNameMysqlRepository
}
