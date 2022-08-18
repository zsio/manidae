package main

import (
	"io"
	"log"
	"net"
)

type client struct {
	conn net.Conn
}

func main() {
	// 等待程序退出
	shutdown := make(chan struct{})

	c := client{}

	go c.connect("127.0.0.1:2022")

	<-shutdown
}

// 和服务端建立连接
func (c *client) connect(addr string) {
	dial, err := net.Dial("tcp", addr)
	if err != nil {
		return
	}
	defer func(dial net.Conn) {
		_ = dial.Close()
	}(dial)
	c.conn = dial
	go c.forward("127.0.0.1:8000")
}

// 转发服务端数据到另一个tcp连接
func (c *client) forward(addr string) {
	dial, err := net.Dial("tcp", addr)
	if err != nil {
		return
	}
	defer func(dial net.Conn) {
		_ = dial.Close()
	}(dial)
	pipe(c.conn, dial)
}

// 管道数据
func pipe(src, dest net.Conn) {
	errChan := make(chan error, 1)
	maxSize := int64(1024) * 64
	defer func() {
		log.Printf("关闭连接 %s", dest.RemoteAddr())
	}()
	go func() {
		for {
			n, err := io.CopyN(src, dest, maxSize)
			log.Printf("in copy %d bytes", n)
			errChan <- err
		}
	}()
	go func() {
		for {
			n, err := io.CopyN(dest, src, maxSize)
			log.Printf("out copy %d bytes", n)
			errChan <- err
		}
	}()
	<-errChan
}
