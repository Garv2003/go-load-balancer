package models

import (
	"fmt"
	"log"
	"net/url"
)

type ServerPool struct {
	Servers []*Server
}

func (serverPool *ServerPool) AddServer(serverUrl string) {
	parseUrl, err := url.Parse(serverUrl)

	if err != nil {
		fmt.Println("error parsing server url!!")
		return
	}

	serverPool.Servers = append(serverPool.Servers, &Server{ServerUrl: *parseUrl})
}

func (serverPool *ServerPool) HealthCheck() {
	for _, b := range serverPool.Servers {
		status := "up"
		alive := b.IsServerAlive()
		b.SetIsAlive(alive)
		if !alive {
			status = "down"
		}
		log.Printf("%s [%s]\n", b.ServerUrl, status)
	}
}
