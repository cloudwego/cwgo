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
	"context"
	"time"

	"github.com/cloudwego/cwgo/platform/server/shared/consts"
	"github.com/cloudwego/cwgo/platform/server/shared/dao/entity"
	"github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/model"
	"gorm.io/gorm"
)

type ITemplateDaoManager interface {
	AddTemplate(ctx context.Context, templateModel model.Template) error

	DeleteTemplate(ctx context.Context, ids []int64) error

	UpdateTemplate(ctx context.Context, templateModel model.Template) error

	GetTemplateList(ctx context.Context, page, limit, order int32, orderBy string) ([]*model.Template, error)

	AddTemplateItem(ctx context.Context, templateItemModel model.TemplateItem) error

	DeleteTemplateItem(ctx context.Context, ids []int64) error

	UpdateTemplateItem(ctx context.Context, templateItemModel model.TemplateItem) error

	GetTemplateItemList(ctx context.Context, templateId int64, page, limit, order int32, orderBy string) ([]*model.TemplateItem, error)
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

func (m *MysqlTemplateManager) AddTemplate(ctx context.Context, templateModel model.Template) error {
	templateEntity := entity.MysqlTemplate{
		Name: templateModel.Name,
		Type: templateModel.Type,
	}

	err := m.db.WithContext(ctx).
		Create(&templateEntity).Error

	return err
}

func (m *MysqlTemplateManager) DeleteTemplate(ctx context.Context, ids []int64) error {
	var templateEntity entity.MysqlTemplate

	err := m.db.WithContext(ctx).
		Delete(&templateEntity, ids).Error

	return err
}

func (m *MysqlTemplateManager) UpdateTemplate(ctx context.Context, templateModel model.Template) error {
	templateEntity := entity.MysqlTemplate{
		ID:   templateModel.Id,
		Name: templateModel.Name,
		Type: templateModel.Type,
	}

	err := m.db.WithContext(ctx).
		Model(&templateEntity).Updates(templateEntity).Error

	return err
}

func (m *MysqlTemplateManager) GetTemplateList(ctx context.Context, page, limit, order int32, orderBy string) ([]*model.Template, error) {
	var templateEntities []*entity.MysqlTemplate

	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	// default sort field to 'update_time' if not provided
	if orderBy == "" {
		orderBy = consts.OrderByUpdateTime
	}

	switch order {
	case consts.OrderNumInc:
		orderBy = orderBy + " " + consts.OrderInc
	case consts.OrderNumDec:
		orderBy = orderBy + " " + consts.OrderDec
	}

	err := m.db.WithContext(ctx).
		Offset(int(offset)).
		Limit(int(limit)).
		Order(orderBy).
		Find(&templateEntities).Error
	if err != nil {
		return nil, err
	}

	templateModels := make([]*model.Template, len(templateEntities))

	for i, templateEntity := range templateEntities {
		templateModels[i] = &model.Template{
			Id:         templateEntity.ID,
			Name:       templateEntity.Name,
			Type:       templateEntity.Type,
			IsDeleted:  false,
			CreateTime: templateEntity.CreateTime.Format(time.DateTime),
			UpdateTime: templateEntity.UpdateTime.Format(time.DateTime),
		}
	}

	return templateModels, nil
}

func (m *MysqlTemplateManager) AddTemplateItem(ctx context.Context, templateItemModel model.TemplateItem) error {
	templateItemEntity := entity.MysqlTemplateItem{
		TemplateID: templateItemModel.TemplateId,
		Name:       templateItemModel.Name,
		Content:    templateItemModel.Content,
	}

	err := m.db.WithContext(ctx).
		Create(&templateItemEntity).Error

	return err
}

func (m *MysqlTemplateManager) DeleteTemplateItem(ctx context.Context, ids []int64) error {
	var templateItemEntity entity.MysqlTemplateItem

	err := m.db.WithContext(ctx).
		Delete(&templateItemEntity, ids).Error

	return err
}

func (m *MysqlTemplateManager) UpdateTemplateItem(ctx context.Context, templateItemModel model.TemplateItem) error {
	templateItemEntity := entity.MysqlTemplateItem{
		ID:      templateItemModel.TemplateId,
		Name:    templateItemModel.Name,
		Content: templateItemModel.Content,
	}

	err := m.db.WithContext(ctx).
		Model(&templateItemEntity).
		Updates(templateItemEntity).Error

	return err
}

func (m *MysqlTemplateManager) GetTemplateItemList(ctx context.Context, templateId int64, page, limit, order int32, orderBy string) ([]*model.TemplateItem, error) {
	var templateItemEntities []*entity.MysqlTemplateItem

	if page < 1 {
		page = 1
	}

	offset := (page - 1) * limit

	// default sort field to 'update_time' if not provided
	if orderBy == "" {
		orderBy = consts.OrderByUpdateTime
	}

	switch order {
	case consts.OrderNumInc:
		orderBy = orderBy + " " + consts.OrderInc
	case consts.OrderNumDec:
		orderBy = orderBy + " " + consts.OrderDec
	}

	err := m.db.WithContext(ctx).
		Where("`template_id` = ?", templateId).
		Offset(int(offset)).
		Limit(int(limit)).
		Order(orderBy).
		Find(&templateItemEntities).Error
	if err != nil {
		return nil, err
	}

	templateItemModels := make([]*model.TemplateItem, len(templateItemEntities))

	for i, templateEntity := range templateItemEntities {
		templateItemModels[i] = &model.TemplateItem{
			Id:         templateEntity.ID,
			TemplateId: templateEntity.TemplateID,
			Name:       templateEntity.Name,
			Content:    templateEntity.Content,
			IsDeleted:  false,
			CreateTime: templateEntity.CreateTime.Format(time.DateTime),
			UpdateTime: templateEntity.UpdateTime.Format(time.DateTime),
		}
	}

	return templateItemModels, nil
}
