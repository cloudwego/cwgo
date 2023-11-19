/*
 *
 *  * Copyright 2022 CloudWeGo Authors
 *  *
 *  * Licensed under the Apache License, Version 2.0 (the "License");
 *  * you may not use this file except in compliance with the License.
 *  * You may obtain a copy of the License at
 *  *
 *  *     http://www.apache.org/licenses/LICENSE-2.0
 *  *
 *  * Unless required by applicable law or agreed to in writing, software
 *  * distributed under the License is distributed on an "AS IS" BASIS,
 *  * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  * See the License for the specific language governing permissions and
 *  * limitations under the License.
 *
 */

package utils

import (
	"fmt"
	"github.com/cloudwego/cwgo/platform/server/shared/consts"
)

func GetRepoFullUrl(repoType int32, repoUrl, ref, filePid string) string {
	switch repoType {
	case consts.RepositoryTypeNumGitLab:
		return GetRepoFullUrlGitLab(repoUrl, ref, filePid)
	case consts.RepositoryTypeNumGithub:
		return GetRepoFullUrlGitHub(repoUrl, ref, filePid)
	default:
		return ""
	}
}

func GetRepoFullUrlGitLab(repoUrl, ref, filePid string) string {
	return fmt.Sprintf("%s/-/blob/%s/%s?ref_type=heads", repoUrl, ref, filePid)
}

func GetRepoFullUrlGitHub(repoUrl, ref, filePid string) string {
	return fmt.Sprintf("%s/blob/%s/%s", repoUrl, ref, filePid)
}
