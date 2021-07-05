package znet

type Message struct {
	ID      uint32 // 消息ID
	DataLen uint32 // 消息长度
	Data    []byte // 消息内容
}

func (m *Message) GetMsgID() uint32 {
	return m.ID
}

func (m *Message) GetMsgLen() uint32 {
	return m.DataLen
}

func (m *Message) GetMsg() []byte {
	return m.Data
}

func (m *Message) SetMsgID(id uint32) {
	m.ID = id
}

func (m *Message) SetData(data []byte) {
	m.Data = data
}

func (m *Message) SetDataLen(len uint32) {
	m.DataLen = len
}

// 创建Message 的方法
func NewMessage(id uint32, data []byte) *Message {
	return &Message{
		ID: id,
		DataLen: uint32(len(data)),
		Data:    data,
	}
}
