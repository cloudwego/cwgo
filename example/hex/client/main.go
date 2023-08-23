package main

import (
	"context"
	"fmt"

	"cwgo/example/hex/kitex_gen/hello/example"
	"cwgo/example/hex/kitex_gen/hello/example/helloservice"
	"github.com/cloudwego/kitex/client"
)

func main() {
	kc, err := helloservice.NewClient("p.s.m", client.WithHostPorts("127.0.0.1:8888"))
	if err != nil {
		panic(err)
	}
	req := &example.HelloReq{Name: "hex"}
	resp, err := kc.HelloMethod(context.Background(), req)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}
