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

package token

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudwego/cwgo/platform/server/shared/consts"
	"github.com/cloudwego/cwgo/platform/server/shared/dao/entity"
	"github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/model"
	"gorm.io/gorm"
)

type ITokenDaoManager interface {
	AddToken(ctx context.Context, tokenModel model.Token) (int64, error)

	DeleteToken(ctx context.Context, ids []int64) error

	GetTokenList(ctx context.Context, tokenModel model.Token, page, limit, order int32, orderBy string) ([]*model.Token, int64, error)
	GetActiveTokenForDomain(ctx context.Context, domain string) ([]*model.Token, error)
	GetTokenById(ctx context.Context, id int64) (*model.Token, error)
}

type MysqlTokenManager struct {
	db *gorm.DB
}

var _ ITokenDaoManager = (*MysqlTokenManager)(nil)

func NewMysqlToken(db *gorm.DB) *MysqlTokenManager {
	return &MysqlTokenManager{
		db: db,
	}
}

func (m *MysqlTokenManager) AddToken(ctx context.Context, tokenModel model.Token) (int64, error) {
	expirationTime, err := time.ParseInLocation("2006-01-02 15:04:05.999999999 -0700 MST", tokenModel.ExpirationTime, consts.TimeZone)
	if err != nil {
		return -1, err
	}

	tokenEntity := entity.MysqlToken{
		Owner:            tokenModel.Owner,
		RepositoryType:   tokenModel.RepositoryType,
		RepositoryDomain: tokenModel.RepositoryDomain,
		Token:            tokenModel.Token,
		Status:           tokenModel.Status,
		ExpirationTime:   expirationTime,
	}

	err = m.db.WithContext(ctx).
		Create(&tokenEntity).Error

	return tokenEntity.ID, err
}

func (m *MysqlTokenManager) DeleteToken(ctx context.Context, ids []int64) error {
	var tokenEntity entity.MysqlToken

	res := m.db.WithContext(ctx).
		Delete(&tokenEntity, ids)
	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return consts.ErrDatabaseRecordNotFound
	}

	return nil
}

func (m *MysqlTokenManager) GetTokenList(ctx context.Context, tokenModel model.Token, page, limit, order int32, orderBy string) ([]*model.Token, int64, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	var total int64

	db := m.db.WithContext(ctx)

	if tokenModel.Owner != "" {
		db = db.Where("`owner` = ?", tokenModel.Owner)
	}
	if tokenModel.RepositoryType != 0 {
		db = db.Where("`repository_type` = ?", tokenModel.RepositoryType)
	}
	if tokenModel.RepositoryDomain != "" {
		db = db.Where("`repository_domain` LIKE ?", fmt.Sprintf("%%%s%%", tokenModel.RepositoryDomain))
	}

	err := db.
		Model(&entity.MysqlToken{}).
		Count(&total).Error
	if err != nil {
		return nil, -1, err
	}

	if int64(offset) >= total {
		return nil, total, nil
	}

	var tokenEntities []*entity.MysqlToken

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
		Find(&tokenEntities).Error
	if err != nil {
		return nil, -1, err
	}

	tokenModels := make([]*model.Token, len(tokenEntities))

	for i, tokenEntity := range tokenEntities {
		tokenModels[i] = &model.Token{
			Id:               tokenEntity.ID,
			Owner:            tokenEntity.Token,
			RepositoryType:   tokenEntity.RepositoryType,
			RepositoryDomain: tokenEntity.RepositoryDomain,
			Token:            tokenEntity.Token,
			Status:           tokenEntity.Status,
			ExpirationTime:   tokenEntity.ExpirationTime.Format(time.DateTime),
			IsDeleted:        false,
			CreateTime:       tokenEntity.CreateTime.Format(time.DateTime),
			UpdateTime:       tokenEntity.UpdateTime.Format(time.DateTime),
		}
	}

	return tokenModels, total, nil
}

func (m *MysqlTokenManager) GetActiveTokenForDomain(ctx context.Context, domain string) ([]*model.Token, error) {
	var tokenEntities []*entity.MysqlToken

	err := m.db.WithContext(ctx).
		Model(&entity.MysqlToken{}).
		Where("`expiration_time` <= ?", time.Now()).
		UpdateColumn("status", consts.TokenStatusNumExpired).Error
	if err != nil {
		return nil, err
	}

	err = m.db.WithContext(ctx).
		Where("`repository_domain` = ? AND `status` = ?", domain, consts.TokenStatusNumValid).
		Find(&tokenEntities).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, consts.ErrDatabaseRecordNotFound
		}
		return nil, nil
	}

	tokenModels := make([]*model.Token, len(tokenEntities))

	for i, tokenEntity := range tokenEntities {
		tokenModels[i] = &model.Token{
			Id:               tokenEntity.ID,
			Owner:            tokenEntity.Token,
			RepositoryType:   tokenEntity.RepositoryType,
			RepositoryDomain: tokenEntity.RepositoryDomain,
			Token:            tokenEntity.Token,
			Status:           tokenEntity.Status,
			ExpirationTime:   tokenEntity.ExpirationTime.Format(time.DateTime),
			IsDeleted:        false,
			CreateTime:       tokenEntity.CreateTime.Format(time.DateTime),
			UpdateTime:       tokenEntity.UpdateTime.Format(time.DateTime),
		}
	}

	return tokenModels, nil
}

func (m *MysqlTokenManager) GetTokenById(ctx context.Context, id int64) (*model.Token, error) {
	var tokenEntity entity.MysqlToken

	err := m.db.WithContext(ctx).
		Where("`id` = ?", id).
		Take(&tokenEntity).Error
	if err != nil {
		return nil, err
	}

	return &model.Token{
		Id:               tokenEntity.ID,
		Owner:            tokenEntity.Token,
		RepositoryType:   tokenEntity.RepositoryType,
		RepositoryDomain: tokenEntity.RepositoryDomain,
		Token:            tokenEntity.Token,
		Status:           tokenEntity.Status,
		ExpirationTime:   tokenEntity.ExpirationTime.Format(time.DateTime),
		IsDeleted:        false,
		CreateTime:       tokenEntity.CreateTime.Format(time.DateTime),
		UpdateTime:       tokenEntity.UpdateTime.Format(time.DateTime),
	}, nil
}
