package znet

import (
	"fmt"
	"net"
	"v1/utils"
	"v1/ziface"
)

// Server IServer接口的实现，定义一个Server服务器模块
type Server struct {
	// 服务器名称
	Name string
	// 服务器绑定的IP版本
	IPVersion string
	// 服务器监听的IP地址
	IP string
	// 服务器监听的Port
	Port int
	// 消管理模块
	MsgHandle ziface.IMsgHandle
	// 该server的连接管理器
	ConnManager ziface.IConnManager
	// 该Server创建连接之后自动调用Hook函数-OnConnStart
	OnConnStart func(ziface.IConnection)
	// 该Server销毁连接之前自动调用Hook函数-OnConnStop
	OnConnStop func(ziface.IConnection)
}

func (s *Server) Start() {
	// TODO 以后这里都是要用一个日志模块里面的
	fmt.Printf("[Zinx]Configuration load successfully:%s %s:%d\n", s.Name, s.IP, utils.GetGlobalObj().Port)
	fmt.Printf("[Zinx]Version:%s\n[Zinx]MaxConn:%d\n[Zinx]MaxPackSize:%d\n", utils.GetGlobalObj().Version, utils.GetGlobalObj().MaxConn, utils.GetGlobalObj().MaxPackageSize)
	// 监听一个TCP地址
	TCPAddr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		fmt.Println("Failed to resolve tcp addr:", err)
		return
	}
	// 异步处理请求
	go func() {
		// 开启消息队列以及Worker工作池
		s.MsgHandle.StartWorkPool()

		// 对解析的地址进行监听
		TCPListener, err := net.ListenTCP(s.IPVersion, TCPAddr)
		if err != nil {
			fmt.Println("Failed to listen a tcp addr:", err)
			return
		}
		fmt.Printf("Congratulates!Zinx is running now.......\n")
		var cid uint32 = 1
		for {
			// 使用Listener进行Accept
			conn, err := TCPListener.AcceptTCP()
			if err != nil {
				fmt.Println("Failed to accept a tcp conn:", err)
				continue
			}

			// 设置最大连接个数的判断，如果超过最大连接，那么则关闭此新的连接
			if s.ConnManager.Count() >= utils.GetGlobalObj().MaxConn {
				// TODO 给响应的客户端超出一个最大连接的错误包
				fmt.Println("too many connection MaxConn=", utils.GetGlobalObj().MaxConn)

				err := conn.Close()
				if err != nil {
					return
				}
				continue
			}

			// 客户端已经建立完毕，进行业务处理
			// 注意这个的handleFunc，以后是由用户自定义
			dealConn := NewConnect(s, conn, cid, s.MsgHandle)
			cid++

			// 尝试启动连接
			go dealConn.Start()
		}
	}()

}

func (s *Server) Stop() {
	// 将一些服务器的资源、状态或者一些已经开辟的链接信息进行停止或者回收
	fmt.Println("[STOP] Zinx server name ", s.Name)
	s.ConnManager.ClearConn()
}

func (s *Server) Serve() {
	s.Start()

	// TODO 可以做一些启动服务之后的额外事务

	// 阻塞状态
	// 注意这里在Server进行阻塞有一定的巧妙性
	select {}
}

func (s *Server) AddRouter(msgId uint32, router ziface.IRouter) {
	s.MsgHandle.AddRouter(msgId, router)
	fmt.Println("Add router success!")
}

// NewServer 初始化Server的方法
func NewServer() ziface.IServer {
	conf := utils.GetGlobalObj()
	s := &Server{
		Name:        conf.Name,
		IPVersion:   "tcp4",
		IP:          conf.Host,
		Port:        conf.Port,
		MsgHandle:   NewMsgHandle(),
		ConnManager: NewConnManager(),
	}
	return s
}

// GetConnMgr 返回当前Server创建连接管理
func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnManager
}

// SetOnConnStar 注册OnConnStart钩子函数
func (s *Server) SetOnConnStar(hookFunc func(ziface.IConnection)) {
	s.OnConnStart = hookFunc
}

// SetOnConnStop 注册OnConnStop钩子函数
func (s *Server) SetOnConnStop(hookFunc func(ziface.IConnection)) {
	s.OnConnStop = hookFunc
}

// CallOnConnStart 调用OnConnStart钩子函数
func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("------> Call OnConnStart().....")
		s.OnConnStart(conn)
	}
}

// CallOnConnStop 调用OnConnStop钩子函数
func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("------> Call OnConnStop().....")
		s.OnConnStop(conn)
	}
}
