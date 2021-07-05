package main

import (
	"net"
	"fmt"
	"time"
)

// 模拟客户端
func main() {
	fmt.Println("client start ...")
	time.Sleep(time.Second * 1)
	// 1. 直接连接远程服务器，得到conn连接
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client dial error: ", err)
		return	
	}

	// 2. 写数据
	for {
		_, err := conn.Write([]byte("Hello Zinx v0.1 ..."))
		if err != nil {
			fmt.Println("write error: ", err)
			continue
		}

		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read buf error")
			continue
		}
		fmt.Printf("server call back: %s, cnt = %d\n", string(buf[:cnt]), cnt)

		time.Sleep(1 * time.Second)
	}
}
