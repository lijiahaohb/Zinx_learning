package znet

import (
	"fmt"
	"zinx/utils"
	"zinx/ziface"
)

type MsgHandler struct {
	// 存放每个 msgID 与 router的映射关系
	Apis map[uint32]ziface.IRouter

	// 负责worker 读取任务的消息队列
	TaskQueue []chan ziface.IRequest
	// 业务工作 worker池的worker数量
	WorkerPoolSize uint32
}

// 创建msgHandler 的方法
func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}

func (m *MsgHandler) DoMsgHandler(request ziface.IRequest) {
	// 从request中找到msgID
	handler, ok := m.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgID = ", request.GetMsgID(), "not found, please register one")
		return
	}
	handler.PreHandler(request)
	handler.Handler(request)
	handler.PostHandler(request)
}

// 为消息添加具体的处理逻辑
func (m *MsgHandler) AddRouter(msgID uint32, router ziface.IRouter) {
	// 判断当前msgID 对应的 api方法是否存在
	if _, ok := m.Apis[msgID]; ok {
		fmt.Println("msgID have registered api")
		return
	}
	m.Apis[msgID] = router
	fmt.Println("add api msgID = ", msgID, "successfully")
}

// 启动一个Worker工作池
func (m *MsgHandler) StartWorkerPool() {
	// 根据workerPoolSize 开启 worker, 每个worker用一个go来承载
	for i := 0; i < int(m.WorkerPoolSize); i++ {
		// 1. 给worker 对应的 channel 消息队列开辟空间
		m.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)

		// 2. 启动当前的worker， 阻塞等待消息从channel传递进来
		go m.startWorker(i, m.TaskQueue[i])
	}
}

// 启动一个worker 工作流程
func (m *MsgHandler) startWorker(workerID int, taskQueue chan ziface.IRequest) {
	fmt.Println("Worker ID = ", workerID, "is starting ...")

	// 不断阻塞等待 worker 对应的消息队列 的消息到来
	for {
		select {
		case request := <-taskQueue:
			m.DoMsgHandler(request)
		}
	}
}

// 将消息交给 taskQueue，由Worker进行处理
func (m *MsgHandler) SendMsgToTaskQueue(request ziface.IRequest) {
	// 将消息平均分配给不同 worker 对应的 taskQueue
	// 根据客户端建立的ConnID来进行分配
	workerID := request.GetConnection().GetConnectionID() % m.WorkerPoolSize
	fmt.Println("Add connID = ", request.GetConnection().GetConnectionID(),
				"request MsgID = ", request.GetMsgID(), "to workerID = ", workerID)

	m.TaskQueue[workerID] <- request
}
