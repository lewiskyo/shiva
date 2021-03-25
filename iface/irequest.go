package iface

// 将客户端请求的链接信息和请求数据包装在一起

type IRequest interface {
	GetConnection() IConnection

	GetData() []byte
}

