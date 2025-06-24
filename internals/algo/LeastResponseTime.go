package algo

import (
	"context"
	"errors"
	"github.com/garv2003/go-load-balancer/internals/models"
	"net/url"
	"time"
)

type LeastResponseTime struct{}

func (lrt *LeastResponseTime) UpdateServerAvgTime(servers []*models.Server, serverUrl url.URL, t time.Duration) {
	for _, s := range servers {
		if s.ServerUrl == serverUrl {
			alpha := 0.7
			newTime := t.Seconds()
			s.AvgTime = alpha*s.AvgTime + (1-alpha)*newTime
			break
		}
	}
}

func (lrt *LeastResponseTime) GetServer(_ context.Context, servers []*models.Server) (url.URL, error) {
	if len(servers) == 0 {
		return url.URL{}, errors.New("no servers in servers list")
	}

	var selected *models.Server
	var minTime float64 = -1

	for _, s := range servers {
		if !s.IsAlive {
			continue
		}

		if selected == nil || s.AvgTime < minTime {
			selected = s
			minTime = s.AvgTime
		}
	}

	if selected == nil {
		return url.URL{}, errors.New("no alive servers available")
	}

	return selected.ServerUrl, nil
}
