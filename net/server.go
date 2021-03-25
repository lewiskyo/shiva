package net

import (
	"fmt"
	"net"
	"shiva/iface"
	"shiva/utils"
)

type Server struct {
	// 服务器名称
	Name string
	// 服务器绑定的ip版本
	IPVersion string
	// 服务器监听的ip
	IP string
	// 服务器监听的端口
	Port int
	// 当前server的消息管理模块, 用来绑定MsgID和对应处理业务API关系
	MsgHandler iface.IMsgHandler
	// 当前server的链接管理器
	ConnMgr iface.IConnManager
	// 该Server创建连接之后自动调用Hook函数 -- OnConnStart
	OnConnStart func(connection iface.IConnection)
	// 该Server销毁链接之前自动调用Hook函数 -- OnConnStop
	OnConnStop func(connection iface.IConnection)
}

func (s *Server) Start() {
	fmt.Printf("[Shiva] ServerName: %s, listen at ip: %s, port: %d\n", utils.GlobalObject.Name, utils.GlobalObject.Host, utils.GlobalObject.TcpPort)
	fmt.Printf("[Shiva] Version: %s, maxconn: %d, maxpackagesize: %d\n", utils.GlobalObject.Version, utils.GlobalObject.MaxConn, utils.GlobalObject.MaxPackageSize)
	fmt.Printf("[Start] Server Listener at IP: %s, port: %d\n", s.IP, s.Port)

	go func() {
		// 开启消息队列以及worker工作池
		s.MsgHandler.StartWorkerPool()

		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))

		if err != nil {
			fmt.Println("resolve tcp addr error: ", err)
			return
		}

		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen: ", s.IPVersion, "err: ", err)
			return
		}

		fmt.Println("start server succ, listening")

		var cid uint32
		cid = 0

		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("accept err: ", err)
				continue
			}

			// 设置最大连接个数的判断, 如果超过最大连接, 那么关闭次新连接
			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				fmt.Println("too many conn")
				conn.Close()
				continue
			}

			// 将conn和connection绑定
			dealConn := NewConnection(s, conn, cid, s.MsgHandler)
			cid++

			dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {
	// 将一些服务器的资源，状态或者一些已经开辟的连接信息进行停止后者回收
	fmt.Println("[STOP] shiva stop ", s.Name)
	s.ConnMgr.ClearConn()
}

func (s *Server) Serve() {
	s.Start()

	// 做一些启动服务器之后的额外业务

	select {}
}

func (s *Server) AddRouter(msgID uint32, router iface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
}

// 初始化Server模块的方法
func NewServer(name string) iface.IServer {
	s := &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		MsgHandler: NewMsgHandle(),
		ConnMgr:    NewConnManager(),
	}

	return s
}

func (s *Server) GetConnMgr() iface.IConnManager {
	return s.ConnMgr
}

// 注册OnConnStart钩子函数方法
func (s *Server) SetOnConnStart(f func(iface.IConnection)) {
	s.OnConnStart = f
}

// 注册OnConnStop钩子函数方法
func (s *Server) SetOnConnStop(f func(iface.IConnection)) {
	s.OnConnStop = f
}

// 调用OnConnStart钩子函数
func (s *Server) CallOnConnStart(connection iface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("call onConnStart!!!")
		s.OnConnStart(connection)
	}
}

// 调用OnConnStop钩子函数
func (s *Server) CallOnConnStop(connection iface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("call onConnStart!!!")
		s.OnConnStop(connection)
	}
}
