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

func deleteCodegen(delete *parse.DeleteParse) []code.Statement {
	if delete.OperateMode == parse.OperateOne {
		return []code.Statement{
			code.DeclColonStmt{
				Left: code.ListCommaStmt{
					code.RawStmt("result"),
					code.RawStmt("err"),
				},
				Right: code.CallStmt{
					Caller:   code.RawStmt("r.collection"),
					CallName: "DeleteOne",
					Args: code.ListCommaStmt{
						code.RawStmt(delete.CtxParamName),
						queryCodegen(delete.Query),
					},
				},
			},
			code.RawStmt("if err != nil {\n\treturn false, err\n}"),
			code.ReturnStmt{
				ListCommaStmt: code.ListCommaStmt{
					code.RawStmt("result.DeletedCount > 0"),
					code.RawStmt("nil"),
				},
			},
		}
	} else {
		return []code.Statement{
			code.DeclColonStmt{
				Left: code.ListCommaStmt{
					code.RawStmt("result"),
					code.RawStmt("err"),
				},
				Right: code.CallStmt{
					Caller:   code.RawStmt("r.collection"),
					CallName: "DeleteMany",
					Args: code.ListCommaStmt{
						code.RawStmt(delete.CtxParamName),
						queryCodegen(delete.Query),
					},
				},
			},
			code.RawStmt("if err != nil {\n\treturn 0, err\n}"),
			code.ReturnStmt{
				ListCommaStmt: code.ListCommaStmt{
					code.RawStmt("int(result.DeletedCount)"),
					code.RawStmt("nil"),
				},
			},
		}
	}
}
