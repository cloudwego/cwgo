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
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/cloudwego/cwgo/config"
	"github.com/cloudwego/cwgo/pkg/common/utils"
)

type PbUsedInfo struct {
	DocArgs  *config.DocArgument
	astFiles []*pbGoFileInfo // key:belongedToDir
}

type pbGoFileInfo struct {
	path          string
	belongedToDir string
	fSet          *token.FileSet
	astFile       *ast.File
	modifiedFile  string
}

func (info *PbUsedInfo) ParsePbIdl() (rawStructs []*IdlExtractStruct, err error) {
	info.astFiles, err = getPbGoFiles(info.DocArgs.ModelDir)
	if err != nil {
		return nil, err
	}

	for _, astFile := range info.astFiles {
		for _, v := range astFile.astFile.Decls {
			if stc, ok := v.(*ast.GenDecl); ok && stc.Tok == token.TYPE {
				hasInterface := false
				if stc.Doc != nil {
					if strings.Contains(stc.Doc.Text(), "mongo.") {
						hasInterface = true
					}
				}
				if hasInterface {
					for _, spec := range stc.Specs {
						if tp, ok := spec.(*ast.TypeSpec); ok {
							if stp, ok := tp.Type.(*ast.StructType); ok {
								if !stp.Struct.IsValid() {
									continue
								}
								rawStruct := newIdlExtractStruct(tp.Name.Name)
								if err = info.extractPbGoStruct(stp, rawStruct, astFile.astFile); err != nil {
									return nil, err
								}

								if len(rawStruct.StructFields) != 0 {
									rawStructs = append(rawStructs, rawStruct)

									if err = rawStruct.recordMongoIfInfo(info.DocArgs.DaoDir); err != nil {
										return nil, err
									}

									tokens, methods, err := getMongoIfTag(stc.Doc.Text())
									if err != nil {
										return nil, err
									}
									ifMethods := ""
									for _, m := range methods {
										ifMethods += m + "\n"
									}
									rawInterface := fmt.Sprintf("package main\ntype %sInterface interface{\n%s\n}", tp.Name.Name, ifMethods)
									if err = extractIdlInterface(rawInterface, rawStruct, tokens); err != nil {
										return nil, err
									}
								}
							}
						}
					}
				}
			}
		}
		buf := new(bytes.Buffer)
		err = printer.Fprint(buf, astFile.fSet, astFile.astFile)
		if err != nil {
			return nil, err
		}
		astFile.modifiedFile = buf.String()
	}
	return
}

func (info *PbUsedInfo) extractPbGoStruct(stNode *ast.StructType, rawStruct *IdlExtractStruct, astFile *ast.File) error {
	for _, field := range stNode.Fields.List {
		if field.Comment != nil {
			if strings.Contains(field.Comment.Text(), "go.tag") &&
				strings.Contains(field.Comment.Text(), "mongo.bson") {
				comment := getMongoStTag(field.Comment.Text())
				if comment == "" {
					return fmt.Errorf("there are grammar errors in %s", field.Comment.Text())
				}

				if field.Tag == nil {
					field.Tag = &ast.BasicLit{Kind: token.STRING, Value: comment}
				} else {
					if !strings.Contains(field.Tag.Value, "mongo.bson") {
						field.Tag = &ast.BasicLit{Kind: token.STRING, Value: field.Tag.Value[0:len(field.Tag.Value)-1] + " " + comment + "`"}
					}
				}

				tag := handleTagOmitempty(comment)

				fieldName := field.Names[0].Name
				t := getType(field.Type, astFile.Name.Name, true)
				if tt, ok := field.Type.(*ast.StarExpr); ok {
					// *Struct
					if ttt, ok := tt.X.(*ast.Ident); ok {
						rs := &IdlExtractStruct{
							Name:         ttt.Name,
							StructFields: make([]*StructField, 0, 10),
						}

						node := getStructNodeByName(astFile, ttt.Name)
						if node == nil {
							return fmt.Errorf("can not find %s field in struct %s", fieldName, rawStruct.Name)
						}
						if err := info.extractPbGoStruct(node, rs, astFile); err != nil {
							return err
						}
						rawStruct.StructFields = append(rawStruct.StructFields, &StructField{
							Name:               fieldName,
							Type:               t,
							Tag:                tag,
							IsBelongedToStruct: true,
							BelongedToStruct:   rs,
						})
					}
					// *pkgName.Struct
					if ttt, ok := tt.X.(*ast.SelectorExpr); ok {
						tttX := ttt.X.(*ast.Ident)
						astFiles := info.getAstFileByDir(tttX.Name)
						// provided by proto
						if len(astFiles) == 0 {
							rawStruct.StructFields = append(rawStruct.StructFields, &StructField{
								Name: fieldName,
								Type: t,
								Tag:  tag,
							})
						} else {
							var node *ast.StructType
							var f *ast.File
							for _, file := range astFiles {
								node = getStructNodeByName(file, ttt.Sel.Name)
								if node != nil {
									f = file
									break
								}
							}
							if node == nil {
								return fmt.Errorf("can not find %s field in struct %s", fieldName, rawStruct.Name)
							}

							rs := &IdlExtractStruct{
								Name:         ttt.Sel.Name,
								StructFields: make([]*StructField, 0, 10),
							}
							if err := info.extractPbGoStruct(node, rs, f); err != nil {
								return err
							}
							rawStruct.StructFields = append(rawStruct.StructFields, &StructField{
								Name:               fieldName,
								Type:               t,
								Tag:                tag,
								IsBelongedToStruct: true,
								BelongedToStruct:   rs,
							})
						}
					}
				} else {
					rawStruct.StructFields = append(rawStruct.StructFields, &StructField{
						Name: fieldName,
						Type: t,
						Tag:  tag,
					})
				}
			}
		}
	}
	return nil
}

