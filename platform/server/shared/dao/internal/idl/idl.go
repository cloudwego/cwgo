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
	"github.com/cloudwego/cwgo/platform/server/shared/consts"
	"github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/model"
	"github.com/cloudwego/cwgo/platform/server/shared/utils"
	"gorm.io/gorm"
)

type IIdlDaoManager interface {
	AddIDL(repoId int64, idlPath, serviceName string) error
	DeleteIDLs(ids []int64) error
	UpdateIDL(id, repoId int64, idlPath, serviceName string) error
	GetIDL(id int64) (model.IDL, error)
	GetIDLs(page, limit, order int32, orderBy string) ([]*model.IDL, error)
	SyncIDLContent(id int64, content string) error
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

func (r *MysqlIDLManager) AddIDL(repoId int64, idlPath, serviceName string) error {
	timeNow := utils.GetCurrentTime()
	idl := model.IDL{
		RepositoryId: repoId,
		MainIdlPath:  idlPath,
		ServiceName:  serviceName,
		LastSyncTime: timeNow,
		CreateTime:   timeNow,
		UpdateTime:   timeNow,
	}
	res := r.db.
		Table(consts.TableNameIDL).
		Create(&idl)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (r *MysqlIDLManager) DeleteIDLs(ids []int64) error {
	var idl model.IDL
	res := r.db.
		Table(consts.TableNameIDL).
		Delete(&idl, ids)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (r *MysqlIDLManager) UpdateIDL(id, repoId int64, idlPath, serviceName string) error {
	timeNow := utils.GetCurrentTime()
	res := r.db.
		Table(consts.TableNameIDL).
		Where("id = ?", id).
		Updates(
			model.IDL{
				Id:           id,
				RepositoryId: repoId,
				MainIdlPath:  idlPath,
				ServiceName:  serviceName,
				UpdateTime:   timeNow,
			},
		)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (r *MysqlIDLManager) GetIDL(id int64) (model.IDL, error) {
	var idl model.IDL
	res := r.db.
		Table(consts.TableNameIDL).
		Where("id = ?", id).
		First(&idl)
	if res.Error != nil {
		return idl, res.Error
	}

	return idl, nil
}

func (r *MysqlIDLManager) GetIDLs(page, limit, order int32, orderBy string) ([]*model.IDL, error) {
	var IDLs []*model.IDL
	offset := (page - 1) * limit

	// Default sort field to 'update_time' if not provided
	if orderBy == "" {
		orderBy = consts.OrderByUpdateTime
	}

	switch order {
	case consts.OrderNumInc:
		orderBy = orderBy + " " + consts.OrderInc
	case consts.OrderNumDec:
		orderBy = orderBy + " " + consts.OrderDec
	default:
		orderBy = orderBy + " " + consts.OrderInc
	}

	res := r.db.
		Table(consts.TableNameIDL).
		Offset(int(offset)).
		Limit(int(limit)).
		Order(orderBy).
		Find(&IDLs)
	if res.Error != nil {
		return nil, res.Error
	}

	return IDLs, nil
}

func (r *MysqlIDLManager) SyncIDLContent(id int64, content string) error {
	timeNow := utils.GetCurrentTime()
	res := r.db.
		Table(consts.TableNameIDL).
		Where("id = ?", id).
		Updates(
			model.IDL{
				Content:      content,
				LastSyncTime: timeNow,
				UpdateTime:   timeNow,
			},
		)
	if res.Error != nil {
		return res.Error
	}

	return nil
}
