package main

import (
	"fmt"
	"shiva/apis"
	"shiva/core"
	"shiva/iface"
	"shiva/net"
)

func OnConnectionAdd(conn iface.IConnection) {
	// 创建一个player对象
	player := core.NewPlayer(conn)

	// 给客户端发送MsgID:1的消息
	player.SyncPid()

	// 给客户端发送MsgID:200的消息
	player.BroadCastStartPosition()

	core.WorldMgrObj.AddPlayer(player)

	conn.SetProperty("pid", player.Pid)

	fmt.Println("player pid: ", player.Pid, " is login!!!")
}

func DoConnectionStop(conn iface.IConnection) {
	fmt.Println("===> DoConnectionStop")
}

func main() {
	s := net.NewServer("[v0.9]")

	s.SetOnConnStart(OnConnectionAdd)
	s.SetOnConnStop(DoConnectionStop)

	s.AddRouter(2, &apis.WorldChatApi{})

	s.Serve()
}
