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
	"fmt"

	"github.com/cloudwego/cwgo/pkg/doc/mongo/code"
	"github.com/cloudwego/cwgo/pkg/doc/mongo/plugin/model"
)

type DeleteParse struct {
	// OperateMode One or Many
	OperateMode OperateMode

	// Query defines the Query information contained in the Delete operation
	Query *Query

	// CtxParamName defines the method's context.Context param name
	CtxParamName string

	// BelongedToMethod defines the method to which Delete belongs
	BelongedToMethod *model.InterfaceMethod
}

func newDeleteParse() *DeleteParse {
	return &DeleteParse{Query: newQuery()}
}

func (dp *DeleteParse) GetOperationName() string {
	return Delete
}

// parseDelete can be called independently or by Bulk or by Transaction, when isCalled = false,  is called independently
//
//	input params description:
//	tokens: it contains all tokens belonging to Delete except for Delete token
//	method: the method to which Delete belongs
//	curParamIndex: current method's param index
//	isCalled: false ==> independently true ==> called by Bulk or Transaction
func (dp *DeleteParse) parseDelete(tokens []string, method *model.InterfaceMethod, curParamIndex *int, isCalled bool) error {
	if !isCalled {
		if err := dp.check(method); err != nil {
			return err
		}
	}

	dp.BelongedToMethod = method

	if err := dp.parseQuery(tokens, method, curParamIndex); err != nil {
		return err
	}

	if !isCalled {
		if *curParamIndex < len(method.Params) {
			return newMethodSyntaxError(method.Name, fmt.Sprintf("too many method parameters written, "+
				"%v and subsequent parameters are useless", method.Params[*curParamIndex].Name))
		}
	}

	return nil
}

func (dp *DeleteParse) check(method *model.InterfaceMethod) error {
	if len(method.Params) < 1 {
		return newMethodSyntaxError(method.Name, "less than one input parameters")
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

	if t, ok := method.Returns[0].(code.IdentType); ok {
		if string(t) == "bool" {
			dp.OperateMode = OperateOne
		} else if string(t) == "int" {
			dp.OperateMode = OperateMany
		} else {
			return newMethodSyntaxError(method.Name, "the first parameter in the return parameters "+
				"should be bool or int")
		}
	} else {
		return newMethodSyntaxError(method.Name, "the first parameter in the return parameters "+
			"should be bool or int")
	}

	dp.CtxParamName = method.Params[0].Name

	return nil
}

func (dp *DeleteParse) parseQuery(tokens []string, method *model.InterfaceMethod, curParamIndex *int) error {
	fqIndex, err := getFirstQueryIndex(tokens)
	if err != nil {
		return newMethodSyntaxError(method.Name, err.Error())
	}
	if err = dp.Query.parseQuery(tokens[fqIndex:], method, curParamIndex); err != nil {
		return err
	}
	return nil
}
