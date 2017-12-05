package myrpc

import (
	"testing"
)

/*
 *import (
 *    "fmt"
 *    "net"
 *    "net/rpc"
 *    "testing"
 *
 *    "hia/rpc/service"
 *)
 */

/*
 *func TestServer(t *testing.T) {
 *    var ms = new(service.MathService) //实例化服务对像
 *    rpc.Register(ms)                  //注册这个服务
 *    fmt.Println("启动服务...")
 *    var address, _ = net.ResolveTCPAddr("tcp", "127.0.0.1:1234") //定义TCP的服务承载地址
 *    listener, err := net.ListenTCP("tcp", address)               //监听TCP连接
 *    if err != nil {
 *        fmt.Println("启动失败！", err)
 *    }
 *    for {
 *        conn, err := listener.Accept() //如果接受到连接
 *        if err != nil {
 *            continue
 *        }
 *        fmt.Println("接收到一个调用请求...")
 *        rpc.ServeConn(conn) //让此rpc绑定到该Tcp连接上。
 *    }
 *}
 */

func TestServer(t *testing.T) {
	/*
	 *TcpRpcServerTest()
	 */
	RpcServerTest()
}
