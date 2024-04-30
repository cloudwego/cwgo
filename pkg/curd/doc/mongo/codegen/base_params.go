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

func GetMInsertOneParams() (res []code.Param) {
	ctx := code.Param{
		Name: "ctx",
		Type: code.IdentType("context.Context"),
	}
	insertOneData := code.Param{
		Name: "insertOneData",
		Type: code.InterfaceType{
			Name: "interface{}",
		},
	}
	res = append(res, ctx, insertOneData)

	return
}

func GetMDeleteOneParams() (res []code.Param) {
	ctx := code.Param{
		Name: "ctx",
		Type: code.IdentType("context.Context"),
	}
	deleteOneData := code.Param{
		Name: "deleteOneData",
		Type: code.InterfaceType{
			Name: "interface{}",
		},
	}
	res = append(res, ctx, deleteOneData)

	return
}

func GetMFindOneParams() (res []code.Param) {
	ctx := code.Param{
		Name: "ctx",
		Type: code.IdentType("context.Context"),
	}
	selector := code.Param{
		Name: "selector",
		Type: code.IdentType("bson.M"),
	}
	result := code.Param{
		Name: "result",
		Type: code.InterfaceType{
			Name: "interface{}",
		},
	}
	res = append(res, ctx, selector, result)

	return
}

func GetMBulkInsertParams() (res []code.Param) {
	ctx := code.Param{
		Name: "ctx",
		Type: code.IdentType("context.Context"),
	}
	batchData := code.Param{
		Name: "batchData",
		Type: code.IdentType("[]interface{}"),
	}

	res = append(res, ctx, batchData)

	return
}

func GetMBulkUpdateParams() (res []code.Param) {
	ctx := code.Param{
		Name: "ctx",
		Type: code.IdentType("context.Context"),
	}
	filter := code.Param{
		Name: "filter",
		Type: code.IdentType("[]interface{}"),
	}
	updater := code.Param{
		Name: "updater",
		Type: code.IdentType("[]interface{}"),
	}

	res = append(res, ctx, filter, updater)

	return
}

func GetMAggregateParams() (res []code.Param) {
	ctx := code.Param{
		Name: "ctx",
		Type: code.IdentType("context.Context"),
	}
	pipeline := code.Param{
		Name: "pipeline",
		Type: code.IdentType("[]bson.M"),
	}
	result := code.Param{
		Name: "result",
		Type: code.IdentType("interface{}"),
	}

	res = append(res, ctx, pipeline, result)

	return
}

func GetMCountParams() (res []code.Param) {
	ctx := code.Param{
		Name: "ctx",
		Type: code.IdentType("context.Context"),
	}
	selector := code.Param{
		Name: "selector",
		Type: code.IdentType("bson.M"),
	}

	res = append(res, ctx, selector)

	return
}

func GetMUpdateOneParams() (res []code.Param) {
	ctx := code.Param{
		Name: "ctx",
		Type: code.IdentType("context.Context"),
	}
	selector := code.Param{
		Name: "selector",
		Type: code.IdentType("bson.M"),
	}
	updater := code.Param{
		Name: "updater",
		Type: code.IdentType("bson.M"),
	}

	res = append(res, ctx, selector, updater)

	return
}

func GetMUpdateManyParams() (res []code.Param) {
	ctx := code.Param{
		Name: "ctx",
		Type: code.IdentType("context.Context"),
	}
	selector := code.Param{
		Name: "selector",
		Type: code.IdentType("bson.M"),
	}
	updater := code.Param{
		Name: "updater",
		Type: code.IdentType("bson.M"),
	}

	res = append(res, ctx, selector, updater)

	return
}

func GetMFindListParams() (res []code.Param) {
	ctx := code.Param{
		Name: "ctx",
		Type: code.IdentType("context.Context"),
	}
	selector := code.Param{
		Name: "selector",
		Type: code.IdentType("bson.M"),
	}
	result := code.Param{
		Name: "result",
		Type: code.InterfaceType{
			Name: "interface{}",
		},
	}
	res = append(res, ctx, selector, result)

	return
}

func GetMFindPageListParams() (res []code.Param) {
	ctx := code.Param{
		Name: "ctx",
		Type: code.IdentType("context.Context"),
	}
	sort := code.Param{
		Name: "sort",
		Type: code.IdentType("string"),
	}
	skip := code.Param{
		Name: "skip",
		Type: code.IdentType("int"),
	}
	limit := code.Param{
		Name: "limit",
		Type: code.IdentType("int"),
	}
	selector := code.Param{
		Name: "selector",
		Type: code.IdentType("bson.M"),
	}
	result := code.Param{
		Name: "result",
		Type: code.InterfaceType{
			Name: "interface{}",
		},
	}
	res = append(res, ctx, selector, sort, skip, limit, result)
	return
}

func GetMFindSortPageListParams() (res []code.Param) {
	ctx := code.Param{
		Name: "ctx",
		Type: code.IdentType("context.Context"),
	}
	skip := code.Param{
		Name: "skip",
		Type: code.IdentType("int"),
	}
	limit := code.Param{
		Name: "limit",
		Type: code.IdentType("int"),
	}
	selector := code.Param{
		Name: "selector",
		Type: code.IdentType("bson.M"),
	}
	result := code.Param{
		Name: "result",
		Type: code.InterfaceType{
			Name: "interface{}",
		},
	}
	sort := code.Param{
		Name: "sorts",
		Type: code.IdentType("...string"),
	}

	res = append(res, ctx, selector, skip, limit, result, sort)
	return
}
