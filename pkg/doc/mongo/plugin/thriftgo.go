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
	"io/ioutil"
	"os"

	"github.com/cloudwego/cwgo/config"
	cwgoMeta "github.com/cloudwego/cwgo/meta"
	"github.com/cloudwego/cwgo/pkg/doc/mongo/code"
	"github.com/cloudwego/cwgo/pkg/doc/mongo/codegen"
	"github.com/cloudwego/cwgo/pkg/doc/mongo/parse"
	"github.com/cloudwego/cwgo/pkg/doc/mongo/plugin/model"
	"github.com/cloudwego/cwgo/pkg/doc/mongo/template"
	"github.com/cloudwego/hertz/cmd/hz/meta"
	"github.com/cloudwego/hertz/cmd/hz/util/logs"
	"github.com/cloudwego/thriftgo/plugin"
	"golang.org/x/tools/imports"
)

type ThriftGoPlugin struct {
	Req     *plugin.Request
	DocArgs *config.DocArgument
}

func pluginRun() int {
	plu := &ThriftGoPlugin{}

	if err := plu.handleRequest(); err != nil {
		logs.Errorf("handle request failed: %s", err.Error())
		return meta.PluginError
	}

	if err := plu.parseArgs(); err != nil {
		logs.Errorf("parse args failed: %s", err.Error())
		return meta.PluginError
	}

	rawStructs, err := parseThriftIdl(plu)
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
	generated, err := plu.buildResponse(rawStructs, methodRenders)
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

func (plu *ThriftGoPlugin) handleRequest() error {
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return fmt.Errorf("read request failed: %s", err.Error())
	}

	req, err := plugin.UnmarshalRequest(data)
	if err != nil {
		return fmt.Errorf("unmarshal request failed: %s", err.Error())
	}

	plu.Req = req
	return nil
}

func (plu *ThriftGoPlugin) parseArgs() error {
	if plu.Req == nil {
		return fmt.Errorf("request is nil")
	}
	args := new(config.DocArgument)
	if err := args.Unpack(plu.Req.PluginParameters); err != nil {
		logs.Errorf("unpack args failed: %s", err.Error())
		return err
	}
	plu.DocArgs = args
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

func (plu *ThriftGoPlugin) buildResponse(structs []*model.IdlExtractStruct, methodRenders [][]*template.MethodRender) (result []*plugin.Generated, err error) {
	for index, struc := range structs {
		pkgName := getPkgName(struc.Name)
		baseRender := &template.BaseRender{
			Version:     cwgoMeta.Version,
			PackageName: pkgName,
			Imports:     codegen.BaseMongoImports,
		}

		fileMongoName, fileIfName := getFileName(struc.Name, plu.DocArgs.DaoDir)
		if struc.Update {
			// build update mongo file
			tplMongo := &template.Template{
				Renders: []template.Render{},
			}
			for _, methodRender := range methodRenders[index] {
				tplMongo.Renders = append(tplMongo.Renders, methodRender)
			}

			buff, err := tplMongo.Build()
			if err != nil {
				return nil, err
			}
			data := string(struc.UpdateMongoFileContent) + "\n" + buff.String()
			formattedCode, err := imports.Process("", []byte(data), nil)
			if err != nil {
				return nil, err
			}
			result = append(result, &plugin.Generated{
				Content: string(formattedCode),
				Name:    &fileMongoName,
			})

			// build update interface file
			tplIf := &template.Template{
				Renders: []template.Render{},
			}
			tplIf.Renders = append(tplIf.Renders, baseRender)

			methods := make(code.InterfaceMethods, 0, 10)
			for _, preMethod := range struc.PreIfMethods {
				methods = append(methods, code.InterfaceMethod{
					Name:    preMethod.Name,
					Params:  preMethod.Params,
					Returns: preMethod.Returns,
				})
			}
			for _, rawMethod := range struc.InterfaceInfo.Methods {
				methods = append(methods, code.InterfaceMethod{
					Name:    rawMethod.Name,
					Params:  rawMethod.Params,
					Returns: rawMethod.Returns,
				})
			}

			ifRender := &template.InterfaceRender{
				Name:    struc.Name + "Repository",
				Methods: methods,
			}
			tplIf.Renders = append(tplIf.Renders, ifRender)

			buff, err = tplIf.Build()
			if err != nil {
				return nil, err
			}
			formattedCode, err = imports.Process("", buff.Bytes(), nil)
			if err != nil {
				return nil, err
			}
			result = append(result, &plugin.Generated{
				Content: string(formattedCode),
				Name:    &fileIfName,
			})
		} else {
			// build new mongo file
			tplMongo := &template.Template{
				Renders: []template.Render{},
			}

			tplMongo.Renders = append(tplMongo.Renders, baseRender)
			tplMongo.Renders = append(tplMongo.Renders, codegen.GetFuncRender(struc))
			tplMongo.Renders = append(tplMongo.Renders, codegen.GetStructRender(struc))
			for _, methodRender := range methodRenders[index] {
				tplMongo.Renders = append(tplMongo.Renders, methodRender)
			}

			buff, err := tplMongo.Build()
			if err != nil {
				return nil, err
			}
			formattedCode, err := imports.Process("", buff.Bytes(), nil)
			if err != nil {
				return nil, err
			}
			result = append(result, &plugin.Generated{
				Content: string(formattedCode),
				Name:    &fileMongoName,
			})

			// build new interface file
			tplIf := &template.Template{
				Renders: []template.Render{},
			}
			tplIf.Renders = append(tplIf.Renders, baseRender)

			methods := make(code.InterfaceMethods, 0, 10)
			for _, rawMethod := range struc.InterfaceInfo.Methods {
				methods = append(methods, code.InterfaceMethod{
					Name:    rawMethod.Name,
					Params:  rawMethod.Params,
					Returns: rawMethod.Returns,
				})
			}
			ifRender := &template.InterfaceRender{
				Name:    struc.Name + "Repository",
				Methods: methods,
			}
			tplIf.Renders = append(tplIf.Renders, ifRender)

			buff, err = tplIf.Build()
			if err != nil {
				return nil, err
			}
			formattedCode, err = imports.Process("", buff.Bytes(), nil)
			if err != nil {
				return nil, err
			}
			result = append(result, &plugin.Generated{
				Content: string(formattedCode),
				Name:    &fileIfName,
			})
		}
	}

	return
}
