package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"zinx/utils"
	"zinx/ziface"
)

// 封包 拆包的具体模块
type DataPack struct {
}

func NewDataPack() *DataPack {
	return &DataPack{}
}

func (d *DataPack) GetHeadLen() uint32 {
	// DataLen uint32(4字节) + ID uint32(4字节)
	return 8
}

func (d *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	// 创建一个byte字节的缓冲
	dataBuf := bytes.NewBuffer([]byte{})

	// 将DataLen写进databuf中
	if err := binary.Write(dataBuf, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		return nil, err
	}

	// 将MsgID写进databuf中
	if err := binary.Write(dataBuf, binary.LittleEndian, msg.GetMsgID()); err != nil {
		return nil, err
	}

	// 将data数据写入到databuf中
	if err := binary.Write(dataBuf, binary.LittleEndian, msg.GetMsg()); err != nil {
		return nil, err
	}

	return dataBuf.Bytes(), nil
}

func (d *DataPack) Unpack(binaryData []byte) (ziface.IMessage, error) {
	// 创建一个从 二进制数据总读取数据的 ioreader
	dataBuf := bytes.NewReader(binaryData)

	msg := &Message{}

	// 读dataLen
	if err := binary.Read(dataBuf, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	// 读MsgID
	if err := binary.Read(dataBuf, binary.LittleEndian, &msg.ID); err != nil {
		return nil, err
	}

	if (utils.GlobalObject.MaxPackageSize >0 && msg.DataLen > utils.GlobalObject.MaxPackageSize) {
		return nil , errors.New("too large msg data received")
	}

	return msg, nil
}
