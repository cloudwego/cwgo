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

package registry

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cloudwego/cwgo/platform/server/shared/consts"
	"github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/registry"
	"github.com/cloudwego/cwgo/platform/server/shared/service"
	"github.com/cloudwego/cwgo/platform/server/shared/utils"
	"github.com/cloudwego/kitex/pkg/discovery"
	kitexregistry "github.com/cloudwego/kitex/pkg/registry"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/patrickmn/go-cache"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type BuiltinRegistry struct {
	cleanInterval time.Duration
	agentCache    *cache.Cache
	rdb           redis.UniversalClient
	syncDuration  time.Duration
	logger        *zap.Logger
}

var _ IRegistry = (*BuiltinRegistry)(nil)

const (
	minCleanInterval  = 100 * time.Millisecond
	defaultExpiration = 60 * time.Second
)

func NewBuiltinRegistry(rdb redis.UniversalClient, logger *zap.Logger) *BuiltinRegistry {
	r := &BuiltinRegistry{
		agentCache:   cache.New(defaultExpiration, 3*time.Second),
		rdb:          rdb,
		syncDuration: 3 * time.Second,
		logger:       logger,
	}

	go r.StartSync()

	return r
}

func (r *BuiltinRegistry) StartSync() {
	for {
		time.Sleep(r.syncDuration)

		cursor := uint64(0)
		var keys []string
		var err error
		var needUpdateServiceId []string
		for {
			keys, cursor, err = r.rdb.Scan(context.Background(), cursor, consts.RdbKeyRegistryService+"*", 10).Result()
			if err != nil {
				r.logger.Error("builtin registry sync agent info in redis failed", zap.Error(err))
				break
			}

			for _, key := range keys {
				serviceId := strings.TrimLeft(key, consts.RdbKeyRegistryService)
				isExist := r.ServiceExists(key)
				if !isExist {
					needUpdateServiceId = append(needUpdateServiceId, serviceId)
				}
			}
		}

		pipe := r.rdb.Pipeline()
		for _, key := range needUpdateServiceId {
			pipe.Get(context.Background(), key)
		}

		cmds, err := pipe.Exec(context.Background())
		if err != nil {
			r.logger.Error("exec redis pipe command failed", zap.Error(err))
			continue
		}

		for i, cmd := range cmds {
			ipPort, err := cmd.(*redis.StringCmd).Result()
			ip, port, err := net.SplitHostPort(ipPort)
			if err != nil {
				r.logger.Error("parse host port string failed", zap.Error(err), zap.String("val", ipPort))
			}

			p, err := strconv.Atoi(port)
			if err != nil {
				r.logger.Error("parse port failed", zap.Error(err), zap.String("val", port))
			}

			err = r.Register(needUpdateServiceId[i], ip, p)
			if err != nil {
				r.logger.Error(
					"register services get by sync progress in redis failed",
					zap.Error(err),
					zap.String("service_id", needUpdateServiceId[i]),
					zap.String("ip", ip),
					zap.Int("port", p),
				)
			}
		}
	}
}

func (r *BuiltinRegistry) Register(serviceId string, host string, port int) error {
	agentService, err := service.NewService(serviceId, host, port)
	if err != nil {
		return err
	}

	r.agentCache.SetDefault(serviceId, agentService)
	err = r.rdb.Set(
		context.Background(),
		fmt.Sprintf(consts.RdbKeyRegistryService, serviceId),
		fmt.Sprintf("%s:%d", host, port),
		defaultExpiration,
	).Err()
	if err != nil {
		r.logger.Error(
			"register service in builtin registry failed",
			zap.Error(err),
			zap.String("service_id", serviceId),
			zap.String("host", host),
			zap.Int("port", port),
		)
	}

	return nil
}

func (r *BuiltinRegistry) Deregister(id string) error {
	r.agentCache.Delete(id)
	err := r.rdb.Del(
		context.Background(),
		fmt.Sprintf(consts.RdbKeyRegistryService, id),
	).Err()
	if err != nil {
		r.logger.Error(
			"deregister service in builtin registry failed",
			zap.Error(err),
			zap.String("service_id", id),
		)
	}

	return nil
}

func (r *BuiltinRegistry) Update(serviceId string) error {
	v, ok := r.agentCache.Get(serviceId)
	if !ok {
		return errors.New("service not found")
	}

	r.agentCache.SetDefault(serviceId, v.(*service.Service))
	err := r.rdb.Expire(
		context.Background(),
		fmt.Sprintf(consts.RdbKeyRegistryService, serviceId),
		defaultExpiration,
	).Err()
	if err != nil {
		if err == redis.Nil {
			return errors.New("service not found")
		}
		r.logger.Error(
			"update service in builtin registry failed",
			zap.Error(err),
			zap.String("service_id", serviceId),
		)
	}
	return nil
}

func (r *BuiltinRegistry) Count() int {
	return r.agentCache.ItemCount()
}

func (r *BuiltinRegistry) GetServiceById(serviceId string) (*service.Service, error) {
	if agentService, ok := r.agentCache.Get(serviceId); !ok {
		return nil, errors.New("service not found")
	} else {
		return agentService.(*service.Service), nil
	}
}

