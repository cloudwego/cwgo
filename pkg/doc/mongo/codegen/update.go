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

func updateCodegen(update *parse.UpdateParse) []code.Statement {
	chainCall := make(code.ChainStmt, 0, 5)
	if update.OperateMode == parse.OperateOne {
		return []code.Statement{
			code.DeclColonStmt{
				Left: code.ListCommaStmt{
					code.RawStmt("result"),
					code.RawStmt("err"),
				},
				Right: code.CallStmt{
					Caller:   code.RawStmt("r.collection"),
					CallName: "UpdateOne",
					Args: code.ListCommaStmt{
						code.RawStmt(update.CtxParamName),
						queryCodegen(update.Query),
						updateFieldsCodegen(update),
						chainCall.ChainCall(code.Chain{
							CallName: "options.Update",
							Args:     code.ListCommaStmt{},
						}).ChainCall(code.Chain{
							CallName: "SetUpsert",
							Args: code.ListCommaStmt{
								upsertCodegen(update.Upsert),
							},
						}),
					},
				},
			},
			code.RawStmt("if err != nil {\n\treturn false, err\n}"),
			code.ReturnStmt{
				ListCommaStmt: code.ListCommaStmt{
					code.RawStmt("result.MatchedCount > 0"),
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
					CallName: "UpdateMany",
					Args: code.ListCommaStmt{
						code.RawStmt(update.CtxParamName),
						queryCodegen(update.Query),
						updateFieldsCodegen(update),
						chainCall.ChainCall(code.Chain{
							CallName: "options.Update",
							Args:     code.ListCommaStmt{},
						}).ChainCall(code.Chain{
							CallName: "SetUpsert",
							Args: code.ListCommaStmt{
								upsertCodegen(update.Upsert),
							},
						}),
					},
				},
			},
			code.RawStmt("if err != nil {\n\treturn 0, err\n}"),
			code.ReturnStmt{
				ListCommaStmt: code.ListCommaStmt{
					code.RawStmt("int(result.MatchedCount)"),
					code.RawStmt("nil"),
				},
			},
		}
	}
}

func updateFieldsCodegen(update *parse.UpdateParse) code.MapStmt {
	if update.UpdateStructObjName == "" {
		mapPairs := make([]code.MapPair, 0, 5)
		for _, field := range update.UpdateFields {
			mapPairs = append(mapPairs, code.MapPair{
				Key:   code.RawStmt(field.MongoFieldName),
				Value: code.RawStmt(field.ParamName),
			})
		}
		return code.MapStmt{
			Name: "bson.M",
			Pair: []code.MapPair{
				{
					Key: code.RawStmt("$set"),
					Value: code.MapStmt{
						Name: "bson.M",
						Pair: mapPairs,
					},
				},
			},
		}
	} else {
		return code.MapStmt{
			Name: "bson.M",
			Pair: []code.MapPair{
				{
					Key:   code.RawStmt("$set"),
					Value: code.RawStmt(update.UpdateStructObjName),
				},
			},
		}
	}
}

func upsertCodegen(upsert bool) code.RawStmt {
	if upsert {
		return "true"
	}
	return "false"
}
