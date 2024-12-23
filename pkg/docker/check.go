package docker

import (
	"errors"
	"github.com/cloudwego/cwgo/config"
	"github.com/cloudwego/cwgo/pkg/common/utils"
	"github.com/cloudwego/cwgo/pkg/consts"
	"path/filepath"
)

func check(c *config.DockerArgument) error {
	if c.Main == "" {
		return errors.New("go main file must be provided")
	}

	if len(c.Version) > 0 {
		c.Version = c.Version + consts.Dash
	}

	if len(c.Main) > 0 {
		isExist, _ := utils.PathExist(c.Main)
		if !isExist {
			return errors.New("go main not exist")
		}
	}

	curPath, err := filepath.Abs(consts.CurrentDir)
	if err != nil {
		return err
	}
	_, _, isOk := utils.SearchGoMod(curPath, true)
	if !isOk {
		return errors.New("not found go mod")
	}

	if c.Port > 0 && c.Port < 1024 || c.Port < 0 || c.Port > 65535 {
		return errors.New("port must between 1024 and 65535")
	}

	isExist, err := utils.PathExist(c.Template)
	if err != nil {
		return err
	}
	if !isExist {
		return errors.New("DockerFile template not exist")
	}
	return nil
}
