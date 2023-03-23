package ziface

import "net"

// IConnection 定义连接的抽象模块
type IConnection interface {
	// Start 启动连接 让当前连接准备开始工作
	Start()
	// Stop 停止连接 让当前连接结束工作
	Stop()
	// GetTcpConnection 获取当前连接所绑定的socket conn
	GetTcpConnection() *net.TCPConn
	// GetConnID 获取当前连接模块的连接ID
	GetConnID() uint32
	// GetRemoteAddr 获取远程客户端的IP和端口
	GetRemoteAddr() net.Addr
	// SendMsg 发送数据将数据发送给远程的客户端
	SendMsg(uint32, []byte) error
	// StartReader 读连接的数据
	StartReader()
	// SetProperty 设置连接属性
	SetProperty(string, []byte)
	// GetProperty 获取连接属性
	GetProperty(string) (interface{}, error)
	// RemoveProperty 移除连接属性
	RemoveProperty(string)
}
