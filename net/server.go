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
}

func (s *Server) Start() {
	fmt.Printf("[Zinx] ServerName: %s, listen at ip: %s, port: %d\n", utils.GlobalObject.Name, utils.GlobalObject.Host, utils.GlobalObject.TcpPort)
	fmt.Printf("[Zinx] Version: %s, maxconn: %d, maxpackagesize: %d\n", utils.GlobalObject.Version, utils.GlobalObject.MaxConn, utils.GlobalObject.MaxPackageSize)
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

			// 将conn和connection绑定
			dealConn := NewConnection(conn, cid, s.MsgHandler)
			cid++

			dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {
	// 将一些服务器的资源，状态或者一些已经开辟的连接信息进行停止后者回收
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
	}

	return s
}
