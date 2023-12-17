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
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/cloudwego/cwgo/platform/server/shared/consts"
	"github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/registry"
	"github.com/cloudwego/cwgo/platform/server/shared/logger"
	"github.com/cloudwego/cwgo/platform/server/shared/service"
	"github.com/cloudwego/cwgo/platform/server/shared/utils"
	"github.com/cloudwego/kitex/pkg/discovery"
	kitexregistry "github.com/cloudwego/kitex/pkg/registry"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/patrickmn/go-cache"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var ErrServiceNotFound = errors.New("service not found")

type BuiltinRegistry struct {
	agentCache   *cache.Cache
	rdb          redis.UniversalClient
	syncDuration time.Duration
}

var _ IRegistry = (*BuiltinRegistry)(nil)

const (
	minCleanInterval     = 100 * time.Millisecond
	defaultCleanInternal = 3 * time.Second
	defaultExpiration    = 60 * time.Second
	defaultSyncDuration  = 30 * time.Second
)

func NewBuiltinRegistry(rdb redis.UniversalClient) *BuiltinRegistry {
	r := &BuiltinRegistry{
		agentCache:   cache.New(defaultExpiration, defaultCleanInternal),
		rdb:          rdb,
		syncDuration: defaultSyncDuration,
	}

	go r.StartSync()

	return r
}

func (r *BuiltinRegistry) StartSync() {
	for {
		time.Sleep(r.syncDuration) // sync every r.syncDuration

		logger.Logger.Debug("start sync agent services")

		// get agent services info in redis
		switch rdb := r.rdb.(type) {
		case *redis.ClusterClient:
			err := rdb.ForEachMaster(context.Background(), func(ctx context.Context, client *redis.Client) error {
				err := r.scanServiceKeysAndUpdate(client)
				if err != nil {
					return err
				}

				return nil
			})
			if err != nil {
				continue
			}
		case *redis.Client:
			err := r.scanServiceKeysAndUpdate(rdb)
			if err != nil {
				continue
			}
		}

		logger.Logger.Debug("sync agent services finished")
	}
}

func (r *BuiltinRegistry) scanServiceKeysAndUpdate(rdb *redis.Client) error {
	var err error
	var keys []string
	var needUpdateServiceId []string
	cursor := uint64(0)

	for {
		logger.Logger.Debug("scanning keys store agent services in redis", zap.Uint64("cursor", cursor))

		// scan service keys in redis
		keys, cursor, err = rdb.Scan(context.Background(), cursor, consts.RdbKeyRegistryService+"*", 100).Result()
		if err != nil {
			logger.Logger.Error("scanning keys store agent services in redis failed", zap.Error(err))
			break
		}

		for _, key := range keys {
			serviceId := strings.TrimLeft(key, consts.RdbKeyRegistryService)
			// check where service in builtin registry
			isExist := r.ServiceExists(key)
			if !isExist {
				// if not exist then need to update
				needUpdateServiceId = append(needUpdateServiceId, serviceId)
			}
		}

		if cursor == 0 {
			break
		}
	}

	if len(needUpdateServiceId) != 0 {
		// get service info (ip and port)
		pipe := rdb.Pipeline()
		for _, key := range needUpdateServiceId {
			pipe.Get(context.Background(), key)
		}

		logger.Logger.Debug("exec redis pipe command")
		cmds, err := pipe.Exec(context.Background())
		if err != nil {
			logger.Logger.Error("exec redis pipe command failed", zap.Error(err))
			return err
		}

		for i, cmd := range cmds {
			ipPort, err := cmd.(*redis.StringCmd).Result()
			if err != nil {
				return err
			}

			ip, port, err := net.SplitHostPort(ipPort)
			if err != nil {
				logger.Logger.Error("parse host port string failed", zap.Error(err), zap.String("val", ipPort))
			}

			p, err := strconv.Atoi(port)
			if err != nil {
				logger.Logger.Error("parse port failed", zap.Error(err), zap.String("val", port))
			}

			// register service in builtin registry
			err = r.Register(needUpdateServiceId[i], ip, p)
			if err != nil {
				logger.Logger.Error(
					"register services get by sync progress in redis failed",
					zap.Error(err),
					zap.String("service_id", needUpdateServiceId[i]),
					zap.String("ip", ip),
					zap.Int("port", p),
				)
			}
		}
	}

	return nil
}

