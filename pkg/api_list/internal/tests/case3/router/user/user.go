package user

import "github.com/cloudwego/hertz/pkg/route"

func InitUserRoutes(g *route.RouterGroup) {
	g.GET("/info", nil)
	g.POST("/nickname", nil)
}
