package znet

import (
	"fmt"
	"sync"
	"v1/ziface"
)

// ConnManager 连接管理模块
type ConnManager struct {
	connections map[uint32]ziface.IConnection // 管理的连接集合
	connLock    *sync.RWMutex                 // 读写连接集合的读写锁
}

// NewConnManager 创建当前连接的方法
func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
		connLock:    new(sync.RWMutex),
	}
}

func (cm *ConnManager) Add(conn ziface.IConnection) {
	// 保护共享资源map，加写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	// 将conn加入
	cm.connections[conn.GetConnID()] = conn
	fmt.Println("ConnId=", conn.GetConnID(), " add to ConnManager successfully")
}

func (cm *ConnManager) Remove(conn ziface.IConnection) {
	// 保护共享资源map，加写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	// 将conn移除
	delete(cm.connections, conn.GetConnID())
	fmt.Println("ConnId=", conn.GetConnID(), " remove from ConnManager successfully")
}

func (cm *ConnManager) Get(connId uint32) (ziface.IConnection, error) {
	// 保护共享资源map，加写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	conn, ok := cm.connections[connId]
	if !ok {
		return nil, fmt.Errorf("there is no ConnId=%v in ConnManage", connId)
	}
	return conn, nil
}

func (cm *ConnManager) Count() int {
	return len(cm.connections)
}

func (cm *ConnManager) ClearConn() {
	// 保护共享资源map，加写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	// 删除conn并停止conn的工作
	for connId, conn := range cm.connections {
		// 停止
		conn.Stop()
		// 删除
		delete(cm.connections, connId)

	}
	fmt.Println("Clear all connection successfully.")
}
