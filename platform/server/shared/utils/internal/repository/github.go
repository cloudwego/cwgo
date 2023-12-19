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

package repository

import (
	"errors"
	"fmt"
	"regexp"
)

func GetRepoFullUrlGitHub(repoUrl, ref, filePid string) string {
	return fmt.Sprintf("%s/blob/%s/%s", repoUrl, ref, filePid)
}

func ParseRepoFileUrlGitHub(url string) (filePid, owner, repoName string, err error) {
	// define a regular expression to parse the GitHub URL.
	regex := regexp.MustCompile(regGitHubURL)

	// use the regular expression to extract relevant components from the URL.
	matches := regex.FindStringSubmatch(url)
	if len(matches) != 5 {
		return "", "", "", errors.New("IDL path format is incorrect; unable to parse the GitHub URL")
	}

	// assign values to the returned variables.
	owner = matches[1]
	repoName = matches[2]
	filePid = matches[4]

	return filePid, owner, repoName, nil
}
