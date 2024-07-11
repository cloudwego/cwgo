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

package cronjob

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/cloudwego/cwgo/meta"

	"github.com/cloudwego/cwgo/config"
	"github.com/cloudwego/cwgo/pkg/common/utils"
	"github.com/cloudwego/cwgo/pkg/consts"
	"github.com/cloudwego/kitex/tool/internal_pkg/log"
)

func Cronjob(c *config.CronJobArgument) error {
	if err := check(c); err != nil {
		return err
	}

	err := generateCronjobFile(c)
	if err != nil {
		return err
	}

	return nil
}

func check(c *config.CronJobArgument) (err error) {
	if len(c.JobName) == 0 {
		return errors.New("job name is empty")
	}

	c.OutDir, err = filepath.Abs(c.OutDir)
	if err != nil {
		return err
	}

	gopath, err := utils.GetGOPATH()
	if err != nil {
		return fmt.Errorf("GET gopath failed: %s", err)
	}
	if gopath == "" {
		return fmt.Errorf("GOPATH is not set")
	}

	gosrc := filepath.Join(gopath, consts.Src)
	gosrc, err = filepath.Abs(gosrc)
	if err != nil {
		log.Warn("Get GOPATH/src path failed:", err.Error())
		os.Exit(1)
	}
	curpath, err := filepath.Abs(consts.CurrentDir)
	if err != nil {
		log.Warn("Get current path failed:", err.Error())
		os.Exit(1)
	}

	if !strings.HasSuffix(c.JobFile, ".go") {
		log.Warn("--job_file must be a go file")
		os.Exit(1)
	}

	for k := range c.JobName {
		c.JobName[k] = utils.CapitalizeFirstLetter(c.JobName[k])
	}

	if strings.HasPrefix(curpath, gosrc) {
		goPkg := ""
		if goPkg, err = filepath.Rel(gosrc, curpath); err != nil {
			log.Warn("Get GOPATH/src relpath failed:", err.Error())
			os.Exit(1)
		}

		if c.GoMod == "" {
			if utils.IsWindows() {
				c.GoMod = strings.ReplaceAll(goPkg, consts.BackSlash, consts.Slash)
			} else {
				c.GoMod = goPkg
			}
		}

		if c.GoMod != "" {
			if utils.IsWindows() {
				goPkgSlash := strings.ReplaceAll(goPkg, consts.BackSlash, consts.Slash)
				if goPkgSlash != c.GoMod {
					return fmt.Errorf("module name: %s is not the same with GoPkg under GoPath: %s", c.GoMod, goPkgSlash)
				}
			} else {
				if c.GoMod != goPkg {
					return fmt.Errorf("module name: %s is not the same with GoPkg under GoPath: %s", c.GoMod, goPkg)
				}
			}
		}
	}

	if !strings.HasPrefix(curpath, gosrc) && c.GoMod == "" {
		log.Warn("Outside of $GOPATH. Please specify a module name with the '-module' flag.")
		os.Exit(1)
	}

	if c.GoMod != "" {
		module, path, ok := utils.SearchGoMod(curpath, true)

		if ok {
			// go.mod exists
			if module != c.GoMod {
				log.Warnf("The module name given by the '-module' option ('%s') is not consist with the name defined in go.mod ('%s' from %s)\n",
					c.GoMod, module, path)
				os.Exit(1)
			}
			if c.PackagePrefix, err = filepath.Rel(path, c.OutDir); err != nil {
				log.Warn("Get package prefix failed:", err.Error())
				os.Exit(1)
			}
			c.PackagePrefix = filepath.Join(c.GoMod, c.PackagePrefix)
		} else {
			if err = utils.InitGoMod(c.GoMod); err != nil {
				log.Warn("Init go mod failed:", err.Error())
				os.Exit(1)
			}
			if c.PackagePrefix, err = filepath.Rel(curpath, c.OutDir); err != nil {
				log.Warn("Get package prefix failed:", err.Error())
				os.Exit(1)
			}
			c.PackagePrefix = filepath.Join(c.GoMod, c.PackagePrefix)
		}
	}

	c.PackagePrefix = strings.ReplaceAll(c.PackagePrefix, consts.BackSlash, consts.Slash)

	return nil
}

type JobsData struct {
	JobInfos []JobInfo
}
type JobInfo struct {
	JobName       string
	GoModule      string
	PackagePrefix string
}

