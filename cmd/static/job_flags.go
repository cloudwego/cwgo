/*
 * Copyright 2024 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
