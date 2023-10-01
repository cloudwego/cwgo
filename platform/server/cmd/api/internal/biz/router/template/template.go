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

package template

import (
	"github.com/cloudwego/hertz/pkg/app/server"

	"github.com/cloudwego/cwgo/platform/server/cmd/api/internal/biz/handler/template"
)

/*
 This file will register all the routes of the services in the master idl.
 And it will update automatically when you use the "update" command for the idl.
 So don't modify the contents of the file, or your code will be deleted when it is updated.
*/

// Register register routes based on the IDL 'api.${HTTP Method}' annotation.
func Register(r *server.Hertz) {

	root := r.Group("/", rootMw()...)
	root.POST("/template", append(_templateMw(), template.AddTemplate)...)
	_template := root.Group("/template", _templateMw()...)
	_template.POST("/item", append(_addtemplateitemMw(), template.AddTemplateItem)...)
	_template.DELETE("/item", append(_deletetemplateitemMw(), template.DeleteTemplateItem)...)
	_template.PATCH("/item", append(_updatetemplateitemMw(), template.UpdateTemplateItem)...)
	_template.GET("/item", append(_gettemplateitemsMw(), template.GetTemplateItems)...)
	root.DELETE("/template", append(_deletetemplateMw(), template.DeleteTemplate)...)
	root.PATCH("/template", append(_updatetemplateMw(), template.UpdateTemplate)...)
	root.GET("/template", append(_gettemplatesMw(), template.GetTemplates)...)
}
