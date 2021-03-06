package net

import (
	"errors"
	"fmt"
	"io"
	"net"
	"shiva/iface"
	"shiva/utils"
	"sync"
)

type Connection struct {
	// 当前conn隶属于server
	Server iface.IServer

	Conn *net.TCPConn

	ConnID uint32

	isClosed bool

	// 告知当前连接已经退出/停止, 由Reader告知Writer退出
	ExitChan chan bool

	// 无缓冲的管道, 用于读写Goroutine之间的消息通信
	msgChan chan []byte

	// 消息的管理MsgID与对应的api处理关系
	MsgHandler iface.IMsgHandler

	// 链接属性集合
	property map[string]interface{}
	// 保护链接属性的锁
	propertyLock sync.RWMutex
}

func NewConnection(server iface.IServer, conn *net.TCPConn, connID uint32, msgHandler iface.IMsgHandler) *Connection {
	c := &Connection{
		Server:     server,
		Conn:       conn,
		ConnID:     connID,
		MsgHandler: msgHandler,
		msgChan:    make(chan []byte),
		isClosed:   false,
		ExitChan:   make(chan bool, 1),
		property:   make(map[string]interface{}),
	}
	c.Server.GetConnMgr().Add(c)
	return c
}

func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running ...")
	defer fmt.Println("connID = ", c.ConnID, "Reader is exit, remote addr is ", c.RemoteAddr().String())
	// 下面for循环只要break, reader就会退出, 关闭reader 和 writer goroutine
	defer c.Stop()

	for {
		dp := NewDataPack()
		headData := make([]byte, dp.GetHeadLen())

		_, err := io.ReadFull(c.Conn, headData)
		if err != nil {
			fmt.Println("io readfull headdata err, ", err)
			break
		}

		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("dataunpack err, ", err)
			break
		}

		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.Conn, data); err != nil {
				fmt.Println("io readfull data err, ", err)
				break
			}
		}

		msg.SetData(data)

		req := Request{
			conn: c,
			msg:  msg,
		}

		if utils.GlobalObject.WorkerPoolSize > 0 {
			// 已经开辟了工作池和消息队列
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			// 根据绑定好的msgID找到对应的api业务执行
			go c.MsgHandler.DoMsgHandler(&req)
		}

	}
}

// 写消息goroutine, 专门发送给客户端消息的模块
func (c *Connection) StartWriter() {
	fmt.Println("writer goroutine is running..")
	defer fmt.Println(c.RemoteAddr().String(), "[conn write exit]")
	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("send data err: ", err)
				return
			}
		case <-c.ExitChan:
			//  代表Reader已经退出, 此时Writer也要退出
			return
		}
	}

}

func (c *Connection) Start() {
	fmt.Println("Conn Start() ... ConnID = ", c.ConnID)

	go c.StartReader()
	go c.StartWriter()

	// 按照开发者传递进来的 创建链接后的业务hook函数
	c.Server.CallOnConnStart(c)
}

func (c *Connection) Stop() {
	fmt.Println("Conn Stop().. ConnID: ", c.ConnID)

	if c.isClosed == true {
		return
	}

	c.isClosed = true

	c.Server.CallOnConnStop(c)

	c.Conn.Close()

	c.ExitChan <- true

	c.Server.GetConnMgr().Remove(c)

	close(c.ExitChan)
	close(c.msgChan)
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) SendMsg(Id uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("Connection closed when send msg")
	}

	// 将data进行封包 MsgDataLen|Id|Data
	dp := NewDataPack()

	binaryMsg, err := dp.Pack(NewMessagePackage(Id, data))
	if err != nil {
		fmt.Println("sendmsg pack err", err)
		return errors.New("pack error msg")
	}

	// 将消息发给writer
	c.msgChan <- binaryMsg

	return nil
}

func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.property[key] = value
}

func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	if value, ok := c.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("no property found")
	}

}

func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}
