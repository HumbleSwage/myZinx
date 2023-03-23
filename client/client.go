package main

import (
	"fmt"
	"io"
	"net"
	"time"
	"v1/znet"
)

func main() {
	// 建议连接
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("Failed to dail:", err)
		return
	}

	for i := 0; i < 100; i++ {
		var msgId uint32
		if i%2 == 0 {
			msgId = 1
		} else {
			msgId = 2
		}
		// 发送封包的message消息
		dp := znet.NewDataPack()
		binaryMsg, err := dp.Pack(znet.NewMsgPackage(msgId, []byte("Zinx client test message")))
		if err != nil {
			fmt.Println("Failed to pack:", err)
			return
		}

		if _, err := conn.Write(binaryMsg); err != nil {
			fmt.Println("Failed to write binaryMsg:", err)
			return
		}

		// 服务器应该给我们响应
		// 先读取流中的head部分 得到Id 和 DataLen
		buffer := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, buffer); err != nil {
			fmt.Println("Failed to StartReader:", err)
			return
		}
		msg, err := dp.UnPack(buffer)
		if err != nil {
			fmt.Println("Failed to unpack")
			return
		}
		if msg.GetMsgLen() > 0 {
			// 再根据DataLen进行读取Data
			msgBuffer := make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(conn, msgBuffer); err != nil {
				fmt.Println("Failed to StartReader:", err)
				return
			}
			msg.SetData(msgBuffer)
			fmt.Println("MsgId=", msg.GetMsgID(), "	MsgData=", string(msg.GetData()))
		}
		// cpu阻塞
		time.Sleep(1 * time.Second)
	}
}
