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
	"github.com/cloudwego/cwgo/platform/server/shared/config/internal/idl"
	"github.com/cloudwego/cwgo/platform/server/shared/config/internal/repository"
)

type IManager interface {
	GetIdlManager() idl.IIdlManager
	GetRepositoryManager() repository.IRepositoryManager
}

type Manager struct {
	IdlManager        idl.IIdlManager
	RepositoryManager repository.IRepositoryManager
}

var manager *Manager

func NewManager() *Manager {

	return &Manager{}
}

func GetManager() *Manager {
	return manager
}

func (m *Manager) GetIdlManager() idl.IIdlManager {
	return m.IdlManager
}

func (m *Manager) GetRepositoryManager() repository.IRepositoryManager {
	return m.RepositoryManager
}
