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

package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/cloudwego/cwgo/platform/server/shared/consts"
	"github.com/cloudwego/cwgo/platform/server/shared/dao/entity"
	"github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IRepositoryDaoManager interface {
	AddRepository(ctx context.Context, repoModel model.Repository) (int64, error)

	DeleteRepository(ctx context.Context, ids []int64) error

	UpdateRepository(ctx context.Context, repoModel model.Repository) error
	Sync(ctx context.Context, repoModel model.Repository) error
	ChangeRepositoryStatus(ctx context.Context, id int64, status int32) error

	GetRepository(ctx context.Context, id int64) (*model.Repository, error)
	GetRepositoryList(ctx context.Context, repositoryModel model.Repository, page, limit, order int32, orderBy string) ([]*model.Repository, int64, error)
	GetAllRepositories(ctx context.Context) ([]*model.Repository, error)
	GetTokenByID(ctx context.Context, id int64) (string, error)
	GetRepoTypeByID(ctx context.Context, id int64) (int32, error)
	IsExist(ctx context.Context, domain, owner, repoName string) (bool, error)
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
		lastUpdateTime, _ = time.ParseInLocation(time.DateTime, repoModel.LastUpdateTime, consts.TimeZone)
	} else {
		lastUpdateTime = time.Now()
	}

	var repoEntity entity.MysqlRepository

	// check if repo exists
	isExist, err := m.IsExist(ctx, repoModel.RepositoryDomain, repoModel.RepositoryOwner, repoModel.RepositoryName)
	if err != nil {
		return -1, err
	}
	if isExist {
		return -1, consts.ErrDatabaseDuplicateRecord
	}

	// create repo if record is not exist or record's `is_deleted` = 1
	repoEntity = entity.MysqlRepository{
		RepositoryType: repoModel.RepositoryType,
		Domain:         repoModel.RepositoryDomain,
		Owner:          repoModel.RepositoryOwner,
		RepositoryName: repoModel.RepositoryName,
		Branch:         repoModel.RepositoryBranch,
		StoreType:      repoModel.StoreType,
		TokenId:        repoModel.TokenId,
		LastUpdateTime: lastUpdateTime,
		LastSyncTime:   time.Now(),
		Status:         consts.RepositoryStatusNumActive,
	}
	err = m.db.WithContext(ctx).
		Clauses(
			clause.OnConflict{
				Columns:   []clause.Column{{Name: "domain"}, {Name: "owner"}, {Name: "repository_name"}},
				UpdateAll: true,
			},
		).
		Create(&repoEntity).Error

	return repoEntity.ID, err
}

func (m *MysqlRepositoryManager) DeleteRepository(ctx context.Context, ids []int64) error {
	var repoEntity entity.MysqlRepository

	res := m.db.WithContext(ctx).
		Delete(&repoEntity, ids)
	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return consts.ErrDatabaseRecordNotFound
	}

	return nil
}

func (m *MysqlRepositoryManager) UpdateRepository(ctx context.Context, repoModel model.Repository) error {
	if repoModel.Status != 0 {
		if _, ok := consts.RepositoryStatusNumMap[int(repoModel.Status)]; !ok {
			return consts.ErrParamRepositoryStatus
		}
	}

	repoEntity := entity.MysqlRepository{
		ID:     repoModel.Id,
		Branch: repoModel.RepositoryBranch,
		Status: repoModel.Status,
	}

	err := m.db.WithContext(ctx).
		Model(&repoEntity).
		Updates(repoEntity).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return consts.ErrDatabaseRecordNotFound
		}
	}

	return err
}

func (m *MysqlRepositoryManager) Sync(ctx context.Context, repoModel model.Repository) error {
	var lastUpdateTime, lastSyncTime time.Time
	if repoModel.LastUpdateTime != "" {
		lastUpdateTime, _ = time.ParseInLocation(time.DateTime, repoModel.LastUpdateTime, consts.TimeZone)
	}
	if repoModel.LastSyncTime != "" {
		lastSyncTime, _ = time.ParseInLocation(time.DateTime, repoModel.LastSyncTime, consts.TimeZone)
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
		if err == gorm.ErrRecordNotFound {
			return nil, consts.ErrDatabaseRecordNotFound
		}
		return nil, err
	}

	return &model.Repository{
		Id:               repoEntity.ID,
		RepositoryType:   repoEntity.RepositoryType,
		RepositoryDomain: repoEntity.Domain,
		RepositoryOwner:  repoEntity.Owner,
		RepositoryName:   repoEntity.RepositoryName,
		RepositoryBranch: repoEntity.Branch,
		StoreType:        repoEntity.StoreType,
		TokenId:          repoEntity.TokenId,
		Status:           repoEntity.Status,
		LastUpdateTime:   repoEntity.UpdateTime.Format(time.DateTime),
		LastSyncTime:     repoEntity.LastSyncTime.Format(time.DateTime),
		IsDeleted:        false,
		CreateTime:       repoEntity.CreateTime.Format(time.DateTime),
		UpdateTime:       repoEntity.UpdateTime.Format(time.DateTime),
	}, nil
}

