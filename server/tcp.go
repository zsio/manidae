package server

import (
	"fmt"
	"log"
	"net"
	"time"
)

type TCPServer struct{}

// StartTCP starts a tcp server
func (t *TCPServer) StartTCP() {
	addr, err := net.ResolveTCPAddr("tcp", ":12022")
	if err != nil {
		return
	}
	conn, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatalf("listen tcp error: %s", err)
		return
	}
	defer conn.Close()
	for {
		conn, err := conn.Accept()
		if err != nil {
			log.Fatalf("accept error: %s", err)
			return
		}
		go t.handleTCPConnection(conn)
	}

}

// handleTCPConnection handles a tcp connection
func (t *TCPServer) handleTCPConnection(conn net.Conn) {
	defer func() {
		_ = conn.Close()
	}()
	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			break
		}
		fmt.Printf("buf: %s\n", buf[:n])
		time.Sleep(time.Second)
		_, err = conn.Write(buf[:n])
		if err != nil {
			break
		}
	}
}
