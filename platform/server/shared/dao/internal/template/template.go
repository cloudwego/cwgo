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

package template

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudwego/cwgo/platform/server/shared/consts"
	"github.com/cloudwego/cwgo/platform/server/shared/dao/entity"
	"github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/model"
	"gorm.io/gorm"
)

type ITemplateDaoManager interface {
	AddTemplate(ctx context.Context, templateModel model.Template) (int64, error)

	DeleteTemplate(ctx context.Context, ids []int64) error

	UpdateTemplate(ctx context.Context, templateModel model.Template) error

	AddTemplateItem(ctx context.Context, templateItemModel model.TemplateItem) (int64, error)

	DeleteTemplateItem(ctx context.Context, ids []int64) error

	UpdateTemplateItem(ctx context.Context, templateItemModel model.TemplateItem) error

	GetTemplate(ctx context.Context, id int64) (*model.TemplateWithInfo, error)
	GetTemplateList(ctx context.Context, templateModel model.Template, page, limit, order int32, orderBy string) ([]*model.TemplateWithInfo, int64, error)

	CheckTemplateIfExist(ctx context.Context, templateId int64) (bool, error)
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

func (m *MysqlTemplateManager) AddTemplate(ctx context.Context, templateModel model.Template) (int64, error) {
	templateEntity := entity.MysqlTemplate{
		Name: templateModel.Name,
		Type: templateModel.Type,
	}

	err := m.db.WithContext(ctx).
		Create(&templateEntity).Error

	return templateEntity.ID, err
}

func (m *MysqlTemplateManager) DeleteTemplate(ctx context.Context, ids []int64) error {
	var templateEntity entity.MysqlTemplate

	err := m.db.WithContext(ctx).Transaction(
		func(tx *gorm.DB) error {
			// update idl that used this template
			err := tx.
				Table(entity.TableNameMysqlIDL).
				Where("`template_id` IN ?", ids).
				UpdateColumn("template_id", 0).Error
			if err != nil {
				return err
			}

			// delete template
			res := tx.Delete(&templateEntity, ids)
			if res.Error != nil {
				return res.Error
			}
			if res.RowsAffected == 0 {
				return consts.ErrDatabaseRecordNotFound
			}

			// delete template item
			err = tx.Where("`template_id` IN ?", ids).
				Delete(&entity.MysqlTemplateItem{}).Error

			return err
		},
	)

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
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return consts.ErrDatabaseRecordNotFound
		}
		return err
	}

	return err
}

func (m *MysqlTemplateManager) AddTemplateItem(ctx context.Context, templateItemModel model.TemplateItem) (int64, error) {
	templateItemEntity := entity.MysqlTemplateItem{
		TemplateID: templateItemModel.TemplateId,
		Name:       templateItemModel.Name,
		Content:    templateItemModel.Content,
	}

	err := m.db.WithContext(ctx).
		Create(&templateItemEntity).Error

	return templateItemEntity.ID, err
}

func (m *MysqlTemplateManager) DeleteTemplateItem(ctx context.Context, ids []int64) error {
	var templateItemEntity entity.MysqlTemplateItem

	res := m.db.WithContext(ctx).
		Delete(&templateItemEntity, ids)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return consts.ErrDatabaseRecordNotFound
	}

	return nil
}

func (m *MysqlTemplateManager) UpdateTemplateItem(ctx context.Context, templateItemModel model.TemplateItem) error {
	templateItemEntity := entity.MysqlTemplateItem{
		ID:      templateItemModel.Id,
		Name:    templateItemModel.Name,
		Content: templateItemModel.Content,
	}

	err := m.db.WithContext(ctx).
		Model(&templateItemEntity).
		Updates(templateItemEntity).Error

	return err
}

func (m *MysqlTemplateManager) GetTemplate(ctx context.Context, templateId int64) (*model.TemplateWithInfo, error) {
	var templateWithInfoEntity entity.MysqlTemplateWithInfo

	err := m.db.WithContext(ctx).
		Preload("Items").
		Where("`template`.`id` = ?", templateId).
		Take(&templateWithInfoEntity).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, consts.ErrDatabaseRecordNotFound
		}
		return nil, err
	}

	itemModels := make([]*model.TemplateItem, len(templateWithInfoEntity.Items))

	for i, templateItemEntity := range templateWithInfoEntity.Items {
		itemModels[i] = &model.TemplateItem{
			Id:         templateItemEntity.ID,
			TemplateId: templateItemEntity.TemplateID,
			Name:       templateItemEntity.Name,
			Content:    templateItemEntity.Content,
			IsDeleted:  false,
			CreateTime: templateItemEntity.CreateTime.Format(time.DateTime),
			UpdateTime: templateItemEntity.UpdateTime.Format(time.DateTime),
		}
	}

	templateModel := &model.TemplateWithInfo{
		Template: &model.Template{
			Id:         templateWithInfoEntity.ID,
			Name:       templateWithInfoEntity.Name,
			Type:       templateWithInfoEntity.Type,
			IsDeleted:  false,
			CreateTime: templateWithInfoEntity.CreateTime.Format(time.DateTime),
			UpdateTime: templateWithInfoEntity.UpdateTime.Format(time.DateTime),
		},
		Items: itemModels,
	}

	return templateModel, nil
}

func (m *MysqlTemplateManager) GetTemplateList(ctx context.Context, templateModel model.Template, page, limit, order int32, orderBy string) ([]*model.TemplateWithInfo, int64, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	var total int64

	db := m.db.WithContext(ctx)

	if templateModel.Name != "" {
		db = db.Where("`template`.`name` LIKE ?", fmt.Sprintf("%%%s%%", templateModel.Name))
	}

	err := db.
		Model(&entity.MysqlTemplate{}).
		Count(&total).Error
	if err != nil {
		return nil, -1, err
	}

	if int64(offset) >= total {
		return nil, total, nil
	}

	var templateEntities []*entity.MysqlTemplateWithInfo

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

	err = db.
		Preload("Items").
		Offset(int(offset)).
		Limit(int(limit)).
		Order(orderBy).
		Find(&templateEntities).Error
	if err != nil {
		return nil, -1, err
	}

	templateModels := make([]*model.TemplateWithInfo, len(templateEntities))

	for i, templateEntity := range templateEntities {
		itemModels := make([]*model.TemplateItem, len(templateEntity.Items))

		for j, templateItemEntity := range templateEntity.Items {
			itemModels[j] = &model.TemplateItem{
				Id:         templateItemEntity.ID,
				TemplateId: templateItemEntity.TemplateID,
				Name:       templateItemEntity.Name,
				Content:    templateItemEntity.Content,
				IsDeleted:  false,
				CreateTime: templateItemEntity.CreateTime.Format(time.DateTime),
				UpdateTime: templateItemEntity.UpdateTime.Format(time.DateTime),
			}
		}

		templateModels[i] = &model.TemplateWithInfo{
			Template: &model.Template{
				Id:         templateEntity.ID,
				Name:       templateEntity.Name,
				Type:       templateEntity.Type,
				IsDeleted:  false,
				CreateTime: templateEntity.CreateTime.Format(time.DateTime),
				UpdateTime: templateEntity.UpdateTime.Format(time.DateTime),
			},
			Items: itemModels,
		}
	}

	return templateModels, total, nil
}

func (m *MysqlTemplateManager) CheckTemplateIfExist(ctx context.Context, templateId int64) (bool, error) {
	var templateEntity entity.MysqlTemplate

	err := m.db.WithContext(ctx).
		Where("`template_id` = ?", templateId).
		Take(&templateEntity).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}

		return false, err
	}

	return true, nil
}
