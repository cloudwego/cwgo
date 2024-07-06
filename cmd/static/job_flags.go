package static

import (
	"github.com/cloudwego/cwgo/pkg/consts"
	"github.com/urfave/cli/v2"
)

func jobFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringSliceFlag{Name: consts.JobName, Usage: "Specify the job name."},
		&cli.StringFlag{Name: consts.Module, Aliases: []string{"mod"}, Usage: "Specify the Go module name to generate go.mod."},
		&cli.StringFlag{Name: consts.OutDir, Usage: "Specify output directory, default is current dir."},
	}
}
