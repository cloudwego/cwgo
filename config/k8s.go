package config

import (
	"github.com/cloudwego/cwgo/pkg/consts"
	"github.com/urfave/cli/v2"
)

type KubeArgument struct {
	Template        string
	Name            string
	Namespace       string
	Port            int
	Image           string
	ImagePullPolicy string
	LimitCpu        int
	LimitMem        int
	MaxReplicas     int
	MinReplicas     int
	Output          string
	Replicas        int
	Branch          string
	RequestCpu      int
	RequestMem      int
	Revisions       int
	Secret          string
	ServiceAccount  string
	NodePort        int
	TargetPort      int
}

func NewKubeArgument() *KubeArgument {
	return &KubeArgument{}
}

func (s *KubeArgument) ParseCli(ctx *cli.Context) error {
	// Parse all the arguments from the CLI context
	s.Name = ctx.String(consts.Name)
	s.Namespace = ctx.String(consts.Namespace)
	s.Port = ctx.Int(consts.Port)
	s.Image = ctx.String(consts.Image)
	s.ImagePullPolicy = ctx.String(consts.ImagePullPolicy)
	s.LimitCpu = ctx.Int(consts.LimitCpu)
	s.LimitMem = ctx.Int(consts.LimitMem)
	s.MaxReplicas = ctx.Int(consts.MaxReplicas)
	s.MinReplicas = ctx.Int(consts.MinReplicas)
	s.Output = ctx.String(consts.Output)
	s.Replicas = ctx.Int(consts.Replicas)
	s.Branch = ctx.String(consts.Branch)
	s.Template = ctx.String(consts.Template)
	s.RequestCpu = ctx.Int(consts.RequestCpu)
	s.RequestMem = ctx.Int(consts.RequestMem)
	s.Revisions = ctx.Int(consts.Revisions)
	s.Secret = ctx.String(consts.Secret)
	s.ServiceAccount = ctx.String(consts.ServiceAccount)
	s.NodePort = ctx.Int(consts.NodePort)
	s.TargetPort = ctx.Int(consts.TargetPort)

	return nil
}
