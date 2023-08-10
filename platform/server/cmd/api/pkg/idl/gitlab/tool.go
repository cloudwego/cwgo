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

package gitlab

import "io/ioutil"

func processFolders(filesM map[string][]byte, sourceDir string, folders ...string) error {
	for _, folder := range folders {
		// 递归处理文件夹
		if err := processFolder(filesM, sourceDir+"/"+folder, folder); err != nil {
			return err
		}
	}
	return nil
}

func processFolder(filesM map[string][]byte, folderPath, targetPath string) error {
	files, err := ioutil.ReadDir(folderPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			// process the subfolders
			subFolder := targetPath + "/" + file.Name()
			if err := processFolder(filesM, folderPath+"/"+file.Name(), subFolder); err != nil {
				return err
			}
		} else {
			// process the file
			filePath := folderPath + "/" + file.Name()
			content, err := ioutil.ReadFile(filePath)
			filesM[targetPath+"/"+file.Name()] = content
			if err != nil {
				return err
			}
		}
	}

	return nil
}
