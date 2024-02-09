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
	"strings"

	"github.com/cloudwego/cwgo/pkg/curd/code"
	"github.com/cloudwego/cwgo/pkg/curd/extract"

	"github.com/fatih/camelcase"
)

// InterfaceOperation is used to store the parsing results of the structure and for use by codegen packages.
type InterfaceOperation struct {
	BelongedToStruct *extract.IdlExtractStruct
	Operations       []Operation
}

const (
	Insert      = "Insert"
	Find        = "Find"
	Update      = "Update"
	Delete      = "Delete"
	Count       = "Count"
	Transaction = "Transaction"
	Bulk        = "Bulk"
)

type OperateMode int

const (
	OperateOne  = OperateMode(0)
	OperateMany = OperateMode(1)
)

const (
	One  = "One"
	Many = "Many"
)

func HandleOperations(structs []*extract.IdlExtractStruct) (result []*InterfaceOperation, err error) {
	for _, st := range structs {
		ifo := newInterfaceOperation()
		if err = ifo.parseInterfaceMethod(st); err != nil {
			return nil, err
		}
		result = append(result, ifo)
	}
	return
}

func newInterfaceOperation() *InterfaceOperation {
	return &InterfaceOperation{Operations: []Operation{}}
}

func (ifo *InterfaceOperation) parseInterfaceMethod(extractStruct *extract.IdlExtractStruct) error {
	for _, method := range extractStruct.InterfaceInfo.Methods {
		tokens := camelcase.Split(method.ParsedTokens)
		switch tokens[0] {
		case Insert:
			curParamIndex := new(int)
			*curParamIndex = 0
			ip := newInsertParse()
			if err := ip.parseInsert(method, curParamIndex, false); err != nil {
				return err
			}
			ifo.BelongedToStruct = extractStruct
			ifo.Operations = append(ifo.Operations, ip)

		case Find:
			curParamIndex := new(int)
			*curParamIndex = 1
			fp := newFindParse()
			if err := fp.parseFind(tokens[1:], method, curParamIndex); err != nil {
				return err
			}
			ifo.BelongedToStruct = extractStruct
			ifo.Operations = append(ifo.Operations, fp)

		case Update:
			curParamIndex := new(int)
			*curParamIndex = 1
			up := newUpdateParse()
			if err := up.parseUpdate(tokens[1:], method, curParamIndex, false); err != nil {
				return err
			}
			ifo.BelongedToStruct = extractStruct
			ifo.Operations = append(ifo.Operations, up)

		case Delete:
			curParamIndex := new(int)
			*curParamIndex = 1
			dp := newDeleteParse()
			if err := dp.parseDelete(tokens[1:], method, curParamIndex, false); err != nil {
				return err
			}
			ifo.BelongedToStruct = extractStruct
			ifo.Operations = append(ifo.Operations, dp)

		case Count:
			curParamIndex := new(int)
			*curParamIndex = 1
			cp := newCountParse()
			if err := cp.parseCount(tokens[1:], method, curParamIndex); err != nil {
				return err
			}
			ifo.BelongedToStruct = extractStruct
			ifo.Operations = append(ifo.Operations, cp)

		case Transaction:
			curParamIndex := new(int)
			*curParamIndex = 2
			tp := newTransactionParse()
			if err := tp.parseTransaction(tokens[1:], method, curParamIndex); err != nil {
				return err
			}
			ifo.BelongedToStruct = extractStruct
			ifo.Operations = append(ifo.Operations, tp)

		case Bulk:
			curParamIndex := new(int)
			*curParamIndex = 1
			bp := newBulkParse()
			if err := bp.parseBulk(tokens[1:], method, curParamIndex, false); err != nil {
				return err
			}
			ifo.BelongedToStruct = extractStruct
			ifo.Operations = append(ifo.Operations, bp)

		default:
			return newMethodSyntaxError(method.Name, "wrong operation name, should be Insert, Find, "+
				"Update, Delete, Count, Transaction, Bulk")
		}
	}

	return nil
}

// getFieldNameType is used to get field names and types in the specified structure.
//
//	input params description:
//	tokens: parsed tokens
//	extractStruct: the structure to which tokens belong
//	curIndex: point to the next token to be parsed
//	isFirst: if it is called in recursion
func getFieldNameType(tokens []string, extractStruct *extract.IdlExtractStruct, curIndex *int, isFirst bool) (names []string,
	types []code.Type, err error,
) {
	if len(tokens) == 0 {
		return nil, nil, errors.New("the length of the field name requested for parsing is empty")
	}

	for i := 0; i < len(tokens); i++ {
		flag := 0
		for _, field := range extractStruct.StructFields {
			if field.Name == tokens[i] || strings.Index(field.Name, tokens[i]) == 0 {
				s := ""
				hasFieldFlag := 0
				for j := i; j < len(tokens); j++ {
					s += tokens[j]
					if s == field.Name {
						hasFieldFlag = 1
						i = j
						*curIndex += i + 1
						break
					}
				}
				if hasFieldFlag == 0 {
					return nil, nil, fmt.Errorf("partially equal but unable to fully locate field name in %v", tokens[i:])
				}

				flag = 1
				if !field.IsBelongedToStruct {
					names = append(names, field.Tag.Get("bson"))
					types = append(types, field.Type)
					break
				} else {
					r, t, err := getFieldNameType(tokens[i+1:], field.BelongedToStruct, curIndex, false)
					// The final result of the structural field
					if err != nil {
						names = append(names, field.Tag.Get("bson"))
						types = append(types, field.Type)
						break
					}
					if len(r) != 1 {
						return nil, nil, fmt.Errorf("no field name corresponding to %v found", tokens[i:])
					}
					i += *curIndex
					names = append(names, field.Tag.Get("bson")+"."+r[0])
					types = append(types, t[0])
					break
				}
			}
		}
		if !isFirst && flag == 1 {
			break
		}
		if flag == 0 {
			return nil, nil, fmt.Errorf("no field name corresponding to %v found", tokens[i:])
		}
	}

	return
}
