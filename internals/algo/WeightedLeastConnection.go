package algo

import (
	"context"
	"errors"
	"github.com/garv2003/go-load-balancer/internals/models"
	"net/url"
)

type WeightedLeastConnection struct{}

func (wlc *WeightedLeastConnection) GetServer(_ context.Context, servers []*models.Server) (url.URL, error) {
	if len(servers) == 0 {
		return url.URL{}, errors.New("no servers available")
	}

	var selected *models.Server
	minRatio := float64(-1)

	for _, server := range servers {
		if !server.IsAlive || server.Weight <= 0 {
			continue
		}

		connCount := server.Connection
		ratio := float64(connCount) / float64(server.Weight)

		if selected == nil || ratio < minRatio {
			selected = server
			minRatio = ratio
		}
	}

	if selected == nil {
		return url.URL{}, errors.New("no alive servers with valid weight")
	}

	selected.IncrementConnection()

	return selected.ServerUrl, nil
}

func (wlc *WeightedLeastConnection) DecrementConnection(servers []*models.Server, server url.URL) {
	for _, v := range servers {
		if v.ServerUrl == server {
			v.DecrementConnection()
		}
	}
}
