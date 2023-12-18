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

package parser

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"github.com/cloudwego/cwgo/platform/server/shared/utils"
)

type ProtoParser struct{}

func NewProtoParser() *ProtoParser {
	return &ProtoParser{}
}

func (*ProtoParser) GetDependentFilePaths(baseDirPath, mainIdlPath string) (string, []string, error) {
	// create maps to keep track of processed and related paths
	processedPaths := make(map[string]bool)
	relatedPaths := make(map[string]bool)
	var resultPaths []string

	var importBaseDirPath string

	// define a function to process each file recursively
	var processFile func(filePath string) error
	processFile = func(filePath string) error {
		// if the file has already been processed, skip it
		if processedPaths[filePath] {
			return nil
		}

		// read the content of the Thrift file
		protoContent, err := ioutil.ReadFile(filePath)
		if err != nil {
			return err
		}
		regex := regexp.MustCompile(importPattern)

		// find all import statements in the Thrift file
		matches := regex.FindAllStringSubmatch(string(protoContent), -1)
		var includePaths []string

		// extract the paths from the import statements
		for _, match := range matches {
			if len(match) >= 2 {
				includePath := match[1]
				// obtain the fields in the import here and process them
				if importBaseDirPath == "" {
					importBaseDirPath = utils.FindRootPath(filePath, includePath)
					if importBaseDirPath == "" {
						continue
					}
				}

				absolutePath := filepath.Clean(filepath.Join(importBaseDirPath, includePath))
				_, err := os.Stat(absolutePath)
				if err != nil {
					continue
				}
				includePaths = append(includePaths, absolutePath)
			}
		}

		// mark the current file as processed
		processedPaths[filePath] = true

		// recursively process the included files
		for _, includePath := range includePaths {
			if !relatedPaths[includePath] {
				relatedPaths[includePath] = true
				resultPaths = append(resultPaths, includePath)
				err := processFile(includePath)
				if err != nil {
					return err
				}
			}
		}

		return nil
	}

	// start the recursive processing with the main IDL file
	mainAbsPath := baseDirPath + mainIdlPath
	err := processFile(mainAbsPath)
	if err != nil {
		return "", nil, err
	}

	// calculate the relative paths to the base dir path
	relativePaths := make([]string, len(resultPaths))
	for i, path := range resultPaths {
		relativePath, _ := filepath.Rel(importBaseDirPath, path)
		relativePaths[i] = filepath.ToSlash(relativePath)
	}

	rel, err := filepath.Rel(baseDirPath, importBaseDirPath)
	if err != nil {
		return "", nil, err
	}

	return rel, relativePaths, nil
}
