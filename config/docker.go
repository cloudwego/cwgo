package config

import (
	"github.com/cloudwego/cwgo/pkg/consts"
	"github.com/urfave/cli/v2"
)

type DockerArgument struct {
	BaseImage   string
	Branch      string
	ExeName     string
	Main        string
	Template    string
	Port        uint
	Tz          string
	Version     string
	Verbose     bool
	Mirrors     []string
	Arguments   []string
	EtcDirs     []string
	RunCmd      []string
	Env         []string
	BuildEnv    []string
	BuildRunCmd []string
}

func NewDockerArgument() *DockerArgument {
	return &DockerArgument{}
}

func (s *DockerArgument) ParseCli(ctx *cli.Context) error {
	s.Template = ctx.String(consts.Template)
	s.BaseImage = ctx.String(consts.Base)
	s.ExeName = ctx.String(consts.Exe)
	s.Branch = ctx.String(consts.Branch)
	s.Verbose = ctx.Bool(consts.Verbose)
	s.Main = ctx.String(consts.Go)
	s.Port = ctx.Uint(consts.Port)
	s.Tz = ctx.String(consts.TZ)
	s.Version = ctx.String(consts.Version)
	s.Mirrors = ctx.StringSlice(consts.Mirror)
	s.Arguments = ctx.StringSlice(consts.Arguments)
	s.EtcDirs = ctx.StringSlice(consts.Etc)
	return nil
}