func (r *BuiltinRegistry) Register(serviceId, host string, port int) error {
	agentService, err := service.NewService(serviceId, host, port)
	if err != nil {
		return err
	}

	r.agentCache.SetDefault(serviceId, agentService) // register service in cache

	// save service info in redis
	err = r.rdb.Set(
		context.Background(),
		fmt.Sprintf(consts.RdbKeyRegistryService, serviceId),
		fmt.Sprintf("%s:%d", host, port),
		defaultExpiration,
	).Err()
	if err != nil {
		logger.Logger.Error(
			"register service in builtin registry failed",
			zap.Error(err),
			zap.String("service_id", serviceId),
			zap.String("host", host),
			zap.Int("port", port),
		)
	}

	logger.Logger.Debug("registered service in builtin registry",
		zap.String("service_id", serviceId),
		zap.String("host", host),
		zap.Int("port", port),
	)

	return nil
}

func (r *BuiltinRegistry) Deregister(serviceId string) error {
	r.agentCache.Delete(serviceId) // deregister service in cache

	// del service info in redis
	err := r.rdb.Del(
		context.Background(),
		fmt.Sprintf(consts.RdbKeyRegistryService, serviceId),
	).Err()
	if err != nil {
		logger.Logger.Error(
			"deregister service in builtin registry failed",
			zap.Error(err),
			zap.String("service_id", serviceId),
		)
	}

	logger.Logger.Debug("deregistered service in builtin registry",
		zap.String("service_id", serviceId),
	)

	return nil
}

func (r *BuiltinRegistry) Update(serviceId string) error {
	v, ok := r.agentCache.Get(serviceId)
	if !ok {
		return ErrServiceNotFound
	}

	r.agentCache.SetDefault(serviceId, v.(*service.Service)) // update service in cache

	// update service in redis
	err := r.rdb.Expire(
		context.Background(),
		fmt.Sprintf(consts.RdbKeyRegistryService, serviceId),
		defaultExpiration,
	).Err()
	if err != nil {
		if err == redis.Nil {
			return errors.New("service not found")
		}
		logger.Logger.Error(
			"update service in builtin registry failed",
			zap.Error(err),
			zap.String("service_id", serviceId),
		)
	}

	logger.Logger.Debug("update service in builtin registry",
		zap.String("service_id", serviceId),
	)

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

func (rc *BuiltinKitexRegistryClient) registry(serviceId, host string, port int) error {
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

	return nil
}

func (rc *BuiltinKitexRegistryClient) Register(info *kitexregistry.Info) error {
	serviceId, ok := info.Tags["service_id"]
	if !ok {
		return ErrServiceNotFound
	}

	host, port, _ := utils.ParseAddr(info.Addr)

	err := rc.registry(serviceId, host, port)
	if err != nil {
		return err
	}

	logger.Logger.Debug("start update service in registry")
	go func() {
		errNum := 0

		for {
			if errNum == 0 {
				time.Sleep(rc.updateInterval)
			} else if errNum <= 6 {
				time.Sleep(time.Duration(int(math.Pow(3, float64(errNum)))) * time.Second)
			} else {
				logger.Logger.Fatal("update service failed more than 6 times, connect to registry fail, stopping agent service")
			}
			select {
			case <-rc.stopChan:
				logger.Logger.Debug("stop update service to registry")
				return
			default:
				logger.Logger.Debug("updating service to registry")
				err = rc.Update(serviceId)
				if err != nil {
					if err == ErrServiceNotFound {
						err = rc.registry(serviceId, host, port)
						if err != nil {
							errNum++
						} else {
							errNum = 0
						}
					} else {
						errNum++
					}
				} else {
					errNum = 0
				}
				logger.Logger.Debug("update service to registry successfully")
			}
		}
	}()

	return nil
}

func (rc *BuiltinKitexRegistryClient) Deregister(info *kitexregistry.Info) error {
	serviceId, ok := info.Tags["service_id"]
	if !ok {
		return ErrServiceNotFound
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
		if j.Msg == ErrServiceNotFound.Error() {
			return ErrServiceNotFound
		}
		return errors.New(j.Msg)
	}

	return nil
}
