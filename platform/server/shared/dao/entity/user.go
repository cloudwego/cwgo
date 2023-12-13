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

const TableNameMysqlUser = "user"

// MysqlUser mapped from table <user>
type MysqlUser struct {
	UserID     int64                 `gorm:"column:user_id;primaryKey;comment:user id" json:"user_id"`                                                    // user id
	Name       string                `gorm:"column:name;not null;comment:user name" json:"name"`                                                          // user name
	IsDeleted  soft_delete.DeletedAt `gorm:"column:is_deleted;not null;softDelete:flag;comment:is deleted" json:"is_deleted"`                             // is deleted
	CreateTime time.Time             `gorm:"column:create_time;autoCreateTime;not null;default:CURRENT_TIMESTAMP;comment:create time" json:"create_time"` // create time
	UpdateTime time.Time             `gorm:"column:update_time;autoUpdateTime;not null;default:CURRENT_TIMESTAMP;comment:update time" json:"update_time"` // update time
}

// TableName MysqlUser's table name
func (*MysqlUser) TableName() string {
	return TableNameMysqlUser
}
