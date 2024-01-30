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

package extract

import (
	"go/ast"
	astParser "go/parser"
	"go/token"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/cloudwego/cwgo/pkg/curd/code"

	"github.com/cloudwego/cwgo/pkg/common/utils"

	"github.com/fatih/camelcase"
)

type IdlExtractStruct struct {
	Name          string
	StructFields  []*StructField
	InterfaceInfo *InterfaceInfo
	UpdateInfo
}

type InterfaceInfo struct {
	Name    string
	Methods []*InterfaceMethod
}

type InterfaceMethod struct {
	Name             string
	ParsedTokens     string
	Params           code.Params
	Returns          code.Returns
	BelongedToStruct *IdlExtractStruct
}

type StructField struct {
	Name               string
	Type               code.Type
	Tag                reflect.StructTag
	IsBelongedToStruct bool
	BelongedToStruct   *IdlExtractStruct
}

type UpdateInfo struct {
	Update                bool
	UpdateCurdFileContent []byte
	UpdateIfFileContent   []byte
	PreMethodNamesMap     map[string]struct{}
	PreIfMethods          []*InterfaceMethod
}

func newIdlExtractStruct(name string) *IdlExtractStruct {
	return &IdlExtractStruct{
		Name:         name,
		StructFields: make([]*StructField, 0, 10),
		InterfaceInfo: &InterfaceInfo{
			Methods: make([]*InterfaceMethod, 0, 10),
		},
		UpdateInfo: UpdateInfo{
			PreMethodNamesMap: map[string]struct{}{},
			PreIfMethods:      []*InterfaceMethod{},
		},
	}
}

func (st *IdlExtractStruct) recordMongoIfInfo(daoDir string) error {
	fileMongoName, fileIfName := GetFileName(st.Name, daoDir)

	isExist, err := utils.PathExist(fileMongoName)
	if err != nil {
		return err
	}

	if isExist {
		isExist, err = utils.PathExist(fileIfName)
		if err != nil {
			return err
		}

		if isExist {
			st.Update = true
			st.UpdateCurdFileContent, err = utils.ReadFileContent(fileMongoName)
			if err != nil {
				return err
			}
			st.UpdateIfFileContent, err = utils.ReadFileContent(fileIfName)
			if err != nil {
				return err
			}

			preMethodNames, err := getInterfaceMethodNames(string(st.UpdateIfFileContent))
			if err != nil {
				return err
			}
			for _, methodName := range preMethodNames {
				st.PreMethodNamesMap[methodName] = struct{}{}
			}
		}
	}

	return nil
}

func extractIdlInterface(rawInterface string, rawStruct *IdlExtractStruct, tokens []string) error {
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
					rawStruct.InterfaceInfo = extractInterfaceType(spec.Name.Name, t, tokens, rawStruct)
				}
			}
		}
	}

	return nil
}

func extractInterfaceType(ifName string, interfaceType *ast.InterfaceType, tokens []string, rawStruct *IdlExtractStruct) *InterfaceInfo {
	intf := &InterfaceInfo{
		Name:    ifName,
		Methods: []*InterfaceMethod{},
	}

	for index, method := range interfaceType.Methods.List {
		funcType, ok := method.Type.(*ast.FuncType)
		if !ok {
			continue
		}

		var name string
		for _, n := range method.Names {
			name = n.Name
			break
		}

		if rawStruct.Update {
			if _, ok = rawStruct.PreMethodNamesMap[name]; !ok {
				meth := extractFunction(name, funcType, tokens[index])
				meth.BelongedToStruct = rawStruct

				intf.Methods = append(intf.Methods, meth)
			} else {
				meth := extractFunction(name, funcType, tokens[index])
				meth.BelongedToStruct = rawStruct

				rawStruct.PreIfMethods = append(rawStruct.PreIfMethods, meth)
			}
		} else {
			meth := extractFunction(name, funcType, tokens[index])
			meth.BelongedToStruct = rawStruct

			intf.Methods = append(intf.Methods, meth)
		}
	}

	return intf
}

func extractFunction(name string, funcType *ast.FuncType, token string) *InterfaceMethod {
	meth := &InterfaceMethod{
		Name:         name,
		ParsedTokens: token,
	}
	for _, param := range funcType.Params.List {
		paramType := getType(param.Type, "", false)

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
			meth.Returns = append(meth.Returns, getType(result.Type, "", false))
		}
	}

	return meth
}

func getType(expr ast.Expr, pkgName string, isPbCall bool) code.Type {
	switch expr := expr.(type) {
	case *ast.Ident:
		if !isPbCall {
			return code.IdentType(expr.Name)
		} else {
			if isGoBaseType(expr.Name) {
				return code.IdentType(expr.Name)
			}
			return code.SelectorExprType{
				X:   pkgName,
				Sel: expr.Name,
			}
		}

	case *ast.SelectorExpr:
		xExpr := expr.X.(*ast.Ident)
		return code.SelectorExprType{
			X:   xExpr.Name,
			Sel: expr.Sel.Name,
		}

	case *ast.StarExpr:
		realType := getType(expr.X, pkgName, isPbCall)
		return code.StarExprType{
			RealType: realType,
		}

	case *ast.ArrayType:
		elementType := getType(expr.Elt, pkgName, isPbCall)
		return code.SliceType{
			ElementType: elementType,
		}

	case *ast.MapType:
		keyType := getType(expr.Key, pkgName, isPbCall)
		valueType := getType(expr.Value, pkgName, isPbCall)
		return code.MapType{KeyType: keyType, ValueType: valueType}

	case *ast.InterfaceType:
		return code.InterfaceType{}
	}

	return nil
}

func GetFileName(structName, prefix string) (fileMongoName, fileIfName string) {
	dir := GetPkgName(structName)
	fileMongoName = filepath.Join(prefix, dir, dir+"_repo_mongo.go")
	fileIfName = filepath.Join(prefix, dir, dir+"_repo.go")
	return
}

func GetPkgName(structName string) string {
	tokens := camelcase.Split(structName)
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
