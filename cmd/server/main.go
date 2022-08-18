package main

import "manidae/server"

func main() {
	// 程序退出信号
	shutdown := make(chan struct{}, 1)

	wsServer := server.WsServer{
		Clients: make(map[string]*server.WsClient),
	}
	go wsServer.StartWS(&shutdown)

	tcpServer := server.TCPServer{}
	go tcpServer.StartTCP()

	<-shutdown
}
