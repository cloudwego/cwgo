/*
 * Copyright 2024 CloudWeGo Authors
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

package codegen

import (
	"bytes"
	"go/parser"
	"go/printer"
	"go/token"

	"github.com/cloudwego/cwgo/pkg/curd/code"
	"github.com/cloudwego/cwgo/pkg/curd/template"

	"golang.org/x/tools/go/ast/astutil"
)

func AddBaseMGoImports(data string) (string, error) {
	fSet := token.NewFileSet()
	file, err := parser.ParseFile(fSet, "", data, parser.ParseComments)
	if err != nil {
		return "", err
	}

	astutil.AddImport(fSet, file, "strings")
	astutil.AddImport(fSet, file, "fmt")
	buf := new(bytes.Buffer)
	if err = printer.Fprint(buf, fSet, file); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func HandleBaseCodegen() []*template.MethodRender {
	var methods []*template.MethodRender
	methods = append(methods, findOneMethod())
	methods = append(methods, findListMethod())
	methods = append(methods, findPageListMethod())
	methods = append(methods, findSortPageListMethod())
	methods = append(methods, insertOneMethod())
	methods = append(methods, updateOneMethod())
	methods = append(methods, updateManyMethod())
	methods = append(methods, deleteOneMethod())
	methods = append(methods, bulkInsertMethod())
	methods = append(methods, bulkUpdateMethod())
	methods = append(methods, aggregateMethod())
	methods = append(methods, countMethod())

	return methods
}

func countMethod() *template.MethodRender {
	return &template.MethodRender{
		Name: "MCount",
		MethodReceiver: code.MethodReceiver{
			Name: "b",
			Type: code.StarExprType{
				RealType: code.IdentType("BaseRepositoryMongo"),
			},
		},
		Params: GetMCountParams(),
		Returns: code.Returns{
			code.IdentType("(int64, error)"),
		},
		MethodBody: countBaseCodegen(),
	}
}

func aggregateMethod() *template.MethodRender {
	return &template.MethodRender{
		Name: "MAggregate",
		MethodReceiver: code.MethodReceiver{
			Name: "b",
			Type: code.StarExprType{
				RealType: code.IdentType("BaseRepositoryMongo"),
			},
		},
		Params: GetMAggregateParams(),
		Returns: code.Returns{
			code.IdentType("error"),
		},
		MethodBody: aggregateBaseCodegen(),
	}
}

func bulkUpdateMethod() *template.MethodRender {
	return &template.MethodRender{
		Name: "MBulkUpdate",
		MethodReceiver: code.MethodReceiver{
			Name: "b",
			Type: code.StarExprType{
				RealType: code.IdentType("BaseRepositoryMongo"),
			},
		},
		Params: GetMBulkUpdateParams(),
		Returns: code.Returns{
			code.IdentType("(*mongo.BulkWriteResult, error)"),
		},
		MethodBody: bulkUpdateBaseCodegen(),
	}
}

func bulkInsertMethod() *template.MethodRender {
	return &template.MethodRender{
		Name: "MBulkInsert",
		MethodReceiver: code.MethodReceiver{
			Name: "b",
			Type: code.StarExprType{
				RealType: code.IdentType("BaseRepositoryMongo"),
			},
		},
		Params: GetMBulkInsertParams(),
		Returns: code.Returns{
			code.IdentType("(*mongo.BulkWriteResult, error)"),
		},
		MethodBody: bulkInsertBaseCodegen(),
	}
}

func updateManyMethod() *template.MethodRender {
	return &template.MethodRender{
		Name: "MUpdateMany",
		MethodReceiver: code.MethodReceiver{
			Name: "b",
			Type: code.StarExprType{
				RealType: code.IdentType("BaseRepositoryMongo"),
			},
		},
		Params: GetMUpdateManyParams(),
		Returns: code.Returns{
			code.IdentType("(*mongo.UpdateResult, error)"),
		},
		MethodBody: updateManyBaseCodegen(),
	}
}

func deleteOneMethod() *template.MethodRender {
	return &template.MethodRender{
		Name: "MDeleteOne",
		MethodReceiver: code.MethodReceiver{
			Name: "b",
			Type: code.StarExprType{
				RealType: code.IdentType("BaseRepositoryMongo"),
			},
		},
		Params: GetMDeleteOneParams(),
		Returns: code.Returns{
			code.IdentType("(*mongo.DeleteResult, error)"),
		},
		MethodBody: deleteOneBaseCodegen(),
	}
}

func updateOneMethod() *template.MethodRender {
	return &template.MethodRender{
		Name: "MUpdateOne",
		MethodReceiver: code.MethodReceiver{
			Name: "b",
			Type: code.StarExprType{
				RealType: code.IdentType("BaseRepositoryMongo"),
			},
		},
		Params: GetMUpdateOneParams(),
		Returns: code.Returns{
			code.IdentType("(*mongo.UpdateResult, error)"),
		},
		MethodBody: updateOneBaseCodegen(),
	}
}

func insertOneMethod() *template.MethodRender {
	return &template.MethodRender{
		Name: "MInsertOne",
		MethodReceiver: code.MethodReceiver{
			Name: "b",
			Type: code.StarExprType{
				RealType: code.IdentType("BaseRepositoryMongo"),
			},
		},
		Params: GetMInsertOneParams(),
		Returns: code.Returns{
			code.IdentType("(*mongo.InsertOneResult, error)"),
		},
		MethodBody: insertOneBaseCodegen(),
	}
}

func findSortPageListMethod() *template.MethodRender {
	return &template.MethodRender{
		Name: "MFindSortPageList",
		MethodReceiver: code.MethodReceiver{
			Name: "b",
			Type: code.StarExprType{
				RealType: code.IdentType("BaseRepositoryMongo"),
			},
		},
		Params: GetMFindSortPageListParams(),
		Returns: code.Returns{
			code.IdentType("error"),
		},
		MethodBody: findSortPageListBaseCodegen(),
	}
}

func findPageListMethod() *template.MethodRender {
	return &template.MethodRender{
		Name: "MFindPageList",
		MethodReceiver: code.MethodReceiver{
			Name: "b",
			Type: code.StarExprType{
				RealType: code.IdentType("BaseRepositoryMongo"),
			},
		},
		Params: GetMFindPageListParams(),
		Returns: code.Returns{
			code.IdentType("error"),
		},
		MethodBody: findPageListBaseCodegen(),
	}
}

func findListMethod() *template.MethodRender {
	return &template.MethodRender{
		Name: "MFindList",
		MethodReceiver: code.MethodReceiver{
			Name: "b",
			Type: code.StarExprType{
				RealType: code.IdentType("BaseRepositoryMongo"),
			},
		},
		Params: GetMFindListParams(),
		Returns: code.Returns{
			code.IdentType("error"),
		},
		MethodBody: findListBaseCodegen(),
	}
}

func findOneMethod() *template.MethodRender {
	return &template.MethodRender{
		Name: "MFindOne",
		MethodReceiver: code.MethodReceiver{
			Name: "b",
			Type: code.StarExprType{
				RealType: code.IdentType("BaseRepositoryMongo"),
			},
		},
		Params: GetMFindOneParams(),
		Returns: code.Returns{
			code.IdentType("error"),
		},
		MethodBody: findBaseCodegen(),
	}
}
