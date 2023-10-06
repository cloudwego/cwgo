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
	"github.com/cloudwego/cwgo/platform/server/shared/consts"
	"github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/template"
	"github.com/cloudwego/cwgo/platform/server/shared/utils"
	"gorm.io/gorm"
)

type ITemplateDaoManager interface {
	AddTemplate(name string, _type int32) error
	DeleteTemplate(ids []int64) error
	UpdateTemplate(id int64, name string) error
	GetTemplates(page, limit, order int32, orderBy string) ([]*template.Template, error)

	AddTemplateItem(templateId int64, name, content string) error
	DeleteTemplateItem(ids []int64) error
	UpdateTemplateItem(id int64, name, content string) error
	GetTemplateItems(page, limit, order int32, orderBy string) ([]*template.TemplateItem, error)
}

type MysqlTemplateManager struct {
	db *gorm.DB
}

var _ ITemplateDaoManager = (*MysqlTemplateManager)(nil)

func NewMysqlTemplate(db *gorm.DB) *MysqlTemplateManager {
	return &MysqlTemplateManager{
		db: db,
	}
}

func (r *MysqlTemplateManager) AddTemplate(name string, _type int32) error {
	timeNow := utils.GetCurrentTime()
	template := template.Template{
		Name:       name,
		Type:       _type,
		CreateTime: timeNow,
		UpdateTime: timeNow,
	}
	res := r.db.Table(consts.TableNameTemplate).Create(&template)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (r *MysqlTemplateManager) DeleteTemplate(ids []int64) error {
	var template template.Template
	res := r.db.Table(consts.TableNameTemplate).Delete(&template, ids)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (r *MysqlTemplateManager) UpdateTemplate(id int64, name string) error {
	timeNow := utils.GetCurrentTime()
	res := r.db.Table(consts.TableNameTemplate).Where("id = ?", id).Updates(template.Template{
		Name:       name,
		UpdateTime: timeNow,
	})
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (r *MysqlTemplateManager) GetTemplates(page, limit, order int32, orderBy string) ([]*template.Template, error) {
	var templates []*template.Template
	offset := (page - 1) * limit

	// Default sort field to 'update_time' if not provided
	if orderBy == "" {
		orderBy = consts.OrderByUpdateTime
	}

	switch order {
	case consts.OrderNumInc:
		orderBy = orderBy + " " + consts.Inc
	case consts.OrderNumDec:
		orderBy = orderBy + " " + consts.Dec
	default:
		orderBy = orderBy + " " + consts.Inc
	}

	res := r.db.Table(consts.TableNameTemplate).Offset(int(offset)).Limit(int(limit)).Order(orderBy).Find(&templates)
	if res.Error != nil {
		return nil, res.Error
	}

	return templates, nil
}

func (r *MysqlTemplateManager) AddTemplateItem(templateId int64, name, content string) error {
	timeNow := utils.GetCurrentTime()
	templateItem := template.TemplateItem{
		TemplateId: templateId,
		Name:       name,
		Content:    content,
		CreateTime: timeNow,
		UpdateTime: timeNow,
	}
	res := r.db.Table(consts.TableNameTemplateItem).Create(&templateItem)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (r *MysqlTemplateManager) DeleteTemplateItem(ids []int64) error {
	var template template.TemplateItem
	res := r.db.Table(consts.TableNameTemplateItem).Delete(&template, ids)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (r *MysqlTemplateManager) UpdateTemplateItem(id int64, name, content string) error {
	timeNow := utils.GetCurrentTime()
	res := r.db.Table(consts.TableNameTemplateItem).Where("id = ?", id).Updates(template.TemplateItem{
		Name:       name,
		Content:    content,
		UpdateTime: timeNow,
	})
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (r *MysqlTemplateManager) GetTemplateItems(page, limit, order int32, orderBy string) ([]*template.TemplateItem, error) {
	var templateItems []*template.TemplateItem
	offset := (page - 1) * limit

	// Default sort field to 'update_time' if not provided
	if orderBy == "" {
		orderBy = consts.OrderByUpdateTime
	}

	switch order {
	case consts.OrderNumInc:
		orderBy = orderBy + " " + consts.Inc
	case consts.OrderNumDec:
		orderBy = orderBy + " " + consts.Dec
	default:
		orderBy = orderBy + " " + consts.Inc
	}

	res := r.db.Table(consts.TableNameTemplateItem).Offset(int(offset)).Limit(int(limit)).Order(orderBy).Find(&templateItems)
	if res.Error != nil {
		return nil, res.Error
	}

	return templateItems, nil
}
