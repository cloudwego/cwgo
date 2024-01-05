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

package server

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/cloudwego/cwgo/config"
	"github.com/cloudwego/cwgo/pkg/common/utils"
	"github.com/cloudwego/cwgo/pkg/consts"
	"github.com/cloudwego/cwgo/tpl"
	hzConfig "github.com/cloudwego/hertz/cmd/hz/config"
	"github.com/cloudwego/hertz/cmd/hz/meta"
	"github.com/cloudwego/kitex"
	kargs "github.com/cloudwego/kitex/tool/cmd/kitex/args"
	"github.com/cloudwego/kitex/tool/internal_pkg/generator"
	"github.com/cloudwego/kitex/tool/internal_pkg/log"
)

func convertKitexArgs(sa *config.ServerArgument, kitexArgument *kargs.Arguments) (err error) {
	f := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	kitexArgument.ModuleName = sa.GoMod
	kitexArgument.ServiceName = sa.Service
	kitexArgument.Includes = sa.SliceParam.ProtoSearchPath
	kitexArgument.Version = kitex.Version
	kitexArgument.RecordCmd = os.Args
	kitexArgument.ThriftOptions = append(kitexArgument.ThriftOptions,
		"naming_style=golint",
		"ignore_initialisms",
		"gen_setter",
		"gen_deep_equal",
		"compatible_names",
		"frugal_tag",
	)
	kitexArgument.IDL = sa.IdlPath

	f.BoolVar(&kitexArgument.NoFastAPI, "no-fast-api", false, "Generate codes without injecting fast method.")
	f.StringVar(&kitexArgument.Use, "use", "",
		"Specify the kitex_gen package to import when generate server side codes.")
	f.BoolVar(&kitexArgument.GenerateInvoker, "invoker", false,
		"Generate invoker side codes when service name is specified.")
	f.StringVar(&kitexArgument.IDLType, "type", "unknown", "Specify the type of IDL: 'thrift' or 'protobuf'.")
	f.Var(&kitexArgument.ThriftOptions, "thrift", "Specify arguments for the thrift go compiler.")
	f.DurationVar(&kitexArgument.ThriftPluginTimeLimit, "thrift-plugin-time-limit", generator.DefaultThriftPluginTimeLimit, "Specify thrift plugin execution time limit.")
	f.Var(&kitexArgument.ThriftPlugins, "thrift-plugin", "Specify thrift plugin arguments for the thrift compiler.")
	f.Var(&kitexArgument.ProtobufPlugins, "protobuf-plugin", "Specify protobuf plugin arguments for the protobuf compiler.(plugin_name:options:out_dir)")
	f.Var(&kitexArgument.ProtobufOptions, "protobuf", "Specify arguments for the protobuf compiler.")
	f.BoolVar(&kitexArgument.CombineService, "combine-service", false,
		"Combine services in root thrift file.")
	f.BoolVar(&kitexArgument.CopyIDL, "copy-idl", false,
		"Copy each IDL file to the output path.")
	f.StringVar(&kitexArgument.ExtensionFile, "template-extension", kitexArgument.ExtensionFile,
		"Specify a file for template extension.")
	f.BoolVar(&kitexArgument.FrugalPretouch, "frugal-pretouch", false,
		"Use frugal to compile arguments and results when new clients and servers.")
	f.BoolVar(&kitexArgument.Record, "record", false, "Record Kitex cmd into kitex-all.sh.")
	f.StringVar(&kitexArgument.GenPath, "gen-path", generator.KitexGenPath,
		"Specify a code gen path.")
	f.StringVar(&kitexArgument.Protocol, "protocol", "", "Specify a protocol for codec.")
	f.Var(&kitexArgument.Hessian2Options, "hessian2", "Specify arguments for the hessian2 codec.")

	f.Usage = func() {
		fmt.Fprintf(os.Stderr, `Version %s
Usage: %s [flags] IDL

Flags:
`, kitexArgument.Version, os.Args[0])
		f.PrintDefaults()
		os.Exit(1)
	}

	err = f.Parse(utils.StringSliceSpilt(sa.SliceParam.Pass))
	if err != nil {
		return
	}

	// Non-standard template
	if strings.HasSuffix(sa.Template, consts.SuffixGit) {
		err = utils.GitClone(sa.Template, path.Join(tpl.KitexDir, consts.Server))
		if err != nil {
			return err
		}
		gitPath, err := utils.GitPath(sa.Template)
		if err != nil {
			return err
		}
		gitPath = path.Join(tpl.KitexDir, consts.Server, gitPath)
		kitexArgument.TemplateDir = gitPath
	} else {
		if len(sa.Template) != 0 {
			kitexArgument.TemplateDir = sa.Template
		} else {
			kitexArgument.TemplateDir = path.Join(tpl.KitexDir, consts.Server, consts.Standard)
		}
	}

	kitexArgument.GenerateMain = false

	return checkKitexArgs(kitexArgument)
}

