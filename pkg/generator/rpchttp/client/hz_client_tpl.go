/*
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
 */

package client

import (
	"github.com/cloudwego/cwgo/pkg/consts"
	"github.com/cloudwego/cwgo/pkg/generator/common/template"
)

var hzClientMVCTemplates = []template.Template{
	{
		Path:   consts.InitGo,
		Delims: [2]string{consts.LeftDelimiter, consts.RightDelimiter},
		UpdateBehavior: template.UpdateBehavior{
			AppendRender: map[string]interface{}{},
			Append: template.Append{
				AppendImport: map[string]string{},
			},
			ReplaceFunc: template.ReplaceFunc{
				ReplaceFuncName:         make([]string, 0, 5),
				ReplaceFuncAppendImport: make([]map[string]string, 0, 10),
				ReplaceFuncDeleteImport: make([]map[string]string, 0, 10),
				ReplaceFuncBody:         make([]string, 0, 5),
			},
		},
		Body: `package {{.InitOptsPackage}}
      import (
		{{range $key, $value := .GoFileImports}}
	    {{if eq $key "init.go"}}
	    {{range $k, $v := $value}}
        {{if ne $k ""}}{{if ne $v ""}}{{$v}} "{{$k}}"{{else}}"{{$k}}"{{end}}{{end}}{{end}}{{end}}{{end}}
	  )

	  func initClientOpts(hostUrl string) (ops []Option, err error) {
		ops = append(ops, withHostUrl(hostUrl))
		
		if err = initResolver(&ops); err != nil {
		  panic(err)
		}

		return
	  }
	  
	  // If you do not use the service resolver function, do not edit this function.
	  // Otherwise, you can customize and modify it.
	  func initResolver(ops *[]Option) (err error) {
		{{if ne .ResolverName ""}}
		{{.ResolverBody}}
		{{else}}
		return
        {{end}}
	  }`,
	},

	{
		Path:   consts.DefaultHZClientDir + consts.Slash + consts.EnvGo,
		Delims: [2]string{"[[", "]]"},
		UpdateBehavior: template.UpdateBehavior{
			AppendRender: map[string]interface{}{},
			Append: template.Append{
				AppendImport: map[string]string{},
			},
		},
		CustomFunc: template.CustomFuncMap,
		Body: `// Code generated by cwgo generator. DO NOT EDIT.

	  package http
	  import (
		[[range $key, $value := .GoFileImports]]
	    [[if eq $key "env.go"]]
	    [[range $k, $v := $value]]
        [[if ne $k ""]][[if ne $v ""]][[$v]] "[[$k]]"[[else]]"[[$k]]"[[end]][[end]][[end]][[end]][[end]]
	  )

      [[if ne .ResolverName ""]]
      func GetResolverAddress() []string {
		e := os.Getenv("GO_HERTZ_REGISTRY_[[ToUpper .ServiceName]]")
	    if len(e) == 0 {
		  return []string{[[$lenSlice := len .ResolverAddress]][[range $key, $value := .ResolverAddress]]"[[$value]]"[[if eq $key (Sub $lenSlice 1)]][[else]], [[end]][[end]]}
	    }
		return strings.Fields(e)
      }
	  [[end]]`,
	},
}