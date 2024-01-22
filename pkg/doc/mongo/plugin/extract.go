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

package plugin

import (
	"go/ast"
	astParser "go/parser"
	"go/token"
	"path/filepath"
	"strings"

	"github.com/cloudwego/cwgo/pkg/doc/mongo/plugin/model"

	"github.com/cloudwego/cwgo/pkg/doc/mongo/code"
	"github.com/fatih/camelcase"
)

func extractIdlInterface(rawInterface string, rawStruc *model.IdlExtractStruct, tokens []string) error {
	fSet := token.NewFileSet()
	f, err := astParser.ParseFile(fSet, "", rawInterface, astParser.ParseComments)
	if err != nil {
		return err
	}

	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}

		for _, spec := range genDecl.Specs {
			switch spec := spec.(type) {
			case *ast.TypeSpec:
				switch t := spec.Type.(type) {
				case *ast.InterfaceType:
					rawStruc.InterfaceInfo = extractInterfaceType(spec.Name.Name, t, tokens, rawStruc)
				}
			}
		}
	}

	return nil
}

func extractInterfaceType(name string, interfaceType *ast.InterfaceType, tokens []string,
	rawStruc *model.IdlExtractStruct,
) *model.InterfaceInfo {
	intf := &model.InterfaceInfo{
		Name:    name,
		Methods: []*model.InterfaceMethod{},
	}

	for index, method := range interfaceType.Methods.List {
		funcType, ok := method.Type.(*ast.FuncType)
		if !ok {
			continue
		}

		var nam string
		for _, n := range method.Names {
			nam = n.Name
			break
		}

		if rawStruc.Update {
			if _, ok = rawStruc.PreMethodNamesMap[nam]; !ok {
				meth := extractFunction(nam, funcType, tokens[index])
				meth.BelongedToStruct = rawStruc

				intf.Methods = append(intf.Methods, meth)
			} else {
				meth := extractFunction(nam, funcType, tokens[index])
				meth.BelongedToStruct = rawStruc

				rawStruc.PreIfMethods = append(rawStruc.PreIfMethods, meth)
			}
		} else {
			meth := extractFunction(nam, funcType, tokens[index])
			meth.BelongedToStruct = rawStruc

			intf.Methods = append(intf.Methods, meth)
		}
	}

	return intf
}

func extractFunction(name string, funcType *ast.FuncType, token string) *model.InterfaceMethod {
	meth := &model.InterfaceMethod{
		Name:         name,
		ParsedTokens: token,
	}
	for _, param := range funcType.Params.List {
		paramType := getType(param.Type)

		if len(param.Names) == 0 {
			meth.Params = append(meth.Params, code.Param{Type: paramType})
			continue
		}

		for _, n := range param.Names {
			meth.Params = append(meth.Params, code.Param{
				Name: n.Name,
				Type: paramType,
			})
		}
	}

	if funcType.Results != nil {
		for _, result := range funcType.Results.List {
			meth.Returns = append(meth.Returns, getType(result.Type))
		}
	}

	return meth
}

func getType(expr ast.Expr) code.Type {
	switch expr := expr.(type) {
	case *ast.Ident:
		return code.IdentType(expr.Name)

	case *ast.SelectorExpr:
		xExpr := expr.X.(*ast.Ident)
		return code.SelectorExprType{
			X:   xExpr.Name,
			Sel: expr.Sel.Name,
		}

	case *ast.StarExpr:
		realType := getType(expr.X)
		return code.StarExprType{
			RealType: realType,
		}

	case *ast.ArrayType:
		elementType := getType(expr.Elt)
		return code.SliceType{
			ElementType: elementType,
		}

	case *ast.MapType:
		keyType := getType(expr.Key)
		valueType := getType(expr.Value)
		return code.MapType{KeyType: keyType, ValueType: valueType}

	case *ast.InterfaceType:
		return code.InterfaceType{}
	}

	return nil
}

func getFileName(strucName, prefix string) (fileMongoName, fileIfName string) {
	dir := getPkgName(strucName)
	fileMongoName = filepath.Join(prefix, dir, dir+"_repo_mongo.go")
	fileIfName = filepath.Join(prefix, dir, dir+"_repo.go")
	return
}

func getPkgName(strucName string) string {
	tokens := camelcase.Split(strucName)
	dir := ""
	for i, toke := range tokens {
		if i != len(tokens)-1 {
			dir += strings.ToLower(toke) + "_"
		} else {
			dir += strings.ToLower(toke)
		}
	}
	return dir
}

func getInterfaceMethodNames(data string) (result []string, err error) {
	fSet := token.NewFileSet()
	f, err := astParser.ParseFile(fSet, "", data, astParser.ParseComments)
	if err != nil {
		return
	}

	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}

		for _, spec := range genDecl.Specs {
			switch spec := spec.(type) {
			case *ast.TypeSpec:
				switch t := spec.Type.(type) {
				case *ast.InterfaceType:
					for _, method := range t.Methods.List {
						_, ok = method.Type.(*ast.FuncType)
						if !ok {
							continue
						}

						for _, n := range method.Names {
							result = append(result, n.Name)
							break
						}
					}
				}
			}
		}
	}

	return
}
