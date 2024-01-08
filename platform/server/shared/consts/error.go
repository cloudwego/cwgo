/*
 *
 * Copyright 2023 CloudWeGo Authors
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
 *
 */

package consts

import (
	"github.com/cloudwego/cwgo/platform/server/shared/errx"
)

const (
	ErrNumCommon = (iota + 1) * 10000
	ErrNumRegistry
	ErrNumRpc
	ErrNumDatabase
	ErrNumParam
	ErrNumToken
	ErrNumRepo
	ErrNumIdl
)

// common err

const (
	ErrNumCommonCreateTempDir = ErrNumCommon + iota + 1
	ErrNumCommonGenerateCode
	ErrNumCommonProcessFolders
	ErrNumCommonMkdir
	ErrNumCommonRepoApiService
	ErrNumCommonJsonMarshal
	ErrNumCommonJsonUnmarshal
)

const (
	ErrMsgCommon               = "common err"
	ErrMsgCommonCreateTempDir  = "create temp dir failed"
	ErrMsgCommonGenerateCode   = "generate code failed"
	ErrMsgCommonProcessFolders = "process folders failed"
	ErrMsgCommonMkdir          = "mkdir failed"
	ErrMsgCommonRepoApiService = "repo api service is down"
	ErrMsgCommonJsonMarshal
	ErrMsgCommonJsonUnmarshal
)

var (
	ErrCommon               = errx.New(ErrNumCommon, ErrMsgCommon)
	ErrCommonCreateTempDir  = errx.New(ErrNumCommonCreateTempDir, ErrMsgCommonCreateTempDir)
	ErrCommonGenerateCode   = errx.New(ErrNumCommonGenerateCode, ErrMsgCommonGenerateCode)
	ErrCommonProcessFolders = errx.New(ErrNumCommonProcessFolders, ErrMsgCommonProcessFolders)
	ErrCommonMkdir          = errx.New(ErrNumCommonMkdir, ErrMsgCommonMkdir)
	ErrCommonRepoApiService = errx.New(ErrNumCommonRepoApiService, ErrMsgCommonRepoApiService)
	ErrCommonJsonMarshal    = errx.New(ErrNumCommonJsonMarshal, ErrMsgCommonJsonMarshal)
	ErrCommonJsonUnmarshal  = errx.New(ErrNumCommonJsonUnmarshal, ErrMsgCommonJsonUnmarshal)
)

// registry err

const (
	ErrNumRegistryServiceNotFound = ErrNumRegistry + iota + 1
	ErrNumRegistryRegisterService
	ErrNumRegistryDeregisterService
)

const (
	ErrMsgRegistry                  = "registry err"
	ErrMsgRegistryServiceNotFound   = "service not found"
	ErrMsgRegistryRegisterService   = "register service failed"
	ErrMsgRegistryDeregisterService = "deregister service failed"
)

var (
	ErrRegistry                  = errx.New(ErrNumRegistry, ErrMsgRegistry)
	ErrRegistryServiceNotFound   = errx.New(ErrNumRegistryServiceNotFound, ErrMsgRegistryServiceNotFound)
	ErrRegistryRegisterService   = errx.New(ErrNumRegistryRegisterService, ErrMsgRegistryRegisterService)
	ErrRegistryDeregisterService = errx.New(ErrNumRegistryDeregisterService, ErrMsgRegistryDeregisterService)
)

// rpc err
const (
	ErrNumRpcGetClient = ErrNumRpc + iota + 1
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
	ErrNumDatabaseRecordNotFound = ErrNumDatabase + iota + 1
	ErrNumDatabaseDuplicateRecord

	ErrNumDatabaseRedisSet
	ErrNumDatabaseRedisPubSubClose
	ErrNumDatabaseRedisRunScript
	ErrNumDatabaseRedisTTL
	ErrNumDatabaseRedisExpire
	ErrNumDatabaseRedisSetNX
	ErrNumDatabaseRedisPublish
	ErrNumDatabaseRedisGet
	ErrNumDatabaseRedisHDel
	ErrNumDatabaseRedisHGelAll
)

const (
	ErrMsgDatabase                 = "database err"
	ErrMsgDatabaseRecordNotFound   = "record not found"
	ErrMsgDatabaseDuplicateRecord  = "duplicate record"
	ErrMsgDatabaseRedisSet         = "redis set failed"
	ErrMsgDatabaseRedisPubSubClose = "redis pubsub connection close failed"
	ErrMsgDatabaseRedisRunScript   = "redis run script failed"
	ErrMsgDatabaseRedisTTL         = "redis exec TTL command failed"
	ErrMsgDatabaseRedisExpire      = "redis exec EXPIRE command failed"
	ErrMsgDatabaseRedisSetNX       = "redis exec SETNX command failed"
	ErrMsgDatabaseRedisPublish     = "redis exec PUBLISH command failed"
	ErrMsgDatabaseRedisGet         = "redis exec GET command failed"
	ErrMsgDatabaseRedisHDel        = "redis exec HDEL command failed"
	ErrMsgDatabaseRedisHGelAll     = "redis exec HGETALL command failed"
)

