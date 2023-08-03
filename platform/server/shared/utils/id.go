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
	typeid "go.jetpack.io/typeid/typed"
)

type servicePrefix struct{}
type taskPrefix struct{}

func (servicePrefix) Type() string { return "service" }
func (taskPrefix) Type() string    { return "task" }

type serviceId struct{ typeid.TypeID[servicePrefix] }
type taskId struct{ typeid.TypeID[taskId] }

func NewServiceId() (string, error) {
	id, err := typeid.New[servicePrefix]()
	return id.String(), err
}

func NewTaskId() (string, error) {
	id, err := typeid.New[taskPrefix]()
	return id.String(), err
}
