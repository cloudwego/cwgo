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

const TableNameMysqlToken = "token"

// MysqlToken mapped from table <token>
type MysqlToken struct {
	ID               int64                 `gorm:"column:id;primaryKey;autoIncrement;comment:id" json:"id"`                                                     // id
	Owner            string                `gorm:"column:owner;not null;comment:repository owner" json:"owner"`                                                 // repository owner
	RepositoryType   int32                 `gorm:"column:repository_type;not null;comment:repository type (1: gitlab, 2: github)" json:"repository_type"`       // repository type (1: gitlab, 2: github)
	RepositoryDomain string                `gorm:"column:repository_domain;not null;comment:repository api domain" json:"repository_domain"`                    // repository api domain
	Token            string                `gorm:"column:token;not null;comment:repository token" json:"token"`                                                 // repository token
	Status           int32                 `gorm:"column:status;not null;comment:token status (0: expired, 1: valid)" json:"status"`                            // token status (0: expired, 1: valid)
	ExpirationTime   time.Time             `gorm:"column:expiration_time;comment:token expiration time" json:"expiration_time"`                                 // token expiration time
	IsDeleted        soft_delete.DeletedAt `gorm:"column:is_deleted;softDelete:flag;not null;comment:is deleted" json:"is_deleted"`                             // is deleted
	CreateTime       time.Time             `gorm:"column:create_time;autoCreateTime;not null;default:CURRENT_TIMESTAMP;comment:create time" json:"create_time"` // create time
	UpdateTime       time.Time             `gorm:"column:update_time;autoUpdateTime;not null;default:CURRENT_TIMESTAMP;comment:update time" json:"update_time"` // update time
}

// TableName MysqlToken's table name
func (*MysqlToken) TableName() string {
	return TableNameMysqlToken
}
