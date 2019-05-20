package main

import (
	"context"
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"thriftserver/gen-go/echo"
)

type EchoServer struct {
}

func (e *EchoServer) Echo(_ context.Context, req *echo.EchoReq, name string) (*echo.EchoRes, error) {
	fmt.Printf("message from client: %v\n ", req.GetMsg())
	fmt.Println("name",name)
	res := &echo.EchoRes{
		Msg: "super name hello world ,this is a long word,this is a long word,this is a long word",
	}

	return res, nil
}

func main() {
	transport, err := thrift.NewTServerSocket(":9898")
	if err != nil {
		panic(err)
	}

	handler := &EchoServer{}
	processor := echo.NewEchoProcessor(handler)

	transportFactory := thrift.NewTBufferedTransportFactory(8192)
	protocolFactory := thrift.NewTCompactProtocolFactory()
	server := thrift.NewTSimpleServer4(
		processor,
		transport,
		transportFactory,
		protocolFactory,
	)

	if err := server.Serve(); err != nil {
		panic(err)
	}
}
