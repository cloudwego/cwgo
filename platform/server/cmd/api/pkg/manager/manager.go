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

package manager

import (
	"sync"
	"time"

	"github.com/cloudwego/kitex/pkg/discovery"

	"github.com/redis/go-redis/v9"

	"github.com/cloudwego/cwgo/platform/server/cmd/api/pkg/dispatcher"
	"github.com/cloudwego/cwgo/platform/server/shared/config"
	"github.com/cloudwego/cwgo/platform/server/shared/dao"
	"github.com/cloudwego/cwgo/platform/server/shared/meta"
	"github.com/cloudwego/cwgo/platform/server/shared/registry"
)

// Manager that
// api manager
type Manager struct {
	sync.Mutex
	updateTaskInterval    time.Duration
	currentUpdateTaskTime time.Time
	lastUpdateTaskTime    time.Time
	agents                []*meta.Agent
	syncAgentInterval     time.Duration
	syncIdlInterval       time.Duration
	// api service id
	apiID string

	// isMasterApi
	// false:
	//   1. will not push task to agent service
	//   2. trying to promote to master
	// true:
	//   1. will push task to agent service
	//   2. trying to maintain the master identity
	isMasterApi bool
	rdb         redis.UniversalClient
	daoManager  *dao.Manager
	dispatcher  dispatcher.IDispatcher
	registry    registry.IRegistry
	resolver    discovery.Resolver
}

func NewApiManager(
	appConf config.AppConfig,
	apiID string,
	rdb redis.UniversalClient,
	daoManager *dao.Manager,
	dispatcher dispatcher.IDispatcher,
	registry registry.IRegistry,
	resolver discovery.Resolver,
) *Manager {
	manager := &Manager{
		Mutex:                 sync.Mutex{},
		rdb:                   rdb,
		isMasterApi:           false,
		updateTaskInterval:    3 * time.Second,
		currentUpdateTaskTime: time.Time{},
		lastUpdateTaskTime:    time.Now(),
		agents:                make([]*meta.Agent, 0),
		syncAgentInterval:     appConf.GetSyncAgentServiceInterval(),
		syncIdlInterval:       appConf.GetSyncIdlInterval(),
		apiID:                 apiID,
		daoManager:            daoManager,
		dispatcher:            dispatcher,
		registry:              registry,
		resolver:              resolver,
	}

	go manager.tryPromoteApiToMasterSrv()
	go manager.watchTaskUpdate()
	go manager.syncTaskFromDB()
	go manager.startUpdate()

	return manager
}
