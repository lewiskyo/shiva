package iface

import "net"

type IConnection interface {
	Start()

	Stop()

	GetTCPConnection() *net.TCPConn

	GetConnID() uint32

	RemoteAddr() net.Addr

	SendMsg(uint32, []byte) error

	SetProperty(string, interface{})

	GetProperty(string) (interface{}, error)

	RemoveProperty(string)
}

// 定义一个处理连接业务的方法
type HandleFunc func(*net.TCPConn, []byte, int) error
