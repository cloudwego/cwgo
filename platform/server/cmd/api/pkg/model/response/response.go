/*
 *
 *  * Copyright 2022 CloudWeGo Authors
 *  *
 *  * Licensed under the Apache License, Version 2.0 (the "License");
 *  * you may not use this file except in compliance with the License.
 *  * You may obtain a copy of the License at
 *  *
 *  *     http://www.apache.org/licenses/LICENSE-2.0
 *  *
 *  * Unless required by applicable law or agreed to in writing, software
 *  * distributed under the License is distributed on an "AS IS" BASIS,
 *  * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  * See the License for the specific language governing permissions and
 *  * limitations under the License.
 *
 */

package response

import (
	"github.com/cloudwego/hertz/pkg/app"
	hertzconsts "github.com/cloudwego/hertz/pkg/protocol/consts"
)

type Response struct {
	Code int32  `json:"code"`
	Msg  string `json:"msg"`
}

type ResponseWithData struct {
	Response
	Data interface{} `json:"data"`
}

func Ok(rCtx *app.RequestContext, msg string) {
	rCtx.JSON(
		hertzconsts.StatusOK,
		Response{
			Code: 0,
			Msg:  msg,
		},
	)
}

func OkWithData(rCtx *app.RequestContext, msg string, data interface{}) {
	rCtx.JSON(
		hertzconsts.StatusOK,
		ResponseWithData{
			Response: Response{
				Code: 0,
				Msg:  msg,
			},
			Data: data,
		},
	)
}

func Fail(rCtx *app.RequestContext, statusCode int, bizCode int32, msg string) {
	rCtx.JSON(
		statusCode,
		Response{
			Code: bizCode,
			Msg:  msg,
		},
	)
}
