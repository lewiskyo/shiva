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
