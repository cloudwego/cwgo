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
)

type QueryMode string

const (
	By  = QueryMode("By")
	All = QueryMode("All")
)

type QueryConnectionOp string

const (
	And = QueryConnectionOp("And")
	Or  = QueryConnectionOp("Or")
)

type QueryComparator string

const (
	Equal            = QueryComparator("Equal")
	NotEqual         = QueryComparator("NotEqual")
	LessThan         = QueryComparator("LessThan")
	LessThanEqual    = QueryComparator("LessThanEqual")
	GreaterThan      = QueryComparator("GreaterThan")
	GreaterThanEqual = QueryComparator("GreaterThanEqual")
	Between          = QueryComparator("Between")
	NotBetween       = QueryComparator("NotBetween")
	In               = QueryComparator("In")
	NotIn            = QueryComparator("NotIn")
	True             = QueryComparator("True")
	False            = QueryComparator("False")
	Exists           = QueryComparator("Exists")
	NotExists        = QueryComparator("NotExists")
)

type Query struct {
	// QueryMode By or All
	QueryMode QueryMode

	// ConnectionOpTree stores query information
	ConnectionOpTree *ConnectionOpTree
}

type ConnectionOpTree struct {
	Name           string // if not leaf, store And || Or, else store QueryComparator(ComparatorName)
	LeftChildren   *ConnectionOpTree
	RightChildren  *ConnectionOpTree
	MongoFieldName string   // if not leaf, empty
	ParamNames     []string // if not leaf, empty
}

const (
	leftBracket  = "Lb"
	rightBracket = "Rb"
)

func newQuery() *Query {
	return &Query{}
}

func (q *Query) parseQuery(methodTokens []string, method *extract.InterfaceMethod, curParamIndex *int) error {
	tokens, err := q.checkQuery(methodTokens)
	if err != nil {
		return newMethodSyntaxError(method.Name, err.Error())
	}

	if q.QueryMode == All {
		return nil
	}

	q.ConnectionOpTree, err = q.createTree(tokens, method, curParamIndex)
	if err != nil {
		return err
	}

	return nil
}

func (q *Query) checkQuery(methodTokens []string) ([]string, error) {
	if len(methodTokens) == 0 {
		return nil, errors.New("no By or All specified")
	}

	switch methodTokens[0] {
	case string(By):
		if len(methodTokens) == 1 {
			return nil, errors.New("by needs to be followed by query tokens")
		}
		q.QueryMode = By
		return methodTokens[1:], nil

	case string(All):
		if len(methodTokens) != 1 {
			return nil, errors.New("there's no need to follow any tokens behind All")
		}
		q.QueryMode = All
		return nil, nil

	default:
		return nil, errors.New("no By or All specified")
	}
}

func (q *Query) createTree(tokens []string, method *extract.InterfaceMethod, curParamIndex *int) (*ConnectionOpTree, error) {
	stack := make([]string, 0, 5)
	for index, token := range tokens {
		if token == leftBracket {
			stack = append(stack, token)
		}
		if token == rightBracket {
			if len(stack) == 0 {
				return nil, errors.New("mismatched parentheses")
			}
			stack = stack[0 : len(stack)-1]
		}

		if (token == string(And) || token == string(Or)) && (len(stack) == 1 || len(stack) == 0) {
			if index == len(tokens)-1 {
				return nil, errors.New("and || or needs to be followed by query tokens")
			}

			var leftTokens []string
			if len(stack) == 1 {
				leftTokens = tokens[1:index]
			} else {
				leftTokens = tokens[:index]
			}
			leftNode, err := q.createTree(leftTokens, method, curParamIndex)
			if err != nil {
				return nil, err
			}

			var rightTokens []string
			if len(stack) == 1 {
				rightTokens = tokens[index+1 : len(tokens)-1]
			} else {
				rightTokens = tokens[index+1:]
			}
			rightNode, err := q.createTree(rightTokens, method, curParamIndex)
			if err != nil {
				return nil, err
			}

			node := &ConnectionOpTree{
				Name:          token,
				LeftChildren:  leftNode,
				RightChildren: rightNode,
			}
			return node, nil
		}
	}

	cpName, fieldName, paramNames, err := q.splitConditionPairs(tokens, method, curParamIndex)
	if err != nil {
		return nil, err
	}

	node := &ConnectionOpTree{
		Name:           cpName,
		LeftChildren:   nil,
		RightChildren:  nil,
		MongoFieldName: fieldName,
		ParamNames:     paramNames,
	}
	return node, nil
}

