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
	"strings"

	"github.com/cloudwego/cwgo/pkg/doc/mongo/plugin/model"
)

type TransactionParse struct {
	// TransactionOperations defines all the operations contained in the Transaction,
	// supports Insert One Many, Update One Many, Delete One Many, Bulk(Mark boundaries with parentheses).
	TransactionOperations []TransactionOperation

	// CtxParamName defines the method's context.Context param name
	CtxParamName string

	// ClientParamName defines the method's *mongo.Client param name
	ClientParamName string

	// collectionParamsMap stores the method's *mongo.Collection param names
	collectionParamsMap map[string]string

	// BelongedToMethod defines the method to which Transaction belongs
	BelongedToMethod *model.InterfaceMethod
}

type TransactionOperation struct {
	// CollectionParamName stores the collection which the operation belongs,
	// if not specified, default is r.collection
	CollectionParamName string
	Operation           Operation
}

func newTransactionParse() *TransactionParse {
	return &TransactionParse{
		TransactionOperations: []TransactionOperation{},
		collectionParamsMap:   map[string]string{},
	}
}

func (tp *TransactionParse) GetOperationName() string {
	return Transaction
}

const (
	defaultCollection = "r.collection"
	collection        = "Collection"
)

// parseTransaction can be called independently.
//
//	input params description:
//	tokens: it contains all tokens belonging to the Transaction except for the Transaction token
//	method: the method to which Transaction belongs
//	curParamIndex: current method's param index
func (tp *TransactionParse) parseTransaction(tokens []string, method *model.InterfaceMethod, curParamIndex *int) error {
	if err := tp.check(method); err != nil {
		return err
	}

	tp.BelongedToMethod = method

	// store method's collection param names
	if len(method.Params) > 2 {
		for _, param := range method.Params {
			if param.Type.RealName() == "*mongo.Collection" {
				collectionParam := ""
				if len(param.Name) == 1 {
					collectionParam = strings.ToUpper(string(param.Name[0]))
				} else {
					collectionParam = strings.ToUpper(string(param.Name[0])) + param.Name[1:]
				}
				tp.collectionParamsMap[collectionParam] = param.Name
				*curParamIndex += 1
			}
		}
	}

	for index := 0; index < len(tokens); index++ {
		if tokens[index] == Find || tokens[index] == Count || tokens[index] == Transaction {
			return newMethodSyntaxError(method.Name, "the Transaction operation does not supports Find, Count, "+
				"Transaction, only supports Insert, Update, Delete, Bulk")
		}

		if tokens[index] == Insert {
			if err := tp.parseTransactionInsert(method, tokens, index, curParamIndex, defaultCollection); err != nil {
				return err
			}
			index += 1
		}

		if tokens[index] == Update {
			noIndex, err := tp.parseTransactionUpdate(method, tokens, index, curParamIndex, defaultCollection, false)
			if err != nil {
				return err
			}
			index = noIndex - 1
		}

		if tokens[index] == Delete {
			noIndex, err := tp.parseTransactionDelete(method, tokens, index, curParamIndex, defaultCollection, false)
			if err != nil {
				return err
			}
			index = noIndex - 1
		}

		if tokens[index] == Bulk {
			noIndex, err := tp.parseTransactionBulk(method, tokens, index, curParamIndex, defaultCollection)
			if err != nil {
				return err
			}
			index = noIndex - 1
		}

		if tokens[index] == collection {
			if index == len(tokens)-1 {
				return newMethodSyntaxError(method.Name, "no tokens specified after Collection")
			}

			belongedToOpIndex := -1
			belongedToOpIndexName := ""
			for i := index + 1; i < len(tokens); i++ {
				if tokens[i] == Insert || tokens[i] == Update || tokens[i] == Delete || tokens[i] == Bulk {
					belongedToOpIndex = i
					belongedToOpIndexName = tokens[i]
					break
				}
			}

			if belongedToOpIndex == -1 {
				return newMethodSyntaxError(method.Name, "there is no Insert, Update, Delete, Bulk "+
					"tokens after the Collection")
			}
			if belongedToOpIndex == index+1 {
				return newMethodSyntaxError(method.Name, "no collection name specified after Collection")
			}

			paramName := strSlice2Str(tokens[index+1 : belongedToOpIndex])
			v, ok := tp.collectionParamsMap[paramName]
			if !ok {
				return newMethodSyntaxError(method.Name, "the collection name specified in tokens was not "+
					"found in the method parameters")
			}

			switch belongedToOpIndexName {
			case Insert:
				if err := tp.parseTransactionInsert(method, tokens, belongedToOpIndex, curParamIndex, v); err != nil {
					return err
				}
				index = belongedToOpIndex + 1

			case Update:
				noIndex, err := tp.parseTransactionUpdate(method, tokens, belongedToOpIndex, curParamIndex, v, true)
				if err != nil {
					return err
				}
				index = noIndex - 1

			case Delete:
				noIndex, err := tp.parseTransactionDelete(method, tokens, belongedToOpIndex, curParamIndex, v, true)
				if err != nil {
					return err
				}
				index = noIndex - 1

			case Bulk:
				noIndex, err := tp.parseTransactionBulk(method, tokens, belongedToOpIndex, curParamIndex, v)
				if err != nil {
					return err
				}
				index = noIndex - 1

			default:
			}
		}
	}

	if *curParamIndex < len(method.Params) {
		return newMethodSyntaxError(method.Name, fmt.Sprintf("too many method parameters written, "+
			"%v and subsequent parameters are useless", method.Params[*curParamIndex].Name))
	}

	return nil
}

