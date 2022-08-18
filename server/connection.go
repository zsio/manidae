package server

import (
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"
)

// 连接的客户端信息
type Connect struct {
	id       string
	ip       net.IP
	Key      string
	Conn     *net.Conn
	createAt time.Time
}

type ConnectManager struct {
	connects map[string]*Connect
}

// NewConnectManager 构造函数
func NewConnectManager() *ConnectManager {
	return &ConnectManager{
		connects: make(map[string]*Connect),
	}
}

// Add 添加连接
func (m *ConnectManager) Add(conn *net.Conn) *Connect {
	remoteAddr := (*conn).RemoteAddr().(*net.TCPAddr)
	ipStr := remoteAddr.String()

	// 从IP:port中提取IP
	ip, _ := url.QueryUnescape(ipStr[:strings.Index(ipStr, ":")])

	c := &Connect{
		// 临时ID，以后需要分配一个唯一的ID
		id:       ipStr,
		ip:       net.ParseIP(ip),
		Conn:     conn,
		createAt: time.Now(),
	}

	fmt.Println("添加连接", ipStr)
	m.connects[c.id] = c
	return c
}

// Remove 移除连接
func (m *ConnectManager) Remove(conn *net.Conn) {
	delete(m.connects, (*conn).RemoteAddr().String())
}

// GetFirst get first connect
func (m *ConnectManager) GetFirst() *Connect {
	for _, c := range m.connects {
		return c
	}
	return nil
}
