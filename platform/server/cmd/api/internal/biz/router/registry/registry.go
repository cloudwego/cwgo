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
	"github.com/cloudwego/hertz/pkg/route"

	registry "github.com/cloudwego/cwgo/platform/server/cmd/api/internal/biz/handler/registry"
)

/*
 This file will register all the routes of the services in the master idl.
 And it will update automatically when you use the "update" command for the idl.
 So don't modify the contents of the file, or your code will be deleted when it is updated.
*/

// Register register routes based on the IDL 'api.${HTTP Method}' annotation.
func Register(r *route.RouterGroup) {
	root := r.Group("/", rootMw()...)
	{
		_api := root.Group("/api", _apiMw()...)
		{
			_registry := _api.Group("/registry", _registryMw()...)
			_registry.GET("/dnregister", append(_deregisterMw(), registry.Deregister)...)
			_registry.GET("/register", append(_registerMw(), registry.Register)...)
			_registry.GET("/update", append(_updateMw(), registry.Update)...)
		}
	}
}
