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

	"github.com/cloudwego/cwgo/pkg/doc/mongo/plugin/model"
)

type CountParse struct {
	// Query defines the Query information contained in the Count operation
	Query *Query

	// CtxParamName defines the method's context.Context param name
	CtxParamName string

	// BelongedToMethod defines the method to which Count belongs
	BelongedToMethod *model.InterfaceMethod
}

func newCountParse() *CountParse {
	return &CountParse{Query: newQuery()}
}

func (cp *CountParse) GetOperationName() string {
	return Count
}

// parseCount can be called independently.
//
//	input params description:
//	tokens: it contains all tokens belonging to Count except for Count token
//	method: the method to which Count belongs
//	curParamIndex: current method's param index
func (cp *CountParse) parseCount(tokens []string, method *model.InterfaceMethod, curParamIndex *int) error {
	if err := cp.check(method); err != nil {
		return err
	}

	cp.BelongedToMethod = method

	if err := cp.parseQuery(tokens, method, curParamIndex); err != nil {
		return err
	}

	if *curParamIndex < len(method.Params) {
		return newMethodSyntaxError(method.Name, fmt.Sprintf("too many method parameters written, "+
			"%v and subsequent parameters are useless", method.Params[*curParamIndex].Name))
	}

	return nil
}

func (cp *CountParse) check(method *model.InterfaceMethod) error {
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

	if method.Returns[0].RealName() != "int" {
		return newMethodSyntaxError(method.Name, "the first parameter in the return parameters "+
			"should be int")
	}

	if method.Returns[1].RealName() != "error" {
		return newMethodSyntaxError(method.Name, "the second parameter in the return parameters "+
			"should be error")
	}

	cp.CtxParamName = method.Params[0].Name

	return nil
}

func (cp *CountParse) parseQuery(tokens []string, method *model.InterfaceMethod, curParamIndex *int) error {
	fqIndex, err := getFirstQueryIndex(tokens)
	if err != nil {
		return newMethodSyntaxError(method.Name, err.Error())
	}
	if err = cp.Query.parseQuery(tokens[fqIndex:], method, curParamIndex); err != nil {
		return err
	}
	return nil
}
