package algo

import (
	"context"
	"errors"
	"github.com/garv2003/go-load-balancer/internals/models"
	"net/url"
	"sync/atomic"
	"time"
)

type WeightedResponseTime struct {
	count atomic.Int64
}

func (wrt *WeightedResponseTime) UpdateServerAvgTime(servers []*models.Server, serverUrl url.URL, t time.Duration) {
	for _, s := range servers {
		if s.ServerUrl == serverUrl {
			alpha := 0.7
			newTime := t.Seconds()
			s.AvgTime = alpha*s.AvgTime + (1-alpha)*newTime
			break
		}
	}
}

func (wrt *WeightedResponseTime) GetServer(_ context.Context, servers []*models.Server) (url.URL, error) {
	if len(servers) == 0 {
		return url.URL{}, errors.New("no servers in server list")
	}

	var selected *models.Server
	minScore := float64(-1)

	for _, s := range servers {
		if !s.IsAlive || s.Weight <= 0 {
			continue
		}

		score := s.AvgTime / s.Weight

		if selected == nil || score < minScore {
			selected = s
			minScore = score
		}
	}

	if selected == nil {
		return url.URL{}, errors.New("no alive servers with positive weight")
	}

	return selected.ServerUrl, nil
}
