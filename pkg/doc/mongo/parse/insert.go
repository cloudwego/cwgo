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

import (
	"github.com/cloudwego/cwgo/pkg/doc/mongo/code"
	"github.com/cloudwego/cwgo/pkg/doc/mongo/plugin/model"
)

type InsertParse struct {
	// OperateMode One or Many
	OperateMode OperateMode

	// MethodParamNames defines the method's param names
	MethodParamNames [2]string

	// BelongedToMethod defines the method to which Insert belongs
	BelongedToMethod *model.InterfaceMethod
}

func newInsertParse() *InsertParse {
	return &InsertParse{MethodParamNames: [2]string{}}
}

func (ip *InsertParse) GetOperationName() string {
	return Insert
}

// parseInsert can be called independently or by Bulk or by Transaction, when isCalled = false,  is called independently
//   input params description:
//   method: the method to which Insert belongs
//   curParamIndex: current method's param index
//   isCalled: false ==> independently true ==> called by Bulk or Transaction
func (ip *InsertParse) parseInsert(method *model.InterfaceMethod, curParamIndex *int, isCalled bool) error {
	if !isCalled {
		if err := ip.check(method); err != nil {
			return err
		}
	}

	if !isCalled {
		ip.MethodParamNames = [2]string{
			method.Params[*curParamIndex].Name,
			method.Params[*curParamIndex+1].Name,
		}
	} else {
		ip.MethodParamNames = [2]string{
			method.Params[*curParamIndex].Name,
		}
		*curParamIndex += 1
	}

	ip.BelongedToMethod = method

	return nil
}

func (ip *InsertParse) check(method *model.InterfaceMethod) error {
	if len(method.Params) != 2 {
		return newMethodSyntaxError(method.Name, "input parameter not equal to 2")
	}

	if len(method.Returns) != 2 {
		return newMethodSyntaxError(method.Name, "return parameter not equal to 2")
	}

	if method.Params[0].Type.RealName() != "context.Context" {
		return newMethodSyntaxError(method.Name, "the first parameter in the input parameters "+
			"should be context.Context")
	}

	if method.Returns[1].RealName() != "error" {
		return newMethodSyntaxError(method.Name, "the second parameter in the return parameters "+
			"should be error")
	}

	if _, ok := method.Params[1].Type.(code.StarExprType); ok {
		if method.Returns[0].RealName() != "interface{}" {
			return newMethodSyntaxError(method.Name, "inconsistent types, the first parameter in the "+
				"return parameters should be interface{}")
		}
		ip.OperateMode = OperateOne
	} else if _, ok = method.Params[1].Type.(code.SliceType); ok {
		if method.Returns[0].RealName() != "[]interface{}" {
			return newMethodSyntaxError(method.Name, "inconsistent types, the first parameter in the "+
				"return parameters should be []interface{}")
		}
		ip.OperateMode = OperateMany
	} else {
		return newMethodSyntaxError(method.Name, "the first parameter in the return parameters "+
			"should be interface{} or []interface{}")
	}

	return nil
}
