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

package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
)

// customizeRegister registers customize routers.
func customizedRegister(r *server.Hertz, staticFilePath string) {
	r.LoadHTMLGlob(staticFilePath + "/index.html")
	fs := &app.FS{Root: staticFilePath}
	r.StaticFS("/assets", fs)
	r.GET("/favicon.ico", fs.NewRequestHandler())
	r.GET("/*static", func(c context.Context, ctx *app.RequestContext) {
		ctx.HTML(http.StatusOK, "index.html", nil)
	})
}

func getPathRewriter(prefix string) app.PathRewriteFunc {
	// cannot have an empty prefix
	if prefix == "" {
		prefix = "/"
	}
	// prefix always start with a '/' or '*'
	if prefix[0] != '/' {
		prefix = "/" + prefix
	}

	// is prefix a direct wildcard
	isStar := prefix == "/*"
	// is prefix a partial wildcard
	if strings.Contains(prefix, "*") {
		isStar = true
		prefix = strings.Split(prefix, "*")[0]
		// fix this later
	}
	prefixLen := len(prefix)
	if prefixLen > 1 && prefix[prefixLen-1:] == "/" {
		// /john/ -> /john
		prefixLen--
		prefix = prefix[:prefixLen]
	}
	return func(ctx *app.RequestContext) []byte {
		path := ctx.Path()
		if len(path) >= prefixLen {
			if isStar && string(path[0:prefixLen]) == prefix {
				path = append(path[0:0], '/')
			} else {
				path = path[prefixLen:]
				if len(path) == 0 || path[len(path)-1] != '/' {
					path = append(path, '/')
				}
			}
		}
		if len(path) > 0 && path[0] != '/' {
			path = append([]byte("/"), path...)
		}
		return path
	}
}
