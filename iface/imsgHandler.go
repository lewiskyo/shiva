package iface

type IMsgHandler interface {
	// 调度/执行对应的Router消息处理方法
	DoMsgHandler(IRequest)

	// 为消息添加具体的处理逻辑
	AddRouter(uint32, IRouter)

	// 启动Worker工作池
	StartWorkerPool()

	SendMsgToTaskQueue(IRequest)
}
