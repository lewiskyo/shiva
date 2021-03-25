package net

import (
	"errors"
	"fmt"
	"shiva/iface"
	"sync"
)

type ConnManager struct {
	connections map[uint32]iface.IConnection // 管理的链接集合
	connLock    sync.RWMutex                 // 保护链接集合的读写锁
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]iface.IConnection),
	}
}

func (connMgr *ConnManager) Add(connection iface.IConnection) {
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	// 将connection加入到ConnManager中
	connMgr.connections[connection.GetConnID()] = connection
	fmt.Println("connection add to connmgr success, conn num = ", connMgr.Len())
}

// 删除链接
func (connMgr *ConnManager) Remove(connection iface.IConnection) {
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	delete(connMgr.connections, connection.GetConnID())
	fmt.Println("connection remove succ, conn num = ", connMgr.Len())
}

// 根据connID获取链接
func (connMgr *ConnManager) Get(connID uint32) (iface.IConnection, error) {
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	if conn, ok := connMgr.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("get connid fail, no exist")
	}
}

// 得到当前总链接数
func (connMgr *ConnManager) Len() int {
	return len(connMgr.connections)
}

// 清除并终止所有连接
func (connMgr *ConnManager) ClearConn() {
	connMgr.connLock.Lock()
	// 删除connection并且停止connection工作
	for connID, conn := range connMgr.connections {
		conn.Stop()
		delete(connMgr.connections, connID)
	}

	fmt.Println("clear all connection success, len = ", connMgr.Len())
}
