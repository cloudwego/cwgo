// Copyright 2024 CloudWeGo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
