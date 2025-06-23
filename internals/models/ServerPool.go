package models

import (
	"fmt"
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
