package models

import (
	"net/url"
	"sync"
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
