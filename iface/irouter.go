package iface

type IRouter interface {
	// 在处理conn业务之前的钩子方法(中间件)
	PreHandle(IRequest)

	// 处理conn业务主方法
	Handle(IRequest)

	// 处理conn业务之后的钩子方法(中间件 next())
	PostHandle(IRequest)
}

