package znet

import (
	"fmt"
	"strconv"
	"v1/utils"
	"v1/ziface"
)

// MsgHandle 消息处理模块的实现
type MsgHandle struct {
	// 消息ID和对应Handle映射
	Apis map[uint32]ziface.IRouter
	// 负责Worker消费的消息队列
	TaskQueue []chan ziface.IRequest
	// 业务工作Worker池的数量
	WorkerPoolSize uint32
}

// NewMsgHandle 创建MsgHandle方法
func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]ziface.IRouter),
		TaskQueue:      make([]chan ziface.IRequest, utils.GetGlobalObj().WorkerPoolSize), // 配置读取
		WorkerPoolSize: utils.GetGlobalObj().WorkerPoolSize,                               // 配置读取
	}
}

// DoMsgHandler 调度
func (mh *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	// 1 从Request中找到msgId
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Printf("Failed to load handler to msgId[%d\n]", request.GetMsgID())
		return
	}
	// 2 根据msgId调度对应的router
	handler.PreHandle(request)
	handler.Handle(request)
	handler.AfterHandle(request)
}

// AddRouter 添加
func (mh *MsgHandle) AddRouter(msgId uint32, router ziface.IRouter) {
	// 1 判断当前Msg绑定的Api处理方法是否已经存在
	if _, ok := mh.Apis[msgId]; ok {
		// 表示该Id已经注册了
		panic("repeat api,msgId=" + strconv.Itoa(int(msgId)))
	}
	// 2 添加Msg与API的绑定关系
	mh.Apis[msgId] = router
	fmt.Println("Add api MsgId=", msgId)
}

// StartWorkPool 启动一个Worker工作池:开启工作池的动作只能开启一次，一个zinx框架只能有一个worker工作池
func (mh *MsgHandle) StartWorkPool() {
	// 根据WorkerPoolSize分别开启Worker，每个Worker用一个go来承载
	for i := 0; i < int(utils.GetGlobalObj().WorkerPoolSize); i++ {
		// 一个worker被启动
		// 1 当前的worker对应的channel消息队列开辟空间 第0个worker就用第0个channel
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GetGlobalObj().QueueSize)
		// 2 启动当前的worker,阻塞等待消息从channel传递进来
		go mh.StartOneWorker(i)
	}
}

// StartOneWorker 启动一个Worker工作池
func (mh *MsgHandle) StartOneWorker(workerId int) {
	fmt.Println("Worker ID=", workerId, "is started...")
	// 不断的阻塞等待对应消息队列的消息
	for {
		select {
		// 如果有消息，出列的就是一个客户端的request，执行当前Request所绑定的业务
		case request := <-mh.TaskQueue[workerId]:
			mh.DoMsgHandler(request)
		}
	}
}

// SendMsgToTaskQueue 将消息发送给消息队列工作池的方法
func (mh *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	// 1 将消息平均分配给不通过的worker：负载均衡
	// 根据客户端建立的ConnId来进行分配
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Println("Add ConnID=", request.GetConnection().GetConnID(),
		" request MsgId=", request.GetMsgID(),
		" to WorkerID=", workerID)

	// 2 将消息发送给对应的worker的TaskQueue即可
	mh.TaskQueue[workerID] <- request
}
