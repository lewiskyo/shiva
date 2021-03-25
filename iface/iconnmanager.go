package iface

type IConnManager interface {
	// 添加链接
	Add(connection IConnection)
	// 删除链接
	Remove(connection IConnection)
	// 根据connID获取链接
	Get(connID uint32) (IConnection, error)
	// 得到当前总链接数
	Len() int
	// 清除并终止所有连接
	ClearConn()
}

