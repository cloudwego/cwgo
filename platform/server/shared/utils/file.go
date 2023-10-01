/*
 * Copyright 2022 CloudWeGo Authors
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

package utils

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const GitlabURLPrefix = "https://gitlab.com/"

func ProcessFolders(fileContentMap map[string][]byte, tempDir string, folders ...string) error {
	for _, folder := range folders {
		err := filepath.Walk(tempDir+"/"+folder, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() {
				relPath, err := filepath.Rel(tempDir, path)
				if err != nil {
					return err
				}

				content, err := ioutil.ReadFile(path)
				if err != nil {
					return err
				}

				fileContentMap[relPath] = content
			}

			return nil
		})

		if err != nil {
			fmt.Printf("Error walking path %s: %v\n", folder, err)
			return err
		}
	}

	return nil
}

func ParseIdlURL(url string) (idlPid, owner, repoName string, err error) {
	var tempPath string
	if strings.HasPrefix(url, GitlabURLPrefix) {
		tempPath = url[len(GitlabURLPrefix):]
		lastQuestionMarkIndex := strings.LastIndex(tempPath, "?")
		if lastQuestionMarkIndex != -1 {
			tempPath = tempPath[:lastQuestionMarkIndex]
		}
	} else {
		return "", "", "", errors.New("idlPath format wrong,do not have prefix:\"https://github.com/\"")
	}
	urlParts := strings.Split(tempPath, "/")
	if len(urlParts) < 5 {
		return "", "", "", errors.New("idlPath format wrong")
	}
	owner = urlParts[0]
	repoName = urlParts[1]
	for i := 5; i < len(urlParts); i++ {
		idlPid = idlPid + "/" + urlParts[i]
	}
	idlPid = idlPid[1:]
	return idlPid, owner, repoName, nil
}

func CheckNotExist(src string) bool {
	_, err := os.Stat(src)

	return os.IsNotExist(err)
}

func MkDir(src string) error {
	err := os.MkdirAll(src, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func IsNotExistMkDir(src string) error {
	if notExist := CheckNotExist(src); notExist == true {
		if err := MkDir(src); err != nil {
			return err
		}
	}

	return nil
}
