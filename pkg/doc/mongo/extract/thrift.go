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

package extract

import (
	"fmt"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/cloudwego/cwgo/config"
	"github.com/cloudwego/hertz/cmd/hz/util"
	"github.com/cloudwego/thriftgo/plugin"

	"github.com/cloudwego/cwgo/pkg/doc/mongo/code"
	"github.com/cloudwego/thriftgo/parser"
)

type ThriftUsedInfo struct {
	Req     *plugin.Request
	DocArgs *config.DocArgument
}

func (info *ThriftUsedInfo) ParseThriftIdl() (rawStructs []*IdlExtractStruct, err error) {
	var getGenGoFilePath func(file *parser.Thrift) error
	getGenGoFilePath = func(file *parser.Thrift) error {
		for _, st := range file.Structs {
			hasInterface := false
			for _, anno := range st.Annotations {
				if strings.Index(anno.Key, "mongo.") == 0 && len(anno.Key) > 6 {
					hasInterface = true
					break
				}
			}
			if hasInterface {
				rawStruct := newIdlExtractStruct(util.CamelString(st.Name))
				if err = extractIdlStruct(st, file, rawStruct); err != nil {
					return err
				}

				if len(rawStruct.StructFields) != 0 {
					rawStructs = append(rawStructs, rawStruct)

					tokens := make([]string, 0, 10)
					methods := ""
					for _, anno := range st.Annotations {
						if strings.Index(anno.Key, "mongo.") == 0 {
							methods += anno.GetValues()[0] + "\n"
							tokens = append(tokens, anno.Key[6:])
						}
					}

					if err = rawStruct.recordMongoIfInfo(info.DocArgs.DaoDir); err != nil {
						return err
					}

					rawInterface := fmt.Sprintf("package main\ntype %sInterface interface{\n%s\n}", st.Name, methods)
					if err = extractIdlInterface(rawInterface, rawStruct, tokens); err != nil {
						return err
					}
				}
			}
		}
		for _, include := range file.Includes {
			if err = getGenGoFilePath(include.Reference); err != nil {
				return err
			}
		}
		return nil
	}
	if err = getGenGoFilePath(info.Req.AST); err != nil {
		return
	}
	return
}

func extractIdlStruct(st *parser.StructLike, file *parser.Thrift, rawStruct *IdlExtractStruct) error {
	for _, field := range st.Fields {
		fag := field.Annotations.Get("go.tag")
		if len(field.Annotations) > 0 && fag != nil && strings.Contains(fag[0], "mongo.bson") {
			tag := handleTagOmitempty(fag[0])

			t := convertThriftType(field.Type, file)
			if t == nil {
				return fmt.Errorf("unsupported type: %s", field.Type.Name)
			}
			if isThriftBaseType(field.Type.Name) || isThriftContainerType(field.Type.Name) {
				sf := &StructField{
					Name: util.CamelString(field.Name),
					Type: t,
					Tag:  tag,
				}
				rawStruct.StructFields = append(rawStruct.StructFields, sf)
			} else if strings.Contains(field.Type.Name, ".") {
				index := strings.Index(field.Type.Name, ".")
				fileName := field.Type.Name[:index]
				structName := field.Type.Name[index+1:]

				var subStruct *parser.StructLike
				var f *parser.Include
				for _, f = range file.Includes {
					name := filepath.Base(f.Reference.Filename)
					if strings.Contains(name, fileName) {
						for _, s := range f.Reference.Structs {
							if s.Name == structName {
								subStruct = s
								break
							}
						}
						break
					}
				}

				// enum
				if subStruct == nil {
					sf := &StructField{
						Name: util.CamelString(field.Name),
						Type: t,
						Tag:  tag,
					}
					rawStruct.StructFields = append(rawStruct.StructFields, sf)
				} else {
					rs := &IdlExtractStruct{
						Name:         subStruct.Name,
						StructFields: make([]*StructField, 0, 10),
					}
					if err := extractIdlStruct(subStruct, f.Reference, rs); err != nil {
						return err
					}
					sf := &StructField{
						Name:               util.CamelString(field.Name),
						Type:               t,
						Tag:                tag,
						IsBelongedToStruct: true,
						BelongedToStruct:   rs,
					}
					rawStruct.StructFields = append(rawStruct.StructFields, sf)
				}
			} else {
				var subStruct *parser.StructLike
				for _, s := range file.Structs {
					if field.Type.Name == s.Name {
						subStruct = s
						break
					}
				}

				// enum
				if subStruct == nil {
					sf := &StructField{
						Name: util.CamelString(field.Name),
						Type: t,
						Tag:  tag,
					}
					rawStruct.StructFields = append(rawStruct.StructFields, sf)
				} else {
					rs := &IdlExtractStruct{
						Name:         subStruct.Name,
						StructFields: make([]*StructField, 0, 10),
					}
					if err := extractIdlStruct(subStruct, file, rs); err != nil {
						return err
					}
					sf := &StructField{
						Name:               util.CamelString(field.Name),
						Type:               t,
						Tag:                tag,
						IsBelongedToStruct: true,
						BelongedToStruct:   rs,
					}
					rawStruct.StructFields = append(rawStruct.StructFields, sf)
				}
			}
		}
	}
	return nil
}

