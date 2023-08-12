/*
 * Copyright 2022 CloudWeGo Authors
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
 */

package config

import "gorm.io/gorm"

type IIdl interface {
	AddIDL(repoId int64, idlPath, idlHash, serviceName string) error
	DeleteIDLs(id int64) error
	UpdateIDL(id, repoId int64, idlPath, serviceName string) error
	GetIDLs(page, limit int32) []IDL
}

type MysqlIDL struct {
	db *gorm.DB
}

func (r *MysqlIDL) AddIDL(repoId int64, idlPath, idlHash, serviceName string) error {
	idl := IDL{
		RepositoryId: repoId,
		MainIdlPath:  idlPath,
		IdlHash:      idlHash,
		ServiceName:  serviceName,
	}
	result := r.db.Create(&idl)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *MysqlIDL) DeleteIDLs(ids []int64) error {
	var idl IDL
	result := r.db.Delete(&idl, ids)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *MysqlIDL) UpdateIDL(id, repoId int64, idlPath, idlHash, serviceName string) error {
	idl := IDL{
		Id:           id,
		RepositoryId: repoId,
		MainIdlPath:  idlPath,
		IdlHash:      idlHash,
		ServiceName:  serviceName,
	}
	result := r.db.Save(&idl)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *MysqlIDL) GetIDLs(page, limit int32) ([]IDL, error) {
	var IDLs []IDL
	offset := (page - 1) * limit
	result := r.db.Offset(int(offset)).Limit(int(limit)).Find(&IDLs)
	if result.Error != nil {
		return nil, result.Error
	}

	return IDLs, nil
}
