// Code generated by Kitex v0.8.0. DO NOT EDIT.

package agentservice

import (
	agent "github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/agent"
	server "github.com/cloudwego/kitex/server"
)

// NewInvoker creates a server.Invoker with the given handler and options.
func NewInvoker(handler agent.AgentService, opts ...server.Option) server.Invoker {
	var options []server.Option

	options = append(options, opts...)

	s := server.NewInvoker(options...)
	if err := s.RegisterService(serviceInfo(), handler); err != nil {
		panic(err)
	}
	if err := s.Init(); err != nil {
		panic(err)
	}
	pretouch()
	return s
}
