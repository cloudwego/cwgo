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
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/cloudwego/cwgo/pkg/common/utils"
	"github.com/cloudwego/cwgo/pkg/consts"
)

type Template struct {
	Path         string
	Delims       [2]string
	Body         string
	IsPathRender bool
	CustomFunc   template.FuncMap
	UpdateBehavior
}

type UpdateBehavior struct {
	Type string // skip/cover/append/replaceFuncBody default:skip
	Append
	ReplaceFunc
	AppendRender map[string]interface{}
}

type Append struct {
	AppendContent string
	AppendImport  []string
}

type ReplaceFunc struct {
	ReplaceFuncName   []string
	ReplaceFuncBody   []string
	ReplaceFuncImport [][]string
}

type TemplateGenerator struct {
	OutputDir string
	files     []File
}

// GenerateServer generate cwgo side server files
func GenerateServer(serverGen *ServerGenerator) error {
	tg := &TemplateGenerator{
		files: make([]File, 0, 10),
	}

	// render template
	if err := tg.renderServer(serverGen); err != nil {
		return err
	}

	// generate files
	if err := tg.persist(); err != nil {
		return err
	}

	// generate .cwgo file
	if err := serverGen.manifest.Persist(tg.OutputDir); err != nil {
		return err
	}

	return nil
}

// GenerateClient generate cwgo side client files
func GenerateClient(clientGen *ClientGenerator) error {
	tg := &TemplateGenerator{
		files: make([]File, 0, 10),
	}

	// render template
	if err := tg.renderClient(clientGen); err != nil {
		return err
	}

	// generate files
	if err := tg.persist(); err != nil {
		return err
	}

	// generate .cwgo file
	if err := clientGen.manifest.Persist(tg.OutputDir); err != nil {
		return err
	}

	return nil
}

func (tg *TemplateGenerator) renderServer(serverGen *ServerGenerator) (err error) {
	var mvcTemplates []Template

	switch serverGen.communicationType {
	case consts.RPC:
		mvcTemplates = kitexServerMVCTemplates
	case consts.HTTP:
		mvcTemplates = hzServerMVCTemplates
	default:
		return errTypeInput
	}

	for _, tpl := range mvcTemplates {
		if err = tg.renderPathBody(&tpl, serverGen.ServerRender); err != nil {
			return err
		}
	}

	return nil
}

func (tg *TemplateGenerator) renderClient(clientGen *ClientGenerator) error {
	var mvcTemplates []Template

	switch clientGen.communicationType {
	case consts.RPC:
		mvcTemplates = kitexClientMVCTemplates
	case consts.HTTP:
		mvcTemplates = hzClientMVCTemplates
	default:
		return errTypeInput
	}

	if clientGen.isNew {
		for index, tpl := range mvcTemplates {
			// render init.go
			if index == consts.FileClientInitIndex {
				bizDir := ""
				if clientGen.communicationType == consts.HTTP {
					bizDir = clientGen.OutDir
				} else {
					bizDir = filepath.Join(clientGen.OutDir, consts.DefaultKitexClientDir)
				}

				subDirs, err := utils.GetSubDirs(bizDir, false)
				if err != nil {
					return err
				}

				for _, name := range subDirs {
					clientGen.InitOptsPackage = filepath.Base(name)

					// render body
					data, err := render(name, tpl.Body, clientGen.ClientRender, &tpl)
					if err != nil {
						return err
					}

					file := File{Path: filepath.Join(name, consts.InitGo), Content: data.Bytes()}
					tg.files = append(tg.files, file)
				}
				continue
			}

			if err := tg.renderPathBody(&tpl, clientGen.ClientRender); err != nil {
				return err
			}
		}
	} else {
		for _, tpl := range mvcTemplates {
			// handle append render
			if tpl.Type == consts.Append || tpl.Type == consts.ReplaceFuncBody {
				tpl.AppendRender["GoModule"] = clientGen.GoModule
				tpl.AppendRender["ServiceName"] = clientGen.ServiceName
				if clientGen.ResolverName != "" {
					tpl.AppendRender["ResolverName"] = clientGen.ResolverName
				}
				if clientGen.ResolverAddress != nil {
					tpl.AppendRender["ResolverAddress"] = clientGen.ResolverAddress
				}
			}

			// skip
			if tpl.Type == consts.Skip || tpl.Type == "" {
				continue
			}

			// cover
			if tpl.Type == consts.Cover {
				if err := tg.renderPathBody(&tpl, clientGen.ClientRender); err != nil {
					return err
				}
				continue
			}

			// append
			if tpl.Type == consts.Append {
				if err := tg.renderAppend(&tpl, clientGen.ClientRender); err != nil {
					return err
				}
				continue
			}

			// replaceFuncBody
			if tpl.Type == consts.ReplaceFuncBody {
				if err := tg.renderReplaceFuncBody(&tpl, clientGen.ClientRender); err != nil {
					return err
				}
				continue
			}
		}
	}

	return nil
}

