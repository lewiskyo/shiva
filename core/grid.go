package core

import (
	"fmt"
	"sync"
)

type Grid struct {
	// 格子id
	GID int
	// 格子左边的边界坐标
	MinX int
	// 格子右边的边界坐标
	MaxX int
	// 格子上边的边界坐标
	MinY int
	// 格子下边的边界坐标
	MaxY int
	// 当前格子内玩家或物体成员的id集合
	playerIDs map[int]bool
	// 保护当前物体集合的锁
	pidLock sync.RWMutex
}

func NewGrid(gid, minX, maxX, minY, maxY int) *Grid {
	return &Grid{
		GID:       gid,
		MinX:      minX,
		MaxX:      maxX,
		MinY:      minY,
		MaxY:      maxY,
		playerIDs: make(map[int]bool),
	}
}

func (g *Grid) Add(pid int) {
	g.pidLock.Lock()
	defer g.pidLock.Unlock()

	g.playerIDs[pid] = true
}

func (g *Grid) Remove(pid int) {
	g.pidLock.Lock()
	defer g.pidLock.Unlock()

	delete(g.playerIDs, pid)
}

func (g *Grid) GetPlayerIds() (playerIds []int) {
	g.pidLock.RLock()
	defer g.pidLock.RUnlock()

	for id := range g.playerIDs {
		playerIds = append(playerIds, id)
	}

	return
}

func (g *Grid) String() string {
	return fmt.Sprintf("gid: %d, minx: %d, maxx: %d, miny: %d, maxy: %d, playerids: %v", g.GID, g.MinX, g.MaxX, g.MinY, g.MaxY, g.playerIDs)
}
