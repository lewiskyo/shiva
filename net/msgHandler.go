package net

import (
	"fmt"
	"shiva/iface"
)

type MsgHandler struct {
	Apis map[uint32]iface.IRouter
}

func NewMsgHandle() *MsgHandler {
	return &MsgHandler{
		Apis: make(map[uint32]iface.IRouter),
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
