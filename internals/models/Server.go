package models

import (
	"log"
	"net"
	"net/url"
	"sync"
	"time"
)

type Server struct {
	ServerUrl  url.URL
	Connection int64
	AvgTime    int64
	IsAlive    bool
	Mut        sync.Mutex
}

func (s *Server) GetAlive() bool {
	s.Mut.Lock()
	defer s.Mut.Unlock()
	return s.IsAlive
}

func (s *Server) SetIsAlive(isAlive bool) {
	s.Mut.Lock()
	defer s.Mut.Unlock()
	s.IsAlive = isAlive
}

func (s *Server) GetServerUrl() url.URL {
	s.Mut.Lock()
	defer s.Mut.Unlock()
	return s.ServerUrl
}

func (s *Server) IsServerAlive() bool {
	timeout := 2 * time.Second
	conn, err := net.DialTimeout("tcp", s.ServerUrl.Host, timeout)
	if err != nil {
		log.Println("server unreachable, error: ", err)
		return false
	}
	_ = conn.Close()
	return true
}
