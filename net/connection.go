package net

import (
	"errors"
	"fmt"
	"io"
	"net"
	"shiva/iface"
)

type Connection struct {
	Conn *net.TCPConn

	ConnID uint32

	isClosed bool

	ExitChan chan bool

	// 该链接处理的方法Router
	Router iface.IRouter
}

func NewConnection(conn *net.TCPConn, connID uint32, router iface.IRouter) *Connection {
	c := &Connection{
		Conn:     conn,
		ConnID:   connID,
		Router:   router,
		isClosed: false,
		ExitChan: make(chan bool, 1),
	}

	return c
}

func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running ...")
	defer fmt.Println("connID = ", c.ConnID, "Reader is exit, remote addr is ", c.RemoteAddr().String())
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

		// 执行注册的路由方法
		go func(request iface.IRequest) {
			// 从路由中找到注册绑定的conn对应router调用
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)
	}
}

func (c *Connection) Start() {
	fmt.Println("Conn Start() ... ConnID = ", c.ConnID)

	go c.StartReader()
}

func (c *Connection) Stop() {
	fmt.Println("Conn Stop().. ConnID: ", c.ConnID)

	if c.isClosed == true {
		return
	}

	c.isClosed = true

	c.Conn.Close()

	close(c.ExitChan)
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

	if _, err := c.Conn.Write(binaryMsg); err != nil {
		fmt.Println("sendmsg write err", err)
		return errors.New("Write msg id")
	}

	return nil
}
