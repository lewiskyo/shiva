package apis

import (
	"fmt"
	"shiva/core"
	"shiva/iface"
	"shiva/net"
	"shiva/pb"
	"shiva/proto"
)

type MoveApi struct {
	net.BaseRouter
}

func (m *MoveApi) Handle(request iface.IRequest) {
	proto_msg := &pb.Position{}
	err := proto.Unmarshal(request.GetData(), proto_msg)
	if err != nil {
		fmt.Println("move api err, ", err)
		return
	}

	pid, err := request.GetConnection().GetProperty("pid")
	if err != nil {
		fmt.Println("get property err", err)
		return
	}

	p := core.WorldMgrObj.GetPlayerByPid(pid.(int32))

	p.UpdatePos(proto_msg.X, proto_msg.Y, proto_msg.Z, proto_msg.V)
}