func checkKitexArgs(a *kargs.Arguments) (err error) {
	// check IDL
	a.IDLType, err = utils.GetIdlType(a.IDL, consts.Protobuf)
	if err != nil {
		return err
	}

	// check service name
	if a.ServiceName == "" {
		if a.Use != "" {
			log.Warn("-use must be used with -service")
			os.Exit(2)
		}
	}

	gopath, err := utils.GetGOPATH()
	if err != nil {
		return fmt.Errorf("get gopath failed: %s", err)
	}
	if gopath == "" {
		return fmt.Errorf("GOPATH is not set")
	}

	gosrc := filepath.Join(gopath, consts.Src)
	gosrc, err = filepath.Abs(gosrc)
	if err != nil {
		log.Warn("Get GOPATH/src path failed:", err.Error())
		os.Exit(1)
	}
	curpath, err := filepath.Abs(consts.CurrentDir)
	if err != nil {
		log.Warn("Get current path failed:", err.Error())
		os.Exit(1)
	}

	if strings.HasPrefix(curpath, gosrc) {
		if a.PackagePrefix, err = filepath.Rel(gosrc, curpath); err != nil {
			log.Warn("Get GOPATH/src relpath failed:", err.Error())
			os.Exit(1)
		}
		a.PackagePrefix = filepath.Join(a.PackagePrefix, generator.KitexGenPath)
	} else {
		if a.ModuleName == "" {
			log.Warn("Outside of $GOPATH. Please specify a module name with the '-module' flag.")
			os.Exit(1)
		}
	}

	if a.ModuleName != "" {
		module, p, ok := utils.SearchGoMod(curpath, true)
		if ok {
			// go.mod exists
			if module != a.ModuleName {
				log.Warnf("The module name given by the '-module' option ('%s') is not consist with the name defined in go.mod ('%s' from %s)\n",
					a.ModuleName, module, p)
				os.Exit(1)
			}
			if a.PackagePrefix, err = filepath.Rel(p, curpath); err != nil {
				log.Warn("Get package prefix failed:", err.Error())
				os.Exit(1)
			}
			a.PackagePrefix = filepath.Join(a.ModuleName, a.PackagePrefix, generator.KitexGenPath)
		} else {
			if err = utils.InitGoMod(a.ModuleName); err != nil {
				log.Warn("Init go mod failed:", err.Error())
				os.Exit(1)
			}
			a.PackagePrefix = filepath.Join(a.ModuleName, generator.KitexGenPath)
		}
	}

	if a.Use != "" {
		a.PackagePrefix = a.Use
	}
	a.OutputPath = curpath
	a.PackagePrefix = strings.ReplaceAll(a.PackagePrefix, consts.BackSlash, consts.Slash)
	return nil
}

func hzArgsForHex(c *config.ServerArgument) (*hzConfig.Argument, error) {
	utils.SetHzVerboseLog(c.Verbose)
	hzArgs := hzConfig.NewArgument()
	err := convertHzArgument(c, hzArgs)
	if err != nil {
		return nil, err
	}
	hzArgs.CmdType = meta.CmdUpdate // update command is enough for hex
	// these options are aligned with the kitex
	if strings.EqualFold(hzArgs.IdlType, consts.Thrift) {
		hzArgs.ThriftOptions = append(hzArgs.ThriftOptions, "naming_style=golint", "ignore_initialisms", "gen_setter", "gen_deep_equal", "compatible_names", "frugal_tag")
		hzArgs.ModelDir = consts.DefaultKitexModelDir
	}
	if strings.EqualFold(hzArgs.IdlType, consts.Proto) {
		hzArgs.Use = fmt.Sprintf("%s/%s", hzArgs.Gomod, consts.DefaultKitexModelDir)
	}
	if hzArgs.CustomizePackage == path.Join(tpl.HertzDir, consts.Server, consts.Standard, consts.PackageLayoutFile) {
		hzArgs.CustomizePackage = "" // disable the default hertz template for hex
	}
	return hzArgs, nil
}

