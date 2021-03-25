package apis

import (
	"fmt"
	"shiva/core"
	"shiva/iface"
	"shiva/pb"
	"shiva/proto"
	"shiva/net"
)

// 世界聊天 路由业务
type WorldChatApi struct {
	net.BaseRouter
}

func (wc *WorldChatApi) Handle(request iface.IRequest) {
	// 1. 解析协议
	proto_msg := &pb.Talk{}
	err := proto.Unmarshal(request.GetData(), proto_msg)
	if err != nil {
		fmt.Println("talk unmarshal error, ", err)
		return
	}

	// 2. 当前聊天数据是哪个玩家发送
	pid, err := request.GetConnection().GetProperty("pid")

	// 3. 根据pid获取玩家对象
	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))

	player.Talk(proto_msg.GetContent())
}
