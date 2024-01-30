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

package code

import "fmt"

type Type interface {
	RealName() string
}

type IdentType string

func (st IdentType) RealName() string {
	return string(st)
}

type SelectorExprType struct {
	X   string
	Sel string
}

func (set SelectorExprType) RealName() string {
	return set.X + "." + set.Sel
}

type InterfaceType struct {
	Name    string
	Methods InterfaceMethods
}

func (it InterfaceType) RealName() string {
	return "interface{}"
}

type SliceType struct {
	ElementType Type
}

func (st SliceType) RealName() string {
	return "[]" + st.ElementType.RealName()
}

type MapType struct {
	KeyType   Type
	ValueType Type
}

func (mt MapType) RealName() string {
	return fmt.Sprintf("map[%s]%s", mt.KeyType.RealName(), mt.ValueType.RealName())
}

type StarExprType struct {
	RealType Type
}

func (set StarExprType) RealName() string {
	return "*" + set.RealType.RealName()
}