func addScheduleNewJobs(data string, jobs []JobInfo) (string, error) {
	fSet := token.NewFileSet()
	file, err := parser.ParseFile(fSet, "", data, parser.ParseComments)
	if err != nil {
		return "", err
	}

	// Extract existing cronjob calls
	existingJobs := make(map[string]bool)
	ast.Inspect(file, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.CallExpr:
			if sel, ok := x.Fun.(*ast.SelectorExpr); ok {
				if ident, ok := sel.X.(*ast.Ident); ok && sel.Sel.Name == "AddFunc" && ident.Name == "c" {
					for _, arg := range x.Args {
						if funcLit, ok := arg.(*ast.FuncLit); ok {
							ast.Inspect(funcLit.Body, func(n ast.Node) bool {
								switch stmt := n.(type) {
								case *ast.ExprStmt:
									if callExpr, ok := stmt.X.(*ast.CallExpr); ok {
										if sel, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
											if id, ok := sel.X.(*ast.Ident); ok && id.Name == "job" {
												existingJobs[sel.Sel.Name] = true
											}
										}
									}
								}
								return true
							})
						}
					}
				}
			}
		}
		return true
	})

	buf := new(bytes.Buffer)

	for _, decl := range file.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok && fn.Name.Name == "Init" {
			for k := range jobs {
				if _, ok := existingJobs[jobs[k].JobName]; !ok {
					jobName := jobs[k].JobName
					addCronJobCall(fn, jobName)
				}
			}
		}
	}

	// Generate the modified code
	if err := printer.Fprint(buf, fSet, file); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func addNewJobs(data string, jobs []JobInfo) (string, error) {
	fSet := token.NewFileSet()
	file, err := parser.ParseFile(fSet, "", data, parser.ParseComments)
	if err != nil {
		return "", err
	}

	// Extract existing cronjob
	existingJobs := make(map[string]bool)
	ast.Inspect(file, func(n ast.Node) bool {
		if fn, ok := n.(*ast.FuncDecl); ok {
			existingJobs[fn.Name.Name] = true
		}
		return true
	})

	for _, job := range jobs {
		if _, ok := existingJobs[job.JobName]; !ok {
			addCronJobFunction(file, job.JobName)
		}
	}

	// Generate the modified code
	buf := new(bytes.Buffer)
	if err = printer.Fprint(buf, fSet, file); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func addCronJobFunction(file *ast.File, jobName string) {
	funcDecl := &ast.FuncDecl{
		Name: ast.NewIdent(jobName),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("ctx")},
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("context"),
							Sel: ast.NewIdent("Context"),
						},
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ExprStmt{
					X: &ast.BasicLit{
						Kind:  token.STRING,
						Value: `// TODO: fill with your own logic`,
					},
				},
			},
		},
		Doc: &ast.CommentGroup{},
	}

	// Append the new function declaration to the AST file
	file.Decls = append(file.Decls, funcDecl)
}

