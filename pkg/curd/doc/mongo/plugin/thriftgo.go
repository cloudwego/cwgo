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

package plugin

import (
	"fmt"
	"go/format"
	"io"
	"os"

	"github.com/cloudwego/cwgo/pkg/curd/code"
	"github.com/cloudwego/cwgo/pkg/curd/doc/mongo/codegen"
	"github.com/cloudwego/cwgo/pkg/curd/extract"
	"github.com/cloudwego/cwgo/pkg/curd/parse"
	"github.com/cloudwego/cwgo/pkg/curd/template"

	"github.com/cloudwego/cwgo/config"
	cwgoMeta "github.com/cloudwego/cwgo/meta"
	"github.com/cloudwego/hertz/cmd/hz/meta"
	"github.com/cloudwego/hertz/cmd/hz/util/logs"
	"github.com/cloudwego/thriftgo/plugin"
)

type thriftGoPlugin struct {
	req     *plugin.Request
	docArgs *config.DocArgument
}

func thriftPluginRun() int {
	plu := &thriftGoPlugin{}

	if err := plu.handleRequest(); err != nil {
		logs.Errorf("handle request failed: %s", err.Error())
		return meta.PluginError
	}

	if err := plu.parseArgs(); err != nil {
		logs.Errorf("parse args failed: %s", err.Error())
		return meta.PluginError
	}

	tfUsedInfo := &extract.ThriftUsedInfo{
		Req:     plu.req,
		DocArgs: plu.docArgs,
	}
	rawStructs, err := tfUsedInfo.ParseThriftIdl()
	if err != nil {
		logs.Errorf("parse thrift idl failed: %s", err.Error())
		return meta.PluginError
	}

	operations, err := parse.HandleOperations(rawStructs)
	if err != nil {
		logs.Error(err.Error())
		return meta.PluginError
	}

	methodRenders := codegen.HandleCodegen(operations)
	generated, err := plu.buildResponse(rawStructs, methodRenders, tfUsedInfo)
	if err != nil {
		logs.Error(err.Error())
		return meta.PluginError
	}

	res := &plugin.Response{
		Contents: generated,
	}
	if err = response(res); err != nil {
		logs.Error(err.Error())
		return meta.PluginError
	}

	return 0
}

func (plu *thriftGoPlugin) handleRequest() error {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return fmt.Errorf("read request failed: %s", err.Error())
	}

	req, err := plugin.UnmarshalRequest(data)
	if err != nil {
		return fmt.Errorf("unmarshal request failed: %s", err.Error())
	}

	plu.req = req
	return nil
}

func (plu *thriftGoPlugin) parseArgs() error {
	if plu.req == nil {
		return fmt.Errorf("request is nil")
	}
	args := new(config.DocArgument)
	if err := args.Unpack(plu.req.PluginParameters); err != nil {
		logs.Errorf("unpack args failed: %s", err.Error())
		return err
	}
	plu.docArgs = args
	return nil
}

func response(res *plugin.Response) error {
	data, err := plugin.MarshalResponse(res)
	if err != nil {
		return fmt.Errorf("marshal response failed: %s", err.Error())
	}
	_, err = os.Stdout.Write(data)
	if err != nil {
		return fmt.Errorf("write response failed: %s", err.Error())
	}
	return nil
}

func (plu *thriftGoPlugin) buildResponse(structs []*extract.IdlExtractStruct, methodRenders [][]*template.MethodRender,
	info *extract.ThriftUsedInfo,
) (result []*plugin.Generated, err error) {
	for index, st := range structs {
		// get base render
		baseRender := getBaseRender(st)
		// get fileMongoName and fileIfName
		fileMongoName, fileIfName := extract.GetFileName(st.Name, plu.docArgs.DaoDir)

		if st.Update {
			// build update mongo file
			formattedCode, err := getUpdateMongoCode(methodRenders[index], string(st.UpdateCurdFileContent))
			if err != nil {
				return nil, err
			}
			formattedCode, err = codegen.AddMongoImports(formattedCode)
			if err != nil {
				return nil, err
			}
			formattedCode, err = extract.AddMongoModelImports(formattedCode, info.ImportPaths)
			if err != nil {
				return nil, err
			}
			result = append(result, &plugin.Generated{
				Content: formattedCode,
				Name:    &fileMongoName,
			})

			// build update interface file
			formattedCode, err = getUpdateIfCode(st, baseRender)
			if err != nil {
				return nil, err
			}
			formattedCode, err = codegen.AddMongoImports(formattedCode)
			if err != nil {
				return nil, err
			}
			formattedCode, err = extract.AddMongoModelImports(formattedCode, info.ImportPaths)
			if err != nil {
				return nil, err
			}
			result = append(result, &plugin.Generated{
				Content: formattedCode,
				Name:    &fileIfName,
			})
		} else {
			// build new mongo file
			formattedCode, err := getNewMongoCode(methodRenders[index], st, baseRender)
			if err != nil {
				return nil, err
			}
			formattedCode, err = codegen.AddMongoImports(formattedCode)
			if err != nil {
				return nil, err
			}
			formattedCode, err = extract.AddMongoModelImports(formattedCode, info.ImportPaths)
			if err != nil {
				return nil, err
			}
			result = append(result, &plugin.Generated{
				Content: formattedCode,
				Name:    &fileMongoName,
			})

			// build new interface file
			formattedCode, err = getNewIfCode(st, baseRender)
			if err != nil {
				return nil, err
			}
			formattedCode, err = codegen.AddMongoImports(formattedCode)
			if err != nil {
				return nil, err
			}
			formattedCode, err = extract.AddMongoModelImports(formattedCode, info.ImportPaths)
			if err != nil {
				return nil, err
			}
			result = append(result, &plugin.Generated{
				Content: formattedCode,
				Name:    &fileIfName,
			})
		}
	}

	return
}

