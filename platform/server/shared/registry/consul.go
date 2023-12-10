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
	"github.com/cloudwego/cwgo/platform/server/shared/service"
)

type ConsulRegistry struct{}

var _ IRegistry = (*ConsulRegistry)(nil)

func NewConsulRegistry() *ConsulRegistry {
	return &ConsulRegistry{}
}

func (r *ConsulRegistry) GetAllService() ([]*service.Service, error) {
	// TODO implement me
	panic("implement me")
}

func (r *ConsulRegistry) Count() int {
	// TODO implement me
	panic("implement me")
}

func (r *ConsulRegistry) GetServiceById(s string) (*service.Service, error) {
	// TODO implement me
	panic("implement me")
}

func (r *ConsulRegistry) ServiceExists(s string) bool {
	// TODO implement me
	panic("implement me")
}
