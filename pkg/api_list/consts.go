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
