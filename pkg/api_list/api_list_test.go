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

package api_list

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/cloudwego/cwgo/pkg/consts"
)

func trimExecutedPath(path string, routers []*RouterParsed) {
	for _, router := range routers {
		router.FilePath, _ = filepath.Rel(path, router.FilePath)
	}
}

func TestApi(t *testing.T) {
	results := map[string][]*RouterParsed{
		// project that generated by hz with default template thrift
		"case1": {
			{
				FilePath:  "biz/router/hello/example/hello.go",
				StartLine: 19,
				EndLine:   19,
				Method:    RouterRegisterFuncNameGET,
				RoutePath: "/hello",
			},
			{
				FilePath:  "router.go",
				StartLine: 12,
				EndLine:   12,
				Method:    RouterRegisterFuncNameGET,
				RoutePath: "/ping",
			},
		},
		"case2": {
			{
				FilePath:  "main.go",
				StartLine: 7,
				EndLine:   7,
				Method:    RouterRegisterFuncNameGET,
				RoutePath: "/server/get",
			},
			{
				FilePath:  "main.go",
				StartLine: 8,
				EndLine:   8,
				Method:    RouterRegisterFuncNamePOST,
				RoutePath: "/server/post",
			},
			{
				FilePath:  "main.go",
				StartLine: 9,
				EndLine:   9,
				Method:    RouterRegisterFuncNamePUT,
				RoutePath: "/server/put",
			},
			{
				FilePath:  "main.go",
				StartLine: 10,
				EndLine:   10,
				Method:    RouterRegisterFuncNameDELETE,
				RoutePath: "/server/delete",
			},
			{
				FilePath:  "main.go",
				StartLine: 11,
				EndLine:   11,
				Method:    RouterRegisterFuncNameHEAD,
				RoutePath: "/server/head",
			},
			{
				FilePath:  "main.go",
				StartLine: 12,
				EndLine:   12,
				Method:    RouterRegisterFuncNamePATCH,
				RoutePath: "/server/patch",
			},
			{
				FilePath:  "main.go",
				StartLine: 13,
				EndLine:   13,
				Method:    RouterRegisterFuncNameOPTIONS,
				RoutePath: "/server/options",
			},
			{
				FilePath:  "main.go",
				StartLine: 14,
				EndLine:   14,
				Method:    RouterRegisterFuncNameGETEX,
				RoutePath: "/server/getex",
			},
			{
				FilePath:  "main.go",
				StartLine: 15,
				EndLine:   15,
				Method:    RouterRegisterFuncNamePOSTEX,
				RoutePath: "/server/postex",
			},
			{
				FilePath:  "main.go",
				StartLine: 16,
				EndLine:   16,
				Method:    RouterRegisterFuncNamePUTEX,
				RoutePath: "/server/putex",
			},
			{
				FilePath:  "main.go",
				StartLine: 17,
				EndLine:   17,
				Method:    RouterRegisterFuncNameDELETEEX,
				RoutePath: "/server/deleteex",
			},
			{
				FilePath:  "main.go",
				StartLine: 18,
				EndLine:   18,
				Method:    RouterRegisterFuncNameHEADEX,
				RoutePath: "/server/headex",
			},
			{
				FilePath:  "main.go",
				StartLine: 19,
				EndLine:   19,
				Method:    RouterRegisterFuncNameAnyEX,
				RoutePath: "/server/anyex",
			},
			{
				FilePath:  "main.go",
				StartLine: 22,
				EndLine:   22,
				Method:    RouterRegisterFuncNameGET,
				RoutePath: "/engine/get",
			},
			{
				FilePath:  "main.go",
				StartLine: 25,
				EndLine:   25,
				Method:    RouterRegisterFuncNameGET,
				RoutePath: "/g1/get",
			},
			{
				FilePath:  "main.go",
				StartLine: 28,
				EndLine:   28,
				Method:    RouterRegisterFuncNameGET,
				RoutePath: "/g2/get",
			},
			{
				FilePath:  "main.go",
				StartLine: 32,
				EndLine:   32,
				Method:    RouterRegisterFuncNameGET,
				RoutePath: "/g3/g1/get",
			},
			{
				FilePath:  "main.go",
				StartLine: 34,
				EndLine:   37,
				Method:    RouterRegisterFuncNamePOST,
				RoutePath: "/g3/g1/post",
			},
		},
	}

	root := "internal/tests"
	dirEntries, err := os.ReadDir(root)
	if err != nil {
		t.Fatal(err)
	}

	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() {
			dirEntryPath := filepath.Join(root, dirEntry.Name())
			t.Run(dirEntry.Name(), func(t *testing.T) {
				parser, err := NewParser(dirEntryPath, consts.HertzRepoDefaultUrl)
				if err != nil {
					t.Fatal(err)
				}

				moduleName, err := getModuleName(dirEntryPath)
				if err != nil {
					t.Fatal(err)
				}

				err = parser.searchFunc(moduleName, "main", make(map[string]*Var), nil)
				if err != nil {
					t.Fatal(err)
				}

				trimExecutedPath(dirEntryPath, parser.routerParsedList)

				if !reflect.DeepEqual(parser.routerParsedList, results[dirEntry.Name()]) {
					t.Errorf("expected: %v, got: %v", results[dirEntry.Name()], parser.routerParsedList)
				}
			})
		}
	}
}
