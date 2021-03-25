package net

import (
"bytes"
"encoding/binary"
"errors"
	"shiva/iface"
	"shiva/utils"
)

// 封包,拆包的具体模块
type DataPack struct {
}

// 拆包封包实例的一个初始化方法
func NewDataPack() *DataPack {
	return &DataPack{}
}

// 获取包的头的长度方法
func (db *DataPack) GetHeadLen() uint32 {
	// Len(uint32) + id(uint32)
	return 8
}

// 封包方法
func (db *DataPack) Pack(msg iface.IMessage) ([]byte, error) {
	// 创建一个存放bytes字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	// 将dataLen写进dataBuff
	err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgLen())
	if err != nil {
		return nil, err
	}

	// 将Id写进dataBuff
	err = binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId())
	if err != nil {
		return nil, err
	}
	// 将data写进dataBuff
	err = binary.Write(dataBuff, binary.LittleEndian, msg.GetData())
	if err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

// 拆包方法
// 将Head信息读出来,之后根据Head信息里面的dataLen长度，再读取data
func (db *DataPack) Unpack(binaryData []byte) (iface.IMessage, error) {
	// 创建一个从输入二进制数据的ioreader
	dataBuff := bytes.NewReader(binaryData)

	// 只解压head信息,得到dataLen和MsgID
	msg := &Message{}
	// 读dataLen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	// 读MsgID
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	// 判断dataLen是否已经超出了我们允许的最大包长度
	maxPackageSize := utils.GlobalObject.MaxPackageSize
	if maxPackageSize > 0 && msg.DataLen > maxPackageSize {
		return nil, errors.New("too lager msg data recv!!")
	}

	return msg, nil
}

