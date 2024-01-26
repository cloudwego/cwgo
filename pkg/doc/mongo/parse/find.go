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
	"errors"
	"fmt"

	"github.com/cloudwego/cwgo/pkg/doc/mongo/extract"

	"github.com/cloudwego/cwgo/pkg/doc/mongo/code"
)

type FindParse struct {
	// OperateMode One or Many
	OperateMode OperateMode

	// Query defines the Query information contained in the Find operation
	Query *Query

	Project        []string
	Order          Order
	SkipParamName  string
	LimitParamName string

	// CtxParamName defines the method's context.Context param name
	CtxParamName string

	// ReturnType defines the method's first return parameter's Type which Find belongs
	ReturnType code.Type

	// BelongedToMethod defines the method to which Find belongs
	BelongedToMethod *extract.InterfaceMethod
}

type Order struct {
	Asc  []string
	Desc []string
}

const (
	order = "Order"
	skip  = "Skip"
	limit = "Limit"
	desc  = "Desc"
)

func newFindParse() *FindParse {
	return &FindParse{
		Project: []string{},
		Order: Order{
			Asc:  []string{},
			Desc: []string{},
		},
		Query: newQuery(),
	}
}

func (fp *FindParse) GetOperationName() string {
	return Find
}

// parseFind can be called independently.
//
//	input params description:
//	tokens: it contains all tokens belonging to Find except for Find token
//	method: the method to which Find belongs
//	curParamIndex: current method's param index
func (fp *FindParse) parseFind(tokens []string, method *extract.InterfaceMethod, curParamIndex *int) error {
	if err := fp.check(method); err != nil {
		return err
	}

	fp.BelongedToMethod = method

	tokenIndex, err := fp.parseProject(tokens, method.BelongedToStruct)
	if err != nil {
		return newMethodSyntaxError(method.Name, err.Error())
	}

	if err = fp.parseFindOptions(tokens[tokenIndex:], method, curParamIndex); err != nil {
		return err
	}

	fqIndex, err := getFirstQueryIndex(tokens[tokenIndex:])
	if err != nil {
		return newMethodSyntaxError(method.Name, err.Error())
	}

	if err = fp.Query.parseQuery(tokens[tokenIndex+fqIndex:], method, curParamIndex); err != nil {
		return err
	}

	if *curParamIndex < len(method.Params) {
		return newMethodSyntaxError(method.Name, fmt.Sprintf("too many method parameters written, "+
			"%v and subsequent parameters are useless", method.Params[*curParamIndex].Name))
	}

	return nil
}

func (fp *FindParse) check(method *extract.InterfaceMethod) error {
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

	if _, ok := method.Returns[0].(code.StarExprType); ok {
		fp.OperateMode = OperateOne
	} else if _, ok = method.Returns[0].(code.SliceType); ok {
		fp.OperateMode = OperateMany
	} else {
		return newMethodSyntaxError(method.Name, "the first parameter in the return parameters input error")
	}

	fp.CtxParamName = method.Params[0].Name
	fp.ReturnType = method.Returns[0]

	return nil
}

func (fp *FindParse) parseProject(tokens []string, extractStruct *extract.IdlExtractStruct) (int, error) {
	tokenIndex, err := getNextTokenIndex(tokens, 0)
	if err != nil {
		return 0, err
	}

	if tokenIndex == 0 {
		return 0, nil
	}

	curIndex := new(int)
	*curIndex = -1
	result, _, err := getFieldNameType(tokens[:tokenIndex], extractStruct, curIndex, true)
	if err != nil {
		return 0, err
	}

	fp.Project = result
	return tokenIndex, nil
}

