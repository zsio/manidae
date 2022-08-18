package main

import (
	"fmt"
	"io"
	"log"
	"manidae/server"
	"net"
)

func main() {
	// create a new server
	connects := server.NewConnectManager()

	// 启动远程客户端监听，等待客户端连接
	go ListenRemoteClient(connects)

	go forward("127.0.0.0:9000", connects)

	// 等待程序退出
	shutdown := make(chan struct{})
	<-shutdown

}

func ListenRemoteClient(connects *server.ConnectManager) {

	ln, err := net.Listen("tcp", ":2022")
	if err != nil {
		fmt.Println("本地监听端口时异常", err)
	}
	for {
		clientConn, err := ln.Accept()
		if err != nil {
			fmt.Println("建立与客户端连接时异常", err)
			return
		}
		go handleRemoteClientConn(clientConn, connects)
	}
}

func handleRemoteClientConn(conn net.Conn, connects *server.ConnectManager) {
	defer func() {
		log.Println("关闭连接", conn.RemoteAddr())
		connects.Remove(&conn)
		_ = conn.Close()
	}()

	// add client to connects
	connects.Add(&conn)

	for {
		// 读取客户端数据
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("读取客户端数据时异常", err)
			return
		}
		fmt.Println("接收到客户端数据：", string(buf[:n]))
	}
}

func pipe(src net.Conn, dest net.Conn) {
	errChan := make(chan error, 1)
	onClose := func(err error) {
		_ = dest.Close()
		_ = src.Close()
	}
	go func() {
		_, err := io.Copy(src, dest)
		errChan <- err
		onClose(err)
	}()
	go func() {
		_, err := io.Copy(dest, src)
		errChan <- err
		onClose(err)
	}()
	<-errChan
}

// 监听一个tcp端口，转发数据到远程客户端
func forward(addr string, connects *server.ConnectManager) {
	dial, err := net.Dial("tcp", addr)
	if err != nil {
		return
	}
	defer func(dial net.Conn) {
		_ = dial.Close()
	}(dial)

	pipe(*connects.GetFirst().Conn, dial)
}
