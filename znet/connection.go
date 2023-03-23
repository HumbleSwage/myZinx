package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"v1/utils"
	"v1/ziface"
)

// Connection 定义连接模块类
type Connection struct {
	// 当前Connection隶属于哪个Server
	TcpServer ziface.IServer
	// 定义当前连接的socket TCP套接字
	Conn *net.TCPConn
	// 连接ID
	ConnId uint32
	// 连接状态
	IsClosed bool
	// 告诉连接已经退出/停止 channel (由Reader告知Writer退出)
	ExitChan chan bool
	// 无缓冲管道，用于读写goroutine之间的通信
	msgChan chan []byte
	// 消息的管理模块和对应的处理业务API关系
	MsgHandle ziface.IMsgHandle
	// 连接属性集合
	Property map[string]interface{}
	// 保护连接属性集合的锁
	PLock sync.RWMutex
}

// NewConnect 初始化连接模块的方法
func NewConnect(server ziface.IServer, conn *net.TCPConn, connId uint32, handle ziface.IMsgHandle) *Connection {
	connection := &Connection{
		TcpServer: server,
		Conn:      conn,
		ConnId:    connId,
		IsClosed:  false,
		ExitChan:  make(chan bool, 1),
		msgChan:   make(chan []byte),
		MsgHandle: handle,
		Property:  make(map[string]interface{}),
	}
	// 将connect加入到ConnManger中:这里未来可以将Server和Connection进行对应
	connection.TcpServer.GetConnMgr().Add(connection)
	return connection
}

// StartWriter 写消息的goroutine，专门发送客户端消息的模块
func (c *Connection) StartWriter() {
	fmt.Println("[Writer goroutine is running]")
	defer fmt.Printf("%v [conn Writer exit]", c.GetRemoteAddr().String())

	// 不断的阻塞等待channel消息
	for {
		select {
		case data := <-c.msgChan:
			// 有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Failed to send data:", err)
				// TODO：这个地方用continue会不会更好一些
				return
			}
		case <-c.ExitChan:
			// 代表Reader退出，此时Writer也要退出
			return
		}
	}
}

func (c *Connection) Start() {
	fmt.Println("Conn Start()...ConnID=", c.ConnId)
	//TODO 启动从当前连接的读数据的业务
	go c.StartReader()
	//TODO 启动从当前连接的写数据的业务
	go c.StartWriter()

	// 按照开发者传递进来的 创建连接之后需要调用的处理业务，执行对应Hook函数
	c.TcpServer.CallOnConnStart(c)
}

func (c *Connection) Stop() {
	fmt.Println("Conn Stop()...ConnID=", c.ConnId)

	// 如果当前连接已经完毕
	if c.IsClosed == true {
		return
	}
	c.IsClosed = true
	// 调用开发者注册的 销毁连接之前 需要执行业务Hook函数

	c.TcpServer.CallOnConnStop(c)
	// 关闭socket的连接
	err := c.Conn.Close()
	if err != nil {
		return
	}
	// 告知Writer关闭
	c.ExitChan <- true
	// 将当前连接从ConnMgr中摘除掉
	c.TcpServer.GetConnMgr().Remove(c)

	// 关闭管道
	close(c.ExitChan)
	close(c.msgChan)
}

func (c *Connection) GetTcpConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnId
}

func (c *Connection) GetRemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.IsClosed {
		return errors.New("connection closed when send msg")
	}
	// 将data进行封包
	dp := NewDataPack()
	msg := NewMsgPackage(msgId, data)
	binaryMsg, err := dp.Pack(msg)
	if err != nil {
		fmt.Println("Failed to pack:", err)
		return err
	}
	// 将数据发送Writer
	c.msgChan <- binaryMsg

	return nil
}

func (c *Connection) StartReader() {
	fmt.Println("[StartReader goroutine is running]")
	defer fmt.Println("[Reader is exit!] connId", c.ConnId, "remote addr is", c.GetRemoteAddr())
	defer c.Stop()
	for {
		//buf := make([]byte, utils.GetGlobalObj().MaxPackageSize)
		//_, err := c.Conn.StartReader(buf) // 阻塞
		//if err != nil && err != io.EOF {
		//	fmt.Println("Failed to read:", err)
		//	continue
		//}

		// 创建一个拆包、解包的对象
		dp := NewDataPack()

		// 读取客户端的Msg Head 二进制流 8个字节
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTcpConnection(), headData); err != nil {
			fmt.Println("Failed to read head:", err)
			break
		}
		//  拆包，得到msgID 和 msgDataLen 放在msg消息中
		msg, err := dp.UnPack(headData)
		// 根据DataLen 再次读取Data 放在msg.Data中
		if err != nil {
			fmt.Println("Failed to unpack:", err)
		}
		if msg.GetMsgLen() > 0 {
			data := make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTcpConnection(), data); err != nil {
				fmt.Println("Failed to get data when unpack:", err)
				break
			}
			msg.SetData(data)
		}
		// 获取当前conn数据的Request
		request := &Request{
			conn: c,
			msg:  msg,
		}
		// 做一个判断
		if utils.GetGlobalObj().WorkerPoolSize > 0 {
			// 已经开启了工作池机制，将消息发送给Worker工作池机制处理即可
			c.MsgHandle.SendMsgToTaskQueue(request)
		} else {
			// 如果没有开启WorkerPool
			// 那么针对每一个请求都开启一个协程序
			// 执行注册的路由的方法
			// 根据绑定好的MsgId找到对应处理api业务执行
			go c.MsgHandle.DoMsgHandler(request)
		}
	}
}

func (c *Connection) SetProperty(key string, data []byte) {
	c.PLock.Lock()
	defer c.PLock.Unlock()
	// 添加一个连接属性
	c.Property[key] = data

}

func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.PLock.Lock()
	defer c.PLock.Unlock()
	if value, ok := c.Property[key]; ok {
		return value, nil
	}
	return nil, fmt.Errorf("there is no key=%v in property", key)
}

func (c *Connection) RemoveProperty(key string) {
	c.PLock.Lock()
	defer c.PLock.Unlock()

	delete(c.Property, key)
}
