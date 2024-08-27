package static

import (
	"github.com/cloudwego/cwgo/pkg/consts"
	"github.com/urfave/cli/v2"
)

func kubeFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:  consts.Name,
			Usage: "Specify the name of the service.",
		},
		&cli.StringFlag{
			Name:  consts.Namespace,
			Usage: "Specify the namespace.",
		},
		&cli.StringFlag{
			Name:  consts.Image,
			Usage: "Specify the image.",
		},
		&cli.StringFlag{
			Name:  consts.Secret,
			Usage: "Specify the secret.",
		},
		&cli.IntFlag{
			Name:        consts.RequestCpu,
			Usage:       "Specify the request cpu.",
			DefaultText: "500",
		},
		&cli.StringFlag{
			Name:        consts.RequestMem,
			Usage:       "Specify the request memory.",
			DefaultText: "512",
		},
		&cli.StringFlag{
			Name:        consts.LimitCpu,
			Usage:       "Specify the limit cpu.",
			DefaultText: "1000",
		},
		&cli.StringFlag{
			Name:        consts.LimitMem,
			Usage:       "Specify the limit memory.",
			DefaultText: "1024",
		},
		&cli.StringFlag{
			Name:  consts.Port,
			Usage: "Specify the port.",
		},
		&cli.StringFlag{
			Name:        consts.MinReplicas,
			Usage:       "Specify the min replicas.",
			DefaultText: "3",
		},
		&cli.StringFlag{
			Name:        consts.MaxReplicas,
			Usage:       "Specify the max replicas.",
			DefaultText: "10",
		},
		&cli.StringFlag{
			Name:  consts.ImagePullPolicy,
			Usage: "Specify the image pull policy. (Always or IfNotPresent or Never)",
		},
		&cli.StringFlag{
			Name:  consts.ServiceAccount,
			Usage: "Specify the service account.",
		},
	}
}
