package iface

type IServer interface {
	// 启动服务器
	Start()
	// 停止服务器
	Stop()
	// 运行服务器
	Serve()

	// 给当前的Server添加一个路由方法, 供客户端连接处理使用
	AddRouter(uint32, IRouter)

	// 获取当前Server的链接管理
	GetConnMgr() IConnManager
}
