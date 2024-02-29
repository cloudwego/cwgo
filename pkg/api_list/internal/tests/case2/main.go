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

package case2

import "github.com/cloudwego/hertz/pkg/app/server"

func main() {
	h := server.Default()
	h.GET("/server/get", nil)
	h.POST("/server/post", nil)
	h.PUT("/server/put", nil)
	h.DELETE("/server/delete", nil)
	h.HEAD("/server/head", nil)
	h.PATCH("/server/patch", nil)
	h.OPTIONS("/server/options", nil)
	h.GETEX("/server/getex", nil, "")
	h.POSTEX("/server/postex", nil, "")
	h.PUTEX("/server/putex", nil, "")
	h.DELETEEX("/server/deleteex", nil, "")
	h.HEADEX("/server/headex", nil, "")
	h.AnyEX("/server/anyex", nil, "")

	e := h.Engine
	e.GET("/engine/get", nil)

	g1 := h.Group("/g1")
	g1.GET("/get", nil)

	g2 := e.Group("/g2")
	g2.GET("/get", nil)

	g3 := h.Group("/g3")
	g31 := g3.Group("/g1")
	g31.GET("/get", nil)

	g31.POST(
		"/post",
		nil,
	)
}
