package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

// 知识负责测试dataPack拆包和封包
func TestDataPack(t *testing.T) {
	/*
		模拟的服务器
	*/
	listener, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("Failed to listen:", err)
		return
	}
	go func() {
		// Accept
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("Failed to accept:", err)
				// 创建一个拆包对象
				break
			}
			go func(conn net.Conn) {
				dp := NewDataPack()
				for {
					// 第一次从conn读，把包的head读出来
					headData := make([]byte, dp.GetHeadLen())
					if _, err := io.ReadFull(conn, headData); err != nil {
						fmt.Println("Failed to read head:", err)
						return
					}
					msgHead, err := dp.UnPack(headData)
					if err != nil {
						fmt.Println("Failed to unpack:", err)
						return
					}
					if msgHead.GetMsgLen() > 0 {
						// msgHead是有数据的，需要进行第二次读取，即读取data内容
						msg := msgHead.(*Message)
						msg.Data = make([]byte, msg.GetMsgLen())
						// 根据dataLen的长度再次从io流中进行读取
						if _, err := io.ReadFull(conn, msg.Data); err != nil {
							fmt.Println("Failed to read data:", err)
							return
						}
						fmt.Println("——> Receive MsgID:", msg.Id, "	dataLen=", msg.DataLen, "	data=", string(msg.Data))
					}
				}

			}(conn)
		}
	}()

	/*
		模式客户端
	*/
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("Failed to dial at 127.0.0.1:8080:", err)
		return
	}
	// 创建一个封包对象
	dp := NewDataPack()

	// 模拟粘包过程，封装两个message一同发送
	// 封装第一个msg1包
	msg1 := &Message{
		Id:      1,
		DataLen: 4,
		Data:    []byte{'z', 'i', 'n', 'x'},
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("client pack msg1 error:", err)
		return
	}

	// 封装第二个msg2包
	msg2 := &Message{
		Id:      1,
		DataLen: 7,
		Data:    []byte{'h', 'i', 'z', 'i', 'n', 'x', '!'},
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("client pack msg2 error:", err)
		return
	}
	// 将两个包粘在一起一次性发送给服务端
	sendData1 = append(sendData1, sendData2...)
	// 一次性发送给服务端
	_, err = conn.Write(sendData1)
	if err != nil {
		return
	}

	// 客户端阻塞
	select {}
}
