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

func countCodegen(count *parse.CountParse) []code.Statement {
	return []code.Statement{
		code.DeclColonStmt{
			Left: code.ListCommaStmt{
				code.RawStmt("result"),
				code.RawStmt("err"),
			},
			Right: code.CallStmt{
				Caller:   code.RawStmt("r.collection"),
				CallName: "CountDocuments",
				Args: code.ListCommaStmt{
					code.RawStmt(count.CtxParamName),
					queryCodegen(count.Query),
				},
			},
		},
		code.RawStmt("if err != nil {\n\treturn 0, err\n}"),
		code.ReturnStmt{
			ListCommaStmt: code.ListCommaStmt{
				code.RawStmt("int(result)"),
				code.RawStmt("nil"),
			},
		},
	}
}
