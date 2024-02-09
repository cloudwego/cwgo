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
	"fmt"

	"github.com/cloudwego/cwgo/pkg/curd/code"
	"github.com/cloudwego/cwgo/pkg/curd/parse"
)

func findCodegen(find *parse.FindParse) []code.Statement {
	if find.OperateMode == parse.OperateOne {
		return []code.Statement{
			code.DeclVarStmt{
				Name: "entity",
				Type: find.ReturnType,
			},
			code.IfBlockStmt{
				Condition: []code.Statement{
					code.RawStmt("err := "),
					code.CallStmt{
						Caller: code.CallStmt{
							Caller:   code.RawStmt("r.collection"),
							CallName: "FindOne",
							Args: code.ListCommaStmt{
								code.RawStmt(find.CtxParamName),
								queryCodegen(find.Query),
								findOptionsCodegen(find),
							},
						},
						CallName: "Decode",
						Args: code.ListCommaStmt{
							code.RawStmt("entity"),
						},
					},
					code.RawStmt("; err != nil "),
				},
				Body: code.Body{
					code.RawStmt("return nil, err"),
				},
			},
			code.ReturnStmt{
				ListCommaStmt: code.ListCommaStmt{
					code.RawStmt("entity"),
					code.RawStmt("nil"),
				},
			},
		}
	} else {
		baseFindStmt := []code.Statement{
			code.DeclColonStmt{
				Left: code.ListCommaStmt{
					code.RawStmt("cursor"),
					code.RawStmt("err"),
				},
				Right: code.CallStmt{
					Caller:   code.RawStmt("r.collection"),
					CallName: "Find",
					Args: code.ListCommaStmt{
						code.RawStmt(find.CtxParamName),
						queryCodegen(find.Query),
						findOptionsCodegen(find),
					},
				},
			},
			code.RawStmt("if err != nil {\n\treturn nil, err\n}"),
			code.DeclVarStmt{
				Name: "entities",
				Type: find.ReturnType,
			},
			code.IfBlockStmt{
				Condition: []code.Statement{
					code.RawStmt("err = "),
					code.CallStmt{
						Caller:   code.RawStmt("cursor"),
						CallName: "All",
						Args: code.ListCommaStmt{
							code.RawStmt(find.CtxParamName),
							code.RawStmt("&entities"),
						},
					},
					code.RawStmt("; err != nil "),
				},
				Body: code.Body{
					code.RawStmt("return nil, err"),
				},
			},
			code.ReturnStmt{
				ListCommaStmt: code.ListCommaStmt{
					code.RawStmt("entities"),
					code.RawStmt("nil"),
				},
			},
		}

		pageRevealStmt := pageRevealCodegen(find)
		if pageRevealStmt != nil {
			findStmt := append([]code.Statement{pageRevealStmt}, baseFindStmt...)
			return findStmt
		} else {
			return baseFindStmt
		}
	}
}

func findOptionsCodegen(find *parse.FindParse) code.Statement {
	chainCall := make(code.ChainStmt, 0, 5)

	if find.OperateMode == parse.OperateOne {
		baseChain := chainCall.ChainCall(code.Chain{
			CallName: "options.FindOne",
			Args:     code.ListCommaStmt{},
		}).ChainCall(code.Chain{
			CallName: "SetSort",
			Args:     code.ListCommaStmt{findOrderCodegen(find.Order)},
		})

		if len(find.Project) != 0 {
			baseChain = baseChain.ChainCall(code.Chain{
				CallName: "SetProjection",
				Args: code.ListCommaStmt{
					findProjectCodegen(find),
				},
			})
		}

		if find.SkipParamName != "" {
			baseChain = baseChain.ChainCall(code.Chain{
				CallName: "SetSkip",
				Args:     code.ListCommaStmt{code.RawStmt(find.SkipParamName)},
			})
		}

		return baseChain
	} else {
		baseChain := chainCall.ChainCall(code.Chain{
			CallName: "options.Find",
			Args:     code.ListCommaStmt{},
		}).ChainCall(code.Chain{
			CallName: "SetSort",
			Args:     code.ListCommaStmt{findOrderCodegen(find.Order)},
		})

		if len(find.Project) != 0 {
			baseChain = baseChain.ChainCall(code.Chain{
				CallName: "SetProjection",
				Args: code.ListCommaStmt{
					findProjectCodegen(find),
				},
			})
		}

		if find.LimitParamName != "" {
			baseChain = baseChain.ChainCall(code.Chain{
				CallName: "SetLimit",
				Args:     code.ListCommaStmt{code.RawStmt(find.LimitParamName)},
			})
		}

		if find.SkipParamName != "" {
			baseChain = baseChain.ChainCall(code.Chain{
				CallName: "SetSkip",
				Args:     code.ListCommaStmt{code.RawStmt(find.SkipParamName)},
			})
		}

		return baseChain
	}
}

func findOrderCodegen(order parse.Order) code.MapStmt {
	mapPairs := make([]code.MapPair, 0, 10)

	for _, field := range order.Asc {
		mapPairs = append(mapPairs, code.MapPair{
			Key:   code.RawStmt(field),
			Value: code.RawStmt("1"),
		})
	}

	for _, field := range order.Desc {
		mapPairs = append(mapPairs, code.MapPair{
			Key:   code.RawStmt(field),
			Value: code.RawStmt("-1"),
		})
	}

	return code.MapStmt{
		Name: "bson.M",
		Pair: mapPairs,
	}
}

func findProjectCodegen(find *parse.FindParse) code.MapStmt {
	mapPairs := make([]code.MapPair, 0, 10)

	for _, field := range find.Project {
		mapPairs = append(mapPairs, code.MapPair{
			Key:   code.RawStmt(field),
			Value: code.RawStmt("1"),
		})
	}

	return code.MapStmt{
		Name: "bson.M",
		Pair: mapPairs,
	}
}

func pageRevealCodegen(find *parse.FindParse) code.Statement {
	if find.LimitParamName != "" {
		return code.IfBlockStmt{
			Condition: []code.Statement{
				code.RawStmt(fmt.Sprintf("%s == 0 ", find.LimitParamName)),
			},
			Body: code.Body{
				code.RawStmt(fmt.Sprintf("%s = 5", find.LimitParamName)),
			},
		}
	} else {
		return nil
	}
}
