package core

import "sync"

/*
	当前游戏的世界总管理模块
*/

type WorldManager struct {
	// AOIManager 当前世界地图AOI管理模块
	AoiMgr *AOIManager
	// 当前全部在线的Players集合
	Players map[int32]*Player
	// 保护集合的读写锁
	pLock sync.RWMutex
}

var WorldMgrObj *WorldManager

func init() {
	WorldMgrObj = &WorldManager{
		AoiMgr:  NewAOIManager(AOI_MIN_X, AOI_MAX_X, AOI_CNTS_X, AOI_MIN_Y, AOI_MAX_Y, AOI_CNTS_Y),
		Players: make(map[int32]*Player),
	}
}

// 添加一个玩家
func (wm *WorldManager) AddPlayer(player *Player) {
	wm.pLock.Lock()
	wm.Players[player.Pid] = player
	wm.pLock.Unlock()

	wm.AoiMgr.AddToGridByPos(int(player.Pid), player.X, player.Z)
}

// 删除一个玩家
func (wm *WorldManager) RemovePlayer(pid int32) {
	player := wm.Players[pid]
	wm.AoiMgr.RemoveFromGridByPos(int(pid), player.X, player.Z)

	wm.pLock.Lock()
	delete(wm.Players, pid)
	wm.pLock.Unlock()
}

// 通过玩家id查询player对象
func (wm *WorldManager) GetPlayerByPid(pid int32) *Player {
	wm.pLock.RLock()
	defer wm.pLock.RUnlock()
	return wm.Players[pid]
}

// 获取全部的在线玩家
func (wm *WorldManager) GetAllPlayers() []*Player {
	players := make([]*Player, 0)

	wm.pLock.RLock()
	defer wm.pLock.RUnlock()

	for _, v := range wm.Players {
		players = append(players, v)
	}

	return players
}
