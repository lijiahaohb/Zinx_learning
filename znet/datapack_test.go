package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

// 负责 datapack 拆包封包的单元测试
func TestDataPack(t *testing.T) {
	// 模拟的客户服务器

	// 1. 创建Socket TCP
	listener, err := net.Listen("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("server Listen error: ", err)
		return
	}

	// 创建一个goroutine来负责从客户端处理业务

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("server accept error: ", err)
			}

			go func(conn net.Conn) {
				// 处理客户端请求
				/*------> 拆包的过程 <------*/
				// 定义一个拆包的对象
				dp := NewDataPack()

				for {
					// 1. 从客户端读取包的head
					headData := make([]byte, dp.GetHeadLen())
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						fmt.Println("read head data error: ", err)
						break
					}

					msgHead, err := dp.Unpack(headData)
					if err != nil {
						fmt.Println("server unpack error: ", err)
						return
					}
					if msgHead.GetMsgLen() > 0 {
						// 2. 根据head中的dataLen来读取数据内容
						msg := msgHead.(*Message)
						msg.Data = make([]byte, msg.GetMsgLen())

						// 根据dataLen的长度再次从io流中读取
						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println("server unpack data error: ", err)
							return
						}

						fmt.Println("Recv MsgID: ", msg.ID, "dataLen: ", msg.DataLen, "data: ", string(msg.Data))
					}
				}
			}(conn)
		}
	}()

	// 模拟客户端
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("dial error: ", err)
		return
	}

	// 创建一个封包对象 dp
	dp := NewDataPack()

	// 模拟黏包过程，封装两个msg一同发送
	// 封装第一个msg
	msg1 := &Message{
		ID:      1,
		DataLen: 4,
		Data:    []byte{'z', 'i', 'n', 'x'},
	}

	// 给第一个msg打包
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("client pack msg1 error: ", err)
	}

	// 封装第二个msg
	msg2 := &Message{
		ID: 2,
		DataLen: 9,
		Data: []byte{'h', 'e','l','l', 'o', 'z', 'i', 'n', 'x'},
	}
	// 给第二个msg打包
	sendData2, err := dp.Pack(msg2) 
	if err != nil {
		fmt.Println("client pack msg2 eror: ", err)
	}

	sendData1 = append(sendData1, sendData2...)

	conn.Write(sendData1)

	// 客户端阻塞
	select{}
}
