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

package repository

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

func GetRepoFullUrlGitLab(repoUrl, ref, filePid string) string {
	return fmt.Sprintf("%s/-/blob/%s/%s?ref_type=heads", repoUrl, ref, filePid)
}

func ParseRepoFileUrlGitLab(url string) (filePid, owner, repoName string, err error) {
	// using regular expressions to match fields
	regex := regexp.MustCompile(regGitLabURL)
	matches := regex.FindStringSubmatch(url)
	if len(matches) != 5 {
		return "", "", "", errors.New("idlPath format wrong, cannot parse gitlab URL")
	}

	owner = matches[1]
	repoName = matches[2]
	filePid = strings.Split(matches[4], "?")[0]

	return filePid, owner, repoName, nil
}
