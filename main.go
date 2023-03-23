package main

import (
	"fmt"
	"v1/ziface"
	"v1/znet"
)

type Router struct {
	znet.BaseRouter
}

type Handler1 struct {
	znet.BaseRouter
}

func (h1 *Handler1) Handle(request ziface.IRequest) {
	fmt.Println("Receive from client:msgId=", request.GetMsgID(),
		" data=", string(request.GetData()))
	if err := request.GetConnection().SendMsg(1, []byte("ping...ping...ping")); err != nil {
		fmt.Println("Failed to sendMsg:", err)
		return
	}
}

type Handler2 struct {
	znet.BaseRouter
}

func (h2 *Handler2) Handle(request ziface.IRequest) {
	fmt.Println("Receive from client:msgId=", request.GetMsgID(),
		" data=", string(request.GetData()))
	if err := request.GetConnection().SendMsg(2, []byte("ha....ha....ha")); err != nil {
		fmt.Println("Failed to sendMsg:", err)
		return
	}
}

func Before(conn ziface.IConnection) {
	fmt.Println("DoConnectionStart is Called!")
	err := conn.SendMsg(3, []byte("Start Hook！"))
	if err != nil {
		return
	}
	// 设置连接的属性:给当前的连接设置一些属性
	fmt.Println("Set conn property....")
	conn.SetProperty("Name", []byte("d_zhao"))
	conn.SetProperty("Github", []byte("https://github.com/HumbleSwage"))
}

func After(conn ziface.IConnection) {
	fmt.Println("DoConnectionStop is Called!")
	// 获取连接属性
	if name, err := conn.GetProperty("Name"); err == nil {
		fmt.Printf("Name=%s\n", name)
	}
	if web, err := conn.GetProperty("Github"); err == nil {
		fmt.Printf("Github=%s\n", web)
	}

}

// 基于zinx框架开发服务端的应用
func main() {
	// 创建一个server句柄，使用Zinx的Api
	s := znet.NewServer()
	s.SetOnConnStar(Before)
	s.SetOnConnStop(After)
	// 自定义router
	h1 := &Handler1{}
	s.AddRouter(1, h1)
	h2 := &Handler2{}
	s.AddRouter(2, h2)
	// 启动server
	s.Serve()

}
