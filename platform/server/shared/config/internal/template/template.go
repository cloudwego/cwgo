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

package template

import (
	"github.com/cloudwego/cwgo/platform/server/shared/config/internal/store"
	"github.com/cloudwego/cwgo/platform/server/shared/utils"
	"gorm.io/gorm"
)

type ITemplateManager interface {
	AddTemplate(name string, _type int8) error
	DeleteTemplate(ids []int64) error
	UpdateTemplate(id int64, name string) error
	GetTemplates(page, limit int32, sortBy string) ([]Template, error)

	AddTemplateItem(templateId int64, name, content string) error
	DeleteTemplateItem(ids []int64) error
	UpdateTemplateItem(id int64, name, content string) error
	GetTemplateItems(page, limit int32, sortBy string) ([]TemplateItem, error)
}

type MysqlTemplateManager struct {
	Db *gorm.DB
}

type Template struct {
	Id         int64
	Name       string
	Type       int8
	CreateTime string
	UpdateTime string
}

type TemplateItem struct {
	Id         int64
	TemplateId int64
	Name       string
	Content    string
	CreateTime string
	UpdateTime string
}

var _ ITemplateManager = (*MysqlTemplateManager)(nil)

func NewMysqlTemplate() (*MysqlTemplateManager, error) {
	db, err := store.GetMysqlDB()
	if err != nil {
		return nil, err
	}

	return &MysqlTemplateManager{
		Db: db,
	}, nil
}

func (r *MysqlTemplateManager) AddTemplate(name string, _type int8) error {
	timeNow := utils.GetCurrentTime()
	template := Template{
		Name:       name,
		Type:       _type,
		CreateTime: timeNow,
		UpdateTime: timeNow,
	}
	res := r.Db.Create(&template)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (r *MysqlTemplateManager) DeleteTemplate(ids []int64) error {
	var template Template
	res := r.Db.Delete(&template, ids)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (r *MysqlTemplateManager) UpdateTemplate(id int64, name string) error {
	timeNow := utils.GetCurrentTime()
	res := r.Db.Where("id = ?", id).Updates(Template{
		Name:       name,
		UpdateTime: timeNow,
	})
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (r *MysqlTemplateManager) GetTemplates(page, limit int32, sortBy string) ([]Template, error) {
	var templates []Template
	offset := (page - 1) * limit

	// Default sort field to 'update_time' if not provided
	if sortBy == "" {
		sortBy = SortByUpdateTime
	}

	res := r.Db.Offset(int(offset)).Limit(int(limit)).Order(sortBy).Find(&templates)
	if res.Error != nil {
		return nil, res.Error
	}

	return templates, nil
}

func (r *MysqlTemplateManager) AddTemplateItem(templateId int64, name, content string) error {
	timeNow := utils.GetCurrentTime()
	templateItem := TemplateItem{
		TemplateId: templateId,
		Name:       name,
		Content:    content,
		CreateTime: timeNow,
		UpdateTime: timeNow,
	}
	res := r.Db.Create(&templateItem)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (r *MysqlTemplateManager) DeleteTemplateItem(ids []int64) error {
	var template Template
	res := r.Db.Delete(&template, ids)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (r *MysqlTemplateManager) UpdateTemplateItem(id int64, name, content string) error {
	timeNow := utils.GetCurrentTime()
	res := r.Db.Where("id = ?", id).Updates(TemplateItem{
		Name:       name,
		Content:    content,
		UpdateTime: timeNow,
	})
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (r *MysqlTemplateManager) GetTemplateItems(page, limit int32, sortBy string) ([]TemplateItem, error) {
	var templateItems []TemplateItem
	offset := (page - 1) * limit

	// Default sort field to 'update_time' if not provided
	if sortBy == "" {
		sortBy = SortByUpdateTime
	}

	res := r.Db.Offset(int(offset)).Limit(int(limit)).Order(sortBy).Find(&templateItems)
	if res.Error != nil {
		return nil, res.Error
	}

	return templateItems, nil
}