func (tg *TemplateGenerator) renderAppend(tpl *Template, renderObj any) error {
	// render path
	if tpl.IsPathRender {
		data, err := render(filepath.Base(tpl.Path)+"path", tpl.Path, renderObj, tpl)
		if err != nil {
			return err
		}

		tpl.Path = data.String()
	}

	// render append body
	if len(tpl.AppendRender) != 0 {
		data, err := render(filepath.Base(tpl.Path)+"body", tpl.AppendContent, tpl.AppendRender, tpl)
		if err != nil {
			return err
		}

		tpl.AppendContent = data.String()
	}

	// read file content
	if isExist, _ := utils.PathExist(tpl.Path); !isExist {
		return fmt.Errorf("file %v does not exist", tpl.Path)
	}
	fileContent, err := utils.ReadFileContent(tpl.Path)
	if err != nil {
		return err
	}

	// append imports
	var importedContent string
	if strings.HasSuffix(tpl.Path, ".go") {
		importedContent, err = appendGoFileImports(string(fileContent), tpl.AppendImport)
		if err != nil {
			return err
		}
	}

	// append body
	importedContent += consts.LineBreak + tpl.AppendContent

	file := File{Path: tpl.Path, Content: []byte(importedContent)}
	tg.files = append(tg.files, file)

	return nil
}

func (tg *TemplateGenerator) renderReplaceFuncBody(tpl *Template, renderObj any) error {
	// render path
	if tpl.IsPathRender {
		data, err := render(filepath.Base(tpl.Path)+"path", tpl.Path, renderObj, tpl)
		if err != nil {
			return err
		}

		tpl.Path = data.String()
	}

	// render append body
	if len(tpl.AppendRender) != 0 {
		for index, body := range tpl.ReplaceFuncBody {
			data, err := render(filepath.Base(tpl.Path)+"body", body, tpl.AppendRender, tpl)
			if err != nil {
				return err
			}

			tpl.ReplaceFuncBody[index] = data.String()
		}
	}

	// read file content
	if isExist, _ := utils.PathExist(tpl.Path); !isExist {
		return fmt.Errorf("file %v does not exist", tpl.Path)
	}
	fileContent, err := utils.ReadFileContent(tpl.Path)
	if err != nil {
		return err
	}

	replaceFuncImpts := make([]string, 0, 10)
	for _, impt := range tpl.ReplaceFuncImport {
		replaceFuncImpts = append(replaceFuncImpts, impt...)
	}

	content, err := appendGoFileImports(string(fileContent), replaceFuncImpts)
	if err != nil {
		return err
	}

	content, err = replaceFuncBody(content, tpl.ReplaceFuncName, tpl.ReplaceFuncBody)
	if err != nil {
		return err
	}

	file := File{Path: tpl.Path, Content: []byte(content)}
	tg.files = append(tg.files, file)

	return nil
}

func (tg *TemplateGenerator) renderPathBody(tpl *Template, renderObj any) (err error) {
	// render path
	if tpl.IsPathRender {
		data, err := render(filepath.Base(tpl.Path)+"path", tpl.Path, renderObj, tpl)
		if err != nil {
			return err
		}

		tpl.Path = data.String()
	}

	// render body
	data, err := render(filepath.Base(tpl.Path)+"body", tpl.Body, renderObj, tpl)
	if err != nil {
		return err
	}

	file := File{Path: tpl.Path, Content: data.Bytes()}
	tg.files = append(tg.files, file)

	return nil
}

func render(name, parseBody string, data any, tpl *Template) (*bytes.Buffer, error) {
	buffer := &bytes.Buffer{}
	var t *template.Template
	if tpl.CustomFunc == nil {
		t = template.Must(template.New(name).Delims(tpl.Delims[0], tpl.Delims[1]).Parse(parseBody))
	} else {
		t = template.Must(template.New(name).Delims(tpl.Delims[0], tpl.Delims[1]).Funcs(tpl.CustomFunc).Parse(parseBody))
	}
	if err := t.Execute(buffer, data); err != nil {
		return nil, err
	}

	return buffer, nil
}

func (tg *TemplateGenerator) persist() error {
	outPath := tg.OutputDir
	if !filepath.IsAbs(outPath) {
		outPath, _ = filepath.Abs(outPath)
	}

	for _, data := range tg.files {
		// lint file
		if err := data.Lint(); err != nil {
			return err
		}

		abPath := filepath.Join(outPath, data.Path)
		if filepath.IsAbs(data.Path) {
			abPath = data.Path
		}
		abDir := filepath.Dir(abPath)
		isExist, err := utils.PathExist(abDir)
		if err != nil {
			return fmt.Errorf("check directory '%s' failed, err: %v", abDir, err.Error())
		}
		if !isExist {
			if err = os.MkdirAll(abDir, os.FileMode(0o744)); err != nil {
				return fmt.Errorf("mkdir %s failed, err: %v", abDir, err.Error())
			}
		}

		if err = utils.CreateFile(abPath, string(data.Content)); err != nil {
			return err
		}

		name := filepath.Base(data.Path)
		if strings.HasSuffix(name, ".go") {
			if err = utils.FormatGoFile(abPath); err != nil {
				return err
			}
		}
	}

	return nil
}
