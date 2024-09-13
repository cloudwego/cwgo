// Copyright 2024 CloudWeGo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package docker

import (
	"errors"
	"github.com/cloudwego/cwgo/config"
	"github.com/cloudwego/cwgo/pkg/common/utils"
	"github.com/cloudwego/cwgo/pkg/consts"
	"github.com/cloudwego/cwgo/tpl"
	"path"
	"path/filepath"
	"strings"
	"text/template"
)

const (
	dockerfileName = "Dockerfile"
	etcDir         = "etc"
	yamlEtx        = ".yaml"
)

type DockerTemplateParams struct {
	GoVersion     string
	EnableGoProxy bool
	GoProxy       string
	HasTimezone   bool
	Timezone      string
	GoRelPath     string
	GoFileName    string
	ExeFile       string
	GoMainFrom    string
	BaseImage     string
	Argument      string
	HasPort       bool
	Port          int
}

func Docker(c *config.DockerArgument) error {
	if strings.HasSuffix(c.Template, consts.SuffixGit) {
		// pull remote template
		err := utils.GitClone(c.Template, tpl.DockerDir)
		if err != nil {
			return err
		}
		gitPath, err := utils.GitPath(c.Template)
		if err != nil {
			return err
		}
		gitPath = path.Join(tpl.DockerDir, gitPath)
		if err = utils.GitCheckout(c.Branch, gitPath); err != nil {
			return err
		}
		c.Template = gitPath
	} else {
		isExist, _ := utils.PathExist(path.Join(c.Template, tpl.DockerFileTpl))
		if !isExist {
			return errors.New("DockerFile template not exist")
		}
	}

	if len(c.GoVersion) > 0 {
		c.GoVersion = c.GoVersion + "-"
	}

	if len(c.GoFileName) > 0 {
		isExist, _ := utils.PathExist(c.GoFileName)
		if !isExist {
			return errors.New("go file not exist")
		}
	}

	var projectPath string
	if len(c.GoFileName) > 0 {
		curpath, err := filepath.Abs(consts.CurrentDir)
		if err != nil {
			return err
		}
		_, projPath, isOk := utils.SearchGoMod(curpath, true)
		if !isOk {
			return errors.New("not found go mod")
		}
		projectPath = projPath
	}

	if len(projectPath) == 0 {
		projectPath = "."
	}

	dockerFile, err := utils.CreateIfNotExist(dockerfileName)
	if err != nil {
		return err
	}
	defer dockerFile.Close()

	text, err := utils.ReadFileContent(path.Join(c.Template, tpl.DockerFileTpl))
	if err != nil {
		return err
	}

	var exeName string
	if len(c.ExeFileName) > 0 {
		exeName = c.ExeFileName
	} else if len(c.GoFileName) > 0 {
		exeName = strings.TrimSuffix(filepath.Base(c.GoFileName), filepath.Ext(c.GoFileName))
	} else {
		absPath, err := filepath.Abs(projectPath)
		if err != nil {
			return err
		}

		exeName = filepath.Base(absPath)
	}

	t := template.Must(template.New("DockerFile").Parse(string(text)))
	return t.Execute(dockerFile, DockerTemplateParams{
		GoVersion:     c.GoVersion,
		EnableGoProxy: c.EnableGoProxy,
		GoProxy:       c.GoProxy,
		HasTimezone:   c.Timezone != "",
		Timezone:      c.Timezone,
		GoRelPath:     projectPath,
		GoFileName:    c.GoFileName,
		ExeFile:       exeName,
		GoMainFrom:    path.Join(projectPath, c.GoFileName),
		BaseImage:     c.BaseImage,
		Argument:      c.RunArgs,
		HasPort:       c.Port > 0,
		Port:          c.Port,
	})
}
