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
	"github.com/cloudwego/cwgo/pkg/doc/mongo/code"
	"github.com/cloudwego/cwgo/pkg/doc/mongo/parse"
)

func insertCodegen(insert *parse.InsertParse) []code.Statement {
	if insert.OperateMode == parse.OperateOne {
		return []code.Statement{
			code.DeclColonStmt{
				Left: code.ListCommaStmt{
					code.RawStmt("result"),
					code.RawStmt("err"),
				},
				Right: code.CallStmt{
					Caller:   code.RawStmt("r.collection"),
					CallName: "InsertOne",
					Args: code.ListCommaStmt{
						code.RawStmt(insert.MethodParamNames[0]),
						code.RawStmt(insert.MethodParamNames[1]),
					},
				},
			},
			code.RawStmt("if err != nil {\n\treturn nil, err\n}"),
			code.ReturnStmt{
				ListCommaStmt: code.ListCommaStmt{
					code.RawStmt("result.InsertedID"),
					code.RawStmt("nil"),
				},
			},
		}
	} else {
		return []code.Statement{
			code.DeclVarStmt{
				Name: "entities",
				Type: code.SliceType{
					ElementType: code.InterfaceType{},
				},
			},
			code.ForRangeBlockStmt{
				RangeName: insert.MethodParamNames[1],
				Value:     "model",
				Body: []code.Statement{
					code.RawStmt("entities = append(entities, model)"),
				},
			},
			code.DeclColonStmt{
				Left: code.ListCommaStmt{
					code.RawStmt("result"),
					code.RawStmt("err"),
				},
				Right: code.CallStmt{
					Caller:   code.RawStmt("r.collection"),
					CallName: "InsertMany",
					Args: code.ListCommaStmt{
						code.RawStmt(insert.MethodParamNames[0]),
						code.RawStmt("entities"),
					},
				},
			},
			code.RawStmt("if err != nil {\n\treturn nil, err\n}"),
			code.ReturnStmt{
				ListCommaStmt: code.ListCommaStmt{
					code.RawStmt("result.InsertedIDs"),
					code.RawStmt("nil"),
				},
			},
		}
	}
}
