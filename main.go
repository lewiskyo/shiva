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
	fmt.Println("Call Router PreHandle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping..."))
	if err != nil {
		fmt.Println("call back before ping error")
	}
}

// Handle
func (this *PingRouter) Handle(request iface.IRequest) {
	fmt.Println("Call Router Handle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping ping ping..."))
	if err != nil {
		fmt.Println("call back ping ping ping error")
	}
}

// PostHandle
func (this *PingRouter) PostHandle(request iface.IRequest) {
	fmt.Println("Call Router PostHandle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping..."))
	if err != nil {
		fmt.Println("call back after ping error")
	}
}

func main() {
	s := net.NewServer("[v0.4]")

	s.AddRouter(&PingRouter{})

	s.Serve()
}
