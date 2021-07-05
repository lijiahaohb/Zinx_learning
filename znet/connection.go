package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"zinx/utils"
	"zinx/ziface"
)

type Connection struct {
	// 当前connection 属于的 Server
	TcpServer ziface.IServer
	// 当前连接的套接字
	Conn *net.TCPConn
	// 连接ID
	ConnID uint32
	// 连接状态
	isClosed bool
	// 告知当前连接被动关闭的channel (reader 通过 ExitChan 告知 Writer 退出)
	ExitChan chan bool
	// 用于读、写Goroutine之间通信的 Channel 无缓冲
	msgChan chan []byte
	// 该连接对应的处理方法
	msgHandler ziface.IMsgHandler
	// 连接属性集合
	properties map[string]interface{}
	// 保护连接属性的锁
	propertiesLock sync.RWMutex
}

// 初始化连接模块的方法
func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandler) *Connection {
	c := &Connection{
		TcpServer:  server,
		Conn:       conn,
		ConnID:     connID,
		isClosed:   false,
		ExitChan:   make(chan bool, 1),
		msgChan:    make(chan []byte),
		msgHandler: msgHandler,
		properties: make(map[string]interface{}),
	}

	// 将新创建的 connection 加入到server中的ConnMgr中
	c.TcpServer.GetConnMgr().Add(c)

	return c
}

// 连接的读业务的方法
func (c *Connection) StartReader() {
	fmt.Println("[Reader goroutine] is running ..., ConnID = ", c.ConnID)
	defer fmt.Println("[Reader Goroutine] is exiting, ConnID = ", c.ConnID, "remote addr is ", c.Conn.RemoteAddr().String())
	defer c.Stop()

	// 循环读数据
	for {
		// 读取客户端的数据到buf中
		// buf := make([]byte, utils.GlobalObject.MaxPackageSize)

		// cnt, err := c.Conn.Read(buf)
		// if cnt == 0 {
		// 	break
		// }
		// if err != nil && err != io.EOF {
		// 	fmt.Println("recv buf error", err)
		// 	continue
		// }

		// 创建一个拆包、解包的对象
		dp := NewDataPack()

		// 读取客户端的Msg head
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read message head error: ", err)
			break
		}

		// 拆包，得到msgID 和msgDataLen 放在消息中
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error: ", err)
			break
		}

		// 根据DataLen 再次读取 data 放在msg.Data 中
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data error: ", err)
				break
			}
		}
		msg.SetData(data)

		// 得到当前连接数据的 Request
		req := Request{
			conn: c,
			msg:  msg,
		}

		// 根据绑定好的 msgID 找到对应的 处理api的业务执行
		if utils.GlobalObject.WorkerPoolSize > 0 {
			// 已经开启工作池机制
			c.msgHandler.SendMsgToTaskQueue(&req)
		} else {
			go c.msgHandler.DoMsgHandler(&req)
		}
	}
}

// 专门将消息发送给客户端的Goroutine
func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine] is running ...")
	defer fmt.Println("[Writer Goroutine] is exiting!", c.Conn.RemoteAddr().String())

	// 不断的阻塞等待channel的消息， 收到消息后写给客户端
	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data error: ", err)
				return
			}
		case <-c.ExitChan:
			// 代表 Reader 已经退出，此时Writer也应该要退出
			return
		}
	}
}

// 以TLV的形式封装数据，之后在进行发包的方法
func (c *Connection) SendMsg(msgID uint32, data []byte) error {
	if c.isClosed {
		return errors.New("Connection closed")
	}

	// 将data进行封包
	dp := NewDataPack()

	binaryMsg, err := dp.Pack(NewMessage(msgID, data))
	if err != nil {
		fmt.Println("Pack error msgid = ", msgID)
		return errors.New("Pack error msg")
	}

	// 将数据发送给 msgChan
	c.msgChan <- binaryMsg

	return nil
}

func (c *Connection) Start() {
	fmt.Println("Conn is starting ..., ConnID = ", c.ConnID)

	// 启动从连接读数据的业务
	go c.StartReader()

	// 启动从当前连接写数据的业务
	go c.StartWriter()

	// 执行开发者传递进来的连接建立之后的钩子函数
	c.TcpServer.CallOnConnStart(c)
}

func (c *Connection) Stop() {
	fmt.Println("Conn is stoping ..., ConnID = ", c.ConnID)

	// 判断当前连接的状态
	if c.isClosed {
		return
	}

	c.isClosed = true
	// 调用开发者注册的连接断开之前需要执行的钩子函数
	c.TcpServer.CallOnConnStop(c)
	// 关闭socket连接
	c.Conn.Close()
	// 告知Writer 关闭
	c.ExitChan <- true
	// 将conn 从 ConnMgr 中移除
	c.TcpServer.GetConnMgr().Remove(c)
	// 关闭管道
	close(c.ExitChan)
	// 关闭msgCHan
	close(c.msgChan)
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnectionID() uint32 {
	return c.ConnID
}

func (c *Connection) GetRemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) Send(data []byte) error {
	return nil
}

// 设置连接属性
func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertiesLock.Lock()
	defer c.propertiesLock.Unlock()

	c.properties[key] = value
}

// 获取连接属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertiesLock.RLock()
	defer c.propertiesLock.RUnlock()

	if value, ok := c.properties[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("property not found")
	}
}

// 移除连接属性
func (c *Connection) RemoveProperty(key string) {
	c.propertiesLock.Lock()
	defer c.propertiesLock.Unlock()

	delete(c.properties, key)
}
