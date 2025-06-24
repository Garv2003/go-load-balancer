package models

import (
	"fmt"
	"log"
	"net/url"
	"sync"
)

type ServerPool struct {
	Servers []*Server
	mu      sync.Mutex
}

func (sp *ServerPool) AddServer(serverUrl string) {
	parseUrl, err := url.Parse(serverUrl)

	if err != nil {
		fmt.Println("error parsing server url!!")
		return
	}

	sp.Servers = append(sp.Servers, &Server{ServerUrl: *parseUrl})
}

func (sp *ServerPool) HealthCheck() {
	for _, b := range sp.Servers {
		status := "up"
		alive := b.IsServerAlive()
		b.SetIsAlive(alive)
		if !alive {
			status = "down"
		}
		log.Printf("%v [%s]\n", b.ServerUrl, status)
	}
}

func (sp *ServerPool) ClearServers() {
	sp.mu.Lock()
	defer sp.mu.Unlock()
	sp.Servers = []*Server{}
}
