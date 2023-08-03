package service

import (
	"context"
	example "cwgo/example/hex/kitex_gen/hello/example"
	"testing"
)

func TestHelloMethod_Run(t *testing.T) {
	ctx := context.Background()
	s := NewHelloMethodService(ctx)
	// init req and assert value

	request := &example.HelloReq{}
	resp, err := s.Run(request)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Errorf("unexpected nil response")
	}
	// todo: edit your unit test

}
