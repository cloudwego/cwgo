package router

import (
	"main/router/user"

	"github.com/cloudwego/hertz/pkg/route"
)

func InitRoutes(e *route.Engine) {
	e.GET("/ping", nil)

	g := e.Group("/api/v1")

	initDefault(g)
	user.InitUserRoutes(g.Group("/user"))
}

func initDefault(g *route.RouterGroup) {
	g.GET("/help", nil)
}
