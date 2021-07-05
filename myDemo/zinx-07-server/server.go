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

// 创建连接建立之后的执行的钩子函数
func DoConnectionBegin(conn ziface.IConnection)  {
	fmt.Println(">>>>>> DoConenctionBegin() is called ... <<<<<<")
	if err := conn.SendMsg(1, []byte("DoConnection BEGIN")); err != nil {
		fmt.Println(err)
	}

	// 给连接设置一些属性
	conn.SetProperty("Name", "lijiahao")
}

// 创建断开连接之前要执行的钩子函数
func DoConnectionStop(conn ziface.IConnection) {
	fmt.Println(">>>>>> DoConenctionStop() is called ... <<<<<<")
	fmt.Println("conn ID = ", conn.GetConnectionID(), "is lost ...")

	// 获取连接属性
	if value, err := conn.GetProperty("Name"); err != nil {
		fmt.Println("name = ", value)
	}
}

func main() {
	// 1. 创建一个server句柄
	server := znet.NewServer("[zinx v0.4]")

	// 2. 注册与连接相关的钩子函数
	server.SetOnConnStart(DoConnectionBegin)
	server.SetOnConnStop(DoConnectionStop)

	// 3. 给当前server添加一个router
	server.AddRouter(0, &PingRouter{})
	server.AddRouter(1, &HelloZinxRouter{})

	// 4. 启动运行server
	server.Serve()
}