func getBaseRender(st *extract.IdlExtractStruct) *template.BaseRender {
	pkgName := extract.GetPkgName(st.Name)
	return &template.BaseRender{
		Version:     cwgoMeta.Version,
		PackageName: pkgName,
		Imports:     codegen.BaseMongoImports,
	}
}

func getUpdateMongoCode(methodRenders []*template.MethodRender, fileContent string) (string, error) {
	tplMongo := &template.Template{
		Renders: []template.Render{},
	}
	for _, methodRender := range methodRenders {
		tplMongo.Renders = append(tplMongo.Renders, methodRender)
	}

	buff, err := tplMongo.Build()
	if err != nil {
		return "", err
	}
	data := fileContent + "\n" + buff.String()
	formattedCode, err := format.Source([]byte(data))
	if err != nil {
		return "", err
	}

	return string(formattedCode), nil
}

func getUpdateIfCode(st *extract.IdlExtractStruct, baseRender *template.BaseRender) (string, error) {
	tplIf := &template.Template{
		Renders: []template.Render{},
	}
	tplIf.Renders = append(tplIf.Renders, baseRender)

	methods := make(code.InterfaceMethods, 0, 10)
	for _, preMethod := range st.PreIfMethods {
		methods = append(methods, code.InterfaceMethod{
			Name:    preMethod.Name,
			Params:  preMethod.Params,
			Returns: preMethod.Returns,
		})
	}
	for _, rawMethod := range st.InterfaceInfo.Methods {
		methods = append(methods, code.InterfaceMethod{
			Name:    rawMethod.Name,
			Params:  rawMethod.Params,
			Returns: rawMethod.Returns,
		})
	}

	ifRender := &template.InterfaceRender{
		Name:    st.Name + "Repository",
		Methods: methods,
	}
	tplIf.Renders = append(tplIf.Renders, ifRender)

	buff, err := tplIf.Build()
	if err != nil {
		return "", err
	}
	formattedCode, err := format.Source(buff.Bytes())
	if err != nil {
		return "", err
	}

	return string(formattedCode), nil
}

func getNewMongoCode(methodRenders []*template.MethodRender, st *extract.IdlExtractStruct, baseRender *template.BaseRender) (string, error) {
	tplMongo := &template.Template{
		Renders: []template.Render{},
	}

	tplMongo.Renders = append(tplMongo.Renders, baseRender)
	tplMongo.Renders = append(tplMongo.Renders, codegen.GetFuncRender(st))
	tplMongo.Renders = append(tplMongo.Renders, codegen.GetStructRender(st))
	for _, methodRender := range methodRenders {
		tplMongo.Renders = append(tplMongo.Renders, methodRender)
	}

	buff, err := tplMongo.Build()
	if err != nil {
		return "", err
	}
	formattedCode, err := format.Source(buff.Bytes())
	if err != nil {
		return "", err
	}

	return string(formattedCode), nil
}

func getNewIfCode(st *extract.IdlExtractStruct, baseRender *template.BaseRender) (string, error) {
	tplIf := &template.Template{
		Renders: []template.Render{},
	}
	tplIf.Renders = append(tplIf.Renders, baseRender)

	methods := make(code.InterfaceMethods, 0, 10)
	for _, rawMethod := range st.InterfaceInfo.Methods {
		methods = append(methods, code.InterfaceMethod{
			Name:    rawMethod.Name,
			Params:  rawMethod.Params,
			Returns: rawMethod.Returns,
		})
	}
	ifRender := &template.InterfaceRender{
		Name:    st.Name + "Repository",
		Methods: methods,
	}
	tplIf.Renders = append(tplIf.Renders, ifRender)

	buff, err := tplIf.Build()
	if err != nil {
		return "", err
	}
	formattedCode, err := format.Source(buff.Bytes())
	if err != nil {
		return "", err
	}

	return string(formattedCode), nil
}
