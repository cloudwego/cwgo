package parser

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

type ProtoFile struct{}

func (_ *ProtoFile) GetDependentFilePaths(mainIdlPath string) ([]string, error) {
	processedPaths := make(map[string]bool)
	relatedPaths := make(map[string]bool)
	var resultPaths []string

	baseDir := filepath.Dir(mainIdlPath)

	var processFile func(filePath string) error
	processFile = func(filePath string) error {
		if processedPaths[filePath] {
			return nil
		}

		thriftContent, err := ioutil.ReadFile(filePath)
		if err != nil {
			return err
		}
		regex := regexp.MustCompile(importPattern)

		matches := regex.FindAllStringSubmatch(string(thriftContent), -1)
		var includePaths []string

		for _, match := range matches {
			if len(match) >= 2 {
				includePath := match[1]
				// Obtain the fields in the import here and process them
				absolutePath := filepath.Clean(filepath.Join(baseDir, includePath))
				_, err := os.Stat(absolutePath)
				if err != nil {
					continue
				}
				includePaths = append(includePaths, absolutePath)
			}
		}

		processedPaths[filePath] = true

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

	err := processFile(mainIdlPath)
	if err != nil {
		return nil, err
	}

	mainIdlDir := filepath.Dir(mainIdlPath)
	relativePaths := make([]string, len(resultPaths))
	for i, path := range resultPaths {
		relativePaths[i], _ = filepath.Rel(mainIdlDir, path)
	}

	return relativePaths, nil
}
