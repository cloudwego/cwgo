package parser

import (
	"os"
	"path/filepath"
	"regexp"
)

type Parser interface {
	GetDependentFilePaths(mainIdlPath string) ([]string, error)
}

func GetDependentFilePaths(mainIdlPath string) ([]string, error) {
	includePattern := `include\s*"([^"]+)"`
	regex := regexp.MustCompile(includePattern)

	content, err := os.ReadFile(mainIdlPath)
	if err != nil {
		return nil, err
	}

	matches := regex.FindAllStringSubmatch(string(content), -1)
	var absolutePaths []string

	for _, match := range matches {
		if len(match) >= 2 {
			relativePath := match[1]
			relPath := filepath.Clean(filepath.Join(filepath.Dir(mainIdlPath), relativePath))
			absolutePaths = append(absolutePaths, relPath)
		}
	}

	return absolutePaths, nil
}
