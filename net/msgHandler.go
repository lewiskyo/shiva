package net

import (
	"fmt"
	"shiva/iface"
	"shiva/log"
	"shiva/utils"
)

type MsgHandler struct {
	// 存放每个msgID 所对应的处理方法
	Apis map[uint32]iface.IRouter

	// 负责worker取任务的消息队列
	TaskQueue []chan iface.IRequest

	// 业务工作worker池的worker数量
	WorkerPoolSize uint32
}

func NewMsgHandle() *MsgHandler {
	return &MsgHandler{
		Apis:           make(map[uint32]iface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueue:      make([]chan iface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}

func (mh *MsgHandler) DoMsgHandler(request iface.IRequest) {
	msgID := request.GetMsgID()

	handler, ok := mh.Apis[msgID]
	if !ok {
		fmt.Println("DoMsgHandler fail, msgid: ", msgID)
		return
	}
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

func (mh *MsgHandler) AddRouter(msgID uint32, router iface.IRouter) {
	if _, ok := mh.Apis[msgID]; ok {
		fmt.Println("AddRouter fail, msgid exist, id: ", msgID)
		return
	}

	mh.Apis[msgID] = router
	fmt.Println("AddRouter success, msgid: ", msgID)
	return
}

// 启动worker池
func (mh *MsgHandler) StartWorkerPool() {
	// 根据workerPoolSize 分别开启Worker, 每个Worker用一个go来承载
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		// 一个worker被启动
		// 1. 当前的worker对应的channel消息队列 开辟空间, 第0个worker用第0个channel...
		mh.TaskQueue[i] = make(chan iface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		// 2. 启动当前的worker, 阻塞等待消息从channel传递进来
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

func (mh *MsgHandler) StartOneWorker(workerID int, taskQueue chan iface.IRequest) {
	log.Info("Worker ID = ", workerID, "is started...")

	for {
		select {
		// 如果有消息过来, 出列的就是一个客户端request, 执行当前request绑定的业务
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

// 将消息交给TaskQueue, 由worker来处理
func (mh *MsgHandler) SendMsgToTaskQueue(request iface.IRequest) {
	// 1. 将消息平均分配给不同的worker
	// 根据客户端建立的ConnID来进行分配
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Println("msg send to workerid", workerID)

	// 2. 将消息发送给对应的worker的TaskQueue即可
	mh.TaskQueue[workerID] <- request
}
