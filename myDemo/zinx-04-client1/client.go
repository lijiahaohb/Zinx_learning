package main

import (
	"fmt"
	"io"
	"net"
	"time"
	"zinx/znet"
)

// 模拟客户端
func main() {
	fmt.Println("client1 start ...")
	time.Sleep(time.Second * 1)
	// 1. 直接连接远程服务器，得到conn连接
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client dial error: ", err)
		return
	}

	// 2. 写数据
	for {
		// 发送TVL格式的数据
		dp := znet.NewDataPack()

		binaryMsg, err := dp.Pack(znet.NewMessage(1, []byte("Zinx v0.5 client1 Test Message")))
		if err != nil {
			fmt.Println("pack error: ", err)
			return
		}

		if _, err := conn.Write(binaryMsg); err != nil {
			fmt.Println("conn write error: ", err)
			return
		}

		// 先读取流中head 部分 得到ID 和 dataLen
		binaryHead := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, binaryHead); err != nil {
			fmt.Println("read head error: ", err)
			break
		}

		// 将二进制的head进行拆包到msg结构体重
		msgHead, err := dp.Unpack(binaryHead)
		if err != nil {
			fmt.Println("client unpack error: ", err)
			break
		}

		if msgHead.GetMsgLen() > 0 {
			// 说明 是有数据的 继续读取
			msg := msgHead.(*znet.Message) 
			msg.Data = make([]byte, msgHead.GetMsgLen())

			if _, err := io.ReadFull(conn, msg.Data); err != nil {
				fmt.Println("read msg data error: ", err)
				return
			}

			fmt.Println("recv server msg, ID = ", msg.ID, " len = ", msg.DataLen, " data = ", string(msg.Data))
		}

		// 再根据dataLen得到数据内容

		time.Sleep(1 * time.Second)
	}
}
