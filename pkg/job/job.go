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

package job

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
	"strconv"
	"strings"
	"text/template"

	"github.com/cloudwego/cwgo/meta"

	"github.com/cloudwego/cwgo/config"
	"github.com/cloudwego/cwgo/pkg/common/utils"
	"github.com/cloudwego/cwgo/pkg/consts"
	"github.com/cloudwego/kitex/tool/internal_pkg/log"

	"golang.org/x/tools/go/ast/astutil"
)

func Job(c *config.JobArgument) error {
	if err := check(c); err != nil {
		return err
	}

	err := generateJobFile(c.GoMod, c.PackagePrefix, c.JobName, c.OutDir)
	if err != nil {
		return err
	}

	return nil
}

func check(c *config.JobArgument) (err error) {
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

func addJobImportsAndRun(data string, jobs []JobInfo) (string, error) {
	fSet := token.NewFileSet()
	file, err := parser.ParseFile(fSet, "", data, parser.ParseComments)
	if err != nil {
		return "", err
	}

	// Extract existing imports and Run function calls
	existingJobs := make(map[string]bool)
	ast.Inspect(file, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.ImportSpec:
			if path, err := strconv.Unquote(x.Path.Value); err == nil {
				for _, job := range jobs {
					if path == fmt.Sprintf(`"%s/%s/job"`, job.PackagePrefix, job.JobName) {
						existingJobs[job.JobName] = true
					}
				}
			}
		case *ast.ExprStmt:
			if call, ok := x.X.(*ast.CallExpr); ok {
				if fun, ok := call.Fun.(*ast.SelectorExpr); ok {
					if ident, ok := fun.X.(*ast.Ident); ok {
						for _, job := range jobs {
							if ident.Name == job.JobName && fun.Sel.Name == "Run" {
								existingJobs[job.JobName] = true
							}
						}
					}
				}
			}
		}
		return true
	})

	// Add missing imports
	for _, job := range jobs {
		if !existingJobs[job.JobName] {
			astutil.AddNamedImport(fSet, file, job.JobName, fmt.Sprintf("%s/%s/job", job.PackagePrefix, job.JobName))
		}
	}

	// Add missing Run calls
	var runFunc *ast.FuncDecl
	for _, decl := range file.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok && fn.Name.Name == "Run" {
			runFunc = fn
			break
		}
	}

	if runFunc != nil && runFunc.Body != nil {
		var newStatements []ast.Stmt
		for _, job := range jobs {
			if !existingJobs[job.JobName] {
				newStatements = append(newStatements, createRunCall(job.JobName)...)
			}
		}

		// Find the wg.Wait() block and insert new statements before it
		var insertIndex int
		for i, stmt := range runFunc.Body.List {
			if exprStmt, ok := stmt.(*ast.ExprStmt); ok {
				if call, ok := exprStmt.X.(*ast.CallExpr); ok {
					if fun, ok := call.Fun.(*ast.SelectorExpr); ok {
						if fun.X.(*ast.Ident).Name == "wg" && fun.Sel.Name == "Wait" {
							insertIndex = i
							break
						}
					}
				}
			}
		}

		// Insert new statements before wg.Wait()
		runFunc.Body.List = append(runFunc.Body.List[:insertIndex], append(newStatements, runFunc.Body.List[insertIndex:]...)...)
	}

	// Generate the modified code
	buf := new(bytes.Buffer)
	if err = printer.Fprint(buf, fSet, file); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func createRunCall(jobName string) []ast.Stmt {
	return []ast.Stmt{
		&ast.ExprStmt{
			X: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   ast.NewIdent("wg"),
					Sel: ast.NewIdent("Add"),
				},
				Args: []ast.Expr{
					&ast.BasicLit{
						Kind:  token.INT,
						Value: "1",
					},
				},
			},
		},
		&ast.GoStmt{
			Call: &ast.CallExpr{
				Fun: &ast.FuncLit{
					Type: &ast.FuncType{
						Params: &ast.FieldList{},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.DeferStmt{
								Call: &ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   ast.NewIdent("wg"),
										Sel: ast.NewIdent("Done"),
									},
								},
							},
							&ast.ExprStmt{
								X: &ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   ast.NewIdent(jobName),
										Sel: ast.NewIdent("Run"),
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func generateJobFile(GoModule, PackagePrefix string, jobNames []string, outDir string) error {
	// Ensure the base output directory exists
	err := os.MkdirAll(outDir, 0o755)
	if err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Create cmd/main.go and overwrite each time
	cmdDir := filepath.Join(outDir, "cmd")
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
	for _, v := range jobNames {
		jobsInfo.JobInfos = append(jobsInfo.JobInfos, JobInfo{
			JobName:       v,
			GoModule:      GoModule,
			PackagePrefix: PackagePrefix,
		})
	}

	var jobFileContent bytes.Buffer
	data := struct {
		PackagePrefix string
		Version       string
	}{
		PackagePrefix: PackagePrefix,
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
	scheduleGoPath := filepath.Join(outDir, "schedule.go")
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

		res, err := addJobImportsAndRun(string(src), jobsInfo.JobInfos)
		if err != nil {
			return err
		}

		err = utils.CreateFile(scheduleGoPath, res)
		if err != nil {
			return err
		}
	}

	// Create or append to run.sh
	scriptsDir := filepath.Join(outDir, "scripts")
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

	// Create job directories and files
	for _, jobName := range jobNames {
		jobDir := filepath.Join(outDir, jobName)
		internalJobDir := filepath.Join(jobDir, "job")

		// Create directories
		err = os.MkdirAll(internalJobDir, 0o755)
		if err != nil {
			return fmt.Errorf("failed to create internal job directory for %s: %w", jobName, err)
		}

		// Create or append to job.go
		jobFilePath := filepath.Join(internalJobDir, "job.go")
		jobTmpl, err := template.New("job").Parse(jobTemplate)
		if err != nil {
			return err
		}

		if exist, _ := utils.PathExist(jobFilePath); !exist {
			jobInfo := struct {
				JobName string
			}{
				JobName: jobName,
			}

			err = jobTmpl.Execute(&jobFileContent, jobInfo)
			if err != nil {
				return err
			}
			err = utils.CreateFile(jobFilePath, jobFileContent.String())
			if err != nil {
				return err
			}
			jobFileContent.Reset()
		}
	}

	return nil
}
