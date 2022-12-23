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
	"fmt"
	"os"
	"path/filepath"

	"github.com/cloudwego/hertz/cmd/hz/meta"
)

// Exists is used to judge whether the path exists in file system.
func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

func GetIdlType(path string, pbName ...string) (string, error) {
	ext := filepath.Ext(path)
	if ext == "" || ext[0] != '.' {
		return "", fmt.Errorf("idl path %s is not a valid file", path)
	}
	ext = ext[1:]
	switch ext {
	case meta.IdlThrift:
		return meta.IdlThrift, nil
	case meta.IdlProto:
		if len(pbName) > 0 {
			return pbName[0], nil
		}
		return meta.IdlProto, nil
	default:
		return "", fmt.Errorf("IDL type %s is not supported", ext)
	}
}
