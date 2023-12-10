/*
 *
 *  * Copyright 2022 CloudWeGo Authors
 *  *
 *  * Licensed under the Apache License, Version 2.0 (the "License");
 *  * you may not use this file except in compliance with the License.
 *  * You may obtain a copy of the License at
 *  *
 *  *     http://www.apache.org/licenses/LICENSE-2.0
 *  *
 *  * Unless required by applicable law or agreed to in writing, software
 *  * distributed under the License is distributed on an "AS IS" BASIS,
 *  * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  * See the License for the specific language governing permissions and
 *  * limitations under the License.
 *
 */

package utils

import (
	"errors"
	"regexp"

	"github.com/cloudwego/cwgo/platform/server/shared/consts"
)

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
