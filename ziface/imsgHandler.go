package ziface

// 消息管理模块

type IMsgHandler interface {
	// 调度消息ID对应的处理方法
	DoMsgHandler(request IRequest)

	// 为消息添加具体的处理逻辑
	AddRouter(msgID uint32, router IRouter)

	// 启动Worker工作池
	StartWorkerPool() 

	// 将消息 发送到 TaskQueue
	SendMsgToTaskQueue(Request IRequest)
}
