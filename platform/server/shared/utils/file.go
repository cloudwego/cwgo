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
	"archive/tar"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"github.com/cloudwego/cwgo/platform/server/shared/consts"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

func ProcessFolders(fileContentMap map[string][]byte, tempDir string, folders ...string) error {
	// Iterate over the specified folders
	for _, folder := range folders {
		// Recursively walk through the folder and its subdirectories
		err := filepath.Walk(tempDir+"/"+folder, func(path string, info os.FileInfo, err error) error {
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

		// Check if an error occurred while walking the folder
		if err != nil {
			fmt.Printf("Error walking path %s: %v\n", folder, err)
			return err
		}
	}

	// Return nil if the processing of folders is successful
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
	if notExist := CheckNotExist(src); notExist == true {
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

func DetermineIdlType(idlPid string) (int32, error) {
	// Define regular expressions to match file extensions for Thrift and Proto files
	thriftRegex := `(?i)\.thrift$` // Case-insensitive match for ".thrift" extension
	protoRegex := `(?i)\.proto$`   // Case-insensitive match for ".proto" extension

	// Check if the input string matches the Thrift or Proto file extensions
	if matched, _ := regexp.MatchString(thriftRegex, idlPid); matched {
		// Matched ".thrift" extension, indicating Thrift file
		return consts.IdlTypeNumThrift, nil // You need to define ThriftIdlType as a constant or return a specific IDL type for Thrift
	} else if matched, _ := regexp.MatchString(protoRegex, idlPid); matched {
		// Matched ".proto" extension, indicating Proto file
		return consts.IdlTypeNumProto, nil // You need to define ProtoIdlType as a constant or return a specific IDL type for Proto
	}

	// Return a default IDL type or error code if no match is found
	return -1, errors.New("incorrect idl type") // You need to define DefaultIdlType as a constant or return an appropriate default value
}