func (r *BuiltinRegistry) GetAllService() ([]*service.Service, error) {
	var services []*service.Service
	for _, svr := range r.agentCache.Items() {
		if !svr.Expired() {
			services = append(services, svr.Object.(*service.Service))
		}
	}

	return services, nil
}

func (r *BuiltinRegistry) ServiceExists(serviceId string) bool {
	_, ok := r.agentCache.Get(serviceId)

	return ok
}

type BuiltinRegistryResolver struct {
	registry *BuiltinRegistry
}

func NewBuiltinRegistryResolver(r *BuiltinRegistry) (discovery.Resolver, error) {
	return &BuiltinRegistryResolver{
		registry: r,
	}, nil
}

func (r *BuiltinRegistryResolver) Target(_ context.Context, target rpcinfo.EndpointInfo) (description string) {
	return consts.ServiceNameAgent
}

func (r *BuiltinRegistryResolver) Resolve(_ context.Context, _ string) (discovery.Result, error) {
	services, _ := r.registry.GetAllService()

	var eps []discovery.Instance

	for _, svr := range services {

		eps = append(eps, discovery.NewInstance(
			"tcp",
			net.JoinHostPort(svr.Host, strconv.Itoa(svr.Port)),
			1,
			map[string]string{"service_id": svr.Id},
		))
	}

	return discovery.Result{
		Cacheable: false,
		CacheKey:  "",
		Instances: eps,
	}, nil
}

func (r *BuiltinRegistryResolver) Diff(cacheKey string, prev, next discovery.Result) (discovery.Change, bool) {
	return discovery.DefaultDiff(cacheKey, prev, next)
}

func (r *BuiltinRegistryResolver) Name() string {
	return "builtin"
}

type BuiltinKitexRegistryClient struct {
	addr           string
	stopChan       chan struct{}
	updateInterval time.Duration
}

func NewBuiltinKitexRegistryClient(addr string) (*BuiltinKitexRegistryClient, error) {
	httpRes, err := http.Get(fmt.Sprintf("http://%s/api/ping", addr))
	if err != nil {
		return nil, err
	}
	defer httpRes.Body.Close()

	body, err := ioutil.ReadAll(httpRes.Body)
	if err != nil {
		return nil, err
	}

	var j registry.RegisterRes

	err = json.Unmarshal(body, &j)
	if err != nil {
		return nil, err
	}

	if j.Code != 0 {
		return nil, errors.New(j.Msg)
	}

	return &BuiltinKitexRegistryClient{
		addr:           addr,
		stopChan:       make(chan struct{}),
		updateInterval: 10 * time.Second,
	}, nil
}

func (rc *BuiltinKitexRegistryClient) Register(info *kitexregistry.Info) error {
	serviceId, ok := info.Tags["service_id"]
	if !ok {
		return errors.New("service_id not found")
	}

	host, port, _ := utils.ParseAddr(info.Addr)

	httpRes, err := http.Get(fmt.Sprintf("http://%s/api/registry/register?service_id=%s&host=%s&port=%d",
		rc.addr,
		serviceId,
		host,
		port,
	))
	if err != nil {
		return err
	}
	defer httpRes.Body.Close()

	body, err := ioutil.ReadAll(httpRes.Body)
	if err != nil {
		return err
	}

	var j registry.RegisterRes

	err = json.Unmarshal(body, &j)
	if err != nil {
		return err
	}

	if j.Code != 0 {
		return errors.New(j.Msg)
	}

	go func() {
		errNum := 0

		for {
			if errNum == 0 {
				time.Sleep(rc.updateInterval)
			} else if errNum <= 6 {
				time.Sleep(time.Second * 3)
			}
			select {
			case <-rc.stopChan:
				return
			default:
				err = rc.Update(serviceId)
				if err != nil {
					errNum++
				}
				errNum = 0
			}
		}
	}()

	return nil
}

func (rc *BuiltinKitexRegistryClient) Deregister(info *kitexregistry.Info) error {
	serviceId, ok := info.Tags["service_id"]
	if !ok {
		return errors.New("service_id not found")
	}

	rc.stopChan <- struct{}{}

	httpRes, err := http.Get(fmt.Sprintf("http://%s/api/registry/deregister?service_id=%s", rc.addr, serviceId))
	if err != nil {
		return err
	}
	defer httpRes.Body.Close()

	body, err := ioutil.ReadAll(httpRes.Body)
	if err != nil {
		return err
	}

	var j registry.RegisterRes

	err = json.Unmarshal(body, &j)
	if err != nil {
		return err
	}

	if j.Code != 0 {
		return errors.New(j.Msg)
	}

	return nil
}

func (rc *BuiltinKitexRegistryClient) Update(serviceId string) error {

	httpRes, err := http.Get(fmt.Sprintf("http://%s/api/registry/update?service_id=%s", rc.addr, serviceId))
	if err != nil {
		return err
	}
	defer httpRes.Body.Close()

	body, err := ioutil.ReadAll(httpRes.Body)
	if err != nil {
		return err
	}

	var j registry.RegisterRes

	err = json.Unmarshal(body, &j)
	if err != nil {
		return err
	}

	if j.Code != 0 {
		return errors.New(j.Msg)
	}

	return nil
}
