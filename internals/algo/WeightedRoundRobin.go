package algo

import (
	"context"
	"errors"
	"github.com/garv2003/go-load-balancer/internals/models"
	"net/url"
	"sync"
)

type WeightedRoundRobin struct {
	currentIndex       int
	currentWeightCount float64
	mu                 sync.Mutex
}

func (rr *WeightedRoundRobin) GetServer(_ context.Context, servers []*models.Server) (url.URL, error) {
	rr.mu.Lock()
	defer rr.mu.Unlock()

	if len(servers) == 0 {
		return url.URL{}, errors.New("no servers in server list")
	}

	var aliveServers []*models.Server
	for _, s := range servers {
		if s.IsAlive && s.Weight > 0 {
			aliveServers = append(aliveServers, s)
		}
	}

	if len(aliveServers) == 0 {
		return url.URL{}, errors.New("no alive servers available")
	}

	if rr.currentIndex >= len(aliveServers) {
		rr.currentIndex = 0
		rr.currentWeightCount = 0
	}

	selected := aliveServers[rr.currentIndex]

	if rr.currentWeightCount < selected.Weight {
		rr.currentWeightCount++
		return selected.ServerUrl, nil
	} else {
		rr.currentIndex = (rr.currentIndex + 1) % len(aliveServers)
		rr.currentWeightCount = 1
		selected = aliveServers[rr.currentIndex]
		return selected.ServerUrl, nil
	}
}
