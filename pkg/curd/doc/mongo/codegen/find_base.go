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

func findBaseCodegen() []code.Statement {
	stmt := `if selector == nil {
		return fmt.Errorf("query param is empty")
	}

	return b.collection.FindOne(ctx, selector).Decode(result)`

	return []code.Statement{
		code.RawStmt(stmt),
	}
}

func findListBaseCodegen() []code.Statement {
	stmt := `if selector == nil {
		return fmt.Errorf("query param is empty")
	}

	cursor, err := b.collection.Find(ctx, selector)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	return cursor.All(ctx, result)`

	return []code.Statement{
		code.RawStmt(stmt),
	}
}

func findPageListBaseCodegen() []code.Statement {
	stmt := `if selector == nil {
		return fmt.Errorf("query param is empty")
	}
	
	if skip < 0 || limit < 0 {
		return fmt.Errorf("skip or limit not correct")
	}	

	ops := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit))
	if sort != "" {
		if strings.Contains(sort, "-") {
			sort = strings.TrimPrefix(sort, "-")
			ops.SetSort(bson.D{bson.E{Key: sort, Value: -1}})
		} else {
			ops.SetSort(bson.D{bson.E{Key: sort, Value: 1}})
		}
	}

	cursor, err := b.collection.Find(ctx, selector, ops)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	return cursor.All(ctx, result)`

	return []code.Statement{
		code.RawStmt(stmt),
	}
}

func findSortPageListBaseCodegen() []code.Statement {
	stmt := `if selector == nil {
		return fmt.Errorf("query param is empty")
	}
	
	if skip < 0 || limit < 0 {
		return fmt.Errorf("skip or limit not correct")
	}

	ops := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit))
	if len(sorts) > 0 {
		sortFields := make(bson.D, len(sorts))
		for i, field := range sorts {
			if strings.Contains(field, "-") {
				field = strings.TrimPrefix(field, "-")
				sortFields[i] = bson.E{Key: field, Value: -1}
			} else {
				sortFields[i] = bson.E{Key: field, Value: 1}
			}
		}
		ops.SetSort(sortFields)
	}

	cursor, err := b.collection.Find(ctx, selector, ops)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	return cursor.All(ctx, result)`

	return []code.Statement{
		code.RawStmt(stmt),
	}
}
