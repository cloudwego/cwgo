package config

import (
	"github.com/cloudwego/cwgo/pkg/consts"
	"github.com/urfave/cli/v2"
)

type ApiArgument struct {
	ProjectPath  string
	HertzRepoUrl string
}

func NewApiArgument() *ApiArgument {
	return &ApiArgument{}
}

func (c *ApiArgument) ParseCli(ctx *cli.Context) error {
	c.ProjectPath = ctx.String(consts.ProjectPath)
	c.HertzRepoUrl = ctx.String(consts.HertzRepoUrl)
	return nil
}
