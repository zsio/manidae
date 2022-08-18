package client

import (
	"github.com/gorilla/websocket"
	"log"
	"time"
)

type tunnelConn struct {
	tun *websocket.Conn
	key string
	url string
}

// Start 创建一个tunnel连接
func Start() {

	dial, _, err := websocket.DefaultDialer.Dial("ws://127.0.0.1:2022/connect?key=12022", nil)
	if err != nil {
		log.Fatalf("dial error: %s", err)
		return
	}

	t := &tunnelConn{
		tun: dial,
	}
	// start read
	go t.onMsg()

}

func (t *tunnelConn) onMsg() {
	defer func() {
		_ = t.tun.Close()
	}()

	t.tun.SetPingHandler(func(string) error {
		log.Printf("ping\n")
		return t.tun.WriteControl(websocket.PongMessage, []byte{}, time.Now().Add(time.Second))
	})

	t.tun.SetPongHandler(func(string) error {
		log.Printf("pong\n")
		return nil
	})

	// write loop for tunnel
	go func() {
		for {
			err := t.tun.WriteMessage(websocket.PingMessage, []byte("ping"))
			if err != nil {
				log.Printf("pong write error: %s", err)
				break
			}
			t.tun.WriteMessage(websocket.TextMessage, []byte("hello server"))
			time.Sleep(time.Second)
		}

	}()

	for {
		_, msg, err := t.tun.ReadMessage()
		if err != nil {
			log.Printf("read error: %s", err)
			return
		}
		log.Printf("msg: %s", msg)
	}

}
