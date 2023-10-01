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
	"github.com/cloudwego/cwgo/platform/server/shared/utils"
	"gorm.io/gorm"
)

type IRepositoryDaoManager interface {
	GetTokenByID(id int64) (string, error)
	GetRepoTypeByID(id int64) (int8, error)

	AddRepository(repoURL, token, status string, repoType int8) error
	DeleteRepository(ids string) error
	UpdateRepository(id, token string) error
	GetRepositories(page, limit int32, sortBy string) ([]Repository, error)
}

type Repository struct {
	Id             int64
	RepositoryUrl  string
	Token          string
	Status         string
	LastUpdateTime string
	LastSyncTime   string
	CreateTime     string
	UpdateTime     string
	RepoType       int8
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

func (r *MysqlRepositoryManager) GetTokenByID(id int64) (string, error) {
	var repo Repository
	result := r.db.Model(&repo).Where("id = ?", id).Take(&repo)
	if result.Error != nil {
		return "", result.Error
	}

	return repo.Token, nil
}

func (r *MysqlRepositoryManager) GetRepoTypeByID(id int64) (int8, error) {
	var repo Repository
	result := r.db.Model(&repo).Where("id = ?", id).Take(&repo)
	if result.Error != nil {
		return 0, result.Error
	}

	return repo.RepoType, nil
}

func (r *MysqlRepositoryManager) AddRepository(repoURL, token, status string, repoType int8) error {
	timeNow := utils.GetCurrentTime()
	repo := Repository{
		RepositoryUrl:  repoURL,
		Token:          token,
		Status:         status,
		RepoType:       repoType,
		LastUpdateTime: "0",
		LastSyncTime:   timeNow,
		CreateTime:     timeNow,
		UpdateTime:     timeNow,
	}
	result := r.db.Create(&repo)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *MysqlRepositoryManager) DeleteRepository(ids string) error {
	var repo Repository
	result := r.db.Delete(&repo, ids)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *MysqlRepositoryManager) UpdateRepository(id, token string) error {
	timeNow := utils.GetCurrentTime()
	result := r.db.Model(&Repository{}).Where("id = ?", id).Updates(Repository{
		Token:      token,
		UpdateTime: timeNow,
	})
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *MysqlRepositoryManager) GetRepositories(page, limit int32, sortBy string) ([]Repository, error) {
	var repos []Repository

	// Default sort field to 'update_time' if not provided
	if sortBy == "" {
		sortBy = SortByUpdateTime
	}

	offset := (page - 1) * limit
	result := r.db.Offset(int(offset)).Limit(int(limit)).Order(sortBy).Find(&repos)
	if result.Error != nil {
		return nil, result.Error
	}

	return repos, nil
}