func (q *Query) splitConditionPairs(methodTokens []string, method *extract.InterfaceMethod, curParamIndex *int) (string, string, []string, error) {
	if len(methodTokens) == 0 || len(methodTokens) == 1 {
		return "", "", nil, newMethodSyntaxError(method.Name, fmt.Sprintf("there are grammar errors in %v", methodTokens))
	}

	for i := len(methodTokens) - 1; i >= 0; i-- {
		if i-1 >= 0 && methodTokens[i] == "Equal" && methodTokens[i-1] != "Not" && methodTokens[i-1] != "Than" {
			return q.parseQueryConditionPair(methodTokens[:i], method, curParamIndex, Equal, 1)
		}

		if i-1 >= 0 && methodTokens[i] == "Equal" && methodTokens[i-1] == "Not" {
			fmt.Printf("%v\n", methodTokens[:i-1])
			return q.parseQueryConditionPair(methodTokens[:i-1], method, curParamIndex, NotEqual, 1)
		}

		if i-1 >= 0 && methodTokens[i] == "Than" && methodTokens[i-1] == "Less" {
			return q.parseQueryConditionPair(methodTokens[:i-1], method, curParamIndex, LessThan, 1)
		}

		if i-2 >= 0 && methodTokens[i] == "Equal" &&
			methodTokens[i-1] == "Than" && methodTokens[i-2] == "Less" {
			return q.parseQueryConditionPair(methodTokens[:i-2], method, curParamIndex, LessThanEqual, 1)
		}

		if i-1 >= 0 && methodTokens[i] == "Than" && methodTokens[i-1] == "Greater" {
			return q.parseQueryConditionPair(methodTokens[:i-1], method, curParamIndex, GreaterThan, 1)
		}

		if i-2 >= 0 && methodTokens[i] == "Equal" &&
			methodTokens[i-1] == "Than" && methodTokens[i-2] == "Greater" {
			return q.parseQueryConditionPair(methodTokens[:i-2], method, curParamIndex, GreaterThanEqual, 1)
		}

		if i-1 >= 0 && methodTokens[i] == "Between" && methodTokens[i-1] != "Not" {
			return q.parseQueryConditionPair(methodTokens[:i], method, curParamIndex, Between, 2)
		}

		if i-1 >= 0 && methodTokens[i] == "Between" && methodTokens[i-1] == "Not" {
			return q.parseQueryConditionPair(methodTokens[:i-1], method, curParamIndex, NotBetween, 2)
		}

		if i-1 >= 0 && methodTokens[i] == "In" && methodTokens[i-1] != "Not" {
			return q.parseQueryConditionPair(methodTokens[:i], method, curParamIndex, In, 1)
		}

		if i-1 >= 0 && methodTokens[i] == "In" && methodTokens[i-1] == "Not" {
			return q.parseQueryConditionPair(methodTokens[:i-1], method, curParamIndex, NotIn, 1)
		}

		if methodTokens[i] == "True" {
			return q.parseQueryConditionPair(methodTokens[:i], method, curParamIndex, True, 0)
		}

		if methodTokens[i] == "False" {
			return q.parseQueryConditionPair(methodTokens[:i], method, curParamIndex, False, 0)
		}

		if i-1 >= 0 && methodTokens[i] == "Exists" && methodTokens[i-1] != "Not" {
			return q.parseQueryConditionPair(methodTokens[:i], method, curParamIndex, Exists, 0)
		}

		if i-1 >= 0 && methodTokens[i] == "Exists" && methodTokens[i-1] == "Not" {
			return q.parseQueryConditionPair(methodTokens[:i-1], method, curParamIndex, NotExists, 0)
		}
	}

	return "", "", nil, newMethodSyntaxError(method.Name, fmt.Sprintf("there are grammar errors in %v, "+
		"not including Equal, NotEqual, LessThan, LessThanEqual, GreaterThan, GreaterThanEqual, Between, NotBetween,"+
		"In, NotIn, True, False, Exists, NotExists", methodTokens))
}

// parseQueryConditionPair is used to parse query's condition pair
//
//	return params description:
//	1. string(queryComparator) 2. field name in structure
//	3. input parameter values corresponding to field names
//	4. error
func (q *Query) parseQueryConditionPair(methodTokens []string, method *extract.InterfaceMethod, curParamIndex *int,
	queryComparator QueryComparator, paramCount int,
) (string, string, []string, error) {
	if len(methodTokens) == 0 {
		return "", "", nil, newMethodSyntaxError(method.Name, fmt.Sprintf("there are grammar errors in %v", methodTokens))
	}

	curIndex := new(int)
	*curIndex = -1
	result, t, err := getFieldNameType(methodTokens, method.BelongedToStruct, curIndex, true)
	if err != nil {
		return "", "", nil, err
	}
	if len(result) != 1 {
		return "", "", nil, newMethodSyntaxError(method.Name, "only one field name can be included between And or Or")
	}

	var values []string
	if paramCount > 0 {
		if *curParamIndex+paramCount > len(method.Params) {
			return "", "", nil, newMethodSyntaxError(method.Name, "insufficient number of input parameters")
		}
		for i := *curParamIndex; i < *curParamIndex+paramCount; i++ {
			if method.Params[i].Type.RealName() != t[0].RealName() {
				return "", "", nil, newMethodSyntaxError(method.Name,
					fmt.Sprintf("the field type in the parameter transfer: %s, the actual required field type: %s",
						method.Params[i].Type.RealName(), t[0].RealName()))
			}
			values = append(values, method.Params[i].Name)
		}
		*curParamIndex += paramCount
	}

	return string(queryComparator), result[0], values, nil
}

func getFirstQueryIndex(tokens []string) (int, error) {
	firstIndex := -1
	for index, token := range tokens {
		if token == string(By) || token == string(All) {
			firstIndex = index
			break
		}
	}
	if firstIndex == -1 {
		return 0, errors.New("no By or All specified")
	}
	return firstIndex, nil
}
