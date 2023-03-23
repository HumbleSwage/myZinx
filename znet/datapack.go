package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"v1/utils"
	"v1/ziface"
)

// DataPack 封包、拆包的具体实例
type DataPack struct{}

// NewDataPack 封包封包实例的初始化方法
func NewDataPack() *DataPack {
	return &DataPack{}
}

// GetHeadLen 获取包的头长度方法
func (d *DataPack) GetHeadLen() uint32 {
	// DataLen uint32 4字节 + ID uint32 4字节
	return 8
}

// Pack 封包方法:|dataLen|msgId|data|
func (d *DataPack) Pack(message ziface.IMessage) ([]byte, error) {
	// 创建一个方法bytes字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})
	// 将dataLen写进dataBuff
	if err := binary.Write(dataBuff, binary.LittleEndian, message.GetMsgLen()); err != nil {
		return nil, err
	}
	// 将MsgId写进dataBuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, message.GetMsgID()); err != nil {
		return nil, err
	}
	// 将data数据写进dataBuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, message.GetData()); err != nil {
		return nil, err
	}
	return dataBuff.Bytes(), nil
}

// UnPack 拆包方法：只需要将包的Head信息提取出来，之后再根据head信息里的data的长度再进行一次读
func (d *DataPack) UnPack(binaryData []byte) (ziface.IMessage, error) {
	// 创建一个从输入二进制的io.Reader
	dataBuff := bytes.NewReader(binaryData)

	// 只解压head信息，得到dataLen和MsgId
	msg := &Message{}

	// 读dataLen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}
	// 读msgId
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	// 判断dataLen是否已经超出了我们允许的最大包的长度
	if utils.GetGlobalObj().MaxPackageSize > 0 && msg.DataLen > utils.GetGlobalObj().MaxPackageSize {
		return nil, errors.New("too large msg data receive!\n")

	}
	return msg, nil
}
