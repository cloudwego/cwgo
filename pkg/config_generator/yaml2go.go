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

import (
	"fmt"
	"go/format"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
)

// New creates Yaml2Go object
func New(key, desc, kind string, configValueType ConfigValueType) Yaml2Go {
	return Yaml2Go{
		StructsMeta: &ConfigGenerateMeta{
			Key:             key,
			Kind:            kind,
			Desc:            desc,
			ConfigValueType: configValueType,
			structTree:      make(map[string]Struct),
		},
	}
}

type line struct {
	structName string
	line       string
}

// Yaml2Go to store converted result
type Yaml2Go struct {
	visited     map[line]bool
	structMap   map[string]string
	StructsMeta *ConfigGenerateMeta
}

// NewStruct creates new entry in structMap result
func (yg *Yaml2Go) NewStruct(structName, parent string) string {
	// If struct already present with the same name
	// rename struct to Parent Struct name
	if _, ok := yg.structMap[structName]; ok {
		structName = goKeyFormat(parent) + structName
	}
	yg.StructsMeta.structTree[structName] = Struct{
		StructName: structName,
	}
	yg.AppendResult(structName, fmt.Sprintf("// %s\n", structName))
	l := fmt.Sprintf("type %s struct {\n", structName)
	yg.visited[line{structName: structName, line: l}] = true
	yg.structMap[structName] += l
	return structName
}

// AppendResult add lines to the result
func (yg *Yaml2Go) AppendResult(structName, l string, s ...Struct) {
	if _, ok := yg.visited[line{structName, l}]; !ok {
		yg.structMap[structName] += l
	}
	yg.visited[line{structName, l}] = true
	if len(s) > 0 {
		v, ok := yg.StructsMeta.structTree[structName]
		if ok {
			v.Fields = append(v.Fields, s[0].Fields...)
			yg.StructsMeta.structTree[structName] = v
		} else {
			yg.StructsMeta.structTree[structName] = Struct{
				Fields: s[0].Fields,
			}
		}
	}
}

// removeUnderscores and camelize string
func goKeyFormat(key string) string {
	var st string
	strList := strings.FieldsFunc(key, func(r rune) bool {
		return r == ' ' || r == '_' || r == '-'
	})
	for _, str := range strList {
		st += strings.Title(str)
	}
	if len(st) == 0 {
		st = key
	}
	return st
}

// Convert transforms map[string]interface{} to go struct
func (yg *Yaml2Go) Convert(structName string, data []byte) (string, error) {
	structName = convertToGoStructName(structName)
	yg.visited = make(map[line]bool)
	yg.structMap = make(map[string]string)

	// Unmarshal to map[string]interface{}
	var obj map[string]interface{}
	err := yaml.Unmarshal(data, &obj)
	if err != nil {
		return "", err
	}

	yg.NewStruct(structName, "")
	for k, v := range obj {
		yg.Structify(structName, k, v, false)
	}
	yg.AppendResult(structName, "}\n")
	var result string
	for _, value := range yg.structMap {
		result += fmt.Sprintf("%s\n", value)
	}

	// Convert result into go format
	goFormat, err := format.Source([]byte(result))
	if err != nil {
		return "", err
	}
	return string(goFormat), nil
}

