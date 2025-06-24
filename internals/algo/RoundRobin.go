package algo

import (
	"context"
	"errors"
	"github.com/garv2003/go-load-balancer/internals/models"
	"net/url"
	"sync/atomic"
)

type RoundRobin struct {
	count atomic.Int64
}

func (rr *RoundRobin) GetCount() int {
	return int(rr.count.Load())
}

func (rr *RoundRobin) IncrementCount() {
	rr.count.Add(1)
}

func (rr *RoundRobin) GetServer(_ context.Context, servers []*models.Server) (url.URL, error) {
	if len(servers) == 0 {
		return url.URL{}, errors.New("no servers in list")
	}

	// Filter alive servers
	var aliveServers []*models.Server
	for _, server := range servers {
		if server.IsAlive {
			aliveServers = append(aliveServers, server)
		}
	}

	if len(aliveServers) == 0 {
		return url.URL{}, errors.New("no alive servers available")
	}

	index := rr.GetCount()
	selected := aliveServers[index%len(aliveServers)]

	return selected.ServerUrl, nil
}
