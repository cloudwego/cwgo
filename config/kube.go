package config

import (
	"github.com/cloudwego/cwgo/pkg/consts"
	"github.com/urfave/cli/v2"
)

type KubeArgument struct {
	Name            string
	Namespace       string
	Image           string
	Secret          string
	RequestCpu      int
	RequestMem      string
	LimitCpu        string
	LimitMem        string
	MinReplicas     int
	MaxReplicas     int
	ImagePullPolicy string
	ServiceAccount  string
}

func NewKubeArgument() *KubeArgument {
	return &KubeArgument{}
}

func (c *KubeArgument) ParseCli(ctx *cli.Context) error {
	c.Name = ctx.String(consts.Name)
	c.Namespace = ctx.String(consts.Namespace)
	c.Image = ctx.String(consts.Image)
	c.Secret = ctx.String(consts.Secret)
	c.RequestCpu = ctx.Int(consts.RequestCpu)
	c.RequestMem = ctx.String(consts.RequestMem)
	c.LimitCpu = ctx.String(consts.LimitCpu)
	c.LimitMem = ctx.String(consts.LimitMem)
	c.MinReplicas = ctx.Int(consts.MinReplicas)
	c.MaxReplicas = ctx.Int(consts.MaxReplicas)
	c.ImagePullPolicy = ctx.String(consts.ImagePullPolicy)
	c.ServiceAccount = ctx.String(consts.ServiceAccount)
	return nil
}
