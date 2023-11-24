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
	Path   string
	Delims [2]string
	Body   string
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

func (tg *TemplateGenerator) renderServer(serverGen *ServerGenerator) error {
	var mvcTemplates []Template

	switch serverGen.communicationType {
	case consts.RPC:
		mvcTemplates = kitexServerMVCTemplates
	case consts.HTTP:
		mvcTemplates = hzServerMVCTemplates
	default:
		return typeInputErr
	}

	for _, tpl := range mvcTemplates {
		buffer := &bytes.Buffer{}
		name := filepath.Base(tpl.Path)
		t := template.Must(template.New(name).Parse(tpl.Body))
		if err := t.Execute(buffer, serverGen.ServerRender); err != nil {
			return err
		}

		file := File{Path: tpl.Path, Content: buffer.Bytes()}
		tg.files = append(tg.files, file)
	}

	return nil
}

func (tg *TemplateGenerator) renderClient(clientGen *ClientGenerator) error {
	var mvcTemplates []Template

	switch clientGen.communicationType {
	case consts.RPC:
		mvcTemplates = kitexClientMVCTemplates
	case consts.HTTP:
		mvcTemplates = hzClientTemplates
	default:
		return typeInputErr
	}

	for index, tpl := range mvcTemplates {
		// render init.go
		if index == consts.HzClientInitFileIndex && clientGen.communicationType == consts.HTTP {
			for _, name := range clientGen.SnakeServiceNames {
				clientGen.CurrentIDLServiceName = name

				// render path
				bufferPath := &bytes.Buffer{}
				t := template.Must(template.New(name).Parse(tpl.Path))
				if err := t.Execute(bufferPath, clientGen); err != nil {
					return err
				}

				// render body
				bufferBody := &bytes.Buffer{}
				t = template.Must(template.New(name + "body").Parse(tpl.Body))
				if err := t.Execute(bufferBody, clientGen); err != nil {
					return err
				}

				file := File{Path: bufferPath.String(), Content: bufferBody.Bytes()}
				tg.files = append(tg.files, file)
			}
			continue
		}

		buffer := &bytes.Buffer{}
		name := filepath.Base(tpl.Path)
		t := template.Must(template.New(name).Parse(tpl.Body))
		if err := t.Execute(buffer, clientGen); err != nil {
			return err
		}

		file := File{Path: tpl.Path, Content: buffer.Bytes()}
		tg.files = append(tg.files, file)
	}

	return nil
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

		// create rendered file
		abPath := filepath.Join(outPath, data.Path)
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

		err = func() error {
			file, err := os.OpenFile(abPath, os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.FileMode(0o755))
			defer file.Close()
			if err != nil {
				return fmt.Errorf("open file '%s' failed, err: %v", abPath, err.Error())
			}
			if _, err = file.Write(data.Content); err != nil {
				return fmt.Errorf("write file '%s' failed, err: %v", abPath, err.Error())
			}

			return nil
		}()
		if err != nil {
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
