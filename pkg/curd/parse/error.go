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

package parse

import "fmt"

// newMethodSyntaxError creates syntaxError
func newMethodSyntaxError(methodName, errReason string) error {
	return methodSyntaxError{
		methodName: methodName,
		errReason:  errReason,
	}
}

type methodSyntaxError struct {
	methodName string
	errReason  string
}

func (err methodSyntaxError) Error() string {
	return fmt.Sprintf("method %s has syntax errors, specific reasons: %s", err.methodName, err.errReason)
}
