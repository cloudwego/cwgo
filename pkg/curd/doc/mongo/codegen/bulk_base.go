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

import "github.com/cloudwego/cwgo/pkg/curd/code"

func bulkInsertBaseCodegen() []code.Statement {
	stmt := `if len(batchData) == 0 {
		return nil, fmt.Errorf("batch param is empty")
	}

	models := make([]mongo.WriteModel, len(batchData))
	for i, doc := range batchData {
		models[i] = mongo.NewInsertOneModel().SetDocument(doc)
	}

	return b.collection.BulkWrite(ctx, models)`

	return []code.Statement{
		code.RawStmt(stmt),
	}
}

func bulkUpdateBaseCodegen() []code.Statement {
	stmt := `if len(filter) != len(updater) {
		return nil, fmt.Errorf("filter and updater must have the same length")
	}

	models := make([]mongo.WriteModel, len(filter))
	for i := range filter {
		models[i] = mongo.NewUpdateOneModel().SetFilter(filter[i]).SetUpdate(updater[i])
	}

	return b.collection.BulkWrite(ctx, models)`

	return []code.Statement{
		code.RawStmt(stmt),
	}
}
