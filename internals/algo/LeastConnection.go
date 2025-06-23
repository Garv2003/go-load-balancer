package algo

import (
	"fmt"
	"github.com/garv2003/go-load-balancer/internals/models"
	"net/url"
)

type LeastConnection struct {
}

func (lc *LeastConnection) GetServer(servers []*models.Server) (url.URL, error) {
	if len(servers) == 0 {
		fmt.Println("there is no servers in server list!!!")
		return url.URL{}, nil
	}

	minServer := servers[0]

	for _, v := range servers {
		if minServer.Connection > v.Connection {
			if v.IsAlive {
				minServer = v
			}
		}
	}

	return minServer.ServerUrl, nil
}
