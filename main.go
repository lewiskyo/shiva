package main

import (
	"shiva/apis"
	"shiva/core"
	"shiva/iface"
	"shiva/log"
	"shiva/net"
)

func OnConnectionAdd(conn iface.IConnection) {
	// 创建一个player对象
	player := core.NewPlayer(conn)

	// 给客户端发送MsgID:1的消息,同步当前player的id给客户端
	player.SyncPid()

	// 给客户端发送MsgID:200的消息,同步当前player的初始位置给客户端
	player.BroadCastStartPosition()

	// 将当前新上线的玩家添加到WorldManager中
	core.WorldMgrObj.AddPlayer(player)

	// 将该链接绑定一个Pid 玩家id的属性
	conn.SetProperty("pid", player.Pid)

	// 同步周边的玩家, 告知他们当前玩家已经上线, 广播当前玩家的位置信息
	player.SyncSurrounding()

	log.Info("player pid: ", player.Pid, "is login!!!")
}

func OnConnectionStop(conn iface.IConnection) {
	pid, err := conn.GetProperty("pid")
	if err != nil{
		return
	}
	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))
	player.Offline()
}

func main() {
	s := net.NewServer("[shiva v0.9]")

	s.SetOnConnStart(OnConnectionAdd)
	s.SetOnConnStop(OnConnectionStop)

	s.AddRouter(2, &apis.WorldChatApi{})
	s.AddRouter(3, &apis.MoveApi{})

	s.Serve()
}
