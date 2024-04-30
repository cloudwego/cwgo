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

package code

func GetMInsertOneParams() (res []Param) {
	ctx := Param{
		Name: "ctx",
		Type: IdentType("context.Context"),
	}
	insertOneData := Param{
		Name: "insertOneData",
		Type: InterfaceType{
			Name: "interface{}",
		},
	}
	res = append(res, ctx, insertOneData)

	return
}

func GetMDeleteOneParams() (res []Param) {
	ctx := Param{
		Name: "ctx",
		Type: IdentType("context.Context"),
	}
	deleteOneData := Param{
		Name: "deleteOneData",
		Type: InterfaceType{
			Name: "interface{}",
		},
	}
	res = append(res, ctx, deleteOneData)

	return
}

func GetMFindOneParams() (res []Param) {
	ctx := Param{
		Name: "ctx",
		Type: IdentType("context.Context"),
	}
	selector := Param{
		Name: "selector",
		Type: IdentType("bson.M"),
	}
	result := Param{
		Name: "result",
		Type: InterfaceType{
			Name: "interface{}",
		},
	}
	res = append(res, ctx, selector, result)

	return
}

func GetMBulkInsertParams() (res []Param) {
	ctx := Param{
		Name: "ctx",
		Type: IdentType("context.Context"),
	}
	batchData := Param{
		Name: "batchData",
		Type: IdentType("[]interface{}"),
	}

	res = append(res, ctx, batchData)

	return
}

func GetMBulkUpdateParams() (res []Param) {
	ctx := Param{
		Name: "ctx",
		Type: IdentType("context.Context"),
	}
	filter := Param{
		Name: "filter",
		Type: IdentType("[]interface{}"),
	}
	updater := Param{
		Name: "updater",
		Type: IdentType("[]interface{}"),
	}

	res = append(res, ctx, filter, updater)

	return
}

func GetMAggregateParams() (res []Param) {
	ctx := Param{
		Name: "ctx",
		Type: IdentType("context.Context"),
	}
	pipeline := Param{
		Name: "pipeline",
		Type: IdentType("[]bson.M"),
	}
	result := Param{
		Name: "result",
		Type: IdentType("interface{}"),
	}

	res = append(res, ctx, pipeline, result)

	return
}

func GetMCountParams() (res []Param) {
	ctx := Param{
		Name: "ctx",
		Type: IdentType("context.Context"),
	}
	selector := Param{
		Name: "selector",
		Type: IdentType("bson.M"),
	}

	res = append(res, ctx, selector)

	return
}

func GetMUpdateOneParams() (res []Param) {
	ctx := Param{
		Name: "ctx",
		Type: IdentType("context.Context"),
	}
	selector := Param{
		Name: "selector",
		Type: IdentType("bson.M"),
	}
	updater := Param{
		Name: "updater",
		Type: IdentType("bson.M"),
	}

	res = append(res, ctx, selector, updater)

	return
}

func GetMUpdateManyParams() (res []Param) {
	ctx := Param{
		Name: "ctx",
		Type: IdentType("context.Context"),
	}
	selector := Param{
		Name: "selector",
		Type: IdentType("bson.M"),
	}
	updater := Param{
		Name: "updater",
		Type: IdentType("bson.M"),
	}

	res = append(res, ctx, selector, updater)

	return
}

func GetMFindListParams() (res []Param) {
	ctx := Param{
		Name: "ctx",
		Type: IdentType("context.Context"),
	}
	selector := Param{
		Name: "selector",
		Type: IdentType("bson.M"),
	}
	result := Param{
		Name: "result",
		Type: InterfaceType{
			Name: "interface{}",
		},
	}
	res = append(res, ctx, selector, result)

	return
}

func GetMFindPageListParams() (res []Param) {
	ctx := Param{
		Name: "ctx",
		Type: IdentType("context.Context"),
	}
	sort := Param{
		Name: "sort",
		Type: IdentType("string"),
	}
	skip := Param{
		Name: "skip",
		Type: IdentType("int"),
	}
	limit := Param{
		Name: "limit",
		Type: IdentType("int"),
	}
	selector := Param{
		Name: "selector",
		Type: IdentType("bson.M"),
	}
	result := Param{
		Name: "result",
		Type: InterfaceType{
			Name: "interface{}",
		},
	}
	res = append(res, ctx, selector, sort, skip, limit, result)
	return
}

func GetMFindSortPageListParams() (res []Param) {
	ctx := Param{
		Name: "ctx",
		Type: IdentType("context.Context"),
	}
	skip := Param{
		Name: "skip",
		Type: IdentType("int"),
	}
	limit := Param{
		Name: "limit",
		Type: IdentType("int"),
	}
	selector := Param{
		Name: "selector",
		Type: IdentType("bson.M"),
	}
	result := Param{
		Name: "result",
		Type: InterfaceType{
			Name: "interface{}",
		},
	}
	sort := Param{
		Name: "sorts",
		Type: IdentType("...string"),
	}

	res = append(res, ctx, selector, skip, limit, result, sort)
	return
}
