package ziface

// IDataPack 定义解决TCP粘包问题的封包和拆包的模块
type IDataPack interface {
	// GetHeadLen 获取包的头长度的方法
	GetHeadLen() uint32
	// Pack 针对Message进行TLV格式的封装
	Pack(IMessage) ([]byte, error)
	// UnPack 针对Message进行TLV格式的拆包
	UnPack([]byte) (IMessage, error)
}
