/*
 * Copyright 2024 CloudWeGo Authors
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

package model

import (
	"reflect"

	"github.com/cloudwego/cwgo/pkg/doc/mongo/code"
)

type IdlExtractStruct struct {
	Name          string
	StructFields  []*StructField
	InterfaceInfo *InterfaceInfo
	UpdateInfo
}

type InterfaceInfo struct {
	Name    string
	Methods []*InterfaceMethod
}

type InterfaceMethod struct {
	Name             string
	ParsedTokens     string
	Params           code.Params
	Returns          code.Returns
	BelongedToStruct *IdlExtractStruct
}

type StructField struct {
	Name               string
	Type               code.Type
	Tag                reflect.StructTag
	IsBelongedToStruct bool
	BelongedToStruct   *IdlExtractStruct
}

type UpdateInfo struct {
	Update                 bool
	UpdateMongoFileContent []byte
	UpdateIfFileContent    []byte
	PreMethodNamesMap      map[string]struct{}
	PreIfMethods           []*InterfaceMethod
}
