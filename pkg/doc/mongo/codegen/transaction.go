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

func taCodegen(transaction *parse.TransactionParse) []code.Statement {
	taOperations := taOperationsCodegen(transaction)
	body := code.Body{
		code.RawStmt("if err := sessionContext.StartTransaction(); err != nil {\n\treturn err\n}\n"),
	}
	for _, taOperation := range taOperations {
		body = append(body, taOperation)
		body = append(body, code.RawStmt("\n"))
	}
	body = append(body, code.RawStmt("return sessionContext.CommitTransaction(context.Background())"))

	return []code.Statement{
		code.IfBlockStmt{
			Condition: []code.Statement{
				code.DeclColonStmt{
					Left: code.ListCommaStmt{
						code.RawStmt("err"),
					},
					Right: code.CallStmt{
						Caller:   code.RawStmt(transaction.ClientParamName),
						CallName: "UseSession",
						Args: code.ListCommaStmt{
							code.RawStmt(transaction.CtxParamName),
							code.AnonymousFuncStmt{
								Params: code.Params{
									code.Param{
										Name: "sessionContext",
										Type: code.SelectorExprType{
											X:   "mongo",
											Sel: "SessionContext",
										},
									},
								},
								Returns: code.Returns{
									code.IdentType("error"),
								},
								Body: body,
							},
						},
					},
				},
				code.RawStmt("; err != nil "),
			},
			Body: code.Body{
				code.RawStmt("return err"),
			},
		},
		code.RawStmt("return nil"),
	}
}

func taOperationsCodegen(transaction *parse.TransactionParse) []code.Statement {
	operations := make([]code.Statement, 0, 10)
	for _, operation := range transaction.TransactionOperations {
		if operation.Operation.GetOperationName() == parse.Insert {
			operations = append(operations, taInsertCodegen(operation)...)
		}
		if operation.Operation.GetOperationName() == parse.Update {
			operations = append(operations, taUpdateCodegen(operation))
		}
		if operation.Operation.GetOperationName() == parse.Delete {
			operations = append(operations, taDeleteCodegen(operation))
		}
		if operation.Operation.GetOperationName() == parse.Bulk {
			operations = append(operations, taBulkCodegen(operation)...)
		}
	}
	return operations
}

func taInsertCodegen(tsOperation parse.TransactionOperation) []code.Statement {
	insert := tsOperation.Operation.(*parse.InsertParse)
	if insert.OperateMode == parse.OperateOne {
		return getInsertCode(tsOperation, insert, "InsertOne", insert.MethodParamNames[0])
	} else {
		return getInsertCode(tsOperation, insert, "InsertMany", "entities")
	}
}

func getInsertCode(tsOperation parse.TransactionOperation, insert *parse.InsertParse, callName, param string) []code.Statement {
	baseInsertCode := code.IfBlockStmt{
		Condition: []code.Statement{
			code.DeclColonStmt{
				Left: code.ListCommaStmt{
					code.RawStmt("_"),
					code.RawStmt("err"),
				},
				Right: code.CallStmt{
					Caller:   code.RawStmt(tsOperation.CollectionParamName),
					CallName: callName,
					Args: code.ListCommaStmt{
						code.RawStmt("sessionContext"),
						code.RawStmt(param),
					},
				},
			},
			code.RawStmt("; err != nil "),
		},
		Body: code.Body{
			code.RawStmt(abortTa),
		},
	}

	if insert.OperateMode == parse.OperateOne {
		return []code.Statement{baseInsertCode}
	} else {
		return []code.Statement{
			code.DeclVarStmt{
				Name: "entities",
				Type: code.SliceType{
					ElementType: code.InterfaceType{},
				},
			},
			code.ForRangeBlockStmt{
				RangeName: insert.MethodParamNames[0],
				Value:     "model",
				Body: []code.Statement{
					code.RawStmt("entities = append(entities, model)"),
				},
			},
			baseInsertCode,
		}
	}
}

func taUpdateCodegen(tsOperation parse.TransactionOperation) code.Statement {
	update := tsOperation.Operation.(*parse.UpdateParse)
	if update.OperateMode == parse.OperateOne {
		return getUpdateCode(tsOperation, update, "UpdateOne")
	} else {
		return getUpdateCode(tsOperation, update, "UpdateMany")
	}
}

func getUpdateCode(tsOperation parse.TransactionOperation, update *parse.UpdateParse, callName string) code.Statement {
	chainCall := make(code.ChainStmt, 0, 5)
	return code.IfBlockStmt{
		Condition: []code.Statement{
			code.DeclColonStmt{
				Left: code.ListCommaStmt{
					code.RawStmt("_"),
					code.RawStmt("err"),
				},
				Right: code.CallStmt{
					Caller:   code.RawStmt(tsOperation.CollectionParamName),
					CallName: callName,
					Args: code.ListCommaStmt{
						code.RawStmt("sessionContext"),
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
			code.RawStmt("; err != nil "),
		},
		Body: code.Body{
			code.RawStmt(abortTa),
		},
	}
}

func taDeleteCodegen(tsOperation parse.TransactionOperation) code.Statement {
	del := tsOperation.Operation.(*parse.DeleteParse)
	if del.OperateMode == parse.OperateOne {
		return getTaDeleteCode(tsOperation, del, "DeleteOne")
	} else {
		return getTaDeleteCode(tsOperation, del, "DeleteMany")
	}
}

func getTaDeleteCode(tsOperation parse.TransactionOperation, del *parse.DeleteParse, callName string) code.Statement {
	return code.IfBlockStmt{
		Condition: []code.Statement{
			code.DeclColonStmt{
				Left: code.ListCommaStmt{
					code.RawStmt("_"),
					code.RawStmt("err"),
				},
				Right: code.CallStmt{
					Caller:   code.RawStmt(tsOperation.CollectionParamName),
					CallName: callName,
					Args: code.ListCommaStmt{
						code.RawStmt("sessionContext"),
						queryCodegen(del.Query),
					},
				},
			},
			code.RawStmt("; err != nil "),
		},
		Body: code.Body{
			code.RawStmt(abortTa),
		},
	}
}

func taBulkCodegen(tsOperation parse.TransactionOperation) []code.Statement {
	bulk := tsOperation.Operation.(*parse.BulkParse)

	return []code.Statement{
		code.DeclVarStmt{
			Name: "models",
			Type: code.SliceType{
				ElementType: code.SelectorExprType{
					X:   "mongo",
					Sel: "WriteModel",
				},
			},
		},
		bulkOperationsCodegen(bulk),
		code.IfBlockStmt{
			Condition: []code.Statement{
				code.DeclColonStmt{
					Left: code.ListCommaStmt{
						code.RawStmt("_"),
						code.RawStmt("err"),
					},
					Right: code.CallStmt{
						Caller:   code.RawStmt(tsOperation.CollectionParamName),
						CallName: "BulkWrite",
						Args: code.ListCommaStmt{
							code.RawStmt("sessionContext"),
							code.RawStmt("models"),
						},
					},
				},
				code.RawStmt("; err != nil "),
			},
			Body: code.Body{
				code.RawStmt(abortTa),
			},
		},
	}
}

var abortTa = `if err = sessionContext.AbortTransaction(context.Background()); err != nil {
    return err
}
return err`
