package ziface

// IConnManager 连接管理模块抽象层
type IConnManager interface {
	// Add 添加连接
	Add(IConnection)
	// Remove 删除连接
	Remove(IConnection)
	// Get 根据connId获取连接
	Get(uint32) (IConnection, error)
	// Count 得到当前的连接总数
	Count() int
	// ClearConn 清除并终止所有的连接
	ClearConn()
}
