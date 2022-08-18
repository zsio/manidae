package server

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

var wsPort = 2022

// WsClient 已经连接的websocket客户端
type WsClient struct {
	Tunnel *websocket.Conn
	Key    string
}

// WsServer websocket服务管理器
type WsServer struct {
	Clients map[string]*WsClient
}

// 根据key获取一个websocket客户端
func (ws *WsServer) getClient(key string) *WsClient {
	return ws.Clients[key]
}

// 添加一个websocket客户端
func (ws *WsServer) addClient(key string, client *WsClient) {
	ws.Clients[key] = client
}

// 删除一个websocket客户端
func (ws *WsServer) removeClient(key string) {
	delete(ws.Clients, key)
}

// StartWS 启动一个websocket服务器
func (ws *WsServer) StartWS(shutdown *chan struct{}) {
	upgrade := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	http.HandleFunc("/connect", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		key := query.Get("key")
		if key == "" {
			log.Printf("key is empty")
			return
		}

		// 判断是否已经连接过
		if _, ok := ws.Clients[key]; ok {
			log.Printf("key: %s already connected\n", key)
			return
		}

		conn, err := upgrade.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("ws upgrade error: %s\n", err)
			return
		}
		// 添加到客户端列表
		ws.addClient(key, &WsClient{
			Tunnel: conn,
			Key:    key,
		})

		// 启动一个goroutine处理websocket
		go ws.handleWSConnection(conn, key)
	})

	err := http.ListenAndServe(fmt.Sprintf(":%d", wsPort), nil)
	if err != nil {
		*shutdown <- struct{}{}
		return
	}
}

// handleWSConnection handle websocket connection
func (ws *WsServer) handleWSConnection(conn *websocket.Conn, key string) {
	defer func() {
		ws.removeClient(key)
		_ = conn.Close()
	}()

	conn.SetPingHandler(func(string) error {
		log.Printf("ping\n")
		return conn.WriteControl(websocket.PongMessage, []byte{}, time.Now().Add(time.Second))
	})

	conn.SetPongHandler(func(string) error {
		log.Printf("pong\n")
		return nil
	})

	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		fmt.Printf("msgType: %d\n", msgType)
		fmt.Printf("msg: %s\n", msg)
		fmt.Printf("url: %s\n", conn.RemoteAddr().String())

		if msgType == websocket.PingMessage {
			_ = conn.WriteMessage(websocket.PongMessage, []byte("pong"))
		} else {
			_ = conn.WriteMessage(msgType, msg)
		}
	}
}
