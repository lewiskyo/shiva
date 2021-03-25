package net

import (
	"fmt"
	"net"
	"shiva/utils"
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
		buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		_, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("recv buff err", err)
			continue
		}

		req := Request{
			conn: c,
			data: buf,
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

func (c *Connection) Send(data []byte) error {
	return nil
}
