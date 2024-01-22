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

type BulkParse struct {
	// Operations defines all the operations contained in the Bulk,
	// supports Insert One(Many not support), Update One Many, Delete One Many
	Operations []Operation

	// CtxParamName defines the method's context.Context param name when Bulk is called independently
	CtxParamName string

	// BelongedToMethod defines the method to which Bulk belongs
	BelongedToMethod *model.InterfaceMethod
}

func newBulkParse() *BulkParse {
	return &BulkParse{Operations: []Operation{}}
}

func (bp *BulkParse) GetOperationName() string {
	return Bulk
}

// parseBulk can be called independently or by Transaction, when isCalled = false,  is called independently
//
//	input params description:
//	tokens: it contains all tokens belonging to the Bulk except for the Bulk token
//	method: the method to which Bulk belongs
//	curParamIndex: current method's param index
//	isCalled: false ==> independently true ==> called by Transaction
func (bp *BulkParse) parseBulk(tokens []string, method *model.InterfaceMethod, curParamIndex *int, isCalled bool) error {
	if !isCalled {
		if err := bp.check(method); err != nil {
			return err
		}
	}

	bp.BelongedToMethod = method

	for index := 0; index < len(tokens); index++ {
		if tokens[index] == Find || tokens[index] == Count || tokens[index] == Bulk || tokens[index] == Transaction {
			return newMethodSyntaxError(method.Name, "the Bulk operation does not supports Find, Count, "+
				"Bulk, Transaction, only supports Insert, Update, Delete")
		}

		if tokens[index] == Insert {
			if index == len(tokens)-1 {
				return newMethodSyntaxError(method.Name, "Insert should be followed by One")
			}

			if tokens[index+1] == Many {
				return newMethodSyntaxError(method.Name, "the Bulk operation does not supports Insert "+
					"Many, only supports Insert One")
			}

			if tokens[index+1] == One {
				ip := newInsertParse()
				if err := ip.parseInsert(method, curParamIndex, true); err != nil {
					return err
				}
				ip.OperateMode = OperateOne
				index += 1
				bp.Operations = append(bp.Operations, ip)
			} else {
				return newMethodSyntaxError(method.Name, "Insert should be followed by One")
			}
		}

		if tokens[index] == Update {
			if index == len(tokens)-1 {
				return newMethodSyntaxError(method.Name, "Update should be followed by One or Many")
			}

			if tokens[index+1] == One || tokens[index+1] == Many {
				if index+1 == len(tokens)-1 {
					return newMethodSyntaxError(method.Name, fmt.Sprintf("there is no content after Update %s",
						tokens[index+1]))
				}

				noIndex := getNextOperationIndex(tokens, index+2, false)
				up := newUpdateParse()
				if err := up.parseUpdate(tokens[index+2:noIndex], method, curParamIndex, true); err != nil {
					return err
				}
				if tokens[index+1] == One {
					up.OperateMode = OperateOne
				} else {
					up.OperateMode = OperateMany
				}
				index = noIndex - 1
				bp.Operations = append(bp.Operations, up)
			} else {
				return newMethodSyntaxError(method.Name, "Update should be followed by One or Many")
			}
		}

		if tokens[index] == Delete {
			if index == len(tokens)-1 {
				return newMethodSyntaxError(method.Name, "Delete should be followed by One or Many")
			}

			if tokens[index+1] == One || tokens[index+1] == Many {
				if index+1 == len(tokens)-1 {
					return newMethodSyntaxError(method.Name, fmt.Sprintf("there is no content after Delete %s",
						tokens[index+1]))
				}
				noIndex := getNextOperationIndex(tokens, index+2, false)
				dp := newDeleteParse()
				if err := dp.parseDelete(tokens[index+2:noIndex], method, curParamIndex, true); err != nil {
					return err
				}
				if tokens[index+1] == One {
					dp.OperateMode = OperateOne
				} else {
					dp.OperateMode = OperateMany
				}
				index = noIndex - 1
				bp.Operations = append(bp.Operations, dp)
			} else {
				return newMethodSyntaxError(method.Name, "Delete should be followed by One or Many")
			}
		}
	}

	if !isCalled {
		if *curParamIndex < len(method.Params) {
			return newMethodSyntaxError(method.Name, fmt.Sprintf("too many method parameters written, "+
				"%v and subsequent parameters are useless", method.Params[*curParamIndex].Name))
		}
	}

	return nil
}

func (bp *BulkParse) check(method *model.InterfaceMethod) error {
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

	if method.Returns[0].RealName() != "*mongo.BulkWriteResult" {
		return newMethodSyntaxError(method.Name, "the first parameter in the return parameters "+
			"should be *mongo.BulkWriteResult")
	}

	if method.Returns[1].RealName() != "error" {
		return newMethodSyntaxError(method.Name, "the second parameter in the return parameters "+
			"should be error")
	}

	bp.CtxParamName = method.Params[0].Name

	return nil
}

// getNextOperationIndex is used to get the next operation index.
// If hasCollection = true, returns the index of the second specified operation obtained through traversal;
// If hasCollection = false, returns the index of the first specified operation obtained through traversal;
func getNextOperationIndex(tokens []string, startIndex int, hasCollection bool) int {
	noIndex := -1
	count := 0
	for i := startIndex; i < len(tokens); i++ {
		if tokens[i] == Insert || tokens[i] == Find || tokens[i] == Update || tokens[i] == Delete ||
			tokens[i] == Count || tokens[i] == Transaction || tokens[i] == Bulk || tokens[i] == collection {
			if !hasCollection && count == 0 {
				noIndex = i
				break
			}
			if hasCollection && count == 1 {
				noIndex = i
				break
			}
			count++
		}
	}
	if noIndex == -1 {
		noIndex = len(tokens)
	}
	return noIndex
}
