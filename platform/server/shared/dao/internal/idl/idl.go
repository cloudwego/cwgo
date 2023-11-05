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
	"github.com/cloudwego/cwgo/platform/server/shared/consts"
	"github.com/cloudwego/cwgo/platform/server/shared/dao/entity"
	"github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/model"
	"gorm.io/gorm"
	"sync"
	"time"
)

type IIdlDaoManager interface {
	AddIDL(ctx context.Context, idlModel model.IDL) error

	DeleteIDLs(ctx context.Context, ids []int64) error

	UpdateIDL(ctx context.Context, idlModel model.IDL) error
	Sync(ctx context.Context, idlModel model.IDL) error

	GetIDL(ctx context.Context, id int64) (*model.IDL, error)
	GetIDLList(ctx context.Context, page, limit, order int32, orderBy string) ([]*model.IDL, error)
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

func (m *MysqlIDLManager) AddIDL(ctx context.Context, idlModel model.IDL) error {
	// check repo id is exists
	var repo entity.MysqlRepository

	err := m.db.WithContext(ctx).
		Take(&repo, idlModel.RepositoryId).Error
	if err != nil {
		return err
	}

	now := time.Now()

	mainIdlEntity := entity.MysqlIDL{
		RepositoryID: idlModel.RepositoryId,
		IdlType:      idlModel.IdlType,
		IdlPath:      idlModel.MainIdlPath,
		CommitHash:   idlModel.CommitHash,
		ServiceName:  idlModel.ServiceName,
		LastSyncTime: now,
	}

	// insert main idlModel
	err = m.db.WithContext(ctx).
		Create(&mainIdlEntity).Error
	if err != nil {
		return err
	}

	// insert import idls
	importedIdlEntities := make([]*entity.MysqlIDL, len(idlModel.ImportIdls))
	for i, importIdl := range idlModel.ImportIdls {
		importedIdlEntities[i] = &entity.MysqlIDL{
			RepositoryID: idlModel.RepositoryId,
			ParentIdlID:  mainIdlEntity.ID,
			IdlType:      idlModel.IdlType,
			IdlPath:      importIdl.IdlPath,
			CommitHash:   importIdl.CommitHash,
			ServiceName:  idlModel.ServiceName,
			LastSyncTime: now,
		}
	}
	err = m.db.WithContext(ctx).
		Create(&importedIdlEntities).Error
	if err != nil {
		return err
	}

	return nil
}

func (m *MysqlIDLManager) DeleteIDLs(ctx context.Context, ids []int64) error {
	var idl entity.MysqlIDL

	err := m.db.WithContext(ctx).
		Delete(&idl, ids).Error

	return err
}

func (m *MysqlIDLManager) UpdateIDL(ctx context.Context, idlModel model.IDL) error {
	// update main idlModel
	mainIdlEntity := entity.MysqlIDL{
		ID:           idlModel.Id,
		RepositoryID: idlModel.RepositoryId,
		ParentIdlID:  0,
		IdlType:      idlModel.IdlType,
		IdlPath:      idlModel.MainIdlPath,
		CommitHash:   idlModel.CommitHash,
		ServiceName:  idlModel.ServiceName,
	}

	err := m.db.WithContext(ctx).
		Model(&mainIdlEntity).Updates(mainIdlEntity).Error
	if err != nil {
		return err
	}

	// update import idls
	if idlModel.ImportIdls != nil {
		importedIdlEntities := make([]*entity.MysqlIDL, len(idlModel.ImportIdls))
		for i, importIdl := range idlModel.ImportIdls {
			importedIdlEntities[i] = &entity.MysqlIDL{
				RepositoryID: idlModel.RepositoryId,
				ParentIdlID:  mainIdlEntity.ID,
				IdlType:      idlModel.IdlType,
				IdlPath:      importIdl.IdlPath,
				CommitHash:   importIdl.CommitHash,
				ServiceName:  idlModel.ServiceName,
			}
		}

		err = m.db.WithContext(ctx).
			Where("`parent_idl_id` = ?", idlModel.Id).
			Delete(&idlModel).Error
		if err != nil {
			return err
		}

		err = m.db.Where(ctx).
			Create(importedIdlEntities).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *MysqlIDLManager) Sync(ctx context.Context, idlModel model.IDL) error {
	// update main idlModel
	mainIdlEntity := entity.MysqlIDL{
		ID:           idlModel.Id,
		RepositoryID: idlModel.RepositoryId,
		ParentIdlID:  0,
		IdlType:      idlModel.IdlType,
		IdlPath:      idlModel.MainIdlPath,
		CommitHash:   idlModel.CommitHash,
		ServiceName:  idlModel.ServiceName,
		LastSyncTime: time.Now(),
	}

	err := m.db.WithContext(ctx).
		Model(&mainIdlEntity).Updates(mainIdlEntity).Error
	if err != nil {
		return err
	}

	// update import idls
	if idlModel.ImportIdls != nil {
		importedIdlEntities := make([]*entity.MysqlIDL, len(idlModel.ImportIdls))
		for i, importIdl := range idlModel.ImportIdls {
			importedIdlEntities[i] = &entity.MysqlIDL{
				RepositoryID: idlModel.RepositoryId,
				ParentIdlID:  mainIdlEntity.ID,
				IdlType:      idlModel.IdlType,
				IdlPath:      importIdl.IdlPath,
				CommitHash:   importIdl.CommitHash,
				ServiceName:  idlModel.ServiceName,
			}
		}

		err = m.db.WithContext(ctx).
			Where("`parent_idl_id` = ?", idlModel.Id).
			Delete(&idlModel).Error
		if err != nil {
			return err
		}

		err = m.db.Where(ctx).
			Create(importedIdlEntities).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *MysqlIDLManager) GetIDL(ctx context.Context, id int64) (*model.IDL, error) {
	var mainIdlEntity entity.MysqlIDL

	err := m.db.WithContext(ctx).
		Where("`id` = ?", id).
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
		Id:           mainIdlEntity.ID,
		RepositoryId: mainIdlEntity.RepositoryID,
		IdlType:      mainIdlEntity.IdlType,
		MainIdlPath:  mainIdlEntity.IdlPath,
		CommitHash:   mainIdlEntity.CommitHash,
		ImportIdls:   importIdlModels,
		ServiceName:  mainIdlEntity.ServiceName,
		LastSyncTime: mainIdlEntity.LastSyncTime.Format(time.DateTime),
		IsDeleted:    false,
		CreateTime:   mainIdlEntity.CreateTime.Format(time.DateTime),
		UpdateTime:   mainIdlEntity.UpdateTime.Format(time.DateTime),
	}, nil
}

func (m *MysqlIDLManager) GetIDLList(ctx context.Context, page, limit, order int32, orderBy string) ([]*model.IDL, error) {
	var idlEntities []*entity.MysqlIDL

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
		Find(&idlEntities).Error
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	idlModels := make([]*model.IDL, len(idlEntities))
	for i, idl := range idlEntities {
		wg.Add(1)
		idlModels[i] = &model.IDL{
			Id:           idl.ID,
			RepositoryId: idl.RepositoryID,
			IdlType:      idl.IdlType,
			MainIdlPath:  idl.IdlPath,
			CommitHash:   idl.CommitHash,
			ImportIdls:   nil,
			ServiceName:  idl.ServiceName,
			LastSyncTime: idl.LastSyncTime.Format(time.DateTime),
			IsDeleted:    false,
			CreateTime:   idl.CreateTime.Format(time.DateTime),
			UpdateTime:   idl.UpdateTime.Format(time.DateTime),
		}
		go func(i int, idl *entity.MysqlIDL) {
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

	return idlModels, nil
}
