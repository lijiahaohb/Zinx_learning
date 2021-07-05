package ziface

// 实际上是把客户端的 连接 和 请求数据 包装到了一个Request中

type IRequest interface {
	// 得到当前连接
	GetConnection() IConnection

	// 得到请求数据
	GetData() []byte

	GetMsgID() uint32
}
