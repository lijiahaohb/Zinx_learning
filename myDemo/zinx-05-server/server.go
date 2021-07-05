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

type HelloZinxRouter struct {
	znet.BaseRouter
}

// Test Handler
func (p *PingRouter) Handler(request ziface.IRequest) {
	fmt.Println("Call PingRouter Handler")
	fmt.Println("recv from client msgID = ",  request.GetMsgID(), " msgData = " +  string(request.GetData()))

	err := request.GetConnection().SendMsg(1, []byte("ping ... ping ... ping ..."))
	if err != nil {
		fmt.Println(err)
	}
}

// HelloZinxrouter test
func (h *HelloZinxRouter) Handler(request ziface.IRequest) {
	fmt.Println("Call HelloZinxRouter Handler")
	fmt.Println("recv from client msgID = ",  request.GetMsgID(), " msgData = " +  string(request.GetData()))

	err := request.GetConnection().SendMsg(201, []byte("hello zinx ... helo zinx ..."))
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	// 1. 创建一个server句柄
	server := znet.NewServer("[zinx v0.4]")

	// 2. 给当前server添加一个router
	server.AddRouter(0, &PingRouter{})
	server.AddRouter(1, &HelloZinxRouter{})

	// 3. 启动运行server
	server.Serve()
}
