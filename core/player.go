package core

import (
	"fmt"
	"math/rand"
	"sync"
	"shiva/pb"
	"shiva/proto"
	"shiva/iface"
)

type Player struct {
	Pid  int32
	Conn iface.IConnection
	X    float32
	Y    float32
	Z    float32
	V    float32
}

var PidGen int32 = 1  // 用来生成玩家ID计数器
var IdLock sync.Mutex // 保护PidGen的Mutex

// 创建玩家的方法
func NewPlayer(conn iface.IConnection) *Player {

	IdLock.Lock()
	id := PidGen
	PidGen++
	IdLock.Unlock()

	p := &Player{
		Pid:  id,
		Conn: conn,
		X:    float32(160 + rand.Intn(10)), // 随机在160坐标点, 基于x轴若干便宜
		Y:    0,
		Z:    float32(140 + rand.Intn(20)), // 随机在140坐标点, 基于y轴若干便宜
		V:    0,
	}

	return p
}

/*
	提供一个发送给客户端消息的方法
	主要是将pb的protobuf数据序列化之后，在调用zinx的SendMsg方法
*/
func (p *Player) SendMsg(msgId uint32, data proto.Message) {
	// 将proto Message结构体序列化, 转换成二进制
	msg, err := proto.Marshal(data)
	if err != nil {
		fmt.Println("marshal msg err, ", err)
		return
	}

	// 将二进制文件 通过zinx框架的sendMsg将数据发送给客户端
	if p.Conn == nil {
		fmt.Println("connection in player is nil")
		return
	}

	if err := p.Conn.SendMsg(msgId, msg); err != nil {
		fmt.Println("Player sendmsg err, ", err)
		return
	}
}

/*
	告知客户端玩家pid, 同步已经生成的玩家ID给客户端
*/
func (p *Player) SyncPid() {
	data := &pb.SyncPid{
		Pid: p.Pid,
	}

	p.SendMsg(1, data)
}

// 广播自己的出生点
func (p *Player) BroadCastStartPosition() {
	msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  2, //TP 2 代表广播坐标
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}

	//发送数据给客户端
	p.SendMsg(200, msg)
}

// 玩家广播聊天消息
func (p *Player) Talk(content string) {
	msg := &pb.BroadCast{
		Pid:  p.Pid,
		Tp:   1,
		Data: &pb.BroadCast_Content{Content: content},
	}

	players := WorldMgrObj.GetAllPlayers()
	for _, p := range players {
		p.SendMsg(200, msg)
	}
}

// 同步玩家上线的位置消息
func (p *Player) SyncSurrounding() {
	// 1. 获取当前玩家周围的玩家有哪些(9宫格)
	players := p.GetSurroundPlayers()

	// 2. 将当前玩家的位置信息通过MsgID:200发给周围的玩家,(让其他玩家看到自己)
	proto_msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  2,
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}
	for _, player := range players {
		player.SendMsg(200, proto_msg)
	}

	// 3. 将周围的全部玩家位置信息发送给当前的玩家 MsgID:202科幻段, (让自己看到其他玩家)
	pb_players := make([]*pb.Player, 0, len(players))
	for _, player := range players {
		pb_players = append(pb_players, &pb.Player{
			Pid: player.Pid,
			P: &pb.Position{
				X: player.X,
				Y: player.Y,
				Z: player.Z,
				V: player.V,
			},
		})
	}

	proto_msg2 := &pb.SyncPlayers{
		Ps: pb_players,
	}
	p.SendMsg(202, proto_msg2)
}

func (p *Player) UpdatePos(x, y, z, v float32) {
	proto_msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  4,
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: x,
				Y: y,
				Z: z,
				V: v,
			},
		},
	}

	players := p.GetSurroundPlayers()
	for _, player := range players {
		player.SendMsg(200, proto_msg)
	}
}

func (p *Player) GetSurroundPlayers() []*Player {
	pids := WorldMgrObj.AoiMgr.GetPidsByPos(p.X, p.Z)
	players := make([]*Player, 0, len(pids))
	for _, pid := range pids {
		players = append(players, WorldMgrObj.GetPlayerByPid(int32(pid)))
	}

	return players
}

func (p *Player) Offline() {
	players := p.GetSurroundPlayers()

	proto_msg := &pb.SyncPid{
		Pid: p.Pid,
	}

	for _, surroundPlayer := range players {
		surroundPlayer.SendMsg(201, proto_msg)
	}

	WorldMgrObj.AoiMgr.RemoveFromGridByPos(int(p.Pid), p.X, p.Z)
	WorldMgrObj.RemovePlayer(p.Pid)
}
