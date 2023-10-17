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
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

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

func UnTar(archiveData []byte, tempDir string, IsTarball bool) (string, error) {
	var tarReader *tar.Reader
	if IsTarball {
		tarballBuffer := bytes.NewReader(archiveData)

		gzipReader, err := gzip.NewReader(tarballBuffer)
		if err != nil {
			return "", err
		}
		defer gzipReader.Close()

		tarReader = tar.NewReader(gzipReader)
	} else {
		tarReader = tar.NewReader(bytes.NewReader(archiveData))
	}
	var rootDirName string

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}

		if header.Typeflag != tar.TypeDir && header.Typeflag != tar.TypeReg {
			continue
		}

		if rootDirName == "" {
			rootDirName = header.Name
		}

		targetPath := filepath.Join(tempDir, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
				return "", err
			}

		case tar.TypeReg:
			file, err := os.Create(targetPath)
			if err != nil {
				return "", err
			}
			defer file.Close()

			if _, err := io.Copy(file, tarReader); err != nil {
				return "", err
			}
		}
	}
	return rootDirName, nil
}
