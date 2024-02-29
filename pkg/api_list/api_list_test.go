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

var (
	results = map[string][]*RouterParsed{
		// project that generated by hz with default template thrift®
		"case1": {
			{
				FilePath:  "biz/router/hello/example/hello.go",
				StartLine: 34,
				EndLine:   34,
				Method:    RouterRegisterFuncNameGET,
				RoutePath: "/hello",
			},
			{
				FilePath:  "router.go",
				StartLine: 27,
				EndLine:   27,
				Method:    RouterRegisterFuncNameGET,
				RoutePath: "/ping",
			},
		},
		"case2": {
			{
				FilePath:  "main.go",
				StartLine: 23,
				EndLine:   23,
				Method:    RouterRegisterFuncNameGET,
				RoutePath: "/server/get",
			},
			{
				FilePath:  "main.go",
				StartLine: 24,
				EndLine:   24,
				Method:    RouterRegisterFuncNamePOST,
				RoutePath: "/server/post",
			},
			{
				FilePath:  "main.go",
				StartLine: 25,
				EndLine:   25,
				Method:    RouterRegisterFuncNamePUT,
				RoutePath: "/server/put",
			},
			{
				FilePath:  "main.go",
				StartLine: 26,
				EndLine:   26,
				Method:    RouterRegisterFuncNameDELETE,
				RoutePath: "/server/delete",
			},
			{
				FilePath:  "main.go",
				StartLine: 27,
				EndLine:   27,
				Method:    RouterRegisterFuncNameHEAD,
				RoutePath: "/server/head",
			},
			{
				FilePath:  "main.go",
				StartLine: 28,
				EndLine:   28,
				Method:    RouterRegisterFuncNamePATCH,
				RoutePath: "/server/patch",
			},
			{
				FilePath:  "main.go",
				StartLine: 29,
				EndLine:   29,
				Method:    RouterRegisterFuncNameOPTIONS,
				RoutePath: "/server/options",
			},
			{
				FilePath:  "main.go",
				StartLine: 30,
				EndLine:   30,
				Method:    RouterRegisterFuncNameGETEX,
				RoutePath: "/server/getex",
			},
			{
				FilePath:  "main.go",
				StartLine: 31,
				EndLine:   31,
				Method:    RouterRegisterFuncNamePOSTEX,
				RoutePath: "/server/postex",
			},
			{
				FilePath:  "main.go",
				StartLine: 32,
				EndLine:   32,
				Method:    RouterRegisterFuncNamePUTEX,
				RoutePath: "/server/putex",
			},
			{
				FilePath:  "main.go",
				StartLine: 33,
				EndLine:   33,
				Method:    RouterRegisterFuncNameDELETEEX,
				RoutePath: "/server/deleteex",
			},
			{
				FilePath:  "main.go",
				StartLine: 34,
				EndLine:   34,
				Method:    RouterRegisterFuncNameHEADEX,
				RoutePath: "/server/headex",
			},
			{
				FilePath:  "main.go",
				StartLine: 35,
				EndLine:   35,
				Method:    RouterRegisterFuncNameAnyEX,
				RoutePath: "/server/anyex",
			},
			{
				FilePath:  "main.go",
				StartLine: 38,
				EndLine:   38,
				Method:    RouterRegisterFuncNameGET,
				RoutePath: "/engine/get",
			},
			{
				FilePath:  "main.go",
				StartLine: 41,
				EndLine:   41,
				Method:    RouterRegisterFuncNameGET,
				RoutePath: "/g1/get",
			},
			{
				FilePath:  "main.go",
				StartLine: 44,
				EndLine:   44,
				Method:    RouterRegisterFuncNameGET,
				RoutePath: "/g2/get",
			},
			{
				FilePath:  "main.go",
				StartLine: 48,
				EndLine:   48,
				Method:    RouterRegisterFuncNameGET,
				RoutePath: "/g3/g1/get",
			},
			{
				FilePath:  "main.go",
				StartLine: 50,
				EndLine:   53,
				Method:    RouterRegisterFuncNamePOST,
				RoutePath: "/g3/g1/post",
			},
		},
		"case3": {
			{
				FilePath:  "router/router.go",
				StartLine: 10,
				EndLine:   10,
				Method:    RouterRegisterFuncNameGET,
				RoutePath: "/ping",
			},
			{
				FilePath:  "router/router.go",
				StartLine: 19,
				EndLine:   19,
				Method:    RouterRegisterFuncNameGET,
				RoutePath: "/api/v1/help",
			},
			{
				FilePath:  "router/user/user.go",
				StartLine: 6,
				EndLine:   6,
				Method:    RouterRegisterFuncNameGET,
				RoutePath: "/api/v1/user/info",
			},
			{
				FilePath:  "router/user/user.go",
				StartLine: 7,
				EndLine:   7,
				Method:    RouterRegisterFuncNamePOST,
				RoutePath: "/api/v1/user/nickname",
			},
		},
	}
)