func (fp *FindParse) parseFindOptions(tokens []string, method *extract.InterfaceMethod, curParamIndex *int) error {
	orderFlag, skipFlag, limitFlag := 0, 0, 0

	for index, token := range tokens {
		if token == order {
			if orderFlag == 1 {
				return newMethodSyntaxError(method.Name, "Order can only be used once")
			}
			if index == len(tokens)-1 {
				return newMethodSyntaxError(method.Name, "there are no sorted fields after the Order")
			}

			tokenIndex, err := getNextTokenIndex(tokens, index+1)
			if err != nil {
				return newMethodSyntaxError(method.Name, err.Error())
			}

			if index+1 == tokenIndex {
				return newMethodSyntaxError(method.Name, "there are no sorted fields after the Order")
			}
			if err = fp.getSortFields(tokens[index+1:tokenIndex], method.BelongedToStruct); err != nil {
				return newMethodSyntaxError(method.Name, err.Error())
			}
			orderFlag = 1
		}

		if token == skip {
			if skipFlag == 1 {
				return newMethodSyntaxError(method.Name, "Skip can only be used once")
			}
			if index == len(tokens)-1 {
				return newMethodSyntaxError(method.Name, "there are no other fields after Skip, such as By or All")
			}

			if method.Params[*curParamIndex].Type.RealName() != "int64" {
				return newMethodSyntaxError(method.Name, "Skip requires passing in a value of type int64")
			}

			fp.SkipParamName = method.Params[*curParamIndex].Name
			*curParamIndex += 1
			skipFlag = 1
		}

		if token == limit {
			if fp.OperateMode == OperateOne {
				return newMethodSyntaxError(method.Name, "Limit operation is not supported in Find One mode")
			}
			if limitFlag == 1 {
				return newMethodSyntaxError(method.Name, "Limit can only be used once")
			}
			if index == len(tokens)-1 {
				return newMethodSyntaxError(method.Name, "there are no other fields after Limit, such as By or All")
			}

			if method.Params[*curParamIndex].Type.RealName() != "int64" {
				return newMethodSyntaxError(method.Name, "Limit requires passing in a value of type int64")
			}

			fp.LimitParamName = method.Params[*curParamIndex].Name
			*curParamIndex += 1
			limitFlag = 1
		}
	}

	return nil
}

func (fp *FindParse) getSortFields(tokens []string, extractStruct *extract.IdlExtractStruct) error {
	preDescIndex := 0
	for i := 0; i < len(tokens); i++ {
		if tokens[i] == desc {
			if i == preDescIndex {
				return errors.New("no field to be specified before desc")
			}

			for j := i - 1; j >= preDescIndex; j-- {
				curIndex := new(int)
				*curIndex = -1
				r, _, err := getFieldNameType(tokens[j:i], extractStruct, curIndex, true)
				if err != nil && j == preDescIndex {
					return err
				}

				if len(r) != 0 {
					repeatFlag := 0
					if j > preDescIndex {
						*curIndex = -1
						r, _, err = getFieldNameType(tokens[preDescIndex:j], extractStruct, curIndex, true)
						if err != nil {
							// To avoid bug: For example, both NameHello and Name are fields of struct,
							// NameHelloDesc can be correctly parsed by this logic.
							for k := j - 1; k >= preDescIndex; k-- {
								r, _, err = getFieldNameType(tokens[k:i], extractStruct, curIndex, true)
								if err != nil && k == preDescIndex {
									return err
								}

								if len(r) != 0 {
									fp.Order.Desc = append(fp.Order.Desc, r...)
									if preDescIndex < k {
										r, _, err = getFieldNameType(tokens[preDescIndex:k], extractStruct, curIndex, true)
										if err != nil {
											return err
										}
										fp.Order.Asc = append(fp.Order.Asc, r...)
									}
									repeatFlag = 1
									break
								}
							}
						} else {
							fp.Order.Asc = append(fp.Order.Asc, r...)
						}
					}

					if repeatFlag == 0 {
						*curIndex = -1
						r, _, err = getFieldNameType(tokens[j:i], extractStruct, curIndex, true)
						if err != nil {
							return err
						}
						fp.Order.Desc = append(fp.Order.Desc, r...)
					}

					preDescIndex = i + 1
					break
				}
			}
		}
	}

	if len(fp.Order.Desc) == 0 {
		curIndex := new(int)
		*curIndex = -1
		r, _, err := getFieldNameType(tokens, extractStruct, curIndex, true)
		if err != nil {
			return err
		}
		fp.Order.Asc = append(fp.Order.Asc, r...)
	}
	return nil
}

func getNextTokenIndex(tokens []string, startIndex int) (int, error) {
	tokenIndex := -1
	for i := startIndex; i < len(tokens); i++ {
		if tokens[i] == order || tokens[i] == skip || tokens[i] == limit ||
			tokens[i] == string(By) || tokens[i] == string(All) {
			tokenIndex = i
			break
		}
	}

	if tokenIndex == -1 {
		return 0, errors.New("no By or All specified")
	}

	return tokenIndex, nil
}
