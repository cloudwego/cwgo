package main

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
	"cwgo/example/hex/biz/router"
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

var httpReg = regexp.MustCompile(`^(?:GET |POST|PUT|DELE|HEAD|OPTI|CONN|TRAC|PATC)$`)

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
	h := hertzServer.New()

	// add a ping route to test
	h.GET("/ping", func(c context.Context, ctx *app.RequestContext) {
		ctx.JSON(consts.StatusOK, utils.H{"ping": "pong"})
	})

	router.GeneratedRegister(h)
	err := h.Engine.Init()
	if err != nil {
		panic(err)
	}
	return h.Engine
}

var hertzEngine *route.Engine

func init() {
	hertzEngine = initHertz()
}

