/*
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
 */

// Code generated by Kitex v0.6.1. DO NOT EDIT.

package helloservice

import (
	"context"
	example "cwgo/example/hex/kitex_gen/hello/example"
	client "github.com/cloudwego/kitex/client"
	kitex "github.com/cloudwego/kitex/pkg/serviceinfo"
)

func serviceInfo() *kitex.ServiceInfo {
	return helloServiceServiceInfo
}

var helloServiceServiceInfo = NewServiceInfo()

func NewServiceInfo() *kitex.ServiceInfo {
	serviceName := "HelloService"
	handlerType := (*example.HelloService)(nil)
	methods := map[string]kitex.MethodInfo{
		"HelloMethod": kitex.NewMethodInfo(helloMethodHandler, newHelloServiceHelloMethodArgs, newHelloServiceHelloMethodResult, false),
	}
	extra := map[string]interface{}{
		"PackageName": "example",
	}
	svcInfo := &kitex.ServiceInfo{
		ServiceName:     serviceName,
		HandlerType:     handlerType,
		Methods:         methods,
		PayloadCodec:    kitex.Thrift,
		KiteXGenVersion: "v0.6.1",
		Extra:           extra,
	}
	return svcInfo
}

func helloMethodHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*example.HelloServiceHelloMethodArgs)
	realResult := result.(*example.HelloServiceHelloMethodResult)
	success, err := handler.(example.HelloService).HelloMethod(ctx, realArg.Request)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newHelloServiceHelloMethodArgs() interface{} {
	return example.NewHelloServiceHelloMethodArgs()
}

func newHelloServiceHelloMethodResult() interface{} {
	return example.NewHelloServiceHelloMethodResult()
}

type kClient struct {
	c client.Client
}

func newServiceClient(c client.Client) *kClient {
	return &kClient{
		c: c,
	}
}

func (p *kClient) HelloMethod(ctx context.Context, request *example.HelloReq) (r *example.HelloResp, err error) {
	var _args example.HelloServiceHelloMethodArgs
	_args.Request = request
	var _result example.HelloServiceHelloMethodResult
	if err = p.c.Call(ctx, "HelloMethod", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}
