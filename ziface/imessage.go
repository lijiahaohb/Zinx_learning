package ziface

// 将请求消息封装到 Message 中，定义抽象接口

type IMessage interface {
	GetMsgID() uint32 
	GetMsgLen() uint32
	GetMsg() []byte

	SetMsgID(uint32)
	SetData([]byte)
	SetDataLen(uint32)
}
