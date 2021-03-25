package net

import "shiva/iface"

// 实现router时, 先嵌入这个BaseRouter基类, 然后根据需要对这个基类的方法进行重写
type BaseRouter struct {
}

// 这里之所以BaseRouter的方法都为空, 因为有的Router不希望有PreHandle, PostHandle两个业务
// 所以Router继承BaseRouter好处是, 不需要实现PreHandle, PostHandle
func (br *BaseRouter) PreHandle(request iface.IRequest) {

}

func (br *BaseRouter) Handle(request iface.IRequest) {

}

func (br *BaseRouter) PostHandle(request iface.IRequest) {

}
