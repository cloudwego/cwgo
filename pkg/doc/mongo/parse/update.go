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

type UpdateParse struct {
	// OperateMode One or Many
	OperateMode OperateMode

	// UpdateStructObjName defines the method's structure point param name,
	// is used when updating the entire structure.
	UpdateStructObjName string

	// Query defines the Query information contained in the Update operation
	Query *Query

	// CtxParamName defines the method's context.Context param name
	CtxParamName string

	// BelongedToMethod defines the method to which Update belongs
	BelongedToMethod *model.InterfaceMethod

	Upsert       bool
	UpdateFields []UpdateField
}

type UpdateField struct {
	MongoFieldName string
	ParamName      string
}

func newUpdateParse() *UpdateParse {
	return &UpdateParse{UpdateFields: []UpdateField{}, Query: newQuery()}
}

func (up *UpdateParse) GetOperationName() string {
	return Update
}

const upsert = "Upsert"

// parseUpdate can be called independently or by Bulk or by Transaction, when isCalled = false,  is called independently
//
//	input params description:
//	tokens: it contains all tokens belonging to Update except for Update token
//	method: the method to which Update belongs
//	curParamIndex: current method's param index
//	isCalled: false ==> independently true ==> called by Bulk or Transaction
func (up *UpdateParse) parseUpdate(tokens []string, method *model.InterfaceMethod, curParamIndex *int, isCalled bool) error {
	if !isCalled {
		if err := up.check(method); err != nil {
			return err
		}
	}

	up.BelongedToMethod = method

	fqIndex, err := getFirstQueryIndex(tokens)
	if err != nil {
		return newMethodSyntaxError(method.Name, err.Error())
	}

	up.parseUpdateOptions(tokens)

	if up.Upsert {
		if err = up.parseUpdateField(tokens[1:fqIndex], method, curParamIndex); err != nil {
			return err
		}
	} else {
		if err = up.parseUpdateField(tokens[:fqIndex], method, curParamIndex); err != nil {
			return err
		}
	}

	if err = up.Query.parseQuery(tokens[fqIndex:], method, curParamIndex); err != nil {
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

func (up *UpdateParse) check(method *model.InterfaceMethod) error {
	if len(method.Params) < 2 {
		return newMethodSyntaxError(method.Name, "less than two input parameters")
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
			up.OperateMode = OperateOne
		} else if string(t) == "int" {
			up.OperateMode = OperateMany
		} else {
			return newMethodSyntaxError(method.Name, "the first parameter in the return parameters "+
				"should be bool or int")
		}
	} else {
		return newMethodSyntaxError(method.Name, "the first parameter in the return parameters "+
			"should be bool or int")
	}

	up.CtxParamName = method.Params[0].Name

	return nil
}

func (up *UpdateParse) parseUpdateOptions(tokens []string) {
	if tokens[0] == upsert {
		up.Upsert = true
	}
}

func (up *UpdateParse) parseUpdateField(tokens []string, method *model.InterfaceMethod, curParamIndex *int) error {
	if len(tokens) == 0 {
		t, ok := method.Params[1].Type.(code.StarExprType)
		if !ok {
			return newMethodSyntaxError(method.Name, "the input when updating the whole structure is not a structure pointer")
		}

		if _, ok = t.RealType.(code.SelectorExprType); !ok {
			return newMethodSyntaxError(method.Name, "the input when updating the whole structure is not in the form of *Package.StructName")
		}

		up.UpdateStructObjName = method.Params[1].Name
		*curParamIndex += 1
		return nil
	}

	curIndex := new(int)
	*curIndex = -1
	result, t, err := getFieldNameType(tokens, method.BelongedToStruct, curIndex, true)
	if err != nil {
		return err
	}

	for i := 0; i < len(result); i++ {
		if i+1 >= len(method.Params) {
			return newMethodSyntaxError(method.Name, "insufficient number of input parameters")
		}
		if method.Params[i+*curParamIndex].Type.RealName() != t[i].RealName() {
			return newMethodSyntaxError(method.Name,
				fmt.Sprintf("the field type in the parameter transfer: %s, the actual required field type: %s",
					method.Params[i].Type.RealName(), t[0].RealName()))
		}
		up.UpdateFields = append(up.UpdateFields, UpdateField{
			MongoFieldName: result[i],
			ParamName:      method.Params[i+1].Name,
		})
	}
	*curParamIndex += len(result)

	return nil
}
