package api_list

import (
	"fmt"
	"path/filepath"

	"github.com/cloudwego/cwgo/config"
	"github.com/cloudwego/cwgo/pkg/consts"
)

func Api(c *config.ApiArgument) error {
	if c.ProjectPath == "" {
		curPath, err := filepath.Abs(".")
		if err != nil {
			return fmt.Errorf("get current path failed, err: %v", err)
		}
		c.ProjectPath = curPath
	}

	if c.HertzRepoUrl == "" {
		c.HertzRepoUrl = consts.HertzRepoDefaultUrl
	}

	parser, err := NewParser(c.ProjectPath, c.HertzRepoUrl)
	if err != nil {
		return err
	}

	moduleName, err := getModuleName(c.ProjectPath)
	if err != nil {
		return err
	}

	fmt.Printf("found module name: %s\n", parser.moduleName)

	err = parser.searchFunc(moduleName, "main", make(map[string]*Var), nil)
	if err != nil {
		return err
	}

	parser.PrintRouters()

	return nil
}
