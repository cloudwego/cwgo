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
	"time"

	"gorm.io/plugin/soft_delete"
)

const TableNameMysqlTemplateItem = "template_item"

// MysqlTemplateItem mapped from table <template_item>
type MysqlTemplateItem struct {
	ID         int64                 `gorm:"column:id;primaryKey;autoIncrement:true;comment:template item id" json:"id"`                                  // template item id
	TemplateID int64                 `gorm:"column:template_id;not null;comment:template id" json:"template_id"`                                          // template id
	Name       string                `gorm:"column:name;not null;comment:template item name" json:"name"`                                                 // template item name
	Content    string                `gorm:"column:content;comment:template content" json:"content"`                                                      // template content
	IsDeleted  soft_delete.DeletedAt `gorm:"column:is_deleted;softDelete:flag;not null;comment:is deleted" json:"is_deleted"`                             // is deleted
	CreateTime time.Time             `gorm:"column:create_time;autoCreateTime;not null;default:CURRENT_TIMESTAMP;comment:create time" json:"create_time"` // create time
	UpdateTime time.Time             `gorm:"column:update_time;autoUpdateTime;not null;default:CURRENT_TIMESTAMP;comment:update time" json:"update_time"` // update time
}

// TableName MysqlTemplateItem's table name
func (*MysqlTemplateItem) TableName() string {
	return TableNameMysqlTemplateItem
}
