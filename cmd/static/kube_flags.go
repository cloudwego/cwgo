package static

import (
	"github.com/cloudwego/cwgo/config"
	"github.com/cloudwego/cwgo/pkg/consts"
	"github.com/urfave/cli/v2"
)

func kubeFlags() []cli.Flag {
	globalArgs := config.GetGlobalArgs()
	return []cli.Flag{
		// Required Flags
		&cli.StringFlag{
			Name:        consts.Name, // Use consts.Name constant
			Required:    true,
			Usage:       "The name of the deployment",
			Destination: &globalArgs.KubeArgument.Name,
		},
		&cli.StringFlag{
			Name:        consts.Namespace, // Use consts.Namespace constant
			Required:    true,
			Usage:       "The namespace of the deployment",
			Destination: &globalArgs.KubeArgument.Namespace,
		},
		&cli.StringFlag{
			Name:        consts.Output, // Use consts.Output constant
			Aliases:     []string{"o"}, // Alias "o" for "output"
			Required:    true,
			Usage:       "The output yaml file",
			Destination: &globalArgs.KubeArgument.Output,
		},
		&cli.IntFlag{
			Name:        consts.Port, // Use consts.Port constant
			Required:    true,
			Usage:       "The port of the deployment to listen on pod",
			Destination: &globalArgs.KubeArgument.Port,
		},

		// Optional Flags with Default Values
		&cli.StringFlag{
			Name:        consts.Branch, // Use consts.Branch constant
			Usage:       "The branch of the remote repo, works with --remote",
			Destination: &globalArgs.KubeArgument.Branch,
		},
		&cli.StringFlag{
			Name:        consts.Template, // Use consts.Template constant
			Value:       consts.Kube,     // Use consts.Kube constant
			Usage:       "The home path of the template",
			Destination: &globalArgs.KubeArgument.Template,
		},
		&cli.StringFlag{
			Name:        consts.Image, // Use consts.Image constant
			Required:    true,
			Usage:       "The docker image for the deployment (required)",
			Destination: &globalArgs.KubeArgument.Image,
		},
		&cli.StringFlag{
			Name:        consts.ImagePullPolicy, // Use consts.ImagePullPolicy constant
			Usage:       "Image pull policy. One of Always, Never, IfNotPresent",
			Destination: &globalArgs.KubeArgument.ImagePullPolicy,
		},
		&cli.IntFlag{
			Name:        consts.LimitCpu, // Use consts.LimitCpu constant
			Usage:       "The CPU limit for deployment (default 1000)",
			Value:       consts.DefaultLimitCpu, // Use consts.DefaultLimitCpu constant
			Destination: &globalArgs.KubeArgument.LimitCpu,
		},
		&cli.IntFlag{
			Name:        consts.LimitMem, // Use consts.LimitMem constant
			Usage:       "The memory limit for deployment (default 1024)",
			Value:       consts.DefaultLimitMem, // Use consts.DefaultLimitCpu constant
			Destination: &globalArgs.KubeArgument.LimitMem,
		},
		&cli.IntFlag{
			Name:        consts.MaxReplicas, // Use consts.MaxReplicas constant
			Usage:       "The maximum number of replicas for deployment (default 10)",
			Value:       consts.DefaultMaxReplicas, // Use consts.DefaultLimitCpu constant
			Destination: &globalArgs.KubeArgument.MaxReplicas,
		},
		&cli.IntFlag{
			Name:        consts.MinReplicas, // Use consts.MinReplicas constant
			Usage:       "The minimum number of replicas for deployment (default 3)",
			Value:       consts.DefaultMinReplicas, // Use consts.DefaultLimitCpu constant
			Destination: &globalArgs.KubeArgument.MinReplicas,
		},

		&cli.IntFlag{
			Name:        consts.Replicas, // Use consts.Replicas constant
			Usage:       "The number of replicas for deployment (default 3)",
			Value:       consts.DefaultReplicas, // Use consts.DefaultLimitCpu constant
			Destination: &globalArgs.KubeArgument.Replicas,
		},
		&cli.IntFlag{
			Name:        consts.RequestCpu, // Use consts.RequestCpu constant
			Usage:       "The requested CPU for deployment (default 500)",
			Value:       consts.DefaultRequestCpu, // Use consts.DefaultLimitCpu constant
			Destination: &globalArgs.KubeArgument.RequestCpu,
		},
		&cli.IntFlag{
			Name:        consts.RequestMem, // Use consts.RequestMem constant
			Usage:       "The requested memory for deployment (default 512)",
			Value:       consts.DefaultRequestMem, // Use consts.DefaultLimitCpu constant
			Destination: &globalArgs.KubeArgument.RequestMem,
		},
		&cli.IntFlag{
			Name:        consts.Revisions, // Use consts.Revisions constant
			Usage:       "The number of revision histories to limit (default 5)",
			Value:       consts.DefaultRevisions, // Use consts.DefaultLimitCpu constant
			Destination: &globalArgs.KubeArgument.Revisions,
		},
		&cli.StringFlag{
			Name:        consts.Secret, // Use consts.Secret constant
			Usage:       "The secret to pull the image from the registry",
			Destination: &globalArgs.KubeArgument.Secret,
		},
		&cli.StringFlag{
			Name:        consts.ServiceAccount, // Use consts.ServiceAccount constant
			Usage:       "The ServiceAccount for the deployment",
			Destination: &globalArgs.KubeArgument.ServiceAccount,
		},
		&cli.IntFlag{
			Name:        consts.NodePort, // Use consts.NodePort constant
			Usage:       "The nodePort for the deployment to expose",
			Destination: &globalArgs.KubeArgument.NodePort,
		},
		&cli.IntFlag{
			Name:        consts.TargetPort, // Use consts.TargetPort constant
			Usage:       "The targetPort for the deployment, default to port",
			Destination: &globalArgs.KubeArgument.TargetPort,
		},
	}
}
