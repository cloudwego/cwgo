package static

import (
	"github.com/cloudwego/cwgo/pkg/consts"
	"github.com/urfave/cli/v2"
)

func apiFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:  consts.ProjectPath,
			Usage: "Specify the project path.",
		},
		&cli.StringFlag{
			Name:        consts.HertzRepoUrl,
			Aliases:     []string{"r"},
			DefaultText: consts.HertzRepoDefaultUrl,
			Usage:       "Specify the url of the hertz repository you want",
		},
	}
}
