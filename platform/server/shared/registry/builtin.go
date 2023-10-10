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
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Manager struct {
	windows     []*service.Service
	currentSize int
	expireTime  time.Duration
	mutex       sync.Mutex
}

func (sw *Manager) add(agentService *service.Service) {
	sw.mutex.Lock()
	defer sw.mutex.Unlock()

	if sw.currentSize == cap(sw.windows) {
		// expand slice if full
		var newAgents []*service.Service
		if cap(sw.windows) == 0 {
			newAgents = make([]*service.Service, 16)
		} else {
			newAgents = make([]*service.Service, cap(sw.windows)<<1)
		}
		copy(newAgents, sw.windows)
		sw.windows = newAgents
	}

	sw.windows[sw.currentSize] = agentService
	sw.currentSize++
}

func (sw *Manager) getExpiredServiceIds() []string {
	sw.mutex.Lock()
	defer sw.mutex.Unlock()

	expiredServiceIds := make([]string, 0)
	for i := 0; i < sw.currentSize; i++ {
		agentService := sw.windows[i]
		if agentService.LastUpdateTime.Add(sw.expireTime).Before(time.Now()) {
			expiredServiceIds = append(expiredServiceIds, agentService.Id)
		} else {
			copy(sw.windows, sw.windows[i:])
			sw.currentSize = sw.currentSize - i
			break
		}
	}

	return expiredServiceIds
}

type BuiltinRegistry struct {
	sync.Mutex
	agents        map[string]*service.Service
	cleanInterval time.Duration
	manager       *Manager
}

var _ IRegistry = (*BuiltinRegistry)(nil)

const (
	minCleanInterval = 100 * time.Millisecond
)

func NewBuiltinRegistry() *BuiltinRegistry {
	r := &BuiltinRegistry{
		Mutex:         sync.Mutex{},
		agents:        make(map[string]*service.Service),
		cleanInterval: 3 * time.Second,
		manager: &Manager{
			windows:     make([]*service.Service, 0),
			currentSize: 0,
			mutex:       sync.Mutex{},
			expireTime:  time.Minute,
		},
	}

	go r.StartCleanUp()

	return r
}

func (r *BuiltinRegistry) Register(serviceId string, host string, port int) error {
	r.Lock()
	defer r.Unlock()

	agentService, err := service.NewService(serviceId, host, port)
	if err != nil {
		return err
	}

	r.agents[serviceId] = agentService

	r.manager.add(agentService)

	return nil
}

func (r *BuiltinRegistry) Deregister(id string) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.agents[id]; !ok {
		return errors.New("service not found")
	}

	delete(r.agents, id)

	return nil
}

func (r *BuiltinRegistry) Update(serviceId string) error {
	r.Lock()
	defer r.Unlock()

	if agentService, ok := r.agents[serviceId]; !ok {
		return errors.New("service not found")
	} else {
		agentService.LastUpdateTime = time.Now()
		r.manager.add(agentService)
		return nil
	}
}

func (r *BuiltinRegistry) StartCleanUp() {
	for {
		time.Sleep(r.cleanInterval)

		r.Lock()
		expiredServiceIds := r.manager.getExpiredServiceIds()

		for _, serviceId := range expiredServiceIds {
			if _, ok := r.agents[serviceId]; ok {
				delete(r.agents, serviceId)
			}
		}
		r.Unlock()
	}
}

func (r *BuiltinRegistry) Count() int {
	return len(r.agents)
}

func (r *BuiltinRegistry) GetServiceById(serviceId string) (*service.Service, error) {
	if agentService, ok := r.agents[serviceId]; !ok {
		return nil, errors.New("service not found")
	} else {
		return agentService, nil
	}
}

func (r *BuiltinRegistry) GetAllService() ([]*service.Service, error) {
	r.Lock()
	defer r.Unlock()

	var services []*service.Service
	for _, svr := range r.agents {
		services = append(services, svr)
	}

	return services, nil
}

func (r *BuiltinRegistry) ServiceExists(serviceId string) bool {
	_, ok := r.agents[serviceId]

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
