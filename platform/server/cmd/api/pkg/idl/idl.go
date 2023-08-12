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
	"github.com/cloudwego/cwgo/platform/server/shared/config"
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

type DeleteIDLsReq struct {
	Ids []int64
}
type DeleteIDLsRes struct {
	Code int32
	Msg  string
}

type UpdateIDLReq struct {
	Id           int64
	RepositoryId int64
	MainIdlPath  string
	ServiceName  string
}
type UpdateIDLRes struct {
	Code int32
	Msg  string
}

type GetIDLsReq struct {
	Page  int32
	Limit int32
}
type GetIDLsRes struct {
	Code int32
	Msg  string
	IDLs []config.IDL
}

type SyncIDLsReq struct {
	Ids []int64
}
type SyncIDLsRes struct {
	Code int32
	Msg  string
}

func AddIDL(req AddIDLReq) (AddIDLRes, error) {
	var err error
	var res AddIDLRes

	//TODO：查数据库
	repoType := 1

	switch repoType {
	case 1:
		err = gitlab.AddIDL(req.RepositoryId, req.MainIdlPath, req.ServiceName)
	case 2:

	}
	if err != nil {
		res.Code = http.StatusBadRequest
		res.Msg = "Internal error"
		return res, err
	}

	return res, nil
}

func DeleteIDLs(req DeleteIDLsReq) (DeleteIDLsRes, error) {
	var err error
	var res DeleteIDLsRes
	repoType := 1 //查数据库
	switch repoType {
	case 1:
		err = gitlab.DeleteIDLs(req.Ids)
	case 2:

	}
	if err != nil {
		res.Code = http.StatusBadRequest
		res.Msg = "Internal error"
		return res, err
	}

	return res, nil
}

func UpdateIDL(req UpdateIDLReq) (UpdateIDLRes, error) {
	var err error
	var res UpdateIDLRes
	repoType := 1 //查数据库
	switch repoType {
	case 1:
		err = gitlab.UpdateIDL(req.Id, req.RepositoryId, req.MainIdlPath, req.ServiceName)
	case 2:

	}
	if err != nil {
		res.Code = http.StatusBadRequest
		res.Msg = "Internal error"
		return res, err
	}

	return res, nil
}

func GetIDLs(req GetIDLsReq) (GetIDLsRes, error) {
	var err error
	var res GetIDLsRes
	repoType := 1 //查数据库
	switch repoType {
	case 1:
		err = gitlab.GetIDLs(req.Limit, req.Page)
	case 2:

	}
	if err != nil {
		res.Code = http.StatusBadRequest
		res.Msg = "Internal error"
		return res, err
	}

	return res, nil
}

func SyncIDLs(req SyncIDLsReq) (SyncIDLsRes, error) {
	var err error
	var res SyncIDLsRes
	repoType := 1 //查数据库
	switch repoType {
	case 1:
		err = gitlab.SyncIDLs(req.Ids)
	case 2:

	}
	if err != nil {
		res.Code = http.StatusBadRequest
		res.Msg = "Internal error"
		return res, err
	}

	return res, nil
}