func addCronJobCall(fn *ast.FuncDecl, jobName string) {
	// Construct the AST nodes for the new cron job call
	addFuncCall := &ast.AssignStmt{
		Lhs: []ast.Expr{
			ast.NewIdent("_, err"),
		},
		Tok: token.ASSIGN,
		Rhs: []ast.Expr{
			&ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   ast.NewIdent("c"),
					Sel: ast.NewIdent("AddFunc"),
				},
				Args: []ast.Expr{
					&ast.BasicLit{
						Kind:  token.STRING,
						Value: `"* * * * *"`,
					},
					&ast.FuncLit{
						Type: &ast.FuncType{
							Params: &ast.FieldList{
								List: []*ast.Field{},
							},
						},
						Body: &ast.BlockStmt{
							List: []ast.Stmt{
								&ast.SelectStmt{
									Body: &ast.BlockStmt{
										List: []ast.Stmt{
											&ast.CommClause{
												Comm: &ast.ExprStmt{
													X: &ast.UnaryExpr{
														Op: token.ARROW,
														X:  ast.NewIdent("ctx.Done()"),
													},
												},
												Body: []ast.Stmt{
													&ast.ExprStmt{
														X: &ast.CallExpr{
															Fun: &ast.SelectorExpr{
																X:   ast.NewIdent("log"),
																Sel: ast.NewIdent("Println"),
															},
															Args: []ast.Expr{
																&ast.BasicLit{
																	Kind:  token.STRING,
																	Value: fmt.Sprintf(`"%s terminated."`, jobName),
																},
															},
														},
													},
													&ast.ReturnStmt{},
												},
											},
											&ast.CommClause{
												Body: []ast.Stmt{
													&ast.ExprStmt{
														X: &ast.CallExpr{
															Fun: &ast.SelectorExpr{
																X:   ast.NewIdent("job"),
																Sel: ast.NewIdent(jobName),
															},
															Args: []ast.Expr{
																ast.NewIdent("ctx"),
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	checkErrCall := &ast.IfStmt{
		Cond: &ast.BinaryExpr{
			X:  ast.NewIdent("err"),
			Op: token.NEQ,
			Y:  ast.NewIdent("nil"),
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent("log"),
							Sel: ast.NewIdent("Fatalf"),
						},
						Args: []ast.Expr{
							&ast.BasicLit{
								Kind:  token.STRING,
								Value: `"Error adding new cron job: %v"`,
							},
							ast.NewIdent("err"),
						},
					},
				},
			},
		},
	}

	dummyCall := &ast.ExprStmt{
		X: &ast.BasicLit{
			Kind:  token.STRING,
			Value: "",
		},
	}

	fn.Body.List = append(fn.Body.List, dummyCall, addFuncCall, checkErrCall)
}

func generateCronjobFile(c *config.CronJobArgument) error {
	// Ensure the base output directory exists
	err := os.MkdirAll(c.OutDir, 0o755)
	if err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Create cmd/main.go and overwrite each time
	cmdDir := filepath.Join(c.OutDir, "cmd")
	err = os.MkdirAll(cmdDir, 0o755)
	if err != nil {
		return fmt.Errorf("failed to create cmd directory: %w", err)
	}

	mainGoPath := filepath.Join(cmdDir, "main.go")
	tmpl, err := template.New("job_main").Parse(jobMainTemplate)
	if err != nil {
		return err
	}

	jobsInfo := &JobsData{
		JobInfos: make([]JobInfo, 0),
	}
	for _, v := range c.JobName {
		jobsInfo.JobInfos = append(jobsInfo.JobInfos, JobInfo{
			JobName:       v,
			GoModule:      c.GoMod,
			PackagePrefix: c.PackagePrefix,
		})
	}

	var jobFileContent bytes.Buffer
	data := struct {
		PackagePrefix string
		Version       string
	}{
		PackagePrefix: c.PackagePrefix,
		Version:       meta.Version,
	}
	err = tmpl.Execute(&jobFileContent, data)
	if err != nil {
		return err
	}
	err = utils.CreateFile(mainGoPath, jobFileContent.String())
	if err != nil {
		return err
	}
	jobFileContent.Reset()

	// Create or append to schedule.go
	internalDir := filepath.Join(c.OutDir, "internal")
	err = os.MkdirAll(internalDir, 0o755)
	if err != nil {
		return fmt.Errorf("failed to create internal directory: %w", err)
	}
	scheduleGoPath := filepath.Join(internalDir, "schedule.go")
	scheduleTmpl, err := template.New("job_schedule").Parse(jobScheduleTemplate)
	if err != nil {
		return err
	}

	if exist, _ := utils.PathExist(scheduleGoPath); !exist {
		err = scheduleTmpl.Execute(&jobFileContent, jobsInfo)
		if err != nil {
			return err
		}
		err = utils.CreateFile(scheduleGoPath, jobFileContent.String())
		if err != nil {
			return fmt.Errorf("failed to write schedule.go: %w", err)
		}
		jobFileContent.Reset()
	} else {
		src, err := utils.ReadFileContent(scheduleGoPath)
		if err != nil {
			return fmt.Errorf("failed to read schedule.go: %w", err)
		}

		res, err := addScheduleNewJobs(string(src), jobsInfo.JobInfos)
		if err != nil {
			return err
		}

		err = utils.CreateFile(scheduleGoPath, res)
		if err != nil {
			return err
		}
	}

	// Create or append to job_file
	jobDir := filepath.Join(internalDir, "job")
	err = os.MkdirAll(jobDir, 0o755)
	if err != nil {
		return fmt.Errorf("failed to create job directory: %w", err)
	}
	jobGoPath := filepath.Join(jobDir, c.JobFile)
	jobTmpl, err := template.New("job_file").Parse(jobTemplate)
	if err != nil {
		return err
	}

	if exist, _ := utils.PathExist(jobGoPath); !exist {
		err = jobTmpl.Execute(&jobFileContent, jobsInfo)
		if err != nil {
			return err
		}
		err = utils.CreateFile(jobGoPath, jobFileContent.String())
		if err != nil {
			return fmt.Errorf("failed to write job_file: %w", err)
		}
		jobFileContent.Reset()
	} else {
		src, err := utils.ReadFileContent(jobGoPath)
		if err != nil {
			return fmt.Errorf("failed to read job_file: %w", err)
		}

		res, err := addNewJobs(string(src), jobsInfo.JobInfos)
		if err != nil {
			return err
		}

		err = utils.CreateFile(jobGoPath, res)
		if err != nil {
			return err
		}
	}

	// Create or append to run.sh
	scriptsDir := filepath.Join(c.OutDir, "scripts")
	err = os.MkdirAll(scriptsDir, 0o755)
	if err != nil {
		return fmt.Errorf("failed to create scripts directory: %w", err)
	}

	// Create run.sh
	runShPath := filepath.Join(scriptsDir, "run.sh")
	scriptTmpl, err := template.New("job_script").Parse(scriptTemplate)
	if err != nil {
		return err
	}

	if exist, _ := utils.PathExist(runShPath); !exist {
		err = scriptTmpl.Execute(&jobFileContent, nil)
		if err != nil {
			return err
		}
		err = utils.CreateFile(runShPath, jobFileContent.String())
		if err != nil {
			return err
		}
		jobFileContent.Reset()
	}

	return nil
}
