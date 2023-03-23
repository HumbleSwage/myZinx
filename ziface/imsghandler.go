package ziface

// IMsgHandle 消息管理抽象层
type IMsgHandle interface {
	// DoMsgHandler 调度对应Router消息处理方法
	DoMsgHandler(IRequest)
	// AddRouter 为消息添加具体的处理逻辑
	AddRouter(uint32, IRouter)
	// StartWorkPool 启动Worker工作池
	StartWorkPool()
	// SendMsgToTaskQueue 将消息发送给消息队列工作池的方法
	SendMsgToTaskQueue(IRequest)
}
