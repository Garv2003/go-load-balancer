package algo

import (
	"context"
	"errors"
	"github.com/garv2003/go-load-balancer/internals/models"
	"math/rand"
	"net/url"
	"sync"
	"time"
)

type Random struct {
	rng  *rand.Rand
	lock sync.Mutex
}

func NewRandom() *Random {
	return &Random{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (r *Random) GetServer(_ context.Context, servers []*models.Server) (url.URL, error) {
	if len(servers) == 0 {
		return url.URL{}, errors.New("no servers in list")
	}

	var aliveServers []*models.Server
	for _, s := range servers {
		if s.IsAlive {
			aliveServers = append(aliveServers, s)
		}
	}

	if len(aliveServers) == 0 {
		return url.URL{}, errors.New("no alive servers available")
	}

	r.lock.Lock()
	index := r.rng.Intn(len(aliveServers))
	r.lock.Unlock()

	return aliveServers[index].ServerUrl, nil
}
