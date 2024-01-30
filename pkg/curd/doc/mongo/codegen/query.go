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

package codegen

import (
	"github.com/cloudwego/cwgo/pkg/curd/code"
	"github.com/cloudwego/cwgo/pkg/curd/parse"
)

func queryCodegen(query *parse.Query) code.Statement {
	if query.QueryMode == parse.All {
		return code.MapStmt{
			Name: "bson.M",
			Pair: []code.MapPair{},
		}
	} else {
		return code.MapStmt{
			Name: "bson.M",
			Pair: []code.MapPair{
				dfsCodegen(query.ConnectionOpTree),
			},
		}
	}
}

func dfsCodegen(node *parse.ConnectionOpTree) code.MapPair {
	// leafs node
	if node.LeftChildren == nil {
		return comparatorCodegen(node)
	} else {
		// none-leafs node
		return code.MapPair{
			Key: code.RawStmt(node.Name),
			Value: code.SliceStmt{
				Name: "[]bson.M",
				Values: []code.MapPair{
					dfsCodegen(node.LeftChildren),
					dfsCodegen(node.RightChildren),
				},
			},
		}
	}
}

func comparatorCodegen(node *parse.ConnectionOpTree) code.MapPair {
	switch parse.QueryComparator(node.Name) {
	case parse.Equal:
		return singleMapCodegen(node.MongoFieldName, node.ParamNames[0])
	case parse.NotEqual:
		return oneMapParamCodegen(node.MongoFieldName, "$ne", node.ParamNames[0])
	case parse.LessThan:
		return oneMapParamCodegen(node.MongoFieldName, "$lt", node.ParamNames[0])
	case parse.LessThanEqual:
		return oneMapParamCodegen(node.MongoFieldName, "$lte", node.ParamNames[0])
	case parse.GreaterThan:
		return oneMapParamCodegen(node.MongoFieldName, "$gt", node.ParamNames[0])
	case parse.GreaterThanEqual:
		return oneMapParamCodegen(node.MongoFieldName, "$gte", node.ParamNames[0])
	case parse.Between:
		return twoMapParamsCodegen(node.MongoFieldName, "$gte", node.ParamNames[0],
			"$lte", node.ParamNames[1])
	case parse.NotBetween:
		return twoMapParamsCodegen(node.MongoFieldName, "$lt", node.ParamNames[0],
			"$gt", node.ParamNames[1])
	case parse.In:
		return oneMapParamCodegen(node.MongoFieldName, "$in", node.ParamNames[0])
	case parse.NotIn:
		return oneMapParamCodegen(node.MongoFieldName, "$n"+"in", node.ParamNames[0])
	case parse.True:
		return singleMapCodegen(node.MongoFieldName, "true")
	case parse.False:
		return singleMapCodegen(node.MongoFieldName, "false")
	case parse.Exists:
		return oneMapParamCodegen(node.MongoFieldName, "$exists", "1")
	case parse.NotExists:
		return oneMapParamCodegen(node.MongoFieldName, "$exists", "0")
	default:
	}

	return code.MapPair{}
}

func singleMapCodegen(key, value string) code.MapPair {
	return code.MapPair{
		Key:   code.RawStmt(key),
		Value: code.RawStmt(value),
	}
}

func oneMapParamCodegen(key, valueKey, valueValue string) code.MapPair {
	return code.MapPair{
		Key: code.RawStmt(key),
		Value: code.MapStmt{
			Name: "bson.M",
			Pair: []code.MapPair{
				{
					Key:   code.RawStmt(valueKey),
					Value: code.RawStmt(valueValue),
				},
			},
		},
	}
}

func twoMapParamsCodegen(key, valueKey1, valueValue1, valueKey2, valueValue2 string) code.MapPair {
	return code.MapPair{
		Key: code.RawStmt(key),
		Value: code.MapStmt{
			Name: "bson.M",
			Pair: []code.MapPair{
				{
					Key:   code.RawStmt(valueKey1),
					Value: code.RawStmt(valueValue1),
				},
				{
					Key:   code.RawStmt(valueKey2),
					Value: code.RawStmt(valueValue2),
				},
			},
		},
	}
}
