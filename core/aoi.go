package core

import "fmt"

/*
	AOI区域管理模块
*/
type AOIManager struct {
	// 区域的左边界坐标
	MinX int
	// 区域的右边界坐标
	MaxX int
	// X方向格子的数量
	CntsX int
	// 区域的上边界坐标
	MinY int
	// 区域的下边界坐标
	MaxY int
	// Y方向格子的数量
	CntsY int
	// 当前区域有哪些格子map-key = 格子id => 格子对象
	grids map[int]*Grid
}

/*
初始化一个AOI区域管理模块
*/
func NewAOIManager(minX, maxX, cntsX, minY, maxY, cntsY int) *AOIManager {
	aoiMgr := &AOIManager{
		MinX:  minX,
		MaxX:  maxX,
		CntsX: cntsX,
		MinY:  minY,
		MaxY:  maxY,
		CntsY: cntsY,
		grids: make(map[int]*Grid),
	}

	// 给aoi初始化区域的格子所有格子进行编号和初始化
	for y := 0; y < cntsY; y++ {
		for x := 0; x < cntsX; x++ {
			// 计算格子id 根据x,y编号
			// id  = idy * cntx + idx
			gid := y*cntsX + x

			// 初始化gid格子
			aoiMgr.grids[gid] = NewGrid(gid,
				aoiMgr.MinX+aoiMgr.gridXLen()*x,
				aoiMgr.MinX+aoiMgr.gridXLen()*(x+1),
				aoiMgr.MinY+aoiMgr.gridYLen()*y,
				aoiMgr.MinY+aoiMgr.gridYLen()*(y+1))
		}
	}

	return aoiMgr
}

// 得到每个格子在x轴方向的长度
func (m *AOIManager) gridXLen() int {
	return (m.MaxX - m.MinX) / m.CntsX
}

// 得到每个格子在y轴方向的长度
func (m *AOIManager) gridYLen() int {
	return (m.MaxY - m.MinY) / m.CntsY
}

// 打印格子信息
func (m *AOIManager) String() string {
	// 打印aoi manager信息
	s := fmt.Sprintf("aoimanager:\nminx: %d, maxx: %d, cntx: %d, miny: %d, maxy: %d, cnty: %d\ngrids in aoimanager\n",
		m.MinX, m.MaxX, m.CntsX, m.MinY, m.MaxY, m.CntsY)

	// 打印全部格子信息
	for _, g := range m.grids {
		s += fmt.Sprintln(g)
	}

	return s
}
