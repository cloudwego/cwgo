package docker

import (
	"fmt"
	"github.com/cloudwego/cwgo/config"
	"github.com/cloudwego/cwgo/pkg/common/utils"
	"github.com/cloudwego/cwgo/pkg/consts"
	"github.com/cloudwego/cwgo/tpl"
	"github.com/cloudwego/kitex/tool/internal_pkg/log"
	"path"
	"path/filepath"
	"strings"
	"text/template"
)

func pullGitTpl(c *config.DockerArgument) error {
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

func Docker(c *config.DockerArgument) error {
	if err := check(c); err != nil {
		return err
	}

	log.Verbose = c.Verbose

	dockerfileInfo := FillInfo(c)

	var tplFile = new(DockerfilesTpl)
	var url string
	if strings.HasSuffix(c.Template, consts.SuffixGit) {
		if err := pullGitTpl(c); err != nil {
			return err
		}
		url = path.Join(c.Template, consts.DefaultDockerfileTpl)
	} else if strings.HasSuffix(c.Template, consts.Yaml) {
		url = c.Template
	} else {
		url = path.Join(c.Template, path.Join(consts.Standard, consts.DefaultDockerfileTpl))
	}
	if err := tplFile.FromYAMLFile(url); err != nil {
		return err
	}
	t := template.Must(template.New(consts.Dockerfile).Parse(tplFile.Body))

	wr, err := utils.WriteFile(tplFile.Path)
	if err != nil {
		return err
	}
	defer wr.Close()

	if err := t.Execute(wr, dockerfileInfo); err != nil {
		return err
	}
	var str, _ = filepath.Abs(consts.CurrentDir)
	fmt.Print("Hint: run \"docker build ...\" command in dir:\n" + "\t" + str + "\nDone.")
	return nil
}
