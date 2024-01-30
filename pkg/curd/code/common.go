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
 *
 * MIT License
 *
 * Copyright (c) 2021 Surawich Laprattanatrai
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package code

import (
	"fmt"
	"reflect"
)

type StructFields []StructField

func (sfs StructFields) GetCode() string {
	result := ""
	for index, field := range sfs {
		if index != len(sfs)-1 {
			result += field.GetCode() + "\n"
		} else {
			result += field.GetCode()
		}
	}
	return result
}

type StructField struct {
	Name string
	Type Type
	Tag  reflect.StructTag
}

func (sf *StructField) GetCode() string {
	return fmt.Sprintf("\t%s %s %s", sf.Name, sf.Type.RealName(), string(sf.Tag))
}

type Params []Param

func (ps Params) GetCode() string {
	return getParamsCode(ps)
}

type Returns []Type

func (rs Returns) GetCode() string {
	if len(rs) == 1 {
		return rs[0].RealName()
	} else {
		result := "("
		for index, param := range rs {
			if index != len(rs)-1 {
				result += param.RealName() + ", "
			} else {
				result += param.RealName() + ")"
			}
		}
		return result
	}
}

func getParamsCode(ps []Param) string {
	result := "("
	for index, param := range ps {
		if index != len(ps)-1 {
			result += param.GetCode() + ", "
		} else {
			result += param.GetCode() + ")"
		}
	}
	return result
}

type Param struct {
	Name string
	Type Type
}

func (p *Param) GetCode() string {
	if p.Name != "" {
		return fmt.Sprintf("%s %s", p.Name, p.Type.RealName())
	} else {
		return p.Type.RealName()
	}
}

type MethodReceiver struct {
	Name string
	Type Type
}

func (mr *MethodReceiver) GetCode() string {
	return fmt.Sprintf("(%s %s)", mr.Name, mr.Type.RealName())
}

type Body []Statement

func (b Body) GetCode() string {
	result := ""
	for index, statement := range b {
		if index != len(b)-1 {
			result += "\t" + statement.Code() + "\n"
		} else {
			result += "\t" + statement.Code()
		}
	}
	return result
}

type InterfaceMethods []InterfaceMethod

func (ims InterfaceMethods) GetCode() string {
	result := ""
	for index, mh := range ims {
		if index != len(ims)-1 {
			result += mh.GetCode() + "\n"
		} else {
			result += mh.GetCode()
		}
	}
	return result
}

type InterfaceMethod struct {
	Comment string
	Name    string
	Params  Params
	Returns Returns
}

func (im *InterfaceMethod) GetCode() string {
	if im.Comment != "" {
		return fmt.Sprintf("\t%s\n\t%s%s %s",
			im.Comment,
			im.Name, im.Params.GetCode(),
			im.Returns.GetCode())
	} else {
		return fmt.Sprintf("\t%s%s %s",
			im.Name, im.Params.GetCode(),
			im.Returns.GetCode())
	}
}
