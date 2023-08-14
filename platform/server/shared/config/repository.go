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
	GetRepoTypeByID(id int64) (int32, error)

	AddRepository(repoURL, lastUpdateTime, lastSyncTime, token, status string, repoType int32) error
	DeleteRepository(ids string) error
	UpdateRepository(id, token string) error
	GetRepositories(page, limit int32) ([]Repository, error)
}

type MysqlRepository struct {
	Db *gorm.DB
}

type Repository struct {
	Id             int64
	RepositoryUrl  string
	LastUpdateTime string
	LastSyncTime   string
	Token          string
	Status         string
	RepoType       int32
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
	result := r.Db.Model(&repo).Where("id = ?", id).Take(&repo)
	if result.Error != nil {
		return "", result.Error
	}

	return repo.Token, nil
}

func (r *MysqlRepository) GetRepoTypeByID(id int64) (int32, error) {
	var repo Repository
	result := r.Db.Model(&repo).Where("id = ?", id).Take(&repo)
	if result.Error != nil {
		return 0, result.Error
	}

	return repo.RepoType, nil
}

func (r *MysqlRepository) AddRepository(repoURL, lastUpdateTime, lastSyncTime, token, status string, repoType int32) error {
	repo := Repository{
		RepositoryUrl:  repoURL,
		LastUpdateTime: lastSyncTime,
		LastSyncTime:   lastSyncTime,
		Token:          token,
		Status:         status,
		RepoType:       repoType,
	}
	result := r.Db.Create(&repo)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *MysqlRepository) DeleteRepository(ids string) error {
	var repo Repository
	result := r.Db.Delete(&repo, ids)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *MysqlRepository) UpdateRepository(id, token string) error {
	result := r.Db.Model(&Repository{}).Where("id = ?", id).Update("token", token)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *MysqlRepository) GetRepositories(page, limit int32, sortBy string) ([]Repository, error) {
	var repos []Repository

	// Default sort field to 'update_time' if not provided
	if sortBy == "" {
		sortBy = "update_time"
	}

	offset := (page - 1) * limit
	result := r.Db.Offset(int(offset)).Limit(int(limit)).Find(&repos)
	if result.Error != nil {
		return nil, result.Error
	}

	return repos, nil
}
