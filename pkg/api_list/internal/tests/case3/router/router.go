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

package router

import (
	"main/router/user"

	"github.com/cloudwego/hertz/pkg/route"
)

func InitRoutes(e *route.Engine) {
	e.GET("/ping", nil)

	g := e.Group("/api/v1")

	initDefault(g)
	user.InitUserRoutes(g.Group("/user"))
}

func initDefault(g *route.RouterGroup) {
	g.GET("/help", nil)
}
