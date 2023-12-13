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

package utils

import (
	"net"
	"strings"
)

// IsTokenError Determine if it is a Token error issue
func IsTokenError(err error) bool {
	return strings.Contains(err.Error(), "401 Unauthorized")
}

// IsNetworkError Determine if it is a network timeout issue
func IsNetworkError(err error) bool {
	if netErr, ok := err.(net.Error); ok {
		if netErr.Timeout() {
			return true
		}
	}
	return false
}

// IsFileNotFoundError Determine if it is file not found error issue
func IsFileNotFoundError(err error) bool {
	return strings.Contains(err.Error(), "A file with this name doesn't exist")
}
