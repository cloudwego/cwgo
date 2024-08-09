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

package config_generator

const (
	Yaml = "yaml"
	Json = "json"
	Text = "text"
)

type ConfigGenerateMeta struct {
	Desc            string          `json:"desc,omitempty"`
	Kind            string          `json:"kind,omitempty"`
	ConfigStruct    Struct          `json:"config_struct,omitempty"`
	ConfigValueType ConfigValueType `json:"config_value_type,omitempty"`
	Key             string          `json:"key,omitempty"`

	structTree map[string]Struct `json:"-"`
}

type FieldTag struct {
	TagKey   string `json:"tag_key"`   // eg: json
	TagValue string `json:"tag_value"` // eg: "name,omitempty"
}

// Struct to store struct fields
type Struct struct {
	StructName string  `json:"struct_name,omitempty"`
	Fields     []Field `json:"fields,omitempty"`
}

type Field struct {
	FieldName   string     `json:"field_name,omitempty"` // FieldName of the field
	Value       string     `json:"value,omitempty"`      // Value of the field (used for primitive types)
	FieldType   string     `json:"field_type,omitempty"` // FieldType of the field
	IsStruct    bool       `json:"is_struct"`
	IsSlice     bool       `json:"is_slice"`
	IsBasicType bool       `json:"is_basic"`
	Tags        []FieldTag `json:"tags,omitempty"`
	Children    []Field    `json:"children,omitempty"`
}
