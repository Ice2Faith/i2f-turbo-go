package main

import (
	"fmt"
	"rpcserver/gorpc"
)

type HelloRpcService struct {
}

func (svc *HelloRpcService) EchoBack(req string, resp *string) error {
	fmt.Println("rpc recv:", req)
	*resp = "echo:" + req
	return nil
}

func main() {
	server := gorpc.GetDefaultRpcServer()
	server.AddService("hello", &HelloRpcService{})
	wg := server.Run()
	wg.Wait()
	server.Stop()
}
