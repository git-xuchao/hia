package main

import (
	"fmt"
	"log"
	"net"
)

func Handle_conn(conn net.Conn) { //这个是在处理客户端会阻塞的代码。
	buf := make([]byte, 10) //定义一个切片的长度是1024。
	n, _ := conn.Read(buf)  //接收到的内容大小。
	fmt.Print(string(buf[:n]))
	conn.Close() //与客户端断开连接。
}

func main() {
	addr := "0.0.0.0:12345" //表示监听本地所有ip的8080端口，也可以这样写：addr := ":8080"
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept() //用conn接收链接
		if err != nil {
			log.Fatal(err)
		}
		go Handle_conn(conn) //开启多个协程。
	}
}
