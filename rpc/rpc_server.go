package myrpc

import (
	"fmt"

	"hia/rpc/service"
)

type RpcServer interface {
	Set(key, value string) (err error)
	AddService(service interface{}) (err error)
	Start() (err error)
}

func RpcServerTest() {
	var serv = new(service.MathService)
	var tcp_server = new(TcpRpcServer)
	var server RpcServer

	server = tcp_server
	fmt.Println("start TcpRpcServer...")

	server.Set("net", "tcp")
	server.Set("addr", "127.0.0.1:1234")

	server.AddService(serv)

	server.Start()
}