func generateHexFile(c *config.ServerArgument) error {
	tmplContent := `package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"regexp"

	"github.com/cloudwego/hertz/pkg/app"
	hertzServer "github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/network"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/cloudwego/kitex/pkg/endpoint"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/remote"
	"github.com/cloudwego/kitex/pkg/remote/trans/detection"
	"github.com/cloudwego/kitex/pkg/remote/trans/netpoll"
	"github.com/cloudwego/kitex/pkg/remote/trans/nphttp2"
	"{{$.ProjPackage}}/biz/router"
)

type mixTransHandlerFactory struct {
	originFactory remote.ServerTransHandlerFactory
}

type transHandler struct {
	remote.ServerTransHandler
}

// SetInvokeHandleFunc is used to set invoke handle func.
func (t *transHandler) SetInvokeHandleFunc(inkHdlFunc endpoint.Endpoint) {
	t.ServerTransHandler.(remote.InvokeHandleFuncSetter).SetInvokeHandleFunc(inkHdlFunc)
}

func (m mixTransHandlerFactory) NewTransHandler(opt *remote.ServerOption) (remote.ServerTransHandler, error) {
	var kitexOrigin remote.ServerTransHandler
	var err error

	if m.originFactory != nil {
		kitexOrigin, err = m.originFactory.NewTransHandler(opt)
	} else {
		// if no customized factory just use the default factory under detection pkg.
		kitexOrigin, err = detection.NewSvrTransHandlerFactory(netpoll.NewSvrTransHandlerFactory(), nphttp2.NewSvrTransHandlerFactory()).NewTransHandler(opt)
	}
	if err != nil {
		return nil, err
	}
	return &transHandler{ServerTransHandler: kitexOrigin}, nil
}

var httpReg = regexp.MustCompile(` + "`^(?:GET |POST|PUT|DELE|HEAD|OPTI|CONN|TRAC|PATC)$`" + `)

func (t *transHandler) OnRead(ctx context.Context, conn net.Conn) error {
	c, ok := conn.(network.Conn)
	if ok {
		pre, _ := c.Peek(4)
		if httpReg.Match(pre) {
			klog.Info("using Hertz to process request")
			err := hertzEngine.Serve(ctx, c)
			if err != nil {
				err = errors.New(fmt.Sprintf("HERTZ: %s", err.Error()))
			}
			return err
		}
	}
	return t.ServerTransHandler.OnRead(ctx, conn)
}

func initHertz() *route.Engine {
	h := hertzServer.New(hertzServer.WithIdleTimeout(0))
	// add a ping route to test
	h.GET("/ping", func(c context.Context, ctx *app.RequestContext) {
		ctx.JSON(consts.StatusOK, utils.H{"ping": "pong"})
	})
	router.GeneratedRegister(h)
	if err := h.Engine.Init(); err != nil {
		panic(err)
	}
	//if err := h.Engine.SetEngineRun(); err != nil {
	//	panic(err)
	//}
	return h.Engine
}

var hertzEngine *route.Engine

func init() {
	hertzEngine = initHertz()
}
`
	exist, err := utils.PathExist("hex_trans_handler.go")
	if err != nil {
		return err
	}
	if exist {
		return nil
	}
	tmpl := template.Must(template.New("hex_trans_handler").Parse(tmplContent))
	file, err := os.Create("hex_trans_handler.go")
	if err != nil {
		return err
	}
	defer file.Close()
	return tmpl.Execute(file, map[string]string{
		"ProjPackage": c.GoMod,
	})
}

func addHexOptions() error {
	filePath := consts.Main
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	if bytes.Contains(content, []byte("server.WithTransHandlerFactory(&mixTransHandlerFactory{nil})")) {
		return nil
	}
	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	found, err := insertCodeInFunction(astFile, "kitexInit", "opts", "append(opts,server.WithTransHandlerFactory(&mixTransHandlerFactory{nil}))")
	if err != nil {
		return err
	}
	if !found {
		return nil
	}
	outputFile, err := os.Create(consts.Main)
	if err != nil {
		return err
	}
	defer outputFile.Close()
	err = printer.Fprint(outputFile, fset, astFile)
	if err != nil {
		return err
	}

	return nil
}

func insertCodeInFunction(file *ast.File, functionName, left, right string) (bool, error) {
	for _, decl := range file.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		if funcDecl.Name.Name == functionName {
			insertedStmt, err := parser.ParseExpr(right)
			if err != nil {
				return false, err
			}

			assignStmt := &ast.AssignStmt{
				Tok: token.ASSIGN,
				Lhs: []ast.Expr{ast.NewIdent(left)},
				Rhs: []ast.Expr{insertedStmt},
			}

			funcDecl.Body.List = append([]ast.Stmt{assignStmt}, funcDecl.Body.List...)
			return true, nil
		}
	}
	return false, nil
}
