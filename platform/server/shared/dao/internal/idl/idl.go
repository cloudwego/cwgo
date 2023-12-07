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

package idl

import (
	"context"
	"fmt"
	"github.com/cloudwego/cwgo/platform/server/shared/consts"
	"github.com/cloudwego/cwgo/platform/server/shared/dao/entity"
	"github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"sync"
	"time"
)

type IIdlDaoManager interface {
	AddIDL(ctx context.Context, idlModel model.IDL) (int64, error)

	DeleteIDLs(ctx context.Context, ids []int64) error

	UpdateIDL(ctx context.Context, idlModel model.IDL) error
	Sync(ctx context.Context, idlModel model.IDL) error

	GetIDL(ctx context.Context, id int64) (*model.IDL, error)
	GetIDLList(ctx context.Context, idlModel model.IDL, page, limit, order int32, orderBy string) ([]*model.IDL, int64, error)
	CheckMainIdlIfExist(ctx context.Context, repositoryId int64, mainIdlPath string) (bool, error)
}

type MysqlIDLManager struct {
	db *gorm.DB
}

var _ IIdlDaoManager = (*MysqlIDLManager)(nil)

func NewMysqlIDL(db *gorm.DB) *MysqlIDLManager {
	return &MysqlIDLManager{
		db: db,
	}
}

func (m *MysqlIDLManager) AddIDL(ctx context.Context, idlModel model.IDL) (int64, error) {
	// check repo id is exists
	var repo entity.MysqlRepository
	var mainIdlEntity entity.MysqlIDL

	err := m.db.WithContext(ctx).Transaction(
		func(tx *gorm.DB) error {
			err := tx.Take(&repo, idlModel.IdlRepositoryId).Error
			if err != nil {
				return err
			}

			now := time.Now()

			mainIdlEntity = entity.MysqlIDL{
				IdlRepositoryID:     idlModel.IdlRepositoryId,
				ServiceRepositoryID: idlModel.ServiceRepositoryId,
				IdlPath:             idlModel.MainIdlPath,
				CommitHash:          idlModel.CommitHash,
				ServiceName:         idlModel.ServiceName,
				LastSyncTime:        now,
			}

			err = tx.Clauses(
				clause.OnConflict{
					Columns:   []clause.Column{{Name: "id"}},
					UpdateAll: true,
				},
			).Create(&mainIdlEntity).Error

			// insert import idls
			err = tx.
				Where("`parent_idl_id` = ?", mainIdlEntity.ID).
				Delete(&entity.MysqlIDL{}).Error
			if err != nil {
				return err
			}
			if len(idlModel.ImportIdls) == 0 {
				return nil
			}
			importedIdlEntities := make([]*entity.MysqlIDL, len(idlModel.ImportIdls))
			for i, importIdl := range idlModel.ImportIdls {
				importedIdlEntities[i] = &entity.MysqlIDL{
					IdlRepositoryID:     idlModel.IdlRepositoryId,
					ServiceRepositoryID: idlModel.ServiceRepositoryId,
					ParentIdlID:         mainIdlEntity.ID,
					IdlPath:             importIdl.IdlPath,
					CommitHash:          importIdl.CommitHash,
					ServiceName:         idlModel.ServiceName,
					LastSyncTime:        now,
				}
			}
			err = tx.WithContext(ctx).
				Create(&importedIdlEntities).Error
			if err != nil {
				return err
			}

			return nil
		},
	)
	if err != nil {
		return -1, err
	}

	return mainIdlEntity.ID, nil
}