var (
	ErrDatabase                 = errx.New(ErrNumDatabase, ErrMsgDatabase)
	ErrDatabaseRecordNotFound   = errx.New(ErrNumDatabaseRecordNotFound, ErrMsgDatabaseRecordNotFound)
	ErrDatabaseDuplicateRecord  = errx.New(ErrNumDatabaseDuplicateRecord, ErrMsgDatabaseDuplicateRecord)
	ErrDatabaseRedisSet         = errx.New(ErrNumDatabaseRedisSet, ErrMsgDatabaseRedisSet)
	ErrDatabaseRedisPubSubClose = errx.New(ErrNumDatabaseRedisPubSubClose, ErrMsgDatabaseRedisPubSubClose)
	ErrDatabaseRedisRunScript   = errx.New(ErrNumDatabaseRedisRunScript, ErrMsgDatabaseRedisRunScript)
	ErrDatabaseRedisTTL         = errx.New(ErrNumDatabaseRedisTTL, ErrMsgDatabaseRedisTTL)
	ErrDatabaseRedisExpire      = errx.New(ErrNumDatabaseRedisExpire, ErrMsgDatabaseRedisExpire)
	ErrDatabaseRedisSetNX       = errx.New(ErrNumDatabaseRedisSetNX, ErrMsgDatabaseRedisSetNX)
	ErrDatabaseRedisPublish     = errx.New(ErrNumDatabaseRedisPublish, ErrMsgDatabaseRedisPublish)
	ErrDatabaseRedisGet         = errx.New(ErrNumDatabaseRedisGet, ErrMsgDatabaseRedisGet)
	ErrDatabaseRedisHDel        = errx.New(ErrNumDatabaseRedisHDel, ErrMsgDatabaseRedisHDel)
	ErrDatabaseRedisHGelAll     = errx.New(ErrNumDatabaseRedisHGelAll, ErrMsgDatabaseRedisHGelAll)
)

// params err
const (
	ErrNumParamOrderNum = ErrNumParam + iota + 1
	ErrNumParamOrderBy
	ErrNumParamRepositoryType
	ErrNumParamStoreType
	ErrNumParamRepositoryUrl
	ErrNumParamDomain
	ErrNumParamMainIdlPath
	ErrNumParamTemplateType
	ErrNumParamUrl
	ErrNumParamRepositoryStatus
	ErrNumParamRepositoryBranch
	ErrNumParamIdlStatus
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
	ErrMsgParamRepositoryBranch = "invalid repository branch"
	ErrMsgParamIdlStatus        = "invalid idl status"
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
	ErrParamRepositoryBranch = errx.New(ErrNumParamRepositoryBranch, ErrMsgParamRepositoryBranch)
)

// token err

const (
	ErrNumTokenInvalid = ErrNumToken + iota + 1
	ErrNumTokenInvalidType
)

const (
	ErrMsgToken            = "token err"
	ErrMsgTokenInvalid     = "token is invalid"
	ErrMsgTokenInvalidType = "invalid token type"
)

var (
	ErrToken            = errx.New(ErrNumToken, ErrMsgToken)
	ErrTokenInvalid     = errx.New(ErrNumTokenInvalid, ErrMsgTokenInvalid)
	ErrTokenInvalidType = errx.New(ErrNumTokenInvalidType, ErrMsgTokenInvalidType)
)

// repo err

const (
	ErrNumRepoGetFile = ErrNumRepo + iota + 1
	ErrNumRepoGetCommitHash
	ErrNumRepoGetArchive
	ErrNumRepoParseArchive
	ErrNumRepoGetPrivacy
	ErrNumRepoCreate
	ErrNumRepoGetClient
	ErrNumRepoPush
	ErrNumRepoValidateBranch
)

const (
	ErrMsgRepo               = "repo err"
	ErrMsgRepoGetFile        = "get file from repo failed"
	ErrMsgRepoGetCommitHash  = "get commit hash failed"
	ErrMsgRepoGetArchive     = "get repo archive failed"
	ErrMsgRepoParseArchive   = "parse repo archive failed"
	ErrMsgRepoGetPrivacy     = "get repo privacy failed"
	ErrMsgRepoCreate         = "create repo failed"
	ErrMsgRepoGetClient      = "get repo client failed"
	ErrMsgRepoPush           = "push files to repo failed"
	ErrMsgRepoValidateBranch = "validate repo branch failed"
)

var (
	ErrRepo               = errx.New(ErrNumRepo, ErrMsgRepo)
	ErrRepoGetFile        = errx.New(ErrNumRepoGetFile, ErrMsgRepoGetFile)
	ErrRepoGetCommitHash  = errx.New(ErrNumRepoGetCommitHash, ErrMsgRepoGetCommitHash)
	ErrRepoGetArchive     = errx.New(ErrNumRepoGetArchive, ErrMsgRepoGetArchive)
	ErrRepoParseArchive   = errx.New(ErrNumRepoParseArchive, ErrMsgRepoParseArchive)
	ErrRepoGetPrivacy     = errx.New(ErrNumRepoGetPrivacy, ErrMsgRepoGetPrivacy)
	ErrRepoCreate         = errx.New(ErrNumRepoCreate, ErrMsgRepoCreate)
	ErrRepoGetClient      = errx.New(ErrNumRepoGetClient, ErrMsgRepoGetClient)
	ErrRepoPush           = errx.New(ErrNumRepoPush, ErrMsgRepoPush)
	ErrRepoValidateBranch = errx.New(ErrNumRepoValidateBranch, ErrMsgRepoValidateBranch)
)

// idl err

const (
	ErrNumIdlAlreadyExist = ErrNumIdl + iota + 1
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
