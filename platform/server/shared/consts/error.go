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

package consts

import (
	"github.com/cloudwego/cwgo/platform/server/shared/errx"
)

// params err
const (
	ErrNumParam = iota + 10000
	ErrNumParamOrderNum
	ErrNumParamOrderBy
	ErrNumParamRepositoryType
	ErrNumParamStoreType
	ErrNumParamRepositoryUrl
	ErrNumParamDomain
	ErrNumParamMainIdlPath
	ErrNumParamTemplateType
	ErrNumParamUrl
	ErrNumParamRepositoryStatus
)

const (
	ErrMsgParam                 = "param err"
	ErrMsgParamOrderNum         = "invalid order num"
	ErrMsgParamOrderBy          = "invalid order by"
	ErrMsgParamRepositoryType   = "invalid repository type"
	ErrMsgParamStoreType        = "invalid store type"
	ErrMsgParamRepositoryUrl    = "invalid repository url"
	ErrMsgParamDomain           = "invalid domain"
	ErrMsgParamMainIdlPath      = "invalid main idl path"
	ErrMsgParamTemplateType     = "invalid template type"
	ErrMsgParamUrl              = "invalid url"
	ErrMsgParamRepositoryStatus = "invalid repository status"
)

var (
	ErrParam                 = errx.New(ErrNumParam, ErrMsgParam)
	ErrParamOrderNum         = errx.New(ErrNumParamOrderNum, ErrMsgParamOrderNum)
	ErrParamOrderBy          = errx.New(ErrNumParamOrderBy, ErrMsgParamOrderBy)
	ErrParamRepositoryType   = errx.New(ErrNumParamRepositoryType, ErrMsgParamRepositoryType)
	ErrParamRepositoryUrl    = errx.New(ErrNumParamStoreType, ErrMsgParamRepositoryUrl)
	ErrParamDomain           = errx.New(ErrNumParamDomain, ErrMsgParamDomain)
	ErrParamMainIdlPath      = errx.New(ErrNumParamMainIdlPath, ErrMsgParamMainIdlPath)
	ErrParamTemplateType     = errx.New(ErrNumParamTemplateType, ErrMsgParamTemplateType)
	ErrParamUrl              = errx.New(ErrNumParamUrl, ErrMsgParamUrl)
	ErrParamRepositoryStatus = errx.New(ErrNumParamRepositoryStatus, ErrMsgParamRepositoryStatus)
)

// rpc err
const (
	ErrNumRpc = iota + 20000
	ErrNumRpcGetClient
	ErrNumRpcConnectClient
)

const (
	ErrMsgRpc              = "rpc err"
	ErrMsgRpcGetClient     = "get rpc client failed"
	ErrMsgRpcConnectClient = "connect to rpc client failed"
)

var (
	ErrRpc              = errx.New(ErrNumRpc, ErrMsgRpc)
	ErrRpcGetClient     = errx.New(ErrNumRpcGetClient, ErrMsgRpcGetClient)
	ErrRpcConnectClient = errx.New(ErrNumRpcConnectClient, ErrMsgRpcConnectClient)
)

// database err

const (
	ErrNumDatabase = iota + 30000
	ErrNumDatabaseRecordNotFound
	ErrNumDatabaseDuplicateRecord
)

const (
	ErrMsgDatabase                = "database err"
	ErrMsgDatabaseRecordNotFound  = "record not found"
	ErrMsgDatabaseDuplicateRecord = "duplicate record"
)

var (
	ErrDatabase                = errx.New(ErrNumDatabase, ErrMsgDatabase)
	ErrDatabaseRecordNotFound  = errx.New(ErrNumDatabaseRecordNotFound, ErrMsgDatabaseRecordNotFound)
	ErrDatabaseDuplicateRecord = errx.New(ErrNumDatabaseDuplicateRecord, ErrMsgDatabaseDuplicateRecord)
)

// token err

const (
	ErrNumToken = iota + 40000
	ErrNumTokenInvalid
)

const (
	ErrMsgToken        = "token err"
	ErrMsgTokenInvalid = "token is invalid"
)

