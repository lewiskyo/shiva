package utils

import (
	"encoding/json"
	"io/ioutil"
	"shiva/iface"
)

// 存储一切有关zinx框架的全局参数,供其他模块使用
// 部分参数通过zinx.json配置

type GlobalObj struct {
	// Server
	TcpServer iface.IServer // 当前zinx全局的Server对象
	Host      string        // 当前服务器主机监听的IP
	TcpPort   int           // 当前服务器主机监听的端口
	Name      string        // 当前服务器名称

	// Zinx
	Version        string // 当前zinx版本号
	MaxConn        int    // 当前服务器主机允许的最大连接数
	MaxPackageSize uint32 // 当前zinx框架数据包的最大值
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
		TcpServer:      nil,
		Host:           "0.0.0.0",
		TcpPort:        8999,
		Name:           "server demo",
		Version:        "V0.4",
		MaxConn:        1000,
		MaxPackageSize: 4096,
	}

	// 应该尝试从conf/zinx.json去加载一些用户自定义的参数
	GlobalObject.Reload()
}
