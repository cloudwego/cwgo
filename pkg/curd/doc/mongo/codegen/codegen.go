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
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"strings"

	"github.com/cloudwego/cwgo/pkg/curd/code"
	"github.com/cloudwego/cwgo/pkg/curd/extract"
	"github.com/cloudwego/cwgo/pkg/curd/parse"
	"github.com/cloudwego/cwgo/pkg/curd/template"

	"golang.org/x/tools/go/ast/astutil"
)

func HandleCodegen(ifOperations []*parse.InterfaceOperation) (methodRenders [][]*template.MethodRender) {
	for _, ifOperation := range ifOperations {
		methods := make([]*template.MethodRender, 0)
		for _, operation := range ifOperation.Operations {
			switch operation.GetOperationName() {
			case parse.Insert:
				insert := operation.(*parse.InsertParse)
				method := &template.MethodRender{
					Name: insert.BelongedToMethod.Name,
					MethodReceiver: code.MethodReceiver{
						Name: "r",
						Type: code.StarExprType{
							RealType: code.IdentType(ifOperation.BelongedToStruct.Name + "RepositoryMongo"),
						},
					},
					Params:     insert.BelongedToMethod.Params,
					Returns:    insert.BelongedToMethod.Returns,
					MethodBody: insertCodegen(insert),
				}
				methods = append(methods, method)

			case parse.Find:
				find := operation.(*parse.FindParse)
				method := &template.MethodRender{
					Name: find.BelongedToMethod.Name,
					MethodReceiver: code.MethodReceiver{
						Name: "r",
						Type: code.StarExprType{
							RealType: code.IdentType(ifOperation.BelongedToStruct.Name + "RepositoryMongo"),
						},
					},
					Params:     find.BelongedToMethod.Params,
					Returns:    find.BelongedToMethod.Returns,
					MethodBody: findCodegen(find),
				}
				methods = append(methods, method)

			case parse.Update:
				update := operation.(*parse.UpdateParse)
				method := &template.MethodRender{
					Name: update.BelongedToMethod.Name,
					MethodReceiver: code.MethodReceiver{
						Name: "r",
						Type: code.StarExprType{
							RealType: code.IdentType(ifOperation.BelongedToStruct.Name + "RepositoryMongo"),
						},
					},
					Params:     update.BelongedToMethod.Params,
					Returns:    update.BelongedToMethod.Returns,
					MethodBody: updateCodegen(update),
				}
				methods = append(methods, method)

			case parse.Delete:
				del := operation.(*parse.DeleteParse)
				method := &template.MethodRender{
					Name: del.BelongedToMethod.Name,
					MethodReceiver: code.MethodReceiver{
						Name: "r",
						Type: code.StarExprType{
							RealType: code.IdentType(ifOperation.BelongedToStruct.Name + "RepositoryMongo"),
						},
					},
					Params:     del.BelongedToMethod.Params,
					Returns:    del.BelongedToMethod.Returns,
					MethodBody: deleteCodegen(del),
				}
				methods = append(methods, method)

			case parse.Count:
				count := operation.(*parse.CountParse)
				method := &template.MethodRender{
					Name: count.BelongedToMethod.Name,
					MethodReceiver: code.MethodReceiver{
						Name: "r",
						Type: code.StarExprType{
							RealType: code.IdentType(ifOperation.BelongedToStruct.Name + "RepositoryMongo"),
						},
					},
					Params:     count.BelongedToMethod.Params,
					Returns:    count.BelongedToMethod.Returns,
					MethodBody: countCodegen(count),
				}
				methods = append(methods, method)

			case parse.Bulk:
				bulk := operation.(*parse.BulkParse)
				method := &template.MethodRender{
					Name: bulk.BelongedToMethod.Name,
					MethodReceiver: code.MethodReceiver{
						Name: "r",
						Type: code.StarExprType{
							RealType: code.IdentType(ifOperation.BelongedToStruct.Name + "RepositoryMongo"),
						},
					},
					Params:     bulk.BelongedToMethod.Params,
					Returns:    bulk.BelongedToMethod.Returns,
					MethodBody: bulkCodegen(bulk),
				}
				methods = append(methods, method)

			case parse.Transaction:
				ta := operation.(*parse.TransactionParse)
				method := &template.MethodRender{
					Name: ta.BelongedToMethod.Name,
					MethodReceiver: code.MethodReceiver{
						Name: "r",
						Type: code.StarExprType{
							RealType: code.IdentType(ifOperation.BelongedToStruct.Name + "RepositoryMongo"),
						},
					},
					Params:     ta.BelongedToMethod.Params,
					Returns:    ta.BelongedToMethod.Returns,
					MethodBody: taCodegen(ta),
				}
				methods = append(methods, method)

			default:
			}
		}
		methodRenders = append(methodRenders, methods)
	}
	return
}

var BaseMongoImports = map[string]string{
	"context": "",
}

func AddMongoImports(data string) (string, error) {
	fSet := token.NewFileSet()
	file, err := parser.ParseFile(fSet, "", data, parser.ParseComments)
	if err != nil {
		return "", err
	}

	flagBson, flagMongo, flagOption := false, false, false
	ast.Inspect(file, func(n ast.Node) bool {
		if importSpec, ok := n.(*ast.ImportSpec); ok && importSpec.Path.Value == "go.mongodb.org/mongo-driver/bson" {
			flagBson = true
			return false
		}
		if importSpec, ok := n.(*ast.ImportSpec); ok && importSpec.Path.Value == "go.mongodb.org/mongo-driver/mongo" {
			flagMongo = true
			return false
		}
		if importSpec, ok := n.(*ast.ImportSpec); ok && importSpec.Path.Value == "go.mongodb.org/mongo-driver/mongo/options" {
			flagOption = true
			return false
		}
		return true
	})

	if strings.Contains(data, "bson") {
		if !flagBson {
			astutil.AddNamedImport(fSet, file, "", "go.mongodb.org/mongo-driver/bson")
		}
	}
	if strings.Contains(data, "mongo") {
		if !flagMongo {
			astutil.AddNamedImport(fSet, file, "", "go.mongodb.org/mongo-driver/mongo")
		}
	}
	if strings.Contains(data, "options") {
		if !flagOption {
			astutil.AddNamedImport(fSet, file, "", "go.mongodb.org/mongo-driver/mongo/options")
		}
	}

	buf := new(bytes.Buffer)
	if err = printer.Fprint(buf, fSet, file); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func GetFuncRender(extractStruct *extract.IdlExtractStruct) *template.FuncRender {
	return &template.FuncRender{
		Name: "New" + extractStruct.Name + "Repository",
		Params: code.Params{
			code.Param{
				Name: "collection",
				Type: code.StarExprType{
					RealType: code.SelectorExprType{
						X:   "mongo",
						Sel: "Collection",
					},
				},
			},
		},
		Returns: code.Returns{
			code.IdentType(extractStruct.Name + "Repository"),
		},
		FuncBody: code.Body{
			code.RawStmt("return &" + extractStruct.Name + "RepositoryMongo{\n\tcollection: collection,\n}"),
		},
	}
}

func GetStructRender(extractStruct *extract.IdlExtractStruct) *template.StructRender {
	return &template.StructRender{
		Name: extractStruct.Name + "RepositoryMongo",
		StructFields: code.StructFields{
			code.StructField{
				Name: "collection",
				Type: code.StarExprType{
					RealType: code.SelectorExprType{
						X:   "mongo",
						Sel: "Collection",
					},
				},
			},
		},
	}
}
