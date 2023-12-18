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

package generator

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"

	"golang.org/x/tools/go/ast/astutil"
)

func appendGoFileImports(src string, appendImports []string) (data string, err error) {
	fSet := token.NewFileSet()
	file, err := parser.ParseFile(fSet, "", src, parser.ParseComments)
	if err != nil {
		return "", err
	}

	extraImports := make([]string, 0, 10)
	for _, impt := range appendImports {
		flag := 0
		ast.Inspect(file, func(n ast.Node) bool {
			if importSpec, ok := n.(*ast.ImportSpec); ok && importSpec.Path.Value == impt {
				flag = 1
				return false
			}
			return true
		})
		if flag == 0 {
			extraImports = append(extraImports, impt)
		}
	}

	for _, imp := range extraImports {
		astutil.AddImport(fSet, file, imp)
	}

	var buf bytes.Buffer
	if err = printer.Fprint(&buf, fSet, file); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func replaceFuncBody(src string, funcName, funcBody []string) (data string, err error) {
	fSet := token.NewFileSet()
	file, err := parser.ParseFile(fSet, "", src, parser.ParseComments)
	if err != nil {
		return "", err
	}

	var targetFunc []*ast.FuncDecl
	for _, name := range funcName {
		ast.Inspect(file, func(n ast.Node) bool {
			if funcDecl, ok := n.(*ast.FuncDecl); ok && funcDecl.Name.Name == name {
				targetFunc = append(targetFunc, funcDecl)
				return false
			}
			return true
		})
	}

	for index, tgFunc := range targetFunc {
		tgFunc.Body = &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ExprStmt{
					X: &ast.BasicLit{
						Kind:  token.STRING,
						Value: funcBody[index],
					},
				},
			},
		}
	}

	var buf bytes.Buffer
	if err = printer.Fprint(&buf, fSet, file); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func isFuncExist(src, funcName string) (isExist bool, err error) {
	fSet := token.NewFileSet()
	file, err := parser.ParseFile(fSet, "", src, parser.ParseComments)
	if err != nil {
		return false, err
	}

	ast.Inspect(file, func(n ast.Node) bool {
		if funcDecl, ok := n.(*ast.FuncDecl); ok && funcDecl.Name.Name == funcName {
			isExist = true
			return false
		}
		return true
	})

	return
}

func isFuncBodyEqual(src, funcName, body string) (isExist bool, err error) {
	fSet := token.NewFileSet()
	file, err := parser.ParseFile(fSet, "", src, parser.ParseComments)
	if err != nil {
		return false, err
	}

	var targetFunc *ast.FuncDecl
	ast.Inspect(file, func(n ast.Node) bool {
		if funcDecl, ok := n.(*ast.FuncDecl); ok && funcDecl.Name.Name == funcName {
			targetFunc = funcDecl
			return false
		}
		return true
	})

	buffer := &bytes.Buffer{}
	if err = printer.Fprint(buffer, fSet, targetFunc.Body); err != nil {
		return false, err
	}

	return buffer.String() == body, nil
}

func getStructNames(src string) (result []string, err error) {
	fSet := token.NewFileSet()
	file, err := parser.ParseFile(fSet, "", src, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	ast.Inspect(file, func(n ast.Node) bool {
		ts, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}
		_, ok = ts.Type.(*ast.StructType)
		if !ok {
			return true
		}
		result = append(result, ts.Name.Name)
		return true
	})

	return
}
