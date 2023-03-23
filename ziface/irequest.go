package ziface

// IRequest 实际上是把客户端请求的连接数据和请求数据包装到一个Request中
type IRequest interface {
	// GetConnection 得到当前连接
	GetConnection() IConnection

	// GetData 得到请求的消息数据
	GetData() []byte

	// GetMsgID 得到消息的类型
	GetMsgID() uint32
}
