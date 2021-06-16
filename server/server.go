package server

import (
	"sync"

	"github.com/souliot/fetcp"
)

type Server struct {
	srv        *fetcp.Server
	clientList map[string]*fetcp.Conn
	mutex      *sync.Mutex
}

func NewServer(opts ...fetcp.SrvOption) (srv *Server) {
	srv = new(Server)
	srv.srv = fetcp.NewServer(newCallback(srv), new(Protocol), opts...)
	srv.clientList = make(map[string]*fetcp.Conn)
	srv.mutex = new(sync.Mutex)
	srv.Start()
	return
}

func (m *Server) Start() {
	go m.srv.Server()
}

func (m *Server) Stop() {
	m.srv.Stop()
	for _, c := range m.clientList {
		c.Close()
	}
}

func (m *Server) AddClient(id string, c *fetcp.Conn, new bool) {
	// 如果id已经在线，强制断开连接
	m.mutex.Lock()
	conn, ok := m.clientList[id]
	m.mutex.Unlock()
	if ok && new && c != conn {
		conn.Close()
		goto NEXT
		return
	}
	// 如果id已经在线，重用之前的连接
	if ok && !new && c != conn {
		c = conn
		return
	}
NEXT:
	m.mutex.Lock()
	m.clientList[id] = c
	m.mutex.Unlock()
	return
}

func (m *Server) DelClient(id string) {
	m.mutex.Lock()
	delete(m.clientList, id)
	m.mutex.Unlock()
}

func (m *Server) GetClientList() (ls map[string]*fetcp.Conn) {
	return m.clientList
}
