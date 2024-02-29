/*
 * Copyright 2022 CloudWeGo Authors
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

package api_list

type VarType int

const (
	VarTypeOther = iota
	VarTypeServerHertz
	VarTypeRouteEngine
	VarTypeRouterGroup
)

const (
	RouterRegisterFuncNameGET     = "GET"
	RouterRegisterFuncNamePOST    = "POST"
	RouterRegisterFuncNamePUT     = "PUT"
	RouterRegisterFuncNameDELETE  = "DELETE"
	RouterRegisterFuncNameHEAD    = "HEAD"
	RouterRegisterFuncNamePATCH   = "PATCH"
	RouterRegisterFuncNameOPTIONS = "OPTIONS"

	RouterRegisterFuncNameGETEX    = "GETEX"
	RouterRegisterFuncNamePOSTEX   = "POSTEX"
	RouterRegisterFuncNamePUTEX    = "PUTEX"
	RouterRegisterFuncNameDELETEEX = "DELETEEX"
	RouterRegisterFuncNameHEADEX   = "HEADEX"
	RouterRegisterFuncNameAnyEX    = "AnyEX"
)

const (
	BuiltinFuncNameAppend  = "append"
	BuiltinFuncNameCopy    = "copy"
	BuiltinFuncNameDelete  = "delete"
	BuiltinFuncNameLen     = "len"
	BuiltinFuncNameCap     = "cap"
	BuiltinFuncNameMake    = "make"
	BuiltinFuncNameNew     = "new"
	BuiltinFuncNameComplex = "complex"
	BuiltinFuncNameReal    = "real"
	BuiltinFuncNameImag    = "imag"
	BuiltinFuncNameClear   = "clear"
	BuiltinFuncNameClose   = "close"
	BuiltinFuncNamePanic   = "panic"
	BuiltinFuncNameRecover = "recover"
	BuiltinFuncNamePrint   = "print"
	BuiltinFuncNamePrintln = "println"
)

var (
	HertzFuncAssignmentFuncOfCoreMap = map[string]struct{}{
		"Default": {},
		"New":     {},
	}

	RouterFuncNameMap = map[string]struct{}{
		RouterRegisterFuncNameGET:      {},
		RouterRegisterFuncNamePOST:     {},
		RouterRegisterFuncNamePUT:      {},
		RouterRegisterFuncNameDELETE:   {},
		RouterRegisterFuncNameHEAD:     {},
		RouterRegisterFuncNamePATCH:    {},
		RouterRegisterFuncNameOPTIONS:  {},
		RouterRegisterFuncNameGETEX:    {},
		RouterRegisterFuncNamePOSTEX:   {},
		RouterRegisterFuncNamePUTEX:    {},
		RouterRegisterFuncNameDELETEEX: {},
		RouterRegisterFuncNameHEADEX:   {},
		RouterRegisterFuncNameAnyEX:    {},
	}

	BuiltinFuncNameMap = map[string]struct{}{
		BuiltinFuncNameAppend:  {},
		BuiltinFuncNameCopy:    {},
		BuiltinFuncNameDelete:  {},
		BuiltinFuncNameLen:     {},
		BuiltinFuncNameCap:     {},
		BuiltinFuncNameMake:    {},
		BuiltinFuncNameNew:     {},
		BuiltinFuncNameComplex: {},
		BuiltinFuncNameReal:    {},
		BuiltinFuncNameImag:    {},
		BuiltinFuncNameClear:   {},
		BuiltinFuncNameClose:   {},
		BuiltinFuncNamePanic:   {},
		BuiltinFuncNameRecover: {},
		BuiltinFuncNamePrint:   {},
		BuiltinFuncNamePrintln: {},
	}
)
