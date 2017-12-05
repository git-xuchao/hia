package myrpc

import (
	"fmt"
	"net"
	"net/rpc"

	"hia/rpc/service"
)

type TcpRpcServer struct {
	net, addr string
}

func (self *TcpRpcServer) Set(key, value string) (err error) {
	switch key {
	case "net":
		self.net = value
	case "addr":
		self.addr = value
	default:
		return nil
	}

	return nil
}

func (self *TcpRpcServer) AddService(service interface{}) (err error) {
	rpc.Register(service)
	return nil
}

func (self *TcpRpcServer) Start() (err error) {
	var address, _ = net.ResolveTCPAddr(self.net, self.addr) //定义TCP的服务承载地址
	listener, err := net.ListenTCP(self.net, address)        //监听TCP连接
	if err != nil {
		fmt.Println("start failed！", err)
	}
	for {
		conn, err := listener.Accept() //如果接受到连接
		if err != nil {
			continue
		}
		fmt.Println("rcv a call of service...")
		rpc.ServeConn(conn) //让此rpc绑定到该Tcp连接上。
	}

	return nil
}

func TcpRpcServerTest() {
	var serv = new(service.MathService)
	var server = new(TcpRpcServer)

	fmt.Println("start TcpRpcServer...")

	server.Set("net", "tcp")
	server.Set("addr", "127.0.0.1:1234")

	server.AddService(serv)

	server.Start()
}
