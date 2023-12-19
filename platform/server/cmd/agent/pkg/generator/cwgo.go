/*
 *
 * Copyright 2023 CloudWeGo Authors
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
 *
 */

package generator

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/cloudwego/cwgo/platform/server/shared/logger"
	"go.uber.org/zap"
)

type CwgoGenerator struct{}

func NewCwgoGenerator() *CwgoGenerator {
	return &CwgoGenerator{}
}

func (g *CwgoGenerator) Generate(repoDomain, repoOwner, idlPath, idlSearchPath, serviceName, generatePath string) error {
	// Fixed options were used to generate code, which can be optimized to be optional in the future
	var build strings.Builder
	build.WriteString(
		fmt.Sprintf("cwgo client "+
			"--idl %s "+
			"--type %s "+
			"--service %s "+
			"--module %s ",
			idlPath, "rpc", serviceName, fmt.Sprintf("%s/%s/%s", repoDomain, repoOwner, serviceName),
		),
	)
	if idlSearchPath != "" {
		build.WriteString(
			fmt.Sprintf("--proto_search_path %s ",
				idlSearchPath,
			),
		)
	}
	build.WriteString("&& go mod tidy")
	cwgoCmd := exec.Command("sh", "-c", build.String())

	cwgoCmd.Dir = generatePath

	logger.Logger.Debug("exec generate command", zap.String("command", cwgoCmd.String()))

	outBytes, err := cwgoCmd.CombinedOutput()

	logger.Logger.Debug("generate command output", zap.String("output", string(outBytes)))

	return err
}
