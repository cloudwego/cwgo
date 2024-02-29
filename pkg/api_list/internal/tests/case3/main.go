package main

import (
	"main/router"

	"github.com/cloudwego/hertz/pkg/app/server"
)

func main() {
	h := server.Default()

	router.InitRoutes(h.Engine)
}
