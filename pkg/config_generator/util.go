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
	"path/filepath"
	"strings"
)

func convertToGoStructName(s string) string {
	var result string
	s = filepath.Base(s)
	if idx := strings.LastIndex(s, "."); idx != -1 {
		s = s[:idx]
	}

	// Split the string by spaces or underscores
	words := strings.FieldsFunc(s, func(r rune) bool {
		return r == ' ' || r == '_' || r == '-'
	})

	// Capitalize the first letter of each word
	for _, word := range words {
		if len(word) > 0 {
			result += strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
		}
	}

	return result
}

// isBasicType checks if the given type is a basic Go type
func isBasicType(typeName string) bool {
	basicTypes := map[string]bool{
		"int":     true,
		"int8":    true,
		"int16":   true,
		"int32":   true,
		"int64":   true,
		"uint":    true,
		"uint8":   true,
		"uint16":  true,
		"uint32":  true,
		"uint64":  true,
		"float32": true,
		"float64": true,
		"bool":    true,
		"string":  true,
		"byte":    true,
		"rune":    true,
	}

	return basicTypes[typeName]
}
