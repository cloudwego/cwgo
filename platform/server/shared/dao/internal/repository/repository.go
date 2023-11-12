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

package repository

import (
	"context"
	"errors"
	"github.com/cloudwego/cwgo/platform/server/shared/consts"
	"github.com/cloudwego/cwgo/platform/server/shared/dao/entity"
	"github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type IRepositoryDaoManager interface {
	AddRepository(ctx context.Context, repoModel model.Repository) (int64, error)

	DeleteRepository(ctx context.Context, ids []int64) error

	UpdateRepository(ctx context.Context, repoModel model.Repository) error
	Sync(ctx context.Context, repoModel model.Repository) error
	ChangeRepositoryStatus(ctx context.Context, id int64, status int32) error

	GetRepository(ctx context.Context, id int64) (*model.Repository, error)
	GetRepositoryList(ctx context.Context, page, limit, order int32, orderBy string) ([]*model.Repository, error)
	GetAllRepositories(ctx context.Context) ([]*model.Repository, error)
	GetTokenByID(ctx context.Context, id int64) (string, error)
	GetRepoTypeByID(ctx context.Context, id int64) (int32, error)
}

type MysqlRepositoryManager struct {
	db *gorm.DB
}

var _ IRepositoryDaoManager = (*MysqlRepositoryManager)(nil)

func NewMysqlRepository(db *gorm.DB) *MysqlRepositoryManager {
	return &MysqlRepositoryManager{
		db: db,
	}
}

func (m *MysqlRepositoryManager) AddRepository(ctx context.Context, repoModel model.Repository) (int64, error) {
	var lastUpdateTime time.Time
	if repoModel.LastUpdateTime != "" {
		lastUpdateTime, _ = time.Parse(time.DateTime, repoModel.LastUpdateTime)
	} else {
		lastUpdateTime = time.Now()
	}

	var repoEntity entity.MysqlRepository

	// check if repo exists
	err := m.db.WithContext(ctx).
		Where("`repository_url` = ? AND `store_type` = ? AND `is_deleted` = 0",
			repoModel.RepositoryUrl,
			repoModel.StoreType,
		).
		Take(&repoEntity).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return -1, err
		}
	} else {
		// repo exists
		return repoEntity.ID, consts.ErrDuplicateRecord
	}

	// create repo if record is not exist or record's `is_deleted` = 1
	repoEntity = entity.MysqlRepository{
		RepositoryType: repoModel.RepositoryType,
		StoreType:      repoModel.StoreType,
		RepositoryURL:  repoModel.RepositoryUrl,
		LastUpdateTime: lastUpdateTime,
		LastSyncTime:   time.Now(),
		Token:          repoModel.Token,
		Status:         consts.RepositoryStatusNumActive,
		IsDeleted:      0,
	}
	err = m.db.WithContext(ctx).
		Clauses(
			clause.OnConflict{
				Columns:   []clause.Column{{Name: "repository_url"}, {Name: "store_type"}},
				UpdateAll: true,
			},
		).
		Create(&repoEntity).Error

	return repoEntity.ID, err
}

func (m *MysqlRepositoryManager) DeleteRepository(ctx context.Context, ids []int64) error {
	var repoEntity entity.MysqlRepository

	res := m.db.WithContext(ctx).Debug().
		Delete(&repoEntity, ids)
	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return consts.ErrRecordNotFound
	}

	return nil
}

func (m *MysqlRepositoryManager) UpdateRepository(ctx context.Context, repoModel model.Repository) error {
	if repoModel.Status != 0 {
		if _, ok := consts.RepositoryStatusNumMap[int(repoModel.Status)]; !ok {
			return errors.New("invalid status")
		}
	}

	repoEntity := entity.MysqlRepository{
		ID:     repoModel.Id,
		Token:  repoModel.Token,
		Status: repoModel.Status,
	}

	err := m.db.WithContext(ctx).
		Model(&repoEntity).
		Updates(repoEntity).Error

	return err
}

func (m *MysqlRepositoryManager) Sync(ctx context.Context, repoModel model.Repository) error {
	var lastUpdateTime, lastSyncTime time.Time
	if repoModel.LastUpdateTime != "" {
		lastUpdateTime, _ = time.Parse(time.DateTime, repoModel.LastUpdateTime)
	}
	if repoModel.LastSyncTime != "" {
		lastSyncTime, _ = time.Parse(time.DateTime, repoModel.LastSyncTime)
	}

	repoEntity := entity.MysqlRepository{
		ID:             repoModel.Id,
		LastUpdateTime: lastUpdateTime,
		LastSyncTime:   lastSyncTime,
	}

	err := m.db.WithContext(ctx).
		Model(&repoEntity).Updates(repoEntity).Error
	if err != nil {
		return err
	}

	return nil
}

