package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

func main() {
	addr := "127.0.0.1:12345"          //定义主机名
	conn, err := net.Dial("tcp", addr) //拨号操作，需要指定协议。
	if err != nil {
		log.Fatal(err)
	}

	n, err := conn.Write([]byte("hello world")) //向服务端发送数据。用n接受返回的数据大小，用err接受错误信息。
	if err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, n) //定义一个切片的长度是1024。

	for {
		n, err = conn.Read(buf) //接收到的内容大小。
		if err == io.EOF {
			conn.Close()
		}
		fmt.Print(string(buf[:n]))
	}

	fmt.Println(string(buf[:n])) //将接受的内容都读取出来。

}
