package algo

import (
	"context"
	"errors"
	"github.com/garv2003/go-load-balancer/internals/models"
	"net/url"
)

type LeastConnection struct{}

func (lc *LeastConnection) GetServer(_ context.Context, servers []*models.Server) (url.URL, error) {
	if len(servers) == 0 {
		return url.URL{}, errors.New("no servers in server list")
	}

	var minServer *models.Server

	for _, server := range servers {
		if !server.IsAlive {
			continue
		}

		if minServer == nil || server.Connection < minServer.Connection {
			minServer = server
		}
	}

	if minServer == nil {
		return url.URL{}, errors.New("no alive servers available")
	}

	return minServer.ServerUrl, nil
}
