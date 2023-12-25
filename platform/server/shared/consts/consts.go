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

package consts

import (
	"os"
	"path/filepath"
)

const (
	ProjectName = "cwgo"

	ServiceNameApi   = ProjectName + "-" + ServerTypeApi
	ServiceNameAgent = ProjectName + "-" + ServerTypeAgent

	AgentMetadataFile = ServiceNameAgent + ".yaml"

	ServerRunModeStandalone = "standalone"
	ServerRunModeCluster    = "cluster"
)

const (
	DefaultLimit = 20
)

const (
	RdbKeyRegistryService = ProjectName + ":" + "registry:service:%s"
)

const (
	TempDirRepo          = "repo"           // repo pull file temp dir
	TempDirTemplate      = "template"       // code template temp dir
	TempDirGeneratedCode = "generated_code" // generated code temp dir
)

var TempDir string

func init() {
	TempDir = filepath.Join(os.TempDir(), ProjectName)
}
