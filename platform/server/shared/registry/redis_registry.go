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

package registry

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/cloudwego/cwgo/platform/server/shared/consts"
	"github.com/cloudwego/cwgo/platform/server/shared/meta"
	"github.com/cloudwego/kitex/pkg/registry"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/cwgo/platform/server/shared/log"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var _ registry.Registry = (*RedisRegistry)(nil)

var (
	registerScript = redis.NewScript(`
		local key = KEYS[1]
		local field = ARGV[1]
		local value = ARGV[2]
		local expireTime = tonumber(ARGV[3])
		local message = ARGV[4]
		
		redis.call('HSET', key, field, value)
		redis.call('EXPIRE', key, expireTime)
    `)
	deregisterScript = redis.NewScript(`
		local key = KEYS[1]
	    local field = ARGV[1]
	    local message = ARGV[2]
	
	    redis.call('HDEL', key, field)
	`)
)

const (
	defaultExpireTime    = 60
	defaultTickerTime    = time.Second * 30
	defaultKeepAliveTime = time.Second * 60
	defaultWeight        = 10

	agentServiceFlagKey = "service_id"

	register   = "register"
	deregister = "deregister"
)

type RedisRegistry struct {
	sync.RWMutex
	sync.WaitGroup
	ctx      context.Context
	cancelFn context.CancelFunc
	rdb      redis.UniversalClient
}

// NewRedisRegistry returns a new RedisRegistry.
func NewRedisRegistry(rdb redis.UniversalClient) *RedisRegistry {
	return &RedisRegistry{
		RWMutex:   sync.RWMutex{},
		WaitGroup: sync.WaitGroup{},
		rdb:       rdb,
	}
}

type registryInfo struct {
	ServiceName string            `json:"service_name"`
	Addr        string            `json:"addr"`
	Weight      int               `json:"weight"`
	Tags        map[string]string `json:"tags"`
}

// convertRegistryInfo convert registry.Info to registryInfo
func convertRegistryInfo(info *registry.Info) *registryInfo {
	return &registryInfo{
		ServiceName: info.ServiceName,
		Addr:        info.Addr.String(),
		Weight:      info.Weight,
		Tags:        info.Tags,
	}
}

// registryHashMap is a struct for redis hash map
type registryHashMap struct {
	// key is a service name
	key string
	// field is a service addr
	field string
	// value is a registryInfo json string
	value string
}

// prepareRegistryMeta prepare the registry meta save in registryHashMap
func prepareRegistryMeta(info *registry.Info) (*registryHashMap, error) {
	registryMeta, err := sonic.Marshal(convertRegistryInfo(info))
	if err != nil {
		return nil, err
	}

	return &registryHashMap{
		key:   generateKey(info.ServiceName),
		field: info.Addr.String(),
		value: string(registryMeta),
	}, nil
}

func (r *RedisRegistry) Register(info *registry.Info) error {
	err := r.rdb.Ping(context.Background()).Err()
	if err != nil {
		log.Error("redis connect fail", zap.Error(err))
		return err
	}

	ctx, cancelFn := context.WithCancel(context.Background())

	r.Lock()
	r.cancelFn = cancelFn
	r.ctx = ctx
	r.Unlock()

	hashTable, err := prepareRegistryMeta(info)
	if err != nil {
		return err
	}
	args := []any{
		hashTable.field, hashTable.value, defaultExpireTime,
		generateMsg(register, info.ServiceName, info.Addr.String()),
	}

	// HSET cwgo:ping "127.0.0.1:8081" registryHashMap
	// PUBLISH cwgo:ping "register-ping-127.0.0.1:8081"
	err = registerScript.Run(ctx, r.rdb, []string{hashTable.key}, args).Err()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			return err
		}
	}

	// go m.monitorTTL(r.ctx, hashTable, info, r)
	go keepAlive(r.ctx, hashTable, r)
	return nil
}

func (r *RedisRegistry) Deregister(info *registry.Info) error {
	hashTable, err := prepareRegistryMeta(info)
	if err != nil {
		return err
	}
	args := []any{
		hashTable.field, generateMsg(deregister, info.ServiceName, info.Addr.String()),
	}
	err = deregisterScript.Run(r.ctx, r.rdb, []string{hashTable.key}, args).Err()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			return err
		}
	}
	r.cancelFn()
	log.Info("deregister agent success", zap.String("addr", info.Addr.String()))
	return nil
}

func (r *RedisRegistry) GetAgents() ([]*meta.Agent, error) {
	var agents []*meta.Agent
	hashTable := r.rdb.HGetAll(context.Background(), generateKey(consts.ServiceNameAgent)).Val()
	for _, v := range hashTable {
		var rInfo registryInfo
		err := sonic.Unmarshal([]byte(v), &rInfo)
		if err != nil {
			return nil, err
		}
		agents = append(agents, &meta.Agent{
			ID:   rInfo.Tags[agentServiceFlagKey],
			Host: rInfo.Addr,
		})
	}
	return agents, nil
}

// keepAlive keep the registry information alive
func keepAlive(ctx context.Context, hash *registryHashMap, r *RedisRegistry) {
	ticker := time.NewTicker(defaultTickerTime)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			r.rdb.Expire(ctx, hash.key, defaultKeepAliveTime)
		case <-ctx.Done():
			break
		}
	}
}

func generateKey(srvName string, srvID ...string) string {
	if len(srvID) > 0 {
		return fmt.Sprintf("%s:%s:%s", consts.ProjectName, srvName, srvID[0])
	}
	return fmt.Sprintf("%s:%s", consts.ProjectName, srvName)
}

func generateMsg(msgType, serviceName, serviceAddr string) string {
	return fmt.Sprintf("%s/%s/%s", msgType, serviceName, serviceAddr)
}
