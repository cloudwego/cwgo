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
	service2 "github.com/cloudwego/cwgo/platform/server/shared/service"
	"sync"
	"time"
)

type BuiltinRegistry struct {
	sync.Mutex
	agents map[string]*service2.BuiltinService
}

var _ IRegistry = (*BuiltinRegistry)(nil)

func (r *BuiltinRegistry) Register(id string, address string, port int) error {
	r.Lock()
	defer r.Unlock()
	// TODO: 连接客户端
	r.agents[id] = &service2.BuiltinService{
		Id:             id,
		LastUpdateTime: time.Now(),
	}
	return nil
}

func (r *BuiltinRegistry) Unregister(id string) error {
	return nil
}

func (r *BuiltinRegistry) Update(id string) error {
	return nil
}

func (r *BuiltinRegistry) CleanUp() error {
	return nil
}

func (r *BuiltinRegistry) Count() int {
	return 0
}

func (r *BuiltinRegistry) GetServiceById(string) service2.IService {
	return nil
}
