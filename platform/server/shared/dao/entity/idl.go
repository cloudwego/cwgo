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

const TableNameMysqlIDL = "idl"

// MysqlIDL mapped from table <idl>
type MysqlIDL struct {
	ID                  int64                 `gorm:"column:id;primaryKey;autoIncrement:true;comment:id" json:"id"`                                                // id
	IdlRepositoryID     int64                 `gorm:"column:idl_repository_id;not null;comment:repository id" json:"idl_repository_id"`                            // idl repository id
	ServiceRepositoryID int64                 `gorm:"column:service_repository_id;not null;comment:repository id" json:"service_repository_id"`                    // service repository id
	ParentIdlID         int64                 `gorm:"column:parent_idl_id;comment:null if main idl else import idl" json:"parent_idl_id"`                          // null if main idl else import idl
	IdlPath             string                `gorm:"column:idl_path;not null;comment:idl path" json:"idl_path"`                                                   // idl path
	CommitHash          string                `gorm:"column:commit_hash;not null;comment:idl file commit hash" json:"commit_hash"`                                 // idl file commit hash
	ServiceName         string                `gorm:"column:service_name;not null;comment:service name" json:"service_name"`                                       // service name
	LastSyncTime        time.Time             `gorm:"column:last_sync_time;not null;default:CURRENT_TIMESTAMP;comment:last update time" json:"last_sync_time"`     // last update time
	IsDeleted           soft_delete.DeletedAt `gorm:"column:is_deleted;softDelete:flag;not null;comment:is deleted" json:"is_deleted"`                             // is deleted
	CreateTime          time.Time             `gorm:"column:create_time;autoCreateTime;not null;default:CURRENT_TIMESTAMP;comment:create time" json:"create_time"` // create time
	UpdateTime          time.Time             `gorm:"column:update_time;autoUpdateTime;not null;default:CURRENT_TIMESTAMP;comment:update time" json:"update_time"` // update time
}

// TableName MysqlIDL's table name
func (*MysqlIDL) TableName() string {
	return TableNameMysqlIDL
}
