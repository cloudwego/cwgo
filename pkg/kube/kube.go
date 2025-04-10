package kube

import (
	"fmt"
	"github.com/cloudwego/cwgo/config"
	"github.com/cloudwego/cwgo/pkg/common/utils"
	"github.com/cloudwego/cwgo/pkg/consts"
	"github.com/cloudwego/cwgo/tpl"
	"path"
	"path/filepath"
	"strings"
	"text/template"
)

func pullGitTpl(c *config.KubeArgument) error {
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
	return nil
}

func Kube(c *config.KubeArgument) error {
	if err := check(c); err != nil {
		return err
	}

	dockerfileInfo := FillInfo(c)

	var tplFile = new(KubeDeployTpl)
	var url string
	if strings.HasSuffix(c.Template, consts.SuffixGit) {
		if err := pullGitTpl(c); err != nil {
			return err
		}
		url = path.Join(c.Template, consts.DefaultDockerfileTpl)
	} else if strings.HasSuffix(c.Template, consts.Yaml) {
		url = c.Template
	} else {
		url = path.Join(c.Template, path.Join(consts.Standard, consts.DefaultKubeDeployTpl))
	}
	if err := tplFile.FromYAMLFile(url); err != nil {
		return err
	}
	t := template.Must(template.New(consts.Dockerfile).Parse(tplFile.Body))

	wr, err := utils.WriteFile(c.Output)
	if err != nil {
		return err
	}
	defer wr.Close()

	if err := t.Execute(wr, dockerfileInfo); err != nil {
		return err
	}
	abs, err := filepath.Abs(filepath.Join(consts.CurrentDir, c.Output))
	if err != nil {
		return err
	}
	fmt.Println("Output to:", abs)
	fmt.Println("Done.")
	return nil
}
