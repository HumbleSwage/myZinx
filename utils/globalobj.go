package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"v1/ziface"
)

// GlobalObj 存储一切有关Zinx框架参数，供其他模块使用
type GlobalObj struct {
	// Server
	TcpServer ziface.IServer `json:"tcpServer"` //当前Zinx全局的server对象
	Host      string         `json:"host"`      // 当前服务器主机监听的IP
	Port      int            `json:"port"`      // 当前服务器监听的端口
	Name      string         `json:"name"`      // 当前服务器的名称

	// Zinx
	Version        string `json:"version"`        // 当前Zinx的版本号
	MaxConn        int    `json:"maxConn"`        // 当前服务器主机的最大连接数量
	MaxPackageSize uint32 `json:"maxPackageSize"` // 当前Zinx框架数据包的最大值
	WorkerPoolSize uint32 `json:"workerPoolSize"` // 当前业务工作Worker池的goroutine
	QueueSize      uint32 `json:"queueSize"`      // Zinx框架允许用户最多开辟多少个Worker(限定条件)
}

// globalObj 定义一个全局对外的GlobalObj
var globalObj *GlobalObj

// GetGlobalObj 向外暴露globalObj
func GetGlobalObj() *GlobalObj {
	return globalObj
}

// 解析配置文件
func init() {
	// 如果没有加载配置文件
	globalObj = &GlobalObj{
		Name:           "ZinxServerApp",
		Version:        "V0.4",
		Port:           8999,
		Host:           "0.0.0.0",
		MaxConn:        1000,
		MaxPackageSize: 4096,
		WorkerPoolSize: 10,
		QueueSize:      1024,
	}
	// 从conf/zinx.json加载
	globalObj.Reload()
}

// Reload 解析配置文件
func (g *GlobalObj) Reload() {
	data, err := os.ReadFile("conf/zinx.json")
	if err != nil {
		fmt.Println("Failed to load conf:", err)
		return
	}
	if err = json.Unmarshal(data, g); err != nil {
		fmt.Println(err)
		panic("Failed to unmarshal")
		return
	}
}