func (m *MysqlIDLManager) DeleteIDLs(ctx context.Context, ids []int64) error {
	var idlEntity entity.MysqlIDL

	err := m.db.WithContext(ctx).Transaction(
		func(tx *gorm.DB) error {
			res := tx.Delete(&idlEntity, ids)
			if res.Error != nil {
				return res.Error
			}

			if res.RowsAffected == 0 {
				return consts.ErrRecordNotFound
			}

			err := tx.
				Where("`parent_idl_id` IN ?", ids).
				Delete(&entity.MysqlIDL{}).Error

			return err
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (m *MysqlIDLManager) UpdateIDL(ctx context.Context, idlModel model.IDL) error {
	var lastSyncTime time.Time
	if idlModel.LastSyncTime != "" {
		lastSyncTime, _ = time.Parse(time.DateTime, idlModel.LastSyncTime)
	} else {
		lastSyncTime = time.Now()
	}

	// update main idlModel
	mainIdlEntity := entity.MysqlIDL{
		ID:           idlModel.Id,
		ParentIdlID:  0,
		CommitHash:   idlModel.CommitHash,
		ServiceName:  idlModel.ServiceName,
		LastSyncTime: lastSyncTime,
	}

	err := m.db.WithContext(ctx).Transaction(
		func(tx *gorm.DB) error {
			// TODO: check
			err := tx.Model(&mainIdlEntity).Updates(mainIdlEntity).Error
			if err != nil {
				return err
			}

			// update import idls
			if idlModel.ImportIdls != nil {
				importedIdlEntities := make([]*entity.MysqlIDL, len(idlModel.ImportIdls))
				for i, importIdl := range idlModel.ImportIdls {
					importedIdlEntities[i] = &entity.MysqlIDL{
						IdlRepositoryID:     mainIdlEntity.IdlRepositoryID,
						ServiceRepositoryID: mainIdlEntity.ServiceRepositoryID,
						ParentIdlID:         mainIdlEntity.ID,
						IdlPath:             importIdl.IdlPath,
						CommitHash:          importIdl.CommitHash,
						ServiceName:         mainIdlEntity.ServiceName,
						LastSyncTime:        lastSyncTime,
					}
				}

				err = tx.
					Where("`parent_idl_id` = ?", idlModel.Id).
					Delete(&idlModel).Error
				if err != nil {
					return err
				}

				err = tx.
					Where(ctx).
					Create(importedIdlEntities).Error
				if err != nil {
					return err
				}
			}

			return nil
		},
	)

	return err
}

func (m *MysqlIDLManager) Sync(ctx context.Context, idlModel model.IDL) error {
	// update main idlModel
	mainIdlEntity := entity.MysqlIDL{
		ID:           idlModel.Id,
		ParentIdlID:  0,
		CommitHash:   idlModel.CommitHash,
		LastSyncTime: time.Now(),
	}

	err := m.db.WithContext(ctx).Transaction(
		func(tx *gorm.DB) error {
			err := tx.Model(&mainIdlEntity).Updates(&mainIdlEntity).Error
			if err != nil {
				return err
			}

			// update import idls
			if len(idlModel.ImportIdls) != 0 {
				importedIdlEntities := make([]*entity.MysqlIDL, len(idlModel.ImportIdls))
				for i, importIdl := range idlModel.ImportIdls {
					importedIdlEntities[i] = &entity.MysqlIDL{
						IdlRepositoryID:     mainIdlEntity.IdlRepositoryID,
						ServiceRepositoryID: mainIdlEntity.ServiceRepositoryID,
						ParentIdlID:         mainIdlEntity.ID,
						IdlPath:             importIdl.IdlPath,
						CommitHash:          importIdl.CommitHash,
						ServiceName:         mainIdlEntity.ServiceName,
						LastSyncTime:        time.Now(),
					}
				}

				err = tx.
					Where("`parent_idl_id` = ?", idlModel.Id).
					Delete(&entity.MysqlIDL{}).Error
				if err != nil {
					return err
				}

				err = tx.Create(importedIdlEntities).Error
				if err != nil {
					return err
				}
			}
			return nil
		},
	)

	return err
}

func (m *MysqlIDLManager) GetIDL(ctx context.Context, id int64) (*model.IDL, error) {
	var mainIdlEntity entity.MysqlIDL

	err := m.db.WithContext(ctx).
		Where("`id` = ? AND `parent_idl_id` = 0", id).
		Take(&mainIdlEntity).Error
	if err != nil {
		return nil, err
	}

	var importIdlEntities []*entity.MysqlIDL

	err = m.db.WithContext(ctx).
		Where("`parent_idl_id` = ?", id).
		Find(&importIdlEntities).Error
	if err != nil {
		return nil, err
	}

	importIdlModels := make([]*model.ImportIDL, len(importIdlEntities))
	for i, importIdl := range importIdlEntities {
		importIdlModels[i] = &model.ImportIDL{
			IdlPath:    importIdl.IdlPath,
			CommitHash: importIdl.CommitHash,
		}
	}

	return &model.IDL{
		Id:                  mainIdlEntity.ID,
		IdlRepositoryId:     mainIdlEntity.IdlRepositoryID,
		ServiceRepositoryId: mainIdlEntity.ServiceRepositoryID,
		MainIdlPath:         mainIdlEntity.IdlPath,
		CommitHash:          mainIdlEntity.CommitHash,
		ImportIdls:          importIdlModels,
		ServiceName:         mainIdlEntity.ServiceName,
		LastSyncTime:        mainIdlEntity.LastSyncTime.Format(time.DateTime),
		IsDeleted:           false,
		CreateTime:          mainIdlEntity.CreateTime.Format(time.DateTime),
		UpdateTime:          mainIdlEntity.UpdateTime.Format(time.DateTime),
	}, nil
}

func (m *MysqlIDLManager) GetIDLList(ctx context.Context, idlModel model.IDL, page, limit, order int32, orderBy string) ([]*model.IDL, int64, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	var total int64

	db := m.db.WithContext(ctx)

	if idlModel.ServiceName != "" {
		db = db.Where("`service_name` LIKE ?", fmt.Sprintf("%%%s%%", idlModel.ServiceName))
	}

	err := db.
		Model(&entity.MysqlIDL{}).
		Count(&total).Error
	if err != nil {
		return nil, -1, err
	}

	if int64(offset) >= total {
		return nil, total, nil
	}

	var idlEntities []*entity.MysqlIDL

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
		Where("`parent_idl_id` = 0").
		Offset(int(offset)).
		Limit(int(limit)).
		Order(orderBy).
		Find(&idlEntities).Error
	if err != nil {
		return nil, -1, err
	}

	var wg sync.WaitGroup
	idlModels := make([]*model.IDL, len(idlEntities))
	for i, idl := range idlEntities {
		wg.Add(1)
		idlModels[i] = &model.IDL{
			Id:                  idl.ID,
			IdlRepositoryId:     idl.IdlRepositoryID,
			ServiceRepositoryId: idl.ServiceRepositoryID,
			MainIdlPath:         idl.IdlPath,
			CommitHash:          idl.CommitHash,
			ImportIdls:          nil,
			ServiceName:         idl.ServiceName,
			LastSyncTime:        idl.LastSyncTime.Format(time.DateTime),
			IsDeleted:           false,
			CreateTime:          idl.CreateTime.Format(time.DateTime),
			UpdateTime:          idl.UpdateTime.Format(time.DateTime),
		}
		go func(i int, idl *entity.MysqlIDL) {
			defer wg.Done()

			var importIdlEntities []*entity.MysqlIDL
			err := m.db.WithContext(ctx).
				Where("`parent_idl_id` = ?", idl.ID).
				Find(&importIdlEntities).Error
			if err != nil {
				return
			}

			importIdlModels := make([]*model.ImportIDL, len(importIdlEntities))
			for j, importIdl := range importIdlEntities {
				importIdlModels[j] = &model.ImportIDL{
					IdlPath:    importIdl.IdlPath,
					CommitHash: importIdl.CommitHash,
				}
			}

			idlModels[i].ImportIdls = importIdlModels
		}(i, idl)
	}
	wg.Wait()

	return idlModels, total, nil
}

func (m *MysqlIDLManager) CheckMainIdlIfExist(ctx context.Context, repositoryId int64, mainIdlPath string) (bool, error) {
	var idlEntity entity.MysqlIDL

	err := m.db.WithContext(ctx).
		Where("`idl_repository_id` = ? AND `parent_idl_id` = 0 AND `idl_path` = ?", repositoryId, mainIdlPath).
		Take(&idlEntity).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}

		return false, err
	}

	return true, nil
}
