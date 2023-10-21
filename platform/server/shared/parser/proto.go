package parser

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

type ProtoFile struct{}

func (_ *ProtoFile) GetDependentFilePaths(mainIdlPath string) ([]string, error) {
	// Create maps to keep track of processed and related paths
	processedPaths := make(map[string]bool)
	relatedPaths := make(map[string]bool)
	var resultPaths []string

	// Get the base directory of the main IDL file
	baseDir := filepath.Dir(mainIdlPath)

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
		regex := regexp.MustCompile(importPattern)

		// Find all import statements in the Thrift file
		matches := regex.FindAllStringSubmatch(string(thriftContent), -1)
		var includePaths []string

		// Extract the paths from the import statements
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
