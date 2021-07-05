package ziface 

// 定义一个服务器接口
type IServer interface {
	// 启动服务器方法
	Start()
	// 停止服务器方法
	Stop()
	// 运行服务器方法
	Serve()
	// 添加路由： 给当前的服务注册一个路由方法 来处理客户端的连接
	AddRouter(uint32, IRouter)
	// 获取 连接管理器 的方法
	GetConnMgr() IConnManager
	// 设置钩子函数的方法
	SetOnConnStart(func(IConnection))
	SetOnConnStop(func(IConnection))
	// 调用钩子函数的方法
	CallOnConnStart(IConnection)
	CallOnConnStop(IConnection)
}
