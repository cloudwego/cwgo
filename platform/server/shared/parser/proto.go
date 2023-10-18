package parser

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
)

type ProtoFile struct{}

func (_ *ProtoFile) GetDependentFilePaths(mainProtoPath string) ([]string, error) {
	visited := make(map[string]struct{})

	importPaths := make([]string, 0)

	dir := filepath.Dir(mainProtoPath)

	var findImports func(string) error
	findImports = func(protoPath string) error {
		if _, ok := visited[protoPath]; ok {
			return nil
		}

		visited[protoPath] = struct{}{}

		file, err := os.Open(protoPath)
		if err != nil {
			return err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			importStatement := regexp.MustCompile(importPattern).FindStringSubmatch(line)
			if len(importStatement) > 1 {
				importPath := importStatement[1]
				fullImportPath := filepath.Join(dir, importPath)
				importPaths = append(importPaths, fullImportPath)
				if err := findImports(fullImportPath); err != nil {
					return err
				}
			}
		}

		return nil
	}

	if err := findImports(mainProtoPath); err != nil {
		return nil, err
	}

	return importPaths, nil
}