// Structify transforms map key values to struct fields
// structName : parent struct name
// k, v       : fields in the struct
func (yg *Yaml2Go) Structify(structName, k string, v interface{}, arrayElem bool) {
	if reflect.TypeOf(v) == nil || len(k) == 0 {
		field := Struct{
			Fields: []Field{
				{
					FieldName: goKeyFormat(k),
					FieldType: "interface{}",
					Tags: []FieldTag{
						{
							TagKey:   "yaml",
							TagValue: k,
						},
						{
							TagKey:   "json",
							TagValue: k,
						},
					},
				},
			},
		}
		yg.AppendResult(structName, fmt.Sprintf("%s interface{} `yaml:\"%s\"`\n", goKeyFormat(k), k), field)
		return
	}

	switch reflect.TypeOf(v).Kind() {
	// If yaml object
	case reflect.Map:
		switch val := v.(type) {
		case map[string]interface{}:
			key := goKeyFormat(k)
			newKey := key
			if !arrayElem {
				// Create new structure
				newKey = yg.NewStruct(key, structName)
				field := Struct{
					StructName: newKey,
					Fields: []Field{
						{
							FieldName: key,
							FieldType: newKey,
							IsStruct:  true,
							Tags: []FieldTag{
								{
									TagKey:   "yaml",
									TagValue: k,
								},
								{
									TagKey:   "json",
									TagValue: k,
								},
							},
						},
					},
				}
				yg.AppendResult(structName, fmt.Sprintf("%s %s `yaml:\"%s\"`\n", key, newKey, k), field)
			}
			// If array of yaml objects
			for k1, v1 := range val {
				yg.Structify(newKey, k1, v1, false)
			}
			if !arrayElem {
				yg.AppendResult(newKey, "}\n")
			}
		}

	// If array
	case reflect.Slice:
		val := v.([]interface{})
		if len(val) == 0 {
			return
		}
		keyFormat := goKeyFormat(k)
		switch val[0].(type) {
		case string, int, bool, float64:
			structOr := Struct{
				StructName: keyFormat,
				Fields: []Field{
					{
						FieldName: goKeyFormat(k),
						FieldType: fmt.Sprintf("[]%s", reflect.TypeOf(val[0])),
						IsStruct:  false,
						IsSlice:   true,
						Tags: []FieldTag{
							{
								TagKey:   "yaml",
								TagValue: k,
							},
							{
								TagKey:   "json",
								TagValue: k,
							},
						},
					},
				},
			}
			for _, v := range val {
				var v1 string
				switch v.(type) {
				case string:
					v1 = fmt.Sprintf("\"%s\"", v)
				case bool:
					v1 = fmt.Sprintf("%t", v)
				default:
					v1 = fmt.Sprintf("%v", v)
				}
				structOr.Fields[0].Children = append(structOr.Fields[0].Children, Field{
					Value:       v1,
					FieldType:   reflect.TypeOf(v).String(),
					IsBasicType: true,
				})
			}
			yg.AppendResult(structName, fmt.Sprintf("%s []%s `yaml:\"%s\"`\n", goKeyFormat(k), reflect.TypeOf(val[0]), k), structOr)

		// if nested object
		case map[string]interface{}:
			key := goKeyFormat(k)
			// Create new structure
			newKey := yg.NewStruct(key, structName)
			field := Struct{
				StructName: newKey,
				Fields: []Field{
					{
						FieldName: key,
						FieldType: fmt.Sprintf("[]%s", newKey),
						IsSlice:   true,
						Tags: []FieldTag{
							{
								TagKey:   "yaml",
								TagValue: k,
							},
							{
								TagKey:   "json",
								TagValue: k,
							},
						},
					},
				},
			}
			yg.AppendResult(structName, fmt.Sprintf("%s []%s `yaml:\"%s\"`\n", key, newKey, k), field)
			for _, v1 := range val {
				yg.Structify(newKey, key, v1, true)
			}
			yg.AppendResult(newKey, "}\n")
		}

	default:
		field := Struct{
			Fields: []Field{
				{
					FieldName:   goKeyFormat(k),
					FieldType:   reflect.TypeOf(v).String(),
					IsBasicType: isBasicType(reflect.TypeOf(v).String()),
					Tags: []FieldTag{
						{
							TagKey:   "yaml",
							TagValue: k,
						},
						{
							TagKey:   "json",
							TagValue: k,
						},
					},
				},
			},
		}
		switch v.(type) {
		case string:
			field.Fields[0].Value = fmt.Sprintf("\"%s\"", v)
		case bool:
			field.Fields[0].Value = fmt.Sprintf("%t", v)
		default:
			field.Fields[0].Value = fmt.Sprintf("%v", v)
		}
		yg.AppendResult(structName, fmt.Sprintf("%s %s `yaml:\"%s\"`\n", goKeyFormat(k), reflect.TypeOf(v).String(), k), field)
	}
}

// organizeStructs organizes structs into a tree structure and removes redundant nodes
func organizeStructs(fileMeta *ConfigGenerateMeta) {
	// Keep track of which structs have been included as children
	included := make(map[string]bool)

	// Iterate over each struct in the map
	for structName, structDef := range fileMeta.structTree {
		// If a struct is already included as a child, skip organizing it as a root
		if included[structName] {
			continue
		}

		// Create a map to track visited structs to avoid circular references
		visited := make(map[string]bool)

		// Recursively organize structs
		fileMeta.structTree[structName] = organizeStructHelper(structDef, fileMeta.structTree, visited, included)
	}

	// Remove redundant nodes from the map
	for structName := range included {
		delete(fileMeta.structTree, structName)
	}
	for _, basic := range fileMeta.structTree {
		fileMeta.ConfigStruct = basic
	}
}

// organizeStructHelper is a helper function to recursively organize structs
func organizeStructHelper(currentStruct Struct, structs map[string]Struct, visited, included map[string]bool) Struct {
	// Iterate over each field in the struct
	for i, field := range currentStruct.Fields {
		// Handle array types by stripping brackets and checking the type
		fieldType := field.FieldType
		isSlice := false
		if strings.HasPrefix(fieldType, "[]") {
			fieldType = strings.TrimPrefix(fieldType, "[]")
			isSlice = true
		}

		// Determine if the field type is a basic type or a struct
		currentStruct.Fields[i].IsBasicType = isBasicType(fieldType)

		if !currentStruct.Fields[i].IsBasicType {
			// Check if the field type matches any struct name
			if childStruct, exists := structs[fieldType]; exists && !visited[fieldType] {
				// Mark the field as a struct
				currentStruct.Fields[i].IsStruct = true

				// Mark the struct as visited
				visited[fieldType] = true

				// Recursively organize the child struct
				childStruct = organizeStructHelper(childStruct, structs, visited, included)

				// Add the fields of the child struct directly as children
				currentStruct.Fields[i].Children = childStruct.Fields

				// Mark the struct as included in a parent
				included[fieldType] = true
			}
		}

		// Mark the field as a slice if applicable
		currentStruct.Fields[i].IsSlice = isSlice
	}
	return currentStruct
}
