package main

import (
	"fmt"
	"shiva/iface"
	"shiva/net"
)

type PingRouter struct {
	net.BaseRouter
}

// PreHandle
func (this *PingRouter) PreHandle(request iface.IRequest) {
	fmt.Println("Call PingRouter PreHandle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping..."))
	if err != nil {
		fmt.Println("call back before ping error")
	}
}

// Handle
func (this *PingRouter) Handle(request iface.IRequest) {
	fmt.Println("Call PingRouter Handle...")

	fmt.Println("recv from client, msgid: ", request.GetMsgID(), " data: ", string(request.GetData()))

	err := request.GetConnection().SendMsg(1, []byte("ping...  ping... ping..."))
	if err != nil {
		fmt.Println("call back ping ping ping error")
	}
}

// PostHandle
func (this *PingRouter) PostHandle(request iface.IRequest) {
	fmt.Println("Call PingRouter PostHandle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping..."))
	if err != nil {
		fmt.Println("call back after ping error")
	}
}

type HelloRouter struct {
	net.BaseRouter
}

// PreHandle
func (this *HelloRouter) PreHandle(request iface.IRequest) {
	fmt.Println("Call HelloRouter PreHandle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping..."))
	if err != nil {
		fmt.Println("call back before ping error")
	}
}

// Handle
func (this *HelloRouter) Handle(request iface.IRequest) {
	fmt.Println("Call HelloRouter Handle...")

	fmt.Println("recv from client, msgid: ", request.GetMsgID(), " data: ", string(request.GetData()))

	err := request.GetConnection().SendMsg(1, []byte("hello...  hello... hello..."))
	if err != nil {
		fmt.Println("call back hello error")
	}
}

// PostHandle
func (this *HelloRouter) PostHandle(request iface.IRequest) {
	fmt.Println("Call HelloRouter PostHandle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping..."))
	if err != nil {
		fmt.Println("call back after hello error")
	}
}

func DoConnectionBegin(conn iface.IConnection) {
	fmt.Println("===> DoConnectionBegin")
	if err := conn.SendMsg(101, []byte("doconnection begin")); err != nil {
		fmt.Println(err)
	}
}

func DoConnectionStop(conn iface.IConnection) {
	fmt.Println("===> DoConnectionStop")
	if err := conn.SendMsg(101, []byte("doconnection stop")); err != nil {
		fmt.Println(err)
	}
}

func main() {
	s := net.NewServer("[v0.7]")

	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionStop)

	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloRouter{})

	s.Serve()
}
