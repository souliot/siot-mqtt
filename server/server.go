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

func (m *Server) AddClient(id string, c *fetcp.Conn) {
	// 如果id已经在线，强制断开连接
	m.mutex.Lock()
	if m.clientList[id] != nil && m.clientList[id] != c {
		m.clientList[id].Close()
	}
	m.clientList[id] = c
	m.mutex.Unlock()
}

func (m *Server) DelClient(id string) {
	m.mutex.Lock()
	delete(m.clientList, id)
	m.mutex.Unlock()
}

func (m *Server) GetClientList() (ls map[string]*fetcp.Conn) {
	return m.clientList
}
