package api_list

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strings"
)

type Parser struct {
	moduleName   string // project module name
	hertzRepoUrl string

	fSet         *token.FileSet
	astFileCache map[string]*ast.File
	funcMap      map[string]map[string]*FuncParsed

	// TODO: should consider external var too
	// globalVarMap map[string]*Var

	routerParsedList []*RouterParsed
}

type RouterParsed struct {
	FilePath  string
	StartLine int
	EndLine   int
	Method    string
	RoutePath string
}

type FuncParsed struct {
	importMap map[string]*ImportParsed
	filePath  string
	funcDecl  *ast.FuncDecl
}

type ImportParsed struct {
	Path                 string
	IsLocalModulePackage bool
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
		astFileCache:     make(map[string]*ast.File),
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
				for fileName, astFile := range astPkg.Files {
					if strings.HasSuffix(fileName, "_test.go") {
						// skip test file
						continue
					}

					p.astFileCache[fileName] = astFile

					// parse imports
					importMap := make(map[string]*ImportParsed)
					for _, importSpec := range astFile.Imports {
						importSpecPath := strings.Trim(importSpec.Path.Value, "\"")
						var importPkgName string
						if importSpec.Name != nil {
							importPkgName = importSpec.Name.Name
						} else {
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

							if _, ok := p.funcMap[fullPkgName]; !ok {
								p.funcMap[fullPkgName] = make(map[string]*FuncParsed)
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

func (p *Parser) searchFunc(packageName, funcName string, localGroupVarMap map[string]*Var, funcParams []*Var) error {
	funcParsed, ok := p.funcMap[packageName][funcName]
	if !ok {
		if _, ok := BuiltinFuncNameMap[funcName]; ok {
			return nil
		}

		return fmt.Errorf("func not found, package_name: %s, func_name: %s", packageName, funcName)
	}

	i := 0
	for _, param := range funcParsed.funcDecl.Type.Params.List {
		for _, name := range param.Names {
			funcParams[i].Name = name.Name
			localGroupVarMap[funcParams[i].Name] = funcParams[i]
			i++
		}
	}

	// traverse stmt in function
	err := p.searchStmts(funcParsed.funcDecl.Body.List, packageName, funcParsed, localGroupVarMap)
	if err != nil {
		return err
	}

	return nil
}

func (p *Parser) searchStmts(stmts []ast.Stmt, packageName string, funcParsed *FuncParsed, localGroupVarMap map[string]*Var) error {
	for _, stmtIface := range stmts {
		switch stmt := stmtIface.(type) {
		case *ast.ExprStmt:
			// if is expr stmt
			// then search for call expr only
			expr, ok := stmt.X.(*ast.CallExpr)
			if !ok {
				continue
			}

			switch fun := expr.Fun.(type) {
			case *ast.Ident:
				// calling func is in current package
				if _, isBuiltinFunc := BuiltinFuncNameMap[fun.Name]; isBuiltinFunc {
					continue
				}

				// get relativePath func param if it has *route.RouterGroup
				params := p.getVarsInArgs(funcParsed.importMap, localGroupVarMap, expr)
				err := p.searchFunc(packageName, fun.Name, make(map[string]*Var), params)
				if err != nil {
					return err
				}

			case *ast.SelectorExpr:
				switch funX := fun.X.(type) {
				case *ast.Ident:
					if funX.Obj == nil {
						// funX is package name
						if pkg, ok := funcParsed.importMap[funX.Name]; ok && pkg.IsLocalModulePackage {
							// calling func is in project module

							// get relativePath func param if it has *route.RouterGroup
							params := p.getVarsInArgs(funcParsed.importMap, localGroupVarMap, expr)
							err := p.searchFunc(pkg.Path, fun.Sel.Name, make(map[string]*Var), params)
							if err != nil {
								return err
							}
						}
					} else {
						if funX.Obj.Kind == ast.Var {
							// funX is var ans is defined in ast
							if v, ok := localGroupVarMap[funX.Obj.Name]; ok && v.Type != VarTypeOther {
								// var is existed in local group var map
								if _, ok := RouterFuncNameMap[fun.Sel.Name]; ok {
									// is calling register func

									// get relativePath in func param
									if len(expr.Args) > 0 {
										switch paramExpr := expr.Args[0].(type) {
										case *ast.BasicLit:
											if paramExpr.Kind == token.STRING {
												relativePath := strings.Trim(paramExpr.Value, "\"")
												fullRouter := filepath.Join(v.Prefix, relativePath)

												startLine := p.fSet.Position(stmt.Pos()).Line
												endLine := p.fSet.Position(stmt.End()).Line

												p.routerParsedList = append(p.routerParsedList, &RouterParsed{
													FilePath:  funcParsed.filePath,
													StartLine: startLine,
													EndLine:   endLine,
													Method:    fun.Sel.Name,
													RoutePath: fullRouter,
												})
											}
										}
									}
								}
							}
						}
					}
				}
			}
		case *ast.AssignStmt:
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
											// if var called with Group Method

											// get relativePath func param
											if len(rhsExpr.Args) > 0 {
												if paramExpr, ok := rhsExpr.Args[0].(*ast.BasicLit); ok {
													if paramExpr.Kind == token.STRING {
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
													// TODO: should consider relativePath var if is not string literal
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

func (p *Parser) checkIsServerHertzVar(importMap map[string]*ImportParsed, objDeclIface interface{}) bool {
	switch objDecl := objDeclIface.(type) {
	case *ast.Field:
		// defined in func param
		if starExpr, ok := objDecl.Type.(*ast.StarExpr); ok {
			if selectorExpr, ok := starExpr.X.(*ast.SelectorExpr); ok {
				if xExpr, ok := selectorExpr.X.(*ast.Ident); ok {
					if xExpr.Name == "server" && selectorExpr.Sel.Name == "Hertz" {
						if imp, ok := importMap["server"]; ok && imp.Path == p.hertzRepoUrl+"/pkg/app/server" {
							return true
						}
					}
				}
			}
		}
	case *ast.AssignStmt:
		// assigned in func body
		if len(objDecl.Rhs) > 0 {
			if callExpr, ok := objDecl.Rhs[0].(*ast.CallExpr); ok {
				if selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
					if xExpr, ok := selectorExpr.X.(*ast.Ident); ok {
						if xExpr.Name == "server" && (selectorExpr.Sel.Name == "Default" || selectorExpr.Sel.Name == "New") {
							if imp, ok := importMap["server"]; ok && imp.Path == p.hertzRepoUrl+"/pkg/app/server" {
								return true
							}
						} else if xExpr.Name == "byted" && selectorExpr.Sel.Name == "Default" {
							if imp, ok := importMap["byted"]; ok && imp.Path == p.hertzRepoUrl+"/byted" {
								return true
							}
						}
					}
				}
			}
		}

	default:
		return false
	}

	return false
}

func (p *Parser) checkIsRouteEngine(importMap map[string]*ImportParsed, objDeclIface interface{}) bool {
	switch objDecl := objDeclIface.(type) {
	case *ast.Field:
		if starExpr, ok := objDecl.Type.(*ast.StarExpr); ok {
			if selectorExpr, ok := starExpr.X.(*ast.SelectorExpr); ok {
				if xExpr, ok := selectorExpr.X.(*ast.Ident); ok {
					if xExpr.Name == "route" && selectorExpr.Sel.Name == "Engine" {
						if imp, ok := importMap["route"]; ok && imp.Path == p.hertzRepoUrl+"/pkg/route" {
							return true
						}
					}
				}
			}
		}
	default:
		return false
	}

	return false
}

func (p *Parser) checkIsRouterGroup(importMap map[string]*ImportParsed, objDeclIface interface{}) bool {
	switch objDecl := objDeclIface.(type) {
	case *ast.Field:
		if starExpr, ok := objDecl.Type.(*ast.StarExpr); ok {
			if selectorExpr, ok := starExpr.X.(*ast.SelectorExpr); ok {
				if xExpr, ok := selectorExpr.X.(*ast.Ident); ok {
					if xExpr.Name == "route" && selectorExpr.Sel.Name == "RouterGroup" {
						if imp, ok := importMap["route"]; ok && imp.Path == "code.byted.org/middleware/hertz/pkg/route" {
							return true
						}
					}
				}
			}
		}

	}

	return false
}

func (p *Parser) getVarType(importMap map[string]*ImportParsed, objDeclIface interface{}) VarType {
	if p.checkIsServerHertzVar(importMap, objDeclIface) {
		return VarTypeServerHertz
	}
	if p.checkIsRouteEngine(importMap, objDeclIface) {
		return VarTypeRouteEngine
	}
	if p.checkIsRouterGroup(importMap, objDeclIface) {
		return VarTypeRouterGroup
	}

	return VarTypeOther
}

func (p *Parser) getVarsInArgs(importMap map[string]*ImportParsed, varMap map[string]*Var, expr *ast.CallExpr) []*Var {
	res := make([]*Var, 0)

	for _, exprArg := range expr.Args {
		switch argExpr := exprArg.(type) {
		case *ast.CallExpr:
			if funSelectorExpr, ok := argExpr.Fun.(*ast.SelectorExpr); ok {
				if xIdent, ok := funSelectorExpr.X.(*ast.Ident); ok {
					if xIdent.Obj != nil && xIdent.Obj.Kind == ast.Var {
						if varType := p.getVarType(importMap, xIdent.Obj.Decl); varType != VarTypeOther {
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
										// TODO: should consider relativePath var if is not string literal
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
	for i, routerParsed := range p.routerParsedList {
		if routerParsed.StartLine == routerParsed.EndLine {
			fmt.Printf(
				"%d -> {%s \"%s\":\"%s#L%d\"}\n",
				i+1,
				routerParsed.Method,
				routerParsed.RoutePath,
				routerParsed.FilePath,
				routerParsed.StartLine,
			)
		} else {
			fmt.Printf(
				"%d -> {%s \"%s\":\"%s#L%d-L%d\"}\n",
				i+1,
				routerParsed.Method,
				routerParsed.RoutePath,
				routerParsed.FilePath,
				routerParsed.StartLine,
				routerParsed.EndLine,
			)
		}
	}
}
