package znet

import (
	"errors"
	"fmt"
	"sync"
	"zinx/ziface"
)

type ConnManager struct {
	// 管理的连接集合
	connections map[uint32]ziface.IConnection
	// 保护连接集合的读写锁
	connLock sync.RWMutex
}

// 创建 连接管理模块 的方法
func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

// 添加连接
func (cm *ConnManager) Add(conn ziface.IConnection) {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	cm.connections[conn.GetConnectionID()] = conn
	fmt.Println("connID = ", conn.GetConnectionID(), "add to ConnManager successfully: conn num = ", cm.Len())
}

// 删除连接
func (cm *ConnManager) Remove(conn ziface.IConnection) {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	delete(cm.connections, conn.GetConnectionID())
	fmt.Println("connID = ", conn.GetConnectionID(), "remove from ConnManager successfully: conn num = ", cm.Len())
}

// 根据connID获取连接
func (cm *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	cm.connLock.RLock()
	defer cm.connLock.RUnlock()

	conn, ok := cm.connections[connID]
	if ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not found")
	}
}

// 得到连接总数
func (cm *ConnManager) Len() int {
	return len(cm.connections)
}

// 清除所有连接
func (cm *ConnManager) ClearConn() {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	// 删除 Connection 并停止 conn 的工作
	for connID, conn := range cm.connections {
		// 停止 conn 工作
		conn.Stop()

		// 从connections中删除conn
		delete(cm.connections, connID)
	}

	fmt.Println("clear all connections successfully, conn num = ", cm.Len())
}
