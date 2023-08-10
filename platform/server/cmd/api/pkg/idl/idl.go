/*
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
 */

package idl

import (
	"github.com/cloudwego/cwgo/platform/server/cmd/api/pkg/idl/gitlab"
	"net/http"
)

type AddIDLReq struct {
	RepositoryId int64
	MainIdlPath  string
	ServiceName  string
}

type AddIDLRes struct {
	Code int32
	Msg  string
}

func AddIDL(req AddIDLReq) (AddIDLRes, error) {
	var err error
	var res AddIDLRes
	repoType := 1 //查数据库
	switch repoType {
	case 1:
		err = gitlab.AddIDL(req.RepositoryId, req.MainIdlPath)
	case 2:

	}
	if err != nil {
		res.Code = http.StatusBadRequest
		res.Msg = "Internal error"
		return res, err
	}

	return res, nil
}
