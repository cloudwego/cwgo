/*
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
 */

package utils

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func ProcessFolders(fileContentMap map[string][]byte, tempDir string, paths ...string) error {
	if len(paths) == 0 {
		// if path is not specified,
		// then walk the entire temp dir
		paths = append(paths, "")
	}
	// Iterate over the specified paths
	for _, path := range paths {
		// Recursively walk through the path and its subdirectories
		err := filepath.Walk(tempDir+"/"+path, func(path string, info os.FileInfo, err error) error {
			// Check if an error occurred during the walk
			if err != nil {
				return err
			}

			// Check if the item is not a directory (i.e., it's a file)
			if !info.IsDir() {
				// Calculate the relative path of the file
				relPath, err := filepath.Rel(tempDir, path)
				if err != nil {
					return err
				}

				// Read the content of the file
				content, err := ioutil.ReadFile(path)
				if err != nil {
					return err
				}

				// Store the file content in the provided map with its relative path as the key
				fileContentMap[relPath] = content
			}

			return nil
		})
		// Check if an error occurred while walking the path
		if err != nil {
			fmt.Printf("Error walking path %s: %v\n", path, err)
			return err
		}
	}

	// Return nil if the processing of paths is successful
	return nil
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
	if notExist := CheckNotExist(src); notExist {
		if err := MkDir(src); err != nil {
			return err
		}
	}

	return nil
}

// UnTar Persist the tar compressed package to the specified disk. If using tarball, gzip decompression is required first
func UnTar(archiveData []byte, tempDir string, IsTarball bool) (string, error) {
	var tarReader *tar.Reader
	if IsTarball {
		// Create a byte reader from the archiveData
		tarballBuffer := bytes.NewReader(archiveData)

		// Create a gzip reader from the tarballBuffer
		gzipReader, err := gzip.NewReader(tarballBuffer)
		if err != nil {
			return "", err
		}
		defer gzipReader.Close()

		// Create a tar reader from the gzip reader
		tarReader = tar.NewReader(gzipReader)
	} else {
		// Create a tar reader directly from the archiveData
		tarReader = tar.NewReader(bytes.NewReader(archiveData))
	}

	var rootDirName string

	for {
		// Read the next header from the tar archive
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}

		// Skip unsupported file types
		if header.Typeflag != tar.TypeDir && header.Typeflag != tar.TypeReg {
			continue
		}

		// Get the root directory name
		if rootDirName == "" {
			rootDirName = header.Name
		}

		// Construct the target path within the temp directory
		targetPath := filepath.Join(tempDir, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			// Create a directory and all parent directories
			if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
				return "", err
			}

		case tar.TypeReg:
			// Create a regular file
			file, err := os.Create(targetPath)
			if err != nil {
				return "", err
			}
			defer file.Close()

			// Copy file content from the tar archive to the created file
			if _, err := io.Copy(file, tarReader); err != nil {
				return "", err
			}
		}
	}
	return rootDirName, nil
}

func FindRootPath(absoluteFilePath, relativeFilePath string) string {
	absRoot := filepath.Dir(absoluteFilePath)
	return findRootPathRecursive(absRoot, relativeFilePath)
}

func findRootPathRecursive(currentDirPath, relativeFilePath string) string {
	filePath := filepath.Join(currentDirPath, relativeFilePath)

	if _, err := os.Stat(filePath); err == nil {
		return currentDirPath
	}

	parentPath := filepath.Dir(currentDirPath)
	if parentPath == currentDirPath {
		return ""
	}

	return findRootPathRecursive(parentPath, relativeFilePath)
}