var (
	ErrToken        = errx.New(ErrNumToken, ErrMsgToken)
	ErrTokenInvalid = errx.New(ErrNumTokenInvalid, ErrMsgTokenInvalid)
)

// idl err

const (
	ErrNumIdl = iota + 50000
	ErrNumIdlAlreadyExist
	ErrNumIdlFileExtension
	ErrNumIdlGetDependentFilePath
)

const (
	ErrMsgIdl                     = "idl err"
	ErrMsgIdlAlreadyExist         = "idl is already exist"
	ErrMsgIdlFileExtension        = "invalid idl file extension"
	ErrMsgIdlGetDependentFilePath = "get dependent file paths from idl failed"
)

var (
	ErrIdl                     = errx.New(ErrNumIdl, ErrMsgIdl)
	ErrIdlAlreadyExist         = errx.New(ErrNumIdlAlreadyExist, ErrMsgIdlAlreadyExist)
	ErrIdlFileExtension        = errx.New(ErrNumIdlFileExtension, ErrMsgIdlFileExtension)
	ErrIdlGetDependentFilePath = errx.New(ErrNumIdlGetDependentFilePath, ErrMsgIdlGetDependentFilePath)
)

// repo err

const (
	ErrNumRepo = iota + 60000
	ErrNumRepoGetFile
	ErrNumRepoGetCommitHash
	ErrNumRepoGetArchive
	ErrNumRepoParseArchive
	ErrNumRepoGetPrivacy
	ErrNumRepoCreate
	ErrNumRepoGetClient
	ErrNumRepoPush
)

const (
	ErrMsgRepo              = "repo err"
	ErrMsgRepoGetFile       = "get file from repo failed"
	ErrMsgRepoGetCommitHash = "get commit hash failed"
	ErrMsgRepoGetArchive    = "get repo archive failed"
	ErrMsgRepoParseArchive  = "parse repo archive failed"
	ErrMsgRepoGetPrivacy    = "get repo privacy failed"
	ErrMsgRepoCreate        = "create repo failed"
	ErrMsgRepoGetClient     = "get repo client failed"
	ErrMsgRepoPush          = "push files to repo failed"
)

var (
	ErrRepo              = errx.New(ErrNumRepo, ErrMsgRepo)
	ErrRepoGetFile       = errx.New(ErrNumRepoGetFile, ErrMsgRepoGetFile)
	ErrRepoGetCommitHash = errx.New(ErrNumRepoGetCommitHash, ErrMsgRepoGetCommitHash)
	ErrRepoGetArchive    = errx.New(ErrNumRepoGetArchive, ErrMsgRepoGetArchive)
	ErrRepoParseArchive  = errx.New(ErrNumRepoParseArchive, ErrMsgRepoParseArchive)
	ErrRepoGetPrivacy    = errx.New(ErrNumRepoGetPrivacy, ErrMsgRepoGetPrivacy)
	ErrRepoCreate        = errx.New(ErrNumRepoCreate, ErrMsgRepoCreate)
	ErrRepoGetClient     = errx.New(ErrNumRepoGetClient, ErrMsgRepoGetClient)
	ErrRepoPush          = errx.New(ErrNumRepoPush, ErrMsgRepoPush)
)

// common err

const (
	ErrNumCommon = iota + 70000
	ErrNumCommonCreateTempDir
	ErrNumCommonGenerateCode
	ErrNumCommonProcessFolders
)

const (
	ErrMsgCommon               = "common err"
	ErrMsgCommonCreateTempDir  = "create temp dir failed"
	ErrMsgCommonGenerateCode   = "generate code failed"
	ErrMsgCommonProcessFolders = "process folders failed"
)

var (
	ErrCommon               = errx.New(ErrNumCommon, ErrMsgCommon)
	ErrCommonCreateTempDir  = errx.New(ErrNumCommonCreateTempDir, ErrMsgCommonCreateTempDir)
	ErrCommonGenerateCode   = errx.New(ErrNumCommonGenerateCode, ErrMsgCommonGenerateCode)
	ErrCommonProcessFolders = errx.New(ErrNumCommonProcessFolders, ErrMsgCommonProcessFolders)
)