func getMongoStTag(s string) (r string) {
	index := strings.Index(s, "go.tag")
	leftIndex, rightIndex := -1, -1
	count := 0
	for i := index; i < len(s); i++ {
		if string(s[i]) == "|" && count == 0 {
			leftIndex = i
			count++
			continue
		}
		if string(s[i]) == "|" && count == 1 {
			rightIndex = i
			count++
		}
	}
	if leftIndex == -1 || rightIndex == -1 {
		return ""
	} else {
		return s[leftIndex+1 : rightIndex]
	}
}

func getMongoIfTag(s string) (tokens, methods []string, err error) {
	if s == "" {
		return
	}

	index := strings.Index(s, "mongo.")
	if index == -1 {
		return
	}

	equalIndex := strings.Index(s, "=")
	if equalIndex == -1 || index+6 >= equalIndex {
		return nil, nil, fmt.Errorf("there are grammar errors in %s", s)
	}
	tokens = append(tokens, strings.Replace(s[index+6:equalIndex], " ", "", -1))

	leftIndex, rightIndex := -1, -1
	count := 0
	for i := index; i < len(s); i++ {
		if string(s[i]) == "|" && count == 0 {
			leftIndex = i
			count++
			continue
		}
		if string(s[i]) == "|" && count == 1 {
			rightIndex = i
			count++
		}
	}
	if leftIndex == -1 || rightIndex == -1 || leftIndex+1 == rightIndex {
		return nil, nil, fmt.Errorf("there are grammar errors in %s", s)
	}
	methods = append(methods, strings.Replace(s[leftIndex+1:rightIndex], "\n", "", -1))

	if rightIndex+1 == len(s) {
		return
	} else {
		ts, ms, err := getMongoIfTag(s[rightIndex+1:])
		if err != nil {
			return nil, nil, err
		}
		tokens = append(tokens, ts...)
		methods = append(methods, ms...)
	}
	return
}

func (info *PbUsedInfo) getAstFileByDir(dir string) (result []*ast.File) {
	for _, astFile := range info.astFiles {
		if astFile.belongedToDir == dir {
			result = append(result, astFile.astFile)
		}
	}
	return
}

func getPbGoFiles(dir string) (result []*pbGoFileInfo, err error) {
	err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if strings.Contains(path, ".pb.go") {
			f := token.NewFileSet()
			p, err := parser.ParseFile(f, path, nil, parser.ParseComments)
			if err != nil {
				return err
			}
			belongDir := filepath.Base(filepath.Dir(path))
			result = append(result, &pbGoFileInfo{
				fSet:          f,
				path:          path,
				belongedToDir: belongDir,
				astFile:       p,
			})
		}
		return nil
	})
	return
}

func isGoBaseType(s string) bool {
	return s == "bool" || s == "int8" || s == "int16" || s == "int32" || s == "int64" || s == "int" ||
		s == "uint8" || s == "uint16" || s == "uint32" || s == "uint64" || s == "uint" || s == "float32" ||
		s == "float64" || s == "string" || s == "byte"
}

func getStructNodeByName(file *ast.File, name string) (result *ast.StructType) {
	ast.Inspect(file, func(n ast.Node) bool {
		ts, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}
		tts, ok := ts.Type.(*ast.StructType)
		if !ok {
			return true
		}
		if ts.Name.Name == name {
			result = tts
		}
		return true
	})
	return
}

func (info *PbUsedInfo) GeneratePbFile() error {
	for _, pbGo := range info.astFiles {
		if err := utils.CreateFile(pbGo.path, pbGo.modifiedFile); err != nil {
			return err
		}
	}
	return nil
}
