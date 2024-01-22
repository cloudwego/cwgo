package model

import (
	"reflect"

	"github.com/cloudwego/cwgo/pkg/doc/mongo/code"
)

type IdlExtractStruct struct {
	Name          string
	StructFields  []*StructField
	InterfaceInfo *InterfaceInfo
	UpdateInfo
}

type InterfaceInfo struct {
	Name    string
	Methods []*InterfaceMethod
}

type InterfaceMethod struct {
	Name             string
	ParsedTokens     string
	Params           code.Params
	Returns          code.Returns
	BelongedToStruct *IdlExtractStruct
}

type StructField struct {
	Name               string
	Type               code.Type
	Tag                reflect.StructTag
	IsBelongedToStruct bool
	BelongedToStruct   *IdlExtractStruct
}

type UpdateInfo struct {
	Update                 bool
	UpdateMongoFileContent []byte
	UpdateIfFileContent    []byte
	PreMethodNamesMap      map[string]struct{}
	PreIfMethods           []*InterfaceMethod
}
