package main

import "zinx/znet"

func main() {
	// 1. 创建一个server句柄
	server := znet.NewServer("[zinx v0.1]")

	// 2. 启动运行server
	server.Serve()
}
