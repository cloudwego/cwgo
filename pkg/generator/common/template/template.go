/*
 * Copyright 2023 CloudWeGo Authors
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

package template

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	geneUtils "github.com/cloudwego/cwgo/pkg/generator/common/utils"

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
	AppendImport  map[string]string
}

type ReplaceFunc struct {
	ReplaceFuncName         []string
	ReplaceFuncBody         []string
	ReplaceFuncAppendImport []map[string]string
	ReplaceFuncDeleteImport []map[string]string
}

type Generator struct {
	OutputDir string
	Files     []File
}

func (tg *Generator) RenderCwgoTemplateFile(tpl *Template, renderObj any) (err error) {
	if isExist, _ := utils.PathExist(tpl.Path); !isExist {
		tpl.Type = consts.Cover
	}

	// skip
	if tpl.Type == consts.Skip || tpl.Type == "" {
		return
	}

	// cover
	if tpl.Type == consts.Cover {
		if err = tg.renderPathBody(tpl, renderObj); err != nil {
			return err
		}
		return
	}

	// append
	if tpl.Type == consts.Append {
		if err = tg.renderAppend(tpl, renderObj); err != nil {
			return err
		}
		return
	}

	// replaceFuncBody
	if tpl.Type == consts.ReplaceFuncBody && strings.HasSuffix(tpl.Path, ".go") {
		if err = tg.renderReplaceFuncBody(tpl, renderObj); err != nil {
			return err
		}
		return
	}

	return
}

func (tg *Generator) renderAppend(tpl *Template, renderObj any) error {
	// render path
	if tpl.IsPathRender {
		data, err := Render(filepath.Base(tpl.Path)+"path", tpl.Path, renderObj, tpl)
		if err != nil {
			return err
		}

		tpl.Path = data.String()
	}

	// render append body
	if len(tpl.AppendRender) != 0 {
		data, err := Render(filepath.Base(tpl.Path)+"body", tpl.AppendContent, tpl.AppendRender, tpl)
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
		importedContent, err = geneUtils.HandleGoFileImports(string(fileContent), tpl.AppendImport, true)
		if err != nil {
			return err
		}
	}

	// append body
	if !strings.HasSuffix(tpl.Path, ".go") {
		importedContent = string(fileContent)
	}
	importedContent += consts.LineBreak + tpl.AppendContent

	file := File{Path: tpl.Path, Content: []byte(importedContent)}
	tg.Files = append(tg.Files, file)

	return nil
}

func (tg *Generator) renderReplaceFuncBody(tpl *Template, renderObj any) error {
	// render path
	if tpl.IsPathRender {
		data, err := Render(filepath.Base(tpl.Path)+"path", tpl.Path, renderObj, tpl)
		if err != nil {
			return err
		}

		tpl.Path = data.String()
	}

	// render append body
	if len(tpl.AppendRender) != 0 {
		for index, body := range tpl.ReplaceFuncBody {
			data, err := Render(filepath.Base(tpl.Path)+"body", body, tpl.AppendRender, tpl)
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

	replaceFuncAppendImpts := make(map[string]string, 10)
	for _, impt := range tpl.ReplaceFuncAppendImport {
		for k, v := range impt {
			replaceFuncAppendImpts[k] = v
		}
	}

	content, err := geneUtils.HandleGoFileImports(string(fileContent), replaceFuncAppendImpts, true)
	if err != nil {
		return err
	}

	replaceFuncDeleteImpts := make(map[string]string, 10)
	for _, impt := range tpl.ReplaceFuncDeleteImport {
		for k, v := range impt {
			replaceFuncDeleteImpts[k] = v
		}
	}

	content, err = geneUtils.HandleGoFileImports(content, replaceFuncDeleteImpts, false)
	if err != nil {
		return err
	}

	content, err = geneUtils.ReplaceFuncBody(content, tpl.ReplaceFuncName, tpl.ReplaceFuncBody)
	if err != nil {
		return err
	}

	file := File{Path: tpl.Path, Content: []byte(content)}
	tg.Files = append(tg.Files, file)

	return nil
}

func (tg *Generator) renderPathBody(tpl *Template, renderObj any) (err error) {
	// render path
	if tpl.IsPathRender {
		data, err := Render(filepath.Base(tpl.Path)+"path", tpl.Path, renderObj, tpl)
		if err != nil {
			return err
		}

		tpl.Path = data.String()
	}

	// render body
	data, err := Render(filepath.Base(tpl.Path)+"body", tpl.Body, renderObj, tpl)
	if err != nil {
		return err
	}

	file := File{Path: tpl.Path, Content: data.Bytes()}
	tg.Files = append(tg.Files, file)

	return nil
}

func Render(name, parseBody string, data any, tpl *Template) (*bytes.Buffer, error) {
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

func (tg *Generator) Persist() error {
	outPath := tg.OutputDir
	if !filepath.IsAbs(outPath) {
		outPath, _ = filepath.Abs(outPath)
	}

	for _, data := range tg.Files {
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
