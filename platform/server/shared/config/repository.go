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

import (
	"gorm.io/gorm"
)

type IRepository interface {
	GetTokenByID(id int64) (string, error)
	GetRepoTypeByID(id int64) (int64, error)

	AddIDL(repoId int64, idlPath, idlHash, serviceName string) error
	DeleteIDLs(id int64) error
	UpdateIDL(id, repoId int64, idlPath, serviceName string) error
	GetIDLs(page, limit int32) []IDL
}

type MysqlRepository struct {
	db *gorm.DB
}

type Repository struct {
	Id             int64
	RepositoryUrl  string
	LastUpdateTime string
	LastSyncTime   string
	Token          string
	Status         string
	RepoType       int64
}

type IDL struct {
	Id           int64
	RepositoryId int64
	MainIdlPath  string
	IdlHash      string
	ServiceName  string
}

func (r *MysqlRepository) GetTokenByID(id int64) (string, error) {
	var repo Repository
	result := r.db.Model(&repo).Where("id = ?", id).First(&repo)
	if result.Error != nil {
		return "", result.Error
	}

	return repo.Token, nil
}

func (r *MysqlRepository) GetRepoTypeByID(id int64) (int64, error) {
	var repo Repository
	result := r.db.Model(&repo).Where("id = ?", id).First(&repo)
	if result.Error != nil {
		return 0, result.Error
	}

	return repo.RepoType, nil
}

func (r *MysqlRepository) AddIDL(repoId int64, idlPath, idlHash, serviceName string) error {
	idl := IDL{
		RepositoryId: repoId,
		MainIdlPath:  idlPath,
		IdlHash:      idlHash,
		ServiceName:  serviceName,
	}
	result := r.db.Model(&idl).Create(&idl)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *MysqlRepository) DeleteIDLs(ids []int64) error {
	var idl IDL
	result := r.db.Delete(&idl, ids)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *MysqlRepository) UpdateIDL(id, repoId int64, idlPath, idlHash, serviceName string) error {
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

func (r *MysqlRepository) GetIDLs(page, limit int32) ([]IDL, error) {
	var IDLs []IDL
	offset := (page - 1) * limit
	result := r.db.Offset(int(offset)).Limit(int(limit)).Find(&IDLs)
	if result.Error != nil {
		return nil, result.Error
	}

	return IDLs, nil
}
