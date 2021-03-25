package net

import (
	"fmt"
	"io"
	"net"
	"testing"
	"time"
)

// 只是负责测试datapack拆包,封包的单元测试
func TestDataPack(t *testing.T) {

	listenner, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("server listen err, ", err)
		return
	}

	go func() {
		for {
			conn, err := listenner.Accept()
			if err != nil {
				fmt.Println("server accept error, ", err)
			}

			go func(conn net.Conn) {
				//  处理客户端的请求
				// 1. 第一次读conn, 把包的head读出来
				dp := NewDataPack()
				for {
					headData := make([]byte, dp.GetHeadLen())
					fmt.Println("before readfull headata")
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						fmt.Println("read head error")
						return
					}
					fmt.Println("after readfull headata")
					// msgHead包含dataLen + Id
					msgHead, err := dp.Unpack(headData)
					if err != nil {
						fmt.Println("server unpack err, ", err)
						return
					}
					if msgHead.GetMsgLen() > 0 {
						// msg是有数据的, 需要进行第二次读取
						// 2. 第二次读conn, 把data内容读出来
						msg := msgHead.(*Message)
						msg.Data = make([]byte, msg.GetMsgLen())

						// 根据dataLen长度再次从io流中读取
						fmt.Println("before readfull msgdata")
						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println("server unpack data err, ", err)
							return
						}
						fmt.Println("after readfull msgdata")
						// 完整的消息读取完毕
						fmt.Println("recv msg, id: ", msg.GetMsgId(), ", datalen: ", msg.GetMsgLen(), "data: ", string(msg.GetData()))
					}
				}
			}(conn)
		}
	}()

	// 模拟客户端
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("dial err, ", err)
		return
	}

	// 创建一个封包对象dp
	dp := NewDataPack()
	// 模拟粘包, 两个包一起发送
	msg1 := &Message{
		Id:      1,
		DataLen: 4,
		Data:    []byte{'z', 'i', 'n', 'x'},
	}

	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("pack1 fail, ", err)
		return
	}

	sendDataLen := len(sendData1)
	conn.Write(sendData1[:4])
	time.Sleep(2 * time.Second)
	conn.Write(sendData1[4:8])
	time.Sleep(2 * time.Second)
	conn.Write(sendData1[8:sendDataLen])

	//msg2 := &Message{
	//	Id:      2,
	//	DataLen: 6,
	//	Data:    []byte{'h', 'e', 'l', 'l', 'l', 'o'},
	//}
	//
	//sendData2, err := dp.Pack(msg2)
	//if err != nil {
	//	fmt.Println("pack2 fail, ", err)
	//	return
	//}
	//
	//sendData1 = append(sendData1, sendData2...)
	//conn.Write(sendData1)

	select {}
}
