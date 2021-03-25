package net

type Message struct {
	Id      uint32
	DataLen uint32
	Data    []byte
}

func NewMessagePackage(id uint32, data []byte) *Message {
	msg := &Message{
		Id:      id,
		DataLen: uint32(len(data)),
		Data:    data,
	}

	return msg
}

func (m *Message) GetMsgId() uint32 {
	return m.Id
}

func (m *Message) GetMsgLen() uint32 {
	return m.DataLen
}

func (m *Message) GetData() []byte {
	return m.Data
}

func (m *Message) SetMsgID(Id uint32) {
	m.Id = Id
}

func (m *Message) SetData(data []byte) {
	m.Data = data
}

func (m *Message) SetDataLen(len uint32) {
	m.DataLen = len
}

