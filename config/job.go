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

package config

import (
	"github.com/cloudwego/cwgo/pkg/consts"
	"github.com/urfave/cli/v2"
)

type JobArgument struct {
	GoMod         string
	PackagePrefix string
	JobName       []string
	OutDir        string
}

func NewJobArgument() *JobArgument {
	return &JobArgument{}
}

func (j *JobArgument) ParseCli(ctx *cli.Context) error {
	j.JobName = ctx.StringSlice(consts.JobName)
	j.GoMod = ctx.String(consts.Module)
	j.OutDir = ctx.String(consts.OutDir)
	return nil
}
