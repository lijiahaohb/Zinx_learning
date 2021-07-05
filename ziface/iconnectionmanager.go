package ziface

// 连接管理模块

type IConnManager interface {
	// 添加连接
	Add(IConnection)

	// 删除连接
	Remove(IConnection)

	// 根据connID获取连接
	Get(connID	uint32) (IConnection, error)

	// 得到连接总数
	Len() int

	// 清除所有连接
	ClearConn()
}
