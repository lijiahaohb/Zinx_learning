package main

import (
	"fmt"
	"zinx/ziface"
	"zinx/znet"
)

// ping test 自定义路由
type PingRouter struct {
	znet.BaseRouter
}

// Test PreHandler
func (p *PingRouter) PreHandler(request ziface.IRequest) {
	fmt.Println("Call Router PreHandler")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping ...\n"))
	if err != nil {
		fmt.Println("call back before ping error")
	}
}

// Test Handler
func (p *PingRouter) Handler(request ziface.IRequest) {
	fmt.Println("Call Router Handler")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping ... ping ... ping ...\n"))
	if err != nil {
		fmt.Println("call back ping ping ping error")
	}
}

// Test PostHandler
func (p *PingRouter) PostHandler(request ziface.IRequest) {
	fmt.Println("Call Router PostHandler")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("post ping ...\n"))
	if err != nil {
		fmt.Println("call back post ping error")
	}
}

func main() {
	// 1. 创建一个server句柄
	server := znet.NewServer("[zinx v0.1]")

	// 2. 给当前server添加一个router
	server.AddRouter(&PingRouter{})

	// 3. 启动运行server
	server.Serve()
}
