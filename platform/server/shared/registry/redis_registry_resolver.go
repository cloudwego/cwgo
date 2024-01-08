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
	"strings"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/cwgo/platform/server/shared/log"
	"github.com/cloudwego/kitex/pkg/discovery"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var _ discovery.Resolver = (*RedisResolver)(nil)

type RedisResolver struct {
	rdb redis.UniversalClient
}

func NewRedisResolver(rdb redis.UniversalClient) discovery.Resolver {
	return &RedisResolver{
		rdb: rdb,
	}
}

func (r *RedisResolver) Target(ctx context.Context, target rpcinfo.EndpointInfo) (description string) {
	var serviceIdentification strings.Builder
	agentServiceValue, ok := target.Tag(agentServiceFlagKey)
	if ok {
		serviceIdentification.WriteString(generateKey(target.ServiceName(), agentServiceValue))
	} else {
		serviceIdentification.WriteString(generateKey(target.ServiceName()))
	}

	return serviceIdentification.String()
}

func (r *RedisResolver) Resolve(ctx context.Context, desc string) (discovery.Result, error) {
	var instances []discovery.Instance
	fvs := r.rdb.HGetAll(ctx, desc).Val()

	for _, hashTable := range fvs {
		var rInfo registryInfo
		err := sonic.Unmarshal([]byte(hashTable), &rInfo)
		if err != nil {
			log.Error("fail to unmarshal", zap.Error(err), zap.String("addr", rInfo.Addr))
			continue
		}
		weight := rInfo.Weight
		if weight <= 0 {
			weight = defaultWeight
		}
		instances = append(instances, discovery.NewInstance("tcp", rInfo.Addr, weight, rInfo.Tags))
	}
	result := discovery.Result{
		Cacheable: true,
		CacheKey:  desc,
		Instances: instances,
	}
	return result, nil
}

func (r *RedisResolver) Diff(cacheKey string, prev, next discovery.Result) (discovery.Change, bool) {
	return discovery.DefaultDiff(cacheKey, prev, next)
}

func (r *RedisResolver) Name() string {
	return "redis"
}
