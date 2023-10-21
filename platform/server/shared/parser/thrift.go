package parser

import (
	"io/ioutil"
	"path/filepath"
	"regexp"
)

type ThriftFile struct{}

func (_ *ThriftFile) GetDependentFilePaths(mainIdlPath string) ([]string, error) {
	// Create maps to keep track of processed and related paths
	processedPaths := make(map[string]bool)
	relatedPaths := make(map[string]bool)
	var resultPaths []string

	// Define a function to process each file recursively
	var processFile func(filePath string) error
	processFile = func(filePath string) error {
		// If the file has already been processed, skip it
		if processedPaths[filePath] {
			return nil
		}

		// Read the content of the Thrift file
		thriftContent, err := ioutil.ReadFile(filePath)
		if err != nil {
			return err
		}
		regex := regexp.MustCompile(includePattern)

		// Find all include statements in the Thrift file
		matches := regex.FindAllStringSubmatch(string(thriftContent), -1)
		var includePaths []string

		baseDir := filepath.Dir(filePath)

		// Extract the paths from the include statements
		for _, match := range matches {
			if len(match) >= 2 {
				includePath := match[1]
				absolutePath := filepath.Clean(filepath.Join(baseDir, includePath))
				includePaths = append(includePaths, absolutePath)
			}
		}

		// Mark the current file as processed
		processedPaths[filePath] = true

		// Recursively process the included files
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

	// Start the recursive processing with the main IDL file
	err := processFile(mainIdlPath)
	if err != nil {
		return nil, err
	}

	// Calculate the relative paths to the main IDL file
	mainIdlDir := filepath.Dir(mainIdlPath)
	relativePaths := make([]string, len(resultPaths))
	for i, path := range resultPaths {
		relativePaths[i], _ = filepath.Rel(mainIdlDir, path)
	}

	return relativePaths, nil
}