func (tp *TransactionParse) check(method *model.InterfaceMethod) error {
	if len(method.Params) < 2 {
		return newMethodSyntaxError(method.Name, "less than two input parameters")
	}

	if len(method.Returns) != 1 {
		return newMethodSyntaxError(method.Name, "return parameter not equal to 1")
	}

	if method.Params[0].Type.RealName() != "context.Context" {
		return newMethodSyntaxError(method.Name, "the first parameter in the input parameters "+
			"should be context.Context")
	}

	if method.Params[1].Type.RealName() != "*mongo.Client" {
		return newMethodSyntaxError(method.Name, "the second parameter in the input parameters "+
			"should be *mongo.Client")
	}

	if method.Returns[0].RealName() != "error" {
		return newMethodSyntaxError(method.Name, "the only parameter in the return parameters "+
			"should be error")
	}

	tp.CtxParamName = method.Params[0].Name
	tp.ClientParamName = method.Params[1].Name

	return nil
}

func (tp *TransactionParse) parseTransactionInsert(method *model.InterfaceMethod, tokens []string,
	index int, curParamIndex *int, collectionParamName string,
) error {
	if index == len(tokens)-1 {
		return newMethodSyntaxError(method.Name, "no tokens specified after Insert")
	}

	if tokens[index+1] == One || tokens[index+1] == Many {
		ip := newInsertParse()
		if err := ip.parseInsert(method, curParamIndex, true); err != nil {
			return err
		}

		if tokens[index+1] == One {
			ip.OperateMode = OperateOne
		} else {
			ip.OperateMode = OperateMany
		}

		tp.TransactionOperations = append(tp.TransactionOperations, TransactionOperation{
			CollectionParamName: collectionParamName,
			Operation:           ip,
		})
		return nil
	} else {
		return newMethodSyntaxError(method.Name, "no One or Many specified after Insert")
	}
}

func (tp *TransactionParse) parseTransactionUpdate(method *model.InterfaceMethod, tokens []string,
	index int, curParamIndex *int, collectionParamName string, hasCollection bool,
) (int, error) {
	if index == len(tokens)-1 {
		return 0, newMethodSyntaxError(method.Name, "no tokens specified after Update")
	}

	if tokens[index+1] == One || tokens[index+1] == Many {
		if index+1 == len(tokens)-1 {
			return 0, newMethodSyntaxError(method.Name, fmt.Sprintf("no tokens specified after Update %s", tokens[index+1]))
		}

		noIndex := getNextOperationIndex(tokens, index+2, hasCollection)
		up := newUpdateParse()
		if err := up.parseUpdate(tokens[index+2:noIndex], method, curParamIndex, true); err != nil {
			return 0, err
		}

		if tokens[index+1] == One {
			up.OperateMode = OperateOne
		} else {
			up.OperateMode = OperateMany
		}

		tp.TransactionOperations = append(tp.TransactionOperations, TransactionOperation{
			CollectionParamName: collectionParamName,
			Operation:           up,
		})
		return noIndex, nil
	} else {
		return 0, newMethodSyntaxError(method.Name, "no One or Many specified after Update")
	}
}

func (tp *TransactionParse) parseTransactionDelete(method *model.InterfaceMethod, tokens []string,
	index int, curParamIndex *int, collectionParamName string, hasCollection bool,
) (int, error) {
	if index == len(tokens)-1 {
		return 0, newMethodSyntaxError(method.Name, "no tokens specified after Delete")
	}

	if tokens[index+1] == One || tokens[index+1] == Many {
		if index+1 == len(tokens)-1 {
			return 0, newMethodSyntaxError(method.Name, fmt.Sprintf("no tokens specified after Delete %s", tokens[index+1]))
		}

		noIndex := getNextOperationIndex(tokens, index+2, hasCollection)
		dp := newDeleteParse()
		if err := dp.parseDelete(tokens[index+2:noIndex], method, curParamIndex, true); err != nil {
			return 0, err
		}
		if tokens[index+1] == One {
			dp.OperateMode = OperateOne
		} else {
			dp.OperateMode = OperateMany
		}

		tp.TransactionOperations = append(tp.TransactionOperations, TransactionOperation{
			CollectionParamName: collectionParamName,
			Operation:           dp,
		})
		return noIndex, nil
	} else {
		return 0, newMethodSyntaxError(method.Name, "no One or Many specified after Delete")
	}
}

func (tp *TransactionParse) parseTransactionBulk(method *model.InterfaceMethod, tokens []string,
	index int, curParamIndex *int, collectionParamName string,
) (int, error) {
	if index == len(tokens)-1 {
		return 0, newMethodSyntaxError(method.Name, "no tokens specified after Bulk")
	}

	bp := newBulkParse()
	if tokens[index+1] != leftBracket {
		return 0, newMethodSyntaxError(method.Name, "parentheses need to be specified after Transaction "+
			"Bulk operation")
	}

	rbIndex := -1
	for i := index + 2; i < len(tokens); i++ {
		if tokens[i] == rightBracket {
			rbIndex = i
			break
		}
	}
	if rbIndex == -1 {
		return 0, newMethodSyntaxError(method.Name, "the Transaction Bulk operation did not specify a "+
			"right parenthesis")
	}

	if err := bp.parseBulk(tokens[index+2:rbIndex], method, curParamIndex, true); err != nil {
		return 0, err
	}
	tp.TransactionOperations = append(tp.TransactionOperations, TransactionOperation{
		CollectionParamName: collectionParamName,
		Operation:           bp,
	})

	if rbIndex == len(tokens)-1 {
		return rbIndex, nil
	} else {
		return rbIndex + 1, nil
	}
}

func strSlice2Str(ss []string) (result string) {
	for _, s := range ss {
		result += s
	}
	return
}
