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
)

type ProtoFile struct{}

func (*ProtoFile) GetDependentFilePaths(mainIdlPath string) ([]string, error) {
	// create maps to keep track of processed and related paths
	processedPaths := make(map[string]bool)
	relatedPaths := make(map[string]bool)
	var resultPaths []string

	// get the base directory of the main IDL file
	baseDir := filepath.Dir(mainIdlPath)

	// define a function to process each file recursively
	var processFile func(filePath string) error
	processFile = func(filePath string) error {
		// if the file has already been processed, skip it
		if processedPaths[filePath] {
			return nil
		}

		// read the content of the Thrift file
		thriftContent, err := ioutil.ReadFile(filePath)
		if err != nil {
			return err
		}
		regex := regexp.MustCompile(importPattern)

		// find all import statements in the Thrift file
		matches := regex.FindAllStringSubmatch(string(thriftContent), -1)
		var includePaths []string

		// extract the paths from the import statements
		for _, match := range matches {
			if len(match) >= 2 {
				includePath := match[1]
				// obtain the fields in the import here and process them
				absolutePath := filepath.Clean(filepath.Join(baseDir, includePath))
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
	err := processFile(mainIdlPath)
	if err != nil {
		return nil, err
	}

	// calculate the relative paths to the main IDL file
	mainIdlDir := filepath.Dir(mainIdlPath)
	relativePaths := make([]string, len(resultPaths))
	for i, path := range resultPaths {
		relativePaths[i], _ = filepath.Rel(mainIdlDir, path)
	}

	return relativePaths, nil
}
