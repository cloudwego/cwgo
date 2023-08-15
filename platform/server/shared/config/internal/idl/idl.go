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
	"github.com/cloudwego/cwgo/platform/server/shared/config/internal/repository"
	"github.com/cloudwego/cwgo/platform/server/shared/config/internal/store"
	"gorm.io/gorm"
)

type IIdlManager interface {
	AddIDL(repoId int64, idlPath, idlHash, serviceName string) error
	DeleteIDLs(ids []int64) error
	UpdateIDL(id, repoId int64, idlPath, idlHash, serviceName string) error
	GetIDLs(page, limit int32, sortBy string) ([]repository.IDL, error)
}

type MysqlIDLManager struct {
	Db *gorm.DB
}

var _ IIdlManager = (*MysqlIDLManager)(nil)

func NewMysqlIDL() (*MysqlIDLManager, error) {
	db, err := store.GetMysqlDB()
	if err != nil {
		return nil, err
	}

	return &MysqlIDLManager{
		Db: db,
	}, nil
}

func (r *MysqlIDLManager) AddIDL(repoId int64, idlPath, idlHash, serviceName string) error {
	idl := repository.IDL{
		RepositoryId: repoId,
		MainIdlPath:  idlPath,
		IdlHash:      idlHash,
		ServiceName:  serviceName,
	}
	res := r.Db.Create(&idl)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (r *MysqlIDLManager) DeleteIDLs(ids []int64) error {
	var idl repository.IDL
	res := r.Db.Delete(&idl, ids)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (r *MysqlIDLManager) UpdateIDL(id, repoId int64, idlPath, idlHash, serviceName string) error {
	idl := repository.IDL{
		Id:           id,
		RepositoryId: repoId,
		MainIdlPath:  idlPath,
		IdlHash:      idlHash,
		ServiceName:  serviceName,
	}
	res := r.Db.Save(&idl)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (r *MysqlIDLManager) GetIDLs(page, limit int32, sortBy string) ([]repository.IDL, error) {
	var IDLs []repository.IDL
	offset := (page - 1) * limit

	// Default sort field to 'update_time' if not provided
	if sortBy == "" {
		sortBy = SortByUpdateTime
	}

	res := r.Db.Offset(int(offset)).Limit(int(limit)).Order(sortBy).Order(sortBy).Find(&IDLs)
	if res.Error != nil {
		return nil, res.Error
	}

	return IDLs, nil
}
