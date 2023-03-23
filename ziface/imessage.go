package ziface

// IMessage 将请求的消息封装到Message中，定义抽象接口
type IMessage interface {
	// GetMsgID 获取消息的Id
	GetMsgID() uint32
	// GetMsgLen 获取消息的长度
	GetMsgLen() uint32
	// GetData 获取消息的内容
	GetData() []byte
	// SetMsgId 设置消息的Id
	SetMsgId(id uint32)
	// SetData 设置消息的内容
	SetData(data []byte)
	// SetMsgLen 设置消息的长度
	SetMsgLen(len uint32)
}
