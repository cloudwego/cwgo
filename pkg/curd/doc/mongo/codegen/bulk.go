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

func bulkCodegen(bulk *parse.BulkParse) []code.Statement {
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
		code.ReturnStmt{
			ListCommaStmt: code.ListCommaStmt{
				code.CallStmt{
					Caller:   code.RawStmt("r.collection"),
					CallName: "BulkWrite",
					Args: code.ListCommaStmt{
						code.RawStmt(bulk.CtxParamName),
						code.RawStmt("models"),
					},
				},
			},
		},
	}
}

func bulkOperationsCodegen(bulk *parse.BulkParse) code.SliceAppendsStmt {
	operations := make([]code.SliceAppendStmt, 0, 10)
	for _, operation := range bulk.Operations {
		if operation.GetOperationName() == parse.Insert {
			operations = append(operations, bulkInsertCodegen(operation.(*parse.InsertParse)))
		}
		if operation.GetOperationName() == parse.Update {
			operations = append(operations, bulkUpdateCodegen(operation.(*parse.UpdateParse)))
		}
		if operation.GetOperationName() == parse.Delete {
			operations = append(operations, bulkDeleteCodegen(operation.(*parse.DeleteParse)))
		}
	}
	return operations
}

func bulkInsertCodegen(insert *parse.InsertParse) code.SliceAppendStmt {
	return code.SliceAppendStmt{
		SliceName: "models",
		AppendData: code.CallStmt{
			Caller:   code.RawStmt("mongo.NewInsertOneModel()"),
			CallName: "SetDocument",
			Args: code.ListCommaStmt{
				code.RawStmt(insert.MethodParamNames[0]),
			},
		},
	}
}

func bulkUpdateCodegen(update *parse.UpdateParse) code.SliceAppendStmt {
	if update.OperateMode == parse.OperateOne {
		return getBulkUpdateCode(update, "mongo.NewUpdateOneModel().SetFilter")
	} else {
		return getBulkUpdateCode(update, "mongo.NewUpdateManyModel().SetFilter")
	}
}

func bulkDeleteCodegen(delete *parse.DeleteParse) code.SliceAppendStmt {
	if delete.OperateMode == parse.OperateOne {
		return getBulkDeleteCode(delete, "mongo.NewDeleteOneModel().SetFilter")
	} else {
		return getBulkDeleteCode(delete, "mongo.NewDeleteManyModel().SetFilter")
	}
}

func getBulkUpdateCode(update *parse.UpdateParse, callName string) code.SliceAppendStmt {
	chainCall := make(code.ChainStmt, 0, 5)
	return code.SliceAppendStmt{
		SliceName: "models",
		AppendData: chainCall.ChainCall(code.Chain{
			CallName: callName,
			Args: code.ListCommaStmt{
				queryCodegen(update.Query),
			},
		}).ChainCall(code.Chain{
			CallName: "SetUpdate",
			Args: code.ListCommaStmt{
				updateFieldsCodegen(update),
			},
		}).ChainCall(code.Chain{
			CallName: "SetUpsert",
			Args: code.ListCommaStmt{
				upsertCodegen(update.Upsert),
			},
		}),
	}
}

func getBulkDeleteCode(delete *parse.DeleteParse, callName string) code.SliceAppendStmt {
	chainCall := make(code.ChainStmt, 0, 5)
	return code.SliceAppendStmt{
		SliceName: "models",
		AppendData: chainCall.ChainCall(code.Chain{
			CallName: callName,
			Args: code.ListCommaStmt{
				queryCodegen(delete.Query),
			},
		}),
	}
}
