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
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/bytedance/sonic"
)

type Parser struct {
	moduleName   string // project module name
	hertzRepoUrl string

	fSet    *token.FileSet
	funcMap map[string]map[string]*FuncParsed

	// TODO: should consider external var too
	// globalVarMap map[string]*Var

	routerParsedList []*RouterParsed
}

type RouterParsed struct {
	FilePath  string `json:"file_path"`
	StartLine int    `json:"start_line"`
	EndLine   int    `json:"end_line"`
	Method    string `json:"method"`
	RoutePath string `json:"route_path"`
}

type FuncParsed struct {
	importMap map[string]*ImportParsed
	filePath  string
	funcDecl  *ast.FuncDecl
}

type ImportParsed struct {
	Path                 string // full package path
	IsLocalModulePackage bool   // is package in current project
}

type Var struct {
	Name   string // variable name TODO: not consider shadowed declaration
	Type   VarType
	Prefix string
}

func NewParser(projectPath, hertzRepoUrl string) (*Parser, error) {
	// get module name
	moduleName, err := getModuleName(projectPath)
	if err != nil {
		return nil, err
	}

	p := &Parser{
		moduleName:       moduleName,
		hertzRepoUrl:     hertzRepoUrl,
		fSet:             token.NewFileSet(),
		funcMap:          make(map[string]map[string]*FuncParsed),
		routerParsedList: make([]*RouterParsed, 0),
	}

	// init func map
	err = filepath.WalkDir(projectPath, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			// parse whole package
			astPkgMap, err := parser.ParseDir(p.fSet, path, nil, parser.ParseComments)
			if err != nil {
				return err
			}

			for _, astPkg := range astPkgMap {
				fullPkgName := strings.Replace(path, projectPath, moduleName, 1)
				p.funcMap[fullPkgName] = make(map[string]*FuncParsed)

				for fileName, astFile := range astPkg.Files {
					if strings.HasSuffix(fileName, "_test.go") {
						// skip test file
						continue
					}

					// parse imports
					importMap := make(map[string]*ImportParsed)
					for _, importSpec := range astFile.Imports {
						importSpecPath := strings.Trim(importSpec.Path.Value, "\"")
						var importPkgName string // package call name
						if importSpec.Name != nil {
							// package alias
							importPkgName = importSpec.Name.Name
						} else {
							// package short name
							importPkgName = filepath.Base(importSpecPath)
						}

						if strings.HasPrefix(importSpecPath, moduleName) {
							importMap[importPkgName] = &ImportParsed{
								Path:                 importSpecPath,
								IsLocalModulePackage: true,
							}
						} else {
							importMap[importPkgName] = &ImportParsed{
								Path:                 importSpecPath,
								IsLocalModulePackage: false,
							}
						}
					}

					// parse funcs
					if astFile.Scope != nil && astFile.Scope.Objects != nil {
						for funcName, object := range astFile.Scope.Objects {
							if object.Kind != ast.Fun {
								continue
							}

							p.funcMap[fullPkgName][funcName] = &FuncParsed{
								importMap: importMap,
								filePath:  fileName,
								funcDecl:  object.Decl.(*ast.FuncDecl),
							}
						}
					}
				}
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return p, nil
}

// searchFunc is a recursive func that traverse func body
// params:
// packageName: func is located in which package
// funcName: func name
// localGroupVarMap: stores local var that can call Group()
// funcParams: var infos of func params that passed by external
func (p *Parser) searchFunc(packageName, funcName string, localGroupVarMap map[string]*Var, funcParams []*Var) error {
	funcParsed, ok := p.funcMap[packageName][funcName]
	if !ok {
		// func not found in parsed func map
		return fmt.Errorf("func not found, package_name: %s, func_name: %s", packageName, funcName)
	}

	// init local group var map by func params passed in
	i := 0
	for _, param := range funcParsed.funcDecl.Type.Params.List {
		for _, name := range param.Names {
			funcParams[i].Name = name.Name
			localGroupVarMap[funcParams[i].Name] = funcParams[i]
			i++
		}
	}

	// traverse stmts in function
	err := p.searchStmts(funcParsed.funcDecl.Body.List, packageName, funcParsed, localGroupVarMap)
	if err != nil {
		return err
	}

	return nil
}

// searchStmts is a recursive func that traverse stmts(func, if, switch, caseClause)
// params:
// stmts: stmt list
// packageName: func is located in which package
// funcParsed: func info which stmts belong
// localGroupVarMap: stores local var that can call Group()
func (p *Parser) searchStmts(stmts []ast.Stmt, packageName string, funcParsed *FuncParsed, localGroupVarMap map[string]*Var) error {
	for _, stmtIface := range stmts {
		switch stmt := stmtIface.(type) {
		case *ast.ExprStmt:
			// if is expr stmt
			// then search for call expr only
			exprXCallExpr, ok := stmt.X.(*ast.CallExpr)
			if !ok {
				continue
			}

			switch exprXCallExprFun := exprXCallExpr.Fun.(type) {
			case *ast.Ident:
				// calling func is in current package
				if _, isBuiltinFunc := BuiltinFuncNameMap[exprXCallExprFun.Name]; isBuiltinFunc {
					// return nil if func is builtin
					continue
				}

				// get relativePath func param if it has *route.RouterGroup
				funcParams := p.getVarsInArgs(localGroupVarMap, exprXCallExpr)
				// recursively search func
				err := p.searchFunc(packageName, exprXCallExprFun.Name, make(map[string]*Var), funcParams)
				if err != nil {
					return err
				}

			case *ast.SelectorExpr:
				switch exprXCallExprFunX := exprXCallExprFun.X.(type) {
				case *ast.Ident:
					if exprXCallExprFunX.Obj == nil {
						// funX is package name
						if pkg, ok := funcParsed.importMap[exprXCallExprFunX.Name]; ok && pkg.IsLocalModulePackage {
							// calling func is in project module
							if _, isBuiltinFunc := BuiltinFuncNameMap[exprXCallExprFun.Sel.Name]; isBuiltinFunc {
								// return nil if func is builtin
								continue
							}

							// get relativePath func param if it has *route.RouterGroup
							funcParams := p.getVarsInArgs(localGroupVarMap, exprXCallExpr)
							// recursively search func
							err := p.searchFunc(pkg.Path, exprXCallExprFun.Sel.Name, make(map[string]*Var), funcParams)
							if err != nil {
								return err
							}
						}
					} else {
						if exprXCallExprFunX.Obj.Kind == ast.Var {
							// funX is var and is defined in ast(represent the var is local var)
							// then search var info in local group var map
							if v, ok := localGroupVarMap[exprXCallExprFunX.Obj.Name]; ok && v.Type != VarTypeOther {
								// var is existed in local group var map
								if _, ok := RouterFuncNameMap[exprXCallExprFun.Sel.Name]; ok {
									// is calling register func

									// get relativePath in func param
									if len(exprXCallExpr.Args) > 0 {
										switch paramExpr := exprXCallExpr.Args[0].(type) {
										// get first param(relativePath) of router func
										case *ast.BasicLit:
											// if param is literal
											if paramExpr.Kind == token.STRING {
												// if parma is string
												relativePath := strings.Trim(paramExpr.Value, "\"")
												fullRouter := filepath.Join(v.Prefix, relativePath)

												startLine := p.fSet.Position(stmt.Pos()).Line
												endLine := p.fSet.Position(stmt.End()).Line

												p.routerParsedList = append(p.routerParsedList, &RouterParsed{
													FilePath:  funcParsed.filePath,
													StartLine: startLine,
													EndLine:   endLine,
													Method:    exprXCallExprFun.Sel.Name,
													RoutePath: fullRouter,
												})
											} else {
												continue
											}
										case *ast.Ident:
											// TODO: if param is var
										}
									}
								}
							}
						}
					}
				}
			}
		case *ast.AssignStmt:
			// only consider assign stmt as follows:
			// 1. h := server.Default()
			// 2. h := server.New()
			// 3. h := byted.Default()
			// 4. e := h.Engine
			// 6. g := h.Group()
			// 7. g := e.Group()
			// 8. g1 := g.Group()
			// 9. g1, g2 := g.Group(). g.Group()
			if len(stmt.Lhs) != len(stmt.Rhs) {
				continue
			}
			for i, lhs := range stmt.Lhs {
				rhs := stmt.Rhs[i]
				switch rhsExpr := rhs.(type) {
				case *ast.CallExpr:
					if callFunSelectorExpr, ok := rhsExpr.Fun.(*ast.SelectorExpr); ok {
						if xExpr, ok := callFunSelectorExpr.X.(*ast.Ident); ok {
							if xExpr.Obj == nil {
								// package name
								if xExpr.Name == "server" && (callFunSelectorExpr.Sel.Name == "Default" || callFunSelectorExpr.Sel.Name == "New") {
									if imp, ok := funcParsed.importMap["server"]; ok && imp.Path == p.hertzRepoUrl+"/pkg/app/server" {
										if lhsIdent, ok := lhs.(*ast.Ident); ok {
											localGroupVarMap[lhsIdent.Name] = &Var{
												Name:   lhsIdent.Name,
												Type:   VarTypeServerHertz,
												Prefix: "",
											}
										}
										continue
									}
								} else if xExpr.Name == "byted" && callFunSelectorExpr.Sel.Name == "Default" {
									if imp, ok := funcParsed.importMap["byted"]; ok && imp.Path == p.hertzRepoUrl+"/byted" {
										if lhsIdent, ok := lhs.(*ast.Ident); ok {
											localGroupVarMap[lhsIdent.Name] = &Var{
												Name:   lhsIdent.Name,
												Type:   VarTypeServerHertz,
												Prefix: "",
											}
										}
										continue
									}
								}
							} else {
								// var
								if xExpr.Obj.Kind == ast.Var {
									if v, ok := localGroupVarMap[xExpr.Obj.Name]; ok && v.Type != VarTypeOther {
										if callFunSelectorExpr.Sel.Name == "Group" {
											// if var call Group()

											// get relativePath func param
											if len(rhsExpr.Args) > 0 {
												if paramExpr, ok := rhsExpr.Args[0].(*ast.BasicLit); ok {
													// get first param(relativePath) of Group()
													if paramExpr.Kind == token.STRING {
														// if parma is string
														if v, ok := localGroupVarMap[xExpr.Name]; ok {
															if lhsIdent, ok := lhs.(*ast.Ident); ok {
																localGroupVarMap[lhsIdent.Name] = &Var{
																	Name:   lhsIdent.Name,
																	Type:   VarTypeRouterGroup,
																	Prefix: filepath.Join(v.Prefix, strings.Trim(paramExpr.Value, "\"")),
																}
																continue
															}
														}
													}
												} else {
													// TODO: if param is var
												}
											}
										}
									}
								}
							}
						}
					}

				case *ast.SelectorExpr:
					if rhsExpr.Sel.Name == "Engine" {
						if rhsExprXIdent, ok := rhsExpr.X.(*ast.Ident); ok {
							if rhsExprXIdent.Obj != nil && rhsExprXIdent.Obj.Kind == ast.Var {
								if v, ok := localGroupVarMap[rhsExprXIdent.Name]; ok && v.Type == VarTypeServerHertz {
									if lhsIdent, ok := lhs.(*ast.Ident); ok {
										localGroupVarMap[lhsIdent.Name] = &Var{
											Name:   lhsIdent.Name,
											Type:   VarTypeRouteEngine,
											Prefix: "",
										}
										continue
									}
								}
							}
						}
					}
				}
			}
		case *ast.IfStmt:
			err := p.searchStmts(stmt.Body.List, packageName, funcParsed, localGroupVarMap)
			if err != nil {
				return err
			}

			if stmt.Else != nil {
				err = p.searchStmts(stmt.Else.(*ast.BlockStmt).List, packageName, funcParsed, localGroupVarMap)
				if err != nil {
					return err
				}
			}
		case *ast.SwitchStmt:
			err := p.searchStmts(stmt.Body.List, packageName, funcParsed, localGroupVarMap)
			if err != nil {
				return err
			}
		case *ast.CaseClause:
			err := p.searchStmts(stmt.Body, packageName, funcParsed, localGroupVarMap)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (p *Parser) getVarsInArgs(varMap map[string]*Var, expr *ast.CallExpr) []*Var {
	res := make([]*Var, 0)

	for _, exprArg := range expr.Args {
		switch argExpr := exprArg.(type) {
		case *ast.CallExpr:
			if funSelectorExpr, ok := argExpr.Fun.(*ast.SelectorExpr); ok {
				if xIdent, ok := funSelectorExpr.X.(*ast.Ident); ok {
					if xIdent.Obj != nil && xIdent.Obj.Kind == ast.Var {
						if v, ok := varMap[xIdent.Obj.Name]; ok && v.Type != VarTypeOther {
							if funSelectorExpr.Sel.Name == "Group" {
								// if var called with Group Method

								// get relativePath func param
								if len(argExpr.Args) > 0 {
									if paramExpr, ok := argExpr.Args[0].(*ast.BasicLit); ok {
										if paramExpr.Kind == token.STRING {
											if v, ok := varMap[xIdent.Name]; ok {
												res = append(res, &Var{
													Type:   VarTypeRouterGroup,
													Prefix: filepath.Join(v.Prefix, strings.Trim(paramExpr.Value, "\"")),
												})
												continue
											}
										}
									} else {
										// TODO: if param is var
									}
								}
							}
						}
					}
				}
			}
		case *ast.Ident:
			if v, ok := varMap[argExpr.Name]; ok {
				res = append(res, &Var{
					Type:   v.Type,
					Prefix: v.Prefix,
				})
				continue
			}
		case *ast.SelectorExpr:
			if argCallExprSelectorXIdent, ok := argExpr.X.(*ast.Ident); ok {
				if v, ok := varMap[argCallExprSelectorXIdent.Name]; ok && v.Type != VarTypeOther && argExpr.Sel.Name == "Engine" {
					res = append(res, &Var{
						Name:   "",
						Type:   VarTypeRouteEngine,
						Prefix: "",
					})
					continue
				}
			}
		}
		res = append(res, &Var{
			Name:   "",
			Type:   VarTypeOther,
			Prefix: "",
		})
	}

	return res
}

func (p *Parser) PrintRouters() {
	j, _ := sonic.Marshal(p.routerParsedList)
	fmt.Println(string(j))
}
