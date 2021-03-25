package main

import (
	"fmt"
	"shiva/iface"
	"shiva/core"
	"shiva/net"
)

func OnConnectionAdd(conn iface.IConnection) {
	// 创建一个player对象
	player := core.NewPlayer(conn)

	// 给客户端发送MsgID:1的消息
	player.SyncPid()

	// 给客户端发送MsgID:200的消息
	player.BroadCastStartPosition()

	fmt.Println("player pid: ", player.Pid, " is login!!!")
}

func DoConnectionStop(conn iface.IConnection) {
	fmt.Println("===> DoConnectionStop")
}

func main() {
	s := net.NewServer("[v0.9]")

	s.SetOnConnStart(OnConnectionAdd)
	s.SetOnConnStop(DoConnectionStop)

	s.Serve()
}
