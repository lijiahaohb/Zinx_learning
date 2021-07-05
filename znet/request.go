package znet

import (
	"zinx/ziface"
)

type Request struct {
	// 和客户端建立好的连接
	conn ziface.IConnection
	// 客户端的请求数据
	msg ziface.IMessage
}

// 获取当前连接
func (r *Request) GetConnection() ziface.IConnection {
	return r.conn
}

// 获取客户端请求数据
func (r *Request) GetData() []byte {
	return r.msg.GetMsg()
}

// 获取当前请求消息的ID
func (r *Request) GetMsgID() uint32 {
	return r.msg.GetMsgID()
}
