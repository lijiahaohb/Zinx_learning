package znet

import (
	"fmt"
	"net"
	"zinx/utils"
	"zinx/ziface"
)

// IServer接口的具体实现，定义一个Server的服务器模块
type Server struct {
	// 服务器名称
	Name string
	// 服务器绑定的IP版本
	IPVersion string
	// 服务器IP
	IP string
	// 服务器端口
	Port int
	// 当前Server的消息管理模块 用来绑定msgID和对应业务API
	MsgHandler ziface.IMsgHandler
	// 连接管理模块
	ConnMgr ziface.IConnManager
	// 该Server创建连接之后自动调用的钩子函数
	OnConnStart func(ziface.IConnection)
	// 该Server销毁连接之前自动调用的钩子函数
	OnConnStop func(ziface.IConnection)
}

func NewServer(name string) ziface.IServer {
	server := &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		MsgHandler: NewMsgHandler(),
		ConnMgr:    NewConnManager(),
	}

	return server
}

func (s *Server) Start() {
	fmt.Printf("[Zinx] Server Name: %s, listener at IP: %s, Port: %d is starting\n",
		utils.GlobalObject.Name, utils.GlobalObject.Host, utils.GlobalObject.TcpPort)
	// 用一个goroutine来承载业务，开启多个goroutine阻塞等待监听
	go func() {
		// 0. 开启 消息队列 和 worker 工作池
		s.MsgHandler.StartWorkerPool()

		// 1. 获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error: ", err)
			return
		}

		// 2. 监听服务器地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen ", s.IPVersion, "error: ", err)
			return
		}
		defer listener.Close()

		fmt.Println("start Zinx server successfully", s.Name, "is listening...")
		var cid uint32 = 0

		// 3. 阻塞监听
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept error: ", err)
				continue
			}

			// 判断 当前连接数 是否等于 最大连接数 如果等于 则关闭新连接
			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				fmt.Println("Too many connections MaxConn = ", utils.GlobalObject.MaxConn)
				conn.Close()
				continue
			}

			// 将处理新连接的业务方法和conn相绑定 得到我们的连接模块
			dealConn := NewConnection(s, conn, cid, s.MsgHandler)
			cid++

			// 启动连接模块
			go dealConn.Start()
		}
	}()

}

func (s *Server) Stop() {
	// TODO 将服务器资源、状态、连接进行回收
	fmt.Println("[Zinx server] is stoping ...")
	s.ConnMgr.ClearConn()
}

func (s *Server) Serve() {
	// 启动server的服务功能
	s.Start()

	// TODO 做一些启动服务器的额外业务

	// 阻塞状态
	select {}

}

func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)

	fmt.Println("Add Router Successfully")
}

func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

// 注册钩子函数的方法
func (s *Server) SetOnConnStart(hookFunc func(ziface.IConnection)) {
	s.OnConnStart = hookFunc
}

func (s *Server) SetOnConnStop(hookFunc func(ziface.IConnection)) {
	s.OnConnStop = hookFunc
}

// 调用钩子函数的方法
func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println(">>>>>> OnConnStart() is called ... <<<<<<")
		s.OnConnStart(conn)
	}
}

func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println(">>>>>> OnConnStop() is called ... <<<<<<")
		s.OnConnStop(conn)
	}
}