func isThriftBaseType(t string) bool {
	return t == "byte" || t == "i8" || t == "i16" || t == "i32" || t == "i64" ||
		t == "bool" || t == "string" || t == "double" || t == "binary"
}

func isThriftContainerType(t string) bool {
	return t == "map" || t == "set" || t == "list"
}

func convertThriftType(node *parser.Type, file *parser.Thrift) code.Type {
	if node == nil {
		return nil
	}

	if node.KeyType == nil && node.ValueType == nil {
		if v, ok := thriftBaseTypeMap[node.Name]; ok {
			return code.IdentType(v)
		} else if node.Name == "binary" {
			return code.SliceType{
				ElementType: code.IdentType("byte"),
			}
		} else if strings.Contains(node.Name, ".") {
			index := strings.Index(node.Name, ".")
			fileName := node.Name[:index]
			structName := node.Name[index+1:]

			var ff *parser.Include
			isSt := false
			isEnum := false
			for _, f := range file.Includes {
				name := filepath.Base(f.Reference.Filename)
				if strings.Contains(name, fileName) {
					for _, s := range f.Reference.Structs {
						if s.Name == structName {
							ff = f
							isSt = true
							break
						}
					}
					for _, s := range f.Reference.Enums {
						if s.Name == structName {
							ff = f
							isEnum = true
							break
						}
					}
					break
				}
			}

			if ff == nil {
				return nil
			}
			includePackageName := strings.Split(ff.Reference.Namespaces[0].Name, ".")
			if isSt {
				return code.StarExprType{
					RealType: code.SelectorExprType{
						X:   includePackageName[len(includePackageName)-1],
						Sel: node.Name[index+1:],
					},
				}
			}
			if isEnum {
				return code.SelectorExprType{
					X:   includePackageName[len(includePackageName)-1],
					Sel: node.Name[index+1:],
				}
			}
		} else {
			curPackageName := strings.Split(file.Namespaces[0].Name, ".")

			isSt := false
			isEnum := false
			for _, s := range file.Structs {
				if s.Name == node.Name {
					isSt = true
					break
				}
			}
			for _, s := range file.Enums {
				if s.Name == node.Name {
					isEnum = true
					break
				}
			}
			if isSt {
				// struct
				return code.StarExprType{
					RealType: code.SelectorExprType{
						X:   curPackageName[len(curPackageName)-1],
						Sel: node.Name,
					},
				}
			} else if isEnum {
				// enum
				return code.SelectorExprType{
					X:   curPackageName[len(curPackageName)-1],
					Sel: node.Name,
				}
			} else {
				return nil
			}
		}
	}

	if node.KeyType == nil && node.ValueType != nil {
		return code.SliceType{
			ElementType: convertThriftType(node.ValueType, file),
		}
	}

	if node.KeyType != nil && node.ValueType != nil {
		return code.MapType{
			KeyType:   convertThriftType(node.KeyType, file),
			ValueType: convertThriftType(node.ValueType, file),
		}
	}

	return nil
}

var thriftBaseTypeMap = map[string]string{
	"byte":   "int8",
	"i8":     "int8",
	"i16":    "int16",
	"i32":    "int32",
	"i64":    "int64",
	"bool":   "bool",
	"string": "string",
	"double": "float64",
}

func handleTagOmitempty(s string) reflect.StructTag {
	commaIndex := strings.Index(s, ",")
	var tag reflect.StructTag
	if commaIndex == -1 {
		tag = reflect.StructTag(s)
	} else {
		tag = reflect.StructTag(s[:commaIndex] + "\"")
	}
	return tag
}
