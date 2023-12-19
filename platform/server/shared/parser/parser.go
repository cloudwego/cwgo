/*
*
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
*
*/

package parser

import "github.com/cloudwego/cwgo/platform/server/shared/consts"

type Parser interface {
	GetDependentFilePaths(baseDirPath, mainIdlPath string) (string, []string, error) // Obtain string slices that rely on main idl
}

func NewParser(idlType int32) Parser {
	switch idlType {
	case consts.IdlTypeNumThrift:
		return NewThriftParser()
	case consts.IdlTypeNumProto:
		return NewProtoParser()
	default:
		return nil
	}
}
