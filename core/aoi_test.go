package core

import (
	"fmt"
	"testing"
)

func TestNewAOIManager(t *testing.T) {
	// 初始化 AOIManager
	aoiMgr := NewAOIManager(0, 250, 5, 0, 250, 5)

	fmt.Println(aoiMgr)
}

func TestAOIManager_GetSurroundGridsByGid(t *testing.T) {
	aoiMgr := NewAOIManager(0, 250, 5, 0, 300, 4)

	ret := aoiMgr.GetSurroundGridsByGid(17)
	for _, v := range ret {
		fmt.Println("gid: ", v.GID)
	}
}