func (m *MysqlRepositoryManager) GetRepositoryList(ctx context.Context, repositoryModel model.Repository, page, limit, order int32, orderBy string) ([]*model.Repository, int64, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	var total int64

	db := m.db.WithContext(ctx)

	if repositoryModel.RepositoryType != 0 {
		db = db.Where("`repository_type` = ?", repositoryModel.RepositoryType)
	}
	if repositoryModel.StoreType != 0 {
		db = db.Where("`store_type` = ?", repositoryModel.StoreType)
	}
	if repositoryModel.RepositoryDomain != "" {
		db = db.Where("`domain` LIKE ?", fmt.Sprintf("%%%s%%", repositoryModel.RepositoryDomain))
	}
	if repositoryModel.RepositoryOwner != "" {
		db = db.Where("`owner` LIKE ?", fmt.Sprintf("%%%s%%", repositoryModel.RepositoryOwner))
	}
	if repositoryModel.RepositoryName != "" {
		db = db.Where("`repository_name` LIKE ?", fmt.Sprintf("%%%s%%", repositoryModel.RepositoryName))
	}

	err := db.
		Model(&entity.MysqlRepository{}).
		Count(&total).Error
	if err != nil {
		return nil, -1, err
	}

	if int64(offset) >= total {
		return nil, total, nil
	}

	var repoEntities []*entity.MysqlRepository

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
		Offset(int(offset)).
		Limit(int(limit)).
		Order(orderBy).
		Find(&repoEntities).Error
	if err != nil {
		return nil, -1, err
	}

	repoModels := make([]*model.Repository, len(repoEntities))

	for i, repoEntity := range repoEntities {
		repoModels[i] = &model.Repository{
			Id:               repoEntity.ID,
			RepositoryType:   repoEntity.RepositoryType,
			RepositoryDomain: repoEntity.Domain,
			RepositoryOwner:  repoEntity.Owner,
			RepositoryName:   repoEntity.RepositoryName,
			RepositoryBranch: repoEntity.Branch,
			StoreType:        repoEntity.StoreType,
			TokenId:          repoEntity.TokenId,
			Status:           repoEntity.Status,
			LastUpdateTime:   repoEntity.UpdateTime.Format(time.DateTime),
			LastSyncTime:     repoEntity.LastSyncTime.Format(time.DateTime),
			IsDeleted:        false,
			CreateTime:       repoEntity.CreateTime.Format(time.DateTime),
			UpdateTime:       repoEntity.UpdateTime.Format(time.DateTime),
		}
	}

	return repoModels, total, nil
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
			Id:               repoEntity.ID,
			RepositoryType:   repoEntity.RepositoryType,
			RepositoryDomain: repoEntity.Domain,
			RepositoryOwner:  repoEntity.Owner,
			RepositoryName:   repoEntity.RepositoryName,
			RepositoryBranch: repoEntity.Branch,
			StoreType:        repoEntity.StoreType,
			TokenId:          repoEntity.TokenId,
			Status:           repoEntity.Status,
			LastUpdateTime:   repoEntity.UpdateTime.Format(time.DateTime),
			LastSyncTime:     repoEntity.LastSyncTime.Format(time.DateTime),
			IsDeleted:        false,
			CreateTime:       repoEntity.CreateTime.Format(time.DateTime),
			UpdateTime:       repoEntity.UpdateTime.Format(time.DateTime),
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

func (m *MysqlRepositoryManager) IsExist(ctx context.Context, domain, owner, repoName string) (bool, error) {
	err := m.db.WithContext(ctx).
		Where("`domain` = ? AND `owner` = ? AND 'repository_name' = ? AND `is_deleted` = 0",
			domain,
			owner,
			repoName,
		).
		Take(&entity.MysqlRepository{}).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return false, err
		}
	} else {
		return true, nil
	}

	return false, nil
}