func (m *MysqlRepositoryManager) ChangeRepositoryStatus(ctx context.Context, id int64, status int32) error {
	if status != 0 {
		if _, ok := consts.RepositoryStatusNumMap[int(status)]; !ok {
			return errors.New("invalid status")
		}
	}

	repoEntity := entity.MysqlRepository{
		ID:     id,
		Status: status,
	}

	err := m.db.WithContext(ctx).
		Model(&repoEntity).
		Updates(repoEntity).Error

	return err
}

func (m *MysqlRepositoryManager) GetRepository(ctx context.Context, id int64) (*model.Repository, error) {
	var repoEntity entity.MysqlRepository

	err := m.db.WithContext(ctx).
		Where("`id` = ?", id).
		Take(&repoEntity).Error
	if err != nil {
		return nil, err
	}

	return &model.Repository{
		Id:             repoEntity.ID,
		RepositoryType: repoEntity.RepositoryType,
		StoreType:      repoEntity.StoreType,
		RepositoryUrl:  repoEntity.RepositoryURL,
		Token:          repoEntity.Token,
		Status:         repoEntity.Status,
		LastUpdateTime: repoEntity.UpdateTime.Format(time.DateTime),
		LastSyncTime:   repoEntity.LastSyncTime.Format(time.DateTime),
		IsDeleted:      false,
		CreateTime:     repoEntity.CreateTime.Format(time.DateTime),
		UpdateTime:     repoEntity.UpdateTime.Format(time.DateTime),
	}, nil
}

func (m *MysqlRepositoryManager) GetRepositoryList(ctx context.Context, page, limit, order int32, orderBy string) ([]*model.Repository, error) {
	var repoEntities []*entity.MysqlRepository

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
		Find(&repoEntities).Error
	if err != nil {
		return nil, err
	}

	repoModels := make([]*model.Repository, len(repoEntities))

	for i, repoEntity := range repoEntities {
		repoModels[i] = &model.Repository{
			Id:             repoEntity.ID,
			RepositoryType: repoEntity.RepositoryType,
			StoreType:      repoEntity.StoreType,
			RepositoryUrl:  repoEntity.RepositoryURL,
			Token:          repoEntity.Token,
			Status:         repoEntity.Status,
			LastUpdateTime: repoEntity.UpdateTime.Format(time.DateTime),
			LastSyncTime:   repoEntity.LastSyncTime.Format(time.DateTime),
			IsDeleted:      false,
			CreateTime:     repoEntity.CreateTime.Format(time.DateTime),
			UpdateTime:     repoEntity.UpdateTime.Format(time.DateTime),
		}
	}

	return repoModels, nil
}

func (m *MysqlRepositoryManager) GetAllRepositories(ctx context.Context) ([]*model.Repository, error) {
	var repoEntities []*entity.MysqlRepository

	err := m.db.WithContext(ctx).
		Find(&repoEntities).Error
	if err != nil {
		return nil, err
	}

	repoModels := make([]*model.Repository, len(repoEntities))

	for i, repoEntity := range repoEntities {
		repoModels[i] = &model.Repository{
			Id:             repoEntity.ID,
			RepositoryType: repoEntity.RepositoryType,
			StoreType:      repoEntity.StoreType,
			RepositoryUrl:  repoEntity.RepositoryURL,
			Token:          repoEntity.Token,
			Status:         repoEntity.Status,
			LastUpdateTime: repoEntity.UpdateTime.Format(time.DateTime),
			LastSyncTime:   repoEntity.LastSyncTime.Format(time.DateTime),
			IsDeleted:      false,
			CreateTime:     repoEntity.CreateTime.Format(time.DateTime),
			UpdateTime:     repoEntity.UpdateTime.Format(time.DateTime),
		}
	}

	return repoModels, nil
}

func (m *MysqlRepositoryManager) GetTokenByID(ctx context.Context, id int64) (string, error) {
	var token string

	err := m.db.WithContext(ctx).
		Table(entity.TableNameMysqlRepository).
		Select("`token`").
		Where("`id = ?`", id).
		Take(&token).Error

	return token, err
}

func (m *MysqlRepositoryManager) GetRepoTypeByID(ctx context.Context, id int64) (int32, error) {
	var repoType int32

	err := m.db.WithContext(ctx).
		Table(entity.TableNameMysqlRepository).
		Select("`repo_type`").
		Where("`id = ?`", id).
		Take(&repoType).Error

	return repoType, err
}
