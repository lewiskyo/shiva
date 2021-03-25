package utils

import (
	"encoding/json"
	"io/ioutil"
	"shiva/iface"
)

// 存储一切有关服务器的全局参数,供其他模块使用
// 部分参数通过shiva.json配置

type GlobalObj struct {
	// Server
	TcpServer iface.IServer // 当前server全局的Server对象
	Host      string        // 当前服务器主机监听的IP
	TcpPort   int           // 当前服务器主机监听的端口
	Name      string        // 当前服务器名称

	// shiva
	Version          string // 当前版本号
	MaxConn          int    // 当前服务器主机允许的最大连接数
	MaxPackageSize   uint32 // 当前服务器数据包的最大值
	WorkerPoolSize   uint32 // 当前业务工作Worker的Goroutine数量
	MaxWorkerTaskLen uint32 // 服务器允许用户最多开辟多少个worker(限定条件)
}

// 定义一个全局的对外的GlobalObj
var GlobalObject *GlobalObj

func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("conf/shiva.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

// 第一次import的时候就会调用init方法
func init() {
	// 如果配置文件没加载的默认值
	GlobalObject = &GlobalObj{
		TcpServer:        nil,
		Host:             "0.0.0.0",
		TcpPort:          8999,
		Name:             "server demo",
		Version:          "V0.7",
		MaxConn:          1000,
		MaxPackageSize:   4096,
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024, // 每个worker对应的消息队列的任务的最大值(channel容量)
	}

	// 应该尝试从conf/zinx.json去加载一些用户自定义的参数
	GlobalObject.Reload()
}
